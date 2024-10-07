package types

import (
	"errors"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/miekg/dns"
	"gorm.io/gorm"
)

type AAAARecord struct {
	gorm.Model
	IPV6     string `gorm:"uniqueIndex:aaaa_dom_idx" bigquery:"ip_v6"`
	DomainID uint   `gorm:"uniqueIndex:aaaa_dom_idx" bigquery:"domain_id"`
}

type AAAARecordBQ struct {
	ID        int       `bigquery:"id"`
	CreatedAt time.Time `bigquery:"created_at"`
	UpdatedAt time.Time `bigquery:"updated_at"`
	IPV6      string    `bigquery:"ip_v6"`
	DomainID  int       `bigquery:"domain_id"`
}

func (a *AAAARecord) ToBQ() AAAARecordBQ {
	return AAAARecordBQ{
		ID:        int(a.ID),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		IPV6:      a.IPV6,
		DomainID:  int(a.DomainID),
	}
}

type ARecord struct {
	gorm.Model
	IP       string `gorm:"uniqueIndex:dom_a_idx" bigquery:"ip"`
	DomainID uint   `gorm:"uniqueIndex:dom_a_idx" bigquery:"domain_id"`
}

type ARecordBQ struct {
	ID        int       `bigquery:"id"`
	CreatedAt time.Time `bigquery:"created_at"`
	UpdatedAt time.Time `bigquery:"updated_at"`
	IP        string    `bigquery:"ip"`
	DomainID  int       `bigquery:"domain_id"`
}

func (a *ARecord) ToBQ() ARecordBQ {
	return ARecordBQ{
		ID:        int(a.ID),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		IP:        a.IP,
		DomainID:  int(a.DomainID),
	}
}

type SOARecord struct {
	gorm.Model
	NS       string `gorm:"uniqueIndex:soa_dom_idx" bigquery:"ns"`
	MBox     string `gorm:"uniqueIndex:soa_dom_idx" bigquery:"mbox"`
	Serial   uint32 `gorm:"uniqueIndex:soa_dom_idx" bigquery:"serial"`
	DomainID uint   `gorm:"uniqueIndex:soa_dom_idx" bigquery:"domain_id"`
}

type SOARecordBQ struct {
	ID        int                 `bigquery:"id"`
	CreatedAt time.Time           `bigquery:"created_at"`
	UpdatedAt time.Time           `bigquery:"updated_at"`
	NS        bigquery.NullString `gorm:"primaryKey" bigquery:"ns"`
	MBox      bigquery.NullString `gorm:"primaryKey" bigquery:"mbox"`
	Serial    bigquery.NullInt64  `gorm:"primaryKey" bigquery:"serial"`
	DomainID  int                 `gorm:"primaryKey" bigquery:"domain_id"`
}

func (s *SOARecord) ToBQ() SOARecordBQ {
	return SOARecordBQ{
		ID:        int(s.ID),
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
		NS:        bigquery.NullString{Valid: s.NS != "", StringVal: s.NS},
		MBox:      bigquery.NullString{Valid: s.MBox != "", StringVal: s.MBox},
		Serial:    bigquery.NullInt64{Valid: s.Serial != 0, Int64: int64(s.Serial)},
		DomainID:  int(s.DomainID),
	}
}

type MXRecord struct {
	gorm.Model
	Mx       string `gorm:"uniqueIndex:mx_dom_idx" bigquery:"mx"`
	DomainID uint   `gorm:"uniqueIndex:mx_dom_idx" bigquery:"domain_id"`
}

type MXRecordBQ struct {
	ID        int       `bigquery:"id"`
	CreatedAt time.Time `bigquery:"created_at"`
	UpdatedAt time.Time `bigquery:"updated_at"`
	Mx        string    `bigquery:"mx"`
	DomainID  int       `bigquery:"domain_id"`
}

func (m *MXRecord) ToBQ() MXRecordBQ {
	return MXRecordBQ{
		ID:        int(m.ID),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		Mx:        m.Mx,
		DomainID:  int(m.DomainID),
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
			mxs = append(mxs, MXRecord{DomainID: d.ID, Mx: a.Mx})
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
			ips = append(ips, ARecord{IP: a.A.String(), DomainID: d.ID})
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
			ips = append(ips, AAAARecord{IPV6: a.AAAA.String(), DomainID: d.ID})
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
			soas = append(soas, SOARecord{NS: a.Ns, MBox: a.Mbox, Serial: a.Serial, DomainID: d.ID})
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
