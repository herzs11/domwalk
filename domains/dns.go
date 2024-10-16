package domains

import (
	"errors"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/miekg/dns"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AAAARecord struct {
	gorm.Model
	IPV6       string `gorm:"uniqueIndex:aaaa_dom_idx" bigquery:"ip_v6"`
	DomainName string `gorm:"uniqueIndex:aaaa_dom_idx" bigquery:"domain_name"`
}

type AAAARecordBQ struct {
	ID         int       `bigquery:"id"`
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	IPV6       string    `bigquery:"ip_v6"`
	DomainName string    `bigquery:"domain_name"`
}

func (d *AAAARecord) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "ip_v6"}, {Name: "domain_name"}},
			DoUpdates: clause.Assignments(
				map[string]interface{}{
					"updated_at": gorm.Expr("MAX(updated_at, excluded.updated_at)"),
				},
			),
		},
	)
	return nil
}

func (d *AAAARecordBQ) ToGorm() AAAARecord {
	return AAAARecord{
		Model: gorm.Model{
			CreatedAt: d.CreatedAt,
			UpdatedAt: d.UpdatedAt,
		},
		IPV6:       d.IPV6,
		DomainName: d.DomainName,
	}
}

func (a *AAAARecord) ToBQ() AAAARecordBQ {
	return AAAARecordBQ{
		ID:         int(a.ID),
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		IPV6:       a.IPV6,
		DomainName: a.DomainName,
	}
}

type ARecord struct {
	gorm.Model
	IP         string `gorm:"uniqueIndex:dom_a_idx" bigquery:"ip"`
	DomainName string `gorm:"uniqueIndex:dom_a_idx" bigquery:"domain_id"`
}

func (d *ARecord) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "ip"}, {Name: "domain_name"}},
			DoUpdates: clause.Assignments(
				map[string]interface{}{
					"updated_at": gorm.Expr("MAX(updated_at, excluded.updated_at)"),
				},
			),
		},
	)
	return nil
}

type ARecordBQ struct {
	ID         int       `bigquery:"id"`
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	IP         string    `bigquery:"ip"`
	DomainName string    `bigquery:"domain_name"`
}

func (a *ARecord) ToBQ() ARecordBQ {
	return ARecordBQ{
		ID:         int(a.ID),
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		IP:         a.IP,
		DomainName: a.DomainName,
	}
}

func (a *ARecordBQ) ToGorm() ARecord {
	return ARecord{
		Model: gorm.Model{
			CreatedAt: a.CreatedAt,
			UpdatedAt: a.UpdatedAt,
		},
		IP:         a.IP,
		DomainName: a.DomainName,
	}
}

type SOARecord struct {
	gorm.Model
	NS         string `gorm:"uniqueIndex:soa_dom_idx" bigquery:"ns"`
	MBox       string `gorm:"uniqueIndex:soa_dom_idx" bigquery:"mbox"`
	Serial     uint32 `gorm:"uniqueIndex:soa_dom_idx" bigquery:"serial"`
	DomainName string `gorm:"uniqueIndex:soa_dom_idx" bigquery:"domain_name"`
}

func (d *SOARecord) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "ns"}, {Name: "m_box"}, {Name: "serial"}, {Name: "domain_name"}},
			DoUpdates: clause.Assignments(
				map[string]interface{}{
					"updated_at": gorm.Expr("MAX(updated_at, excluded.updated_at)"),
				},
			),
		},
	)
	return nil
}

type SOARecordBQ struct {
	ID         int                 `bigquery:"id"`
	CreatedAt  time.Time           `bigquery:"created_at"`
	UpdatedAt  time.Time           `bigquery:"updated_at"`
	NS         bigquery.NullString `gorm:"primaryKey" bigquery:"ns"`
	MBox       bigquery.NullString `gorm:"primaryKey" bigquery:"mbox"`
	Serial     bigquery.NullInt64  `gorm:"primaryKey" bigquery:"serial"`
	DomainName string              `gorm:"primaryKey" bigquery:"domain_name"`
}

func (s *SOARecord) ToBQ() SOARecordBQ {
	return SOARecordBQ{
		ID:         int(s.ID),
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		NS:         bigquery.NullString{Valid: s.NS != "", StringVal: s.NS},
		MBox:       bigquery.NullString{Valid: s.MBox != "", StringVal: s.MBox},
		Serial:     bigquery.NullInt64{Valid: s.Serial != 0, Int64: int64(s.Serial)},
		DomainName: s.DomainName,
	}
}

func (s *SOARecordBQ) ToGorm() SOARecord {
	return SOARecord{
		Model: gorm.Model{
			CreatedAt: s.CreatedAt,
			UpdatedAt: s.UpdatedAt,
		},
		NS:         s.NS.StringVal,
		MBox:       s.MBox.StringVal,
		Serial:     uint32(s.Serial.Int64),
		DomainName: s.DomainName,
	}
}

type MXRecord struct {
	gorm.Model
	Mx         string `gorm:"uniqueIndex:mx_dom_idx" bigquery:"mx"`
	DomainName string `gorm:"uniqueIndex:mx_dom_idx" bigquery:"domain_name"`
}

type MXRecordBQ struct {
	ID         int       `bigquery:"id"`
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	Mx         string    `bigquery:"mx"`
	DomainName string    `bigquery:"domain_name"`
}

func (d *MXRecord) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "mx"}, {Name: "domain_name"}},
			DoUpdates: clause.Assignments(
				map[string]interface{}{
					"updated_at": gorm.Expr("MAX(updated_at, excluded.updated_at)"),
				},
			),
		},
	)
	return nil
}

func (m *MXRecord) ToBQ() MXRecordBQ {
	return MXRecordBQ{
		ID:         int(m.ID),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Mx:         m.Mx,
		DomainName: m.DomainName,
	}
}

func (m *MXRecordBQ) ToGorm() MXRecord {
	return MXRecord{
		Model: gorm.Model{
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		},
		Mx:         m.Mx,
		DomainName: m.DomainName,
	}
}

var (
	DomainClient *dns.Client
	ClientConfig *dns.ClientConfig
)

func (d *Domain) QueryMX() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(d.DomainName), dns.TypeMX)
	r, err := queryAllServers(msg)
	if err != nil {
		return err
	}
	mxs := []MXRecord{}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.MX); ok {
			mxs = append(mxs, MXRecord{DomainName: d.DomainName, Mx: a.Mx})
		}
	}
	d.MXRecords = mxs
	return nil
}

func (d *Domain) QueryA() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(d.DomainName), dns.TypeA)
	r, err := queryAllServers(msg)
	if err != nil {
		return err
	}
	ips := []ARecord{}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			ips = append(ips, ARecord{IP: a.A.String(), DomainName: d.DomainName})
		}
	}
	d.ARecords = ips
	return nil
}

func (d *Domain) QueryAAAA() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(d.DomainName), dns.TypeAAAA)
	r, err := queryAllServers(msg)
	if err != nil {
		return err
	}
	ips := []AAAARecord{}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.AAAA); ok {
			ips = append(ips, AAAARecord{IPV6: a.AAAA.String(), DomainName: d.DomainName})
		}
	}
	d.AAAARecords = ips
	return nil
}

func (d *Domain) QuerySOA() error {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(d.DomainName), dns.TypeSOA)
	r, err := queryAllServers(msg)
	if err != nil {
		return err
	}
	soas := []SOARecord{}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.SOA); ok {
			soas = append(soas, SOARecord{NS: a.Ns, MBox: a.Mbox, Serial: a.Serial, DomainName: d.DomainName})
		}
	}
	d.SOARecords = soas
	return nil
}

func query(msg *dns.Msg, nameserver string) (*dns.Msg, error) {
	r, _, err := DomainClient.Exchange(msg, nameserver)
	return r, err
}

func queryAllServers(msg *dns.Msg) (*dns.Msg, error) {
	for _, ns := range ClientConfig.Servers {
		r, err := query(msg, ns+":"+ClientConfig.Port)
		if err == nil {
			return r, nil
		}
	}
	log.Printf(
		"Failed to query record for domain %s with configured servers, trying with google DNS...\n",
		msg.Question[0].Name,
	)
	r, err := query(msg, "8.8.8.8:53")
	if err == nil {
		return r, nil
	}

	return nil, errors.New("Failed to query all servers")
}

func (d *Domain) GetDNSRecords() []error {
	d.LastRanDns = time.Now()
	if d.NonPublicDomain {
		return []error{errors.New("Non public domain")}
	}
	errs := []error{}
	err := d.QueryA()
	if err != nil {
		errs = append(errs, err)
	}
	err = d.QueryAAAA()
	if err != nil {
		errs = append(errs, err)
	}
	err = d.QueryMX()
	if err != nil {
		errs = append(errs, err)
	}
	err = d.QuerySOA()
	if err != nil {
		errs = append(errs, err)
	}
	return errs
}

func init() {
	var err error
	ClientConfig, err = dns.ClientConfigFromFile("/etc/resolv.conf") // TODO: Make this part of the config
	if err != nil {
		log.Fatal(err)
	}
	DomainClient = new(dns.Client)
}
