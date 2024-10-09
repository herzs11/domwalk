/*
Copyright © 2024 Seth Herz sherz@unum.com
*/
package cmd

import (
	"os"
	"time"

	"domwalk/db"
	"domwalk/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	domainsToExecute []string
	enrichCfg        enrichmentConfig
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "domwalk",
	Short: "CLI tool to find and store domain relationships",
	Long: `domwalk is a CLI tool to find and store domain relationships.
	It is written in Go and uses GORM and a local SQLite backend.
	Data is stored in a local SQLite database and can be pushed to BigQuery.

	Currently, the tool can enrich domains with the following relationships:
	- Certificate Subject Alternative Names (SANs)
	- Web Redirects
	- Sitemap Web Domains
	- Sitemap Contact Page Domains

	The tool can also enrich domains with DNS data. In a future version, this dns data will be used to form additional domain relationships
	`,
	Example: `domwalk -d unum.com,coloniallife.com --workers 20 --cert-sans --web-redirects --sitemaps-web --sitemaps-contact --dns --gorm-db domwalk.db`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		gormDBName, _ := cmd.Flags().GetString("gorm-db")
		if gormDBName == "" {
			cmd.Println("GORM SQLite database name is required (use --gorm-db or set $GORM_SQLITE_NAME)")
			os.Exit(1)
		}
		err := db.GormDBConnect(gormDBName)
		if err != nil {
			color.Set(color.FgRed)
			cmd.Println(err)
			color.Unset()
			os.Exit(1)
		}
		types.CreateTables()
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		domainsToExecute, _ = cmd.Flags().GetStringSlice("domains")
		fname, _ := cmd.Flags().GetString("file")
		if fname != "" {
			var err error
			domainsToExecute, err = readDomainsFromFile(fname)
			if err != nil {
				color.Set(color.FgRed)
				cmd.Println(err)
				color.Unset()
				os.Exit(1)
			}
			if len(domainsToExecute) == 0 {
				color.Set(color.FgRed)
				cmd.Println("No domains found in file")
				color.Unset()
				os.Exit(1)
			}
			if header, _ := cmd.Flags().GetBool("header"); header {
				domainsToExecute = domainsToExecute[1:]
			}
		}
		if foundDomainsOnly, _ := cmd.Flags().GetBool("found-domains-only"); !foundDomainsOnly && len(domainsToExecute) == 0 {
			color.Set(color.FgRed)
			cmd.Println("No domains provided, use --domains, --file, or --found-domains-only")
			color.Unset()
			os.Exit(1)
		}
		cs, _ := cmd.Flags().GetBool("cert-sans")
		wr, _ := cmd.Flags().GetBool("web-redirects")
		sw, _ := cmd.Flags().GetBool("sitemaps-web")
		sc, _ := cmd.Flags().GetBool("sitemaps-contact")
		dns, _ := cmd.Flags().GetBool("dns")
		workers, _ := cmd.Flags().GetInt("workers")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		if workers < 1 {
			color.Red("Workers must be greater than 0\n")
			os.Exit(1)
		}
		if limit < 1 {
			color.Red("Limit must be greater than 0\n")
			os.Exit(1)
		}
		if offset > len(domainsToExecute) && len(domainsToExecute) > 0 {
			color.Red("Offset must be less than the number of domains\n")
			os.Exit(1)
		}
		minDate, _ := cmd.Flags().GetString("min-freshness")
		staleDate, err := time.Parse(time.DateOnly, minDate)
		if err != nil {
			color.Red("Invalid date format for min-freshness: (YYYY-MM-DD)\n")
			os.Exit(1)
		}
		if !cs && !wr && !sw && !sc && !dns {
			cs = true
			wr = true
			sw = true
			sc = true
			dns = true
		}
		enrichCfg = enrichmentConfig{
			CertSans:         cs,
			DNS:              dns,
			SitemapWeb:       sw,
			SitemapContact:   sc,
			WebRedirect:      wr,
			Limit:            limit,
			Offset:           offset,
			NWorkers:         workers,
			MinFreshnessDate: staleDate,
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		if len(domainsToExecute) > 0 {
			enrichDomainNames(domainsToExecute, enrichCfg)
		} else {
			enrichDBDomains(enrichCfg)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSliceP("domains", "d", []string{}, "Domains to process")
	rootCmd.Flags().String("file", "", "File with domains to process")
	rootCmd.Flags().BoolP("found-domains-only", "f", false, "Only process domains from database")

	_ = rootCmd.MarkFlagFilename("file", "csv", "txt")
	rootCmd.MarkFlagsOneRequired("domains", "file", "found-domains-only")
	rootCmd.MarkFlagsMutuallyExclusive("domains", "file", "found-domains-only")

	rootCmd.Flags().Bool("header", true, "File with domains to process has a header row")

	rootCmd.PersistentFlags().String(
		"gorm-db", os.Getenv("GORM_SQLITE_NAME"),
		"GORM SQLite database name, can also set 'GORM_SQLITE_NAME' environment variable",
	)

	// Get minimum date to refresh
	rootCmd.Flags().String("min-freshness", "0001-01-01", "Minimum date to refresh relationships, (YYYY-MM-DD)")

	rootCmd.Flags().Bool("cert-sans", false, "Enrich domains with cert SANs")
	rootCmd.Flags().Bool("web-redirects", false, "Enrich domains with web redirects")
	rootCmd.Flags().Bool("sitemaps-web", false, "Enrich domains with sitemap web domains")
	rootCmd.Flags().Bool("sitemaps-contact", false, "Enrich domains with sitemap contact page scraped domains")
	rootCmd.Flags().Bool("dns", false, "Enrich domains with dns data")

	rootCmd.Flags().IntP("workers", "w", 15, "Number of concurrent workers to use")
	rootCmd.Flags().IntP("limit", "l", 3000, "Limit of domains to process")
	rootCmd.Flags().IntP("offset", "s", 0, "Offset of domains to process")
}
