package types

import (
	"errors"
	"log"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/miekg/dns"
)

type AAAARecord struct {
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	IPV6       string    `gorm:"primaryKey" bigquery:"ip_v6"`
	DomainName string    `gorm:"primaryKey" bigquery:"domain_name"`
}

type ARecord struct {
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	IP         string    `gorm:"primaryKey" bigquery:"ip"`
	DomainName string    `gorm:"primaryKey" bigquery:"domain_name"`
}

type SOARecord struct {
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	NS         string    `gorm:"primaryKey" bigquery:"ns"`
	MBox       string    `gorm:"primaryKey" bigquery:"mbox"`
	Serial     uint32    `gorm:"primaryKey" bigquery:"serial"`
	DomainName string    `gorm:"primaryKey" bigquery:"domain_name"`
}

type SOARecordBQ struct {
	CreatedAt  time.Time           `bigquery:"created_at"`
	UpdatedAt  time.Time           `bigquery:"updated_at"`
	NS         bigquery.NullString `gorm:"primaryKey" bigquery:"ns"`
	MBox       bigquery.NullString `gorm:"primaryKey" bigquery:"mbox"`
	Serial     bigquery.NullInt64  `gorm:"primaryKey" bigquery:"serial"`
	DomainName string              `gorm:"primaryKey" bigquery:"domain_name"`
}

func (s *SOARecord) ToBQ() SOARecordBQ {
	return SOARecordBQ{
		CreatedAt:  s.CreatedAt,
		UpdatedAt:  s.UpdatedAt,
		NS:         bigquery.NullString{Valid: s.NS != "", StringVal: s.NS},
		MBox:       bigquery.NullString{Valid: s.MBox != "", StringVal: s.MBox},
		Serial:     bigquery.NullInt64{Valid: s.Serial != 0, Int64: int64(s.Serial)},
		DomainName: s.DomainName,
	}
}

type MXRecord struct {
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	Mx         string    `gorm:"primaryKey" bigquery:"mx"`
	DomainName string    `gorm:"primaryKey" bigquery:"domain_name"`
}

var (
	DomainClient *dns.Client
	ClientConfig *dns.ClientConfig
)

func QueryMX(domain string) ([]MXRecord, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeMX)
	r, err := queryAllServers(msg)
	if err != nil {
		return nil, err
	}
	mxs := []MXRecord{}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.MX); ok {
			mxs = append(mxs, MXRecord{DomainName: domain, Mx: a.Mx})
		}
	}
	return mxs, nil
}

func QueryA(domain string) ([]ARecord, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	r, err := queryAllServers(msg)
	if err != nil {
		return nil, err
	}
	ips := []ARecord{}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			ips = append(ips, ARecord{IP: a.A.String(), DomainName: domain})
		}
	}
	return ips, nil
}

func QueryAAAA(domain string) ([]AAAARecord, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeAAAA)
	r, err := queryAllServers(msg)
	if err != nil {
		return nil, err
	}
	ips := []AAAARecord{}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.AAAA); ok {
			ips = append(ips, AAAARecord{IPV6: a.AAAA.String(), DomainName: domain})
		}
	}
	return ips, nil
}

func QuerySOA(domain string) ([]SOARecord, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeSOA)
	r, err := queryAllServers(msg)
	if err != nil {
		return nil, err
	}
	soas := []SOARecord{}
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.SOA); ok {
			soas = append(soas, SOARecord{NS: a.Ns, MBox: a.Mbox, Serial: a.Serial, DomainName: domain})
		}
	}
	return soas, nil
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
	var err error
	errs := []error{}
	d.ARecords, err = QueryA(d.DomainName)
	if err != nil {
		errs = append(errs, err)
	}
	d.AAAARecords, err = QueryAAAA(d.DomainName)
	if err != nil {
		errs = append(errs, err)
	}
	d.MXRecords, err = QueryMX(d.DomainName)
	if err != nil {
		errs = append(errs, err)
	}
	d.SOARecords, err = QuerySOA(d.DomainName)
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
