package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"domwalk/db"
	"domwalk/domains"
	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/idtoken"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const SNAPSHOT_JOB_URL = "https://snapshot-domwalk-593149879404.us-central1.run.app"

var dataset *bigquery.Dataset

type syncConfig struct {
	Dataset            string
	Domains            bool
	CertSansDomains    bool
	DNS                bool
	WebRedirectDomains bool
	Sitemaps           bool
}

func backupFile(filePath string) error {
	// Get the directory, filename, and extension of the original file.
	dir, file := filepath.Split(filePath)
	ext := filepath.Ext(file)
	fileName := strings.TrimSuffix(filepath.Base(file), ext)
	// Construct the backup file path.
	backupFilePath := filepath.Join(dir, fmt.Sprintf("%s_bak%s", fileName, ext))

	// if file does not exist, do not backup
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		color.Yellow("File does not exist, skipping backup...\n")
	}
	// Open the original file for reading.
	src, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer src.Close()

	// Create the backup file for writing.
	dst, err := os.Create(backupFilePath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer dst.Close()

	// Copy the contents of the original file to the backup file.
	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}

func pullFromBQ(cfg syncConfig) {

	fmt.Println("Pulling data from BigQuery (this takes a while)...")
	domains.ClearTables()
	domains.CreateTables()
	var pb *progressbar.ProgressBar
	if cfg.Domains {
		doms := []domains.Domain{}
		color.Green("Pulling domains")
		q := db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.domains ORDER BY domain_name", cfg.Dataset))
		rows, err := q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.DomainBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			dom := d.ToGorm()
			doms = append(doms, dom)
		}
		color.Yellow("Creating %d rows", len(doms))
		pb = progressbar.Default(int64(len(doms)), "Creating domains")
		for _, d := range doms {
			err = db.GormDB.Save(&d).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save domain: %v", err)
			}

		}
	}
	if cfg.CertSansDomains {
		csans := []domains.CertSansDomain{}
		color.Green("Pulling cert sans domains")
		q := db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.cert_sans_domains ORDER BY id", cfg.Dataset))
		rows, err := q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.MatchedDomainBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			csan := domains.CertSansDomain{MatchedDomain: d.ToGorm()}
			csans = append(csans, csan)
		}
		color.Yellow("Creating %d rows", len(csans))
		pb = progressbar.Default(int64(len(csans)), "Creating cert sans domains")
		for _, c := range csans {
			err = db.GormDB.Save(&c).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save cert sans domain: %v", err)
			}
		}
	}
	if cfg.DNS {
		mxs := []domains.MXRecord{}
		color.Green("Pulling DNS records")
		q := db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.mx_records ORDER BY id", cfg.Dataset))
		rows, err := q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.MXRecordBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			mxr := d.ToGorm()
			mxs = append(mxs, mxr)
		}
		color.Yellow("Creating %d rows", len(mxs))
		pb = progressbar.Default(int64(len(mxs)), "Creating mx records")
		for _, m := range mxs {
			err = db.GormDB.Save(&m).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save mx record: %v", err)
			}
		}

		ars := []domains.ARecord{}
		q = db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.a_records ORDER BY id", cfg.Dataset))
		rows, err = q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.ARecordBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			ar := d.ToGorm()
			ars = append(ars, ar)
		}
		color.Yellow("Creating %d rows", len(ars))
		pb = progressbar.Default(int64(len(ars)), "Creating a records")
		for _, a := range ars {
			err = db.GormDB.Save(&a).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save a record: %v", err)
			}
		}

		aars := []domains.AAAARecord{}
		q = db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.aaaa_records ORDER BY id", cfg.Dataset))
		rows, err = q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.AAAARecordBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			ar := d.ToGorm()
			aars = append(aars, ar)
		}
		color.Yellow("Creating %d rows", len(aars))
		pb = progressbar.Default(int64(len(aars)), "Creating aaaa records")
		for _, a := range aars {
			err = db.GormDB.Save(&a).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save aaaa record: %v", err)
			}
		}

		soas := []domains.SOARecord{}
		q = db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.soa_records ORDER BY id", cfg.Dataset))
		rows, err = q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.SOARecordBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			ar := d.ToGorm()
			soas = append(soas, ar)
		}
		color.Yellow("Creating %d rows", len(soas))
		pb = progressbar.Default(int64(len(soas)), "Creating soa records")
		for _, s := range soas {
			err = db.GormDB.Save(&s).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save soa record: %v", err)
			}
		}
	}

	if cfg.WebRedirectDomains {
		wred := []domains.WebRedirectDomain{}
		color.Green("Pulling web redirect domains")
		q := db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.web_redirect_domains ORDER BY id", cfg.Dataset))
		rows, err := q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.MatchedDomainBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			wr := domains.WebRedirectDomain{MatchedDomain: d.ToGorm()}
			wred = append(wred, wr)
		}
		color.Yellow("Creating %d rows", len(wred))
		pb = progressbar.Default(int64(len(wred)), "Creating web redirect records")
		for _, w := range wred {
			err = db.GormDB.Save(&w).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save web redirect domain: %v", err)
			}
		}
	}
	if cfg.Sitemaps {
		sms := []domains.Sitemap{}
		color.Green("Pulling sitemaps")
		q := db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.sitemaps ORDER BY id", cfg.Dataset))
		rows, err := q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.SitemapBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			s := d.ToGorm()
			sms = append(sms, s)
		}
		color.Yellow("Creating %d rows", len(sms))
		pb = progressbar.Default(int64(len(sms)), "Creating SitemapLoc records")
		for _, s := range sms {
			err = db.GormDB.Save(&s).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save sitemap: %v", err)
			}
		}

		smwd := []domains.SitemapWebDomain{}
		q = db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.sitemap_web_domains ORDER BY id", cfg.Dataset))
		rows, err = q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.MatchedDomainBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			swd := domains.SitemapWebDomain{MatchedDomain: d.ToGorm()}
			smwd = append(smwd, swd)
		}
		color.Yellow("Creating %d rows", len(smwd))
		pb = progressbar.Default(int64(len(smwd)), "Creating SitemapLoc web domain records")
		for _, s := range smwd {
			err = db.GormDB.Save(&s).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save sitemap web domain: %v", err)
			}
		}

		scds := []domains.SitemapContactDomain{}
		q = db.BQConn.Query(fmt.Sprintf("SELECT * FROM %s.sitemap_contact_domains ORDER BY id", cfg.Dataset))
		rows, err = q.Read(context.Background())
		if err != nil {
			log.Fatalf("Failed to read rows: %v", err)
		}
		for {
			var d domains.MatchedDomainBQ
			err := rows.Next(&d)
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to read row: %v", err)
			}
			scd := domains.SitemapContactDomain{MatchedDomain: d.ToGorm()}
			scds = append(scds, scd)
		}
		color.Yellow("Creating %d rows", len(scds))
		pb = progressbar.Default(int64(len(scds)), "Creating SitemapLoc contact records")
		for _, s := range scds {
			err = db.GormDB.Save(&s).Error
			pb.Add(1)
			if err != nil {
				log.Fatalf("Failed to save sitemap contact domain: %v", err)
			}
		}
	}
}

func pushToBQ(cfg syncConfig) {
	dataset = db.BQConn.Dataset(cfg.Dataset)
	defer db.BQConn.Close()
	var (
		offset    int
		limit     int
		tableName string
	)

	if cfg.Domains {
		tableName = "domains"
		if err := recreateTable(domains.DomainBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		offset = 0
		limit = 1000
		for {
			chunk := []domains.Domain{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			bqds := make([]domains.DomainBQ, len(chunk))
			for i, m := range chunk {
				bqds[i] = m.ToBQ()
			}
			loadToBigQuery(bqds, tableName)
			offset += limit
		}
	}
	if cfg.CertSansDomains {
		offset = 0
		limit = 5000
		tableName = "cert_sans_domains"
		if err := recreateTable(domains.MatchedDomainBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []domains.CertSansDomain{}
			bqr := []domains.MatchedDomainBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}
	}
	if cfg.DNS {
		tableName = "mx_records"
		offset = 0
		limit = 5000
		if err := recreateTable(domains.MXRecordBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		for {
			bqr := []domains.MXRecordBQ{}
			chunk := []domains.MXRecord{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}

		tableName = "a_records"
		offset = 0
		limit = 5000
		if err := recreateTable(domains.ARecordBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		for {
			bqr := []domains.ARecordBQ{}
			chunk := []domains.ARecord{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit

		}

		tableName = "aaaa_records"
		offset = 0
		limit = 1000
		if err := recreateTable(domains.AAAARecordBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			bqr := []domains.AAAARecordBQ{}
			chunk := []domains.AAAARecord{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}

		tableName = "soa_records"
		offset = 0
		limit = 1000
		if err := recreateTable(domains.SOARecordBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			bqr := []domains.SOARecordBQ{}
			chunk := []domains.SOARecord{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit

		}
	}
	if cfg.WebRedirectDomains {
		tableName = "web_redirect_domains"
		offset = 0
		limit = 1000
		if err := recreateTable(domains.MatchedDomainBQ{}, tableName); err != nil {
			log.Printf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []domains.WebRedirectDomain{}
			bqr := []domains.MatchedDomainBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit

		}
	}
	if cfg.Sitemaps {
		tableName = "sitemaps"
		offset = 0
		limit = 1000
		if err := recreateTable(domains.SitemapBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []domains.Sitemap{}
			bqr := []domains.SitemapBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}

		tableName = "sitemap_web_domains"
		offset = 0
		limit = 5000
		if err := recreateTable(domains.MatchedDomainBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []domains.SitemapWebDomain{}
			bqr := []domains.MatchedDomainBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit

		}

		tableName = "sitemap_contact_domains"
		offset = 0
		limit = 5000
		if err := recreateTable(domains.MatchedDomainBQ{}, tableName); err != nil {
			log.Fatalf("Failed to truncate table: %v", err)
		}
		for {
			chunk := []domains.SitemapContactDomain{}
			bqr := []domains.MatchedDomainBQ{}
			err := db.GormDB.Limit(limit).Offset(offset).Find(&chunk).Error
			if err != nil {
				log.Fatalf("Failed to get chunk: %v", err)
			}
			if len(chunk) == 0 {
				break
			}
			for _, m := range chunk {
				bqr = append(bqr, m.ToBQ())
			}
			loadToBigQuery(bqr, tableName)
			offset += limit
		}
	}
}

func loadToBigQuery(model interface{}, tableName string) {

	ctx := context.Background()
	table := dataset.Table(tableName)

	inserter := table.Inserter()
	for {
		if err := inserter.Put(ctx, model); err != nil {
			log.Printf("Failed to insert data into BigQuery table: %v, sleeping for 3 seconds\n", err)
			time.Sleep(3 * time.Second)

		} else {
			break
		}
	}

	fmt.Printf("Data loaded into BigQuery table: %s\n", table.FullyQualifiedName())
}

func recreateTable(model interface{}, tableName string) error {
	ctx := context.Background()
	table := dataset.Table(tableName)
	if err := table.Delete(ctx); err != nil {
		log.Printf("Failed to delete table: %v, trying to create\n", err)
	}
	time.Sleep(5 * time.Second)
	sch, err := bigquery.InferSchema(model)
	if err != nil {
		log.Printf("Failed to infer schema: %v", err)
		return err
	}
	metaData := &bigquery.TableMetadata{
		Schema: sch,
	}
	if err := table.Create(ctx, metaData); err != nil {
		log.Printf("Failed to create table: %v", err)
		return err
	}
	fmt.Printf("Table %s created\n", tableName)
	return nil
}

func snapshotDomains() error {

	// client, err := google.DefaultClient(context.Background(), "https://snapshot-domwalk-593149879404.us-central1.run.app")
	ctx := context.Background()

	// Construct the GoogleCredentials object which obtains the default configuration from your
	// working environment.
	credentials, err := google.FindDefaultCredentials(ctx)
	if err != nil {
		return fmt.Errorf("failed to generate default credentials: %w", err)
	}

	ts, err := idtoken.NewTokenSource(ctx, SNAPSHOT_JOB_URL, option.WithCredentials(credentials))
	if err != nil {
		return fmt.Errorf("failed to create NewTokenSource: %w", err)
	}

	// Get the ID token.
	// Once you've obtained the ID token, you can use it to make an authenticated call
	// to the target audience.
	tok, err := ts.Token()
	if err != nil {
		return fmt.Errorf("failed to receive token: %w", err)
	}
	req, err := http.NewRequest("GET", SNAPSHOT_JOB_URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tok.AccessToken))
	resp, err := http.DefaultClient.Do(req)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))
	return nil
}
