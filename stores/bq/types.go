package bq

import (
	"time"

	"cloud.google.com/go/bigquery"
	"dev.azure.com/Unum/Mkt_Analytics/_git/domwalk/domains"
)

type DomainBQ struct {
	CreatedAt             time.Time           `bigquery:"created_at"`
	UpdatedAt             time.Time           `bigquery:"updated_at"`
	DomainName            string              `bigquery:"domain_name"`
	NonPublicDomain       bool                `bigquery:"non_public_domain"`
	Hostname              bigquery.NullString `bigquery:"hostname"`
	Subdomain             bigquery.NullString `bigquery:"subdomain"`
	Suffix                bigquery.NullString `bigquery:"suffix"`
	SuccessfulWebLanding  bool                `bigquery:"successful_web_landing"`
	WebRedirectURLFinal   bigquery.NullString `bigquery:"web_redirect_url_final"`
	LastRanWebRedirect    time.Time           `bigquery:"last_ran_web_redirect"`
	LastRanDns            time.Time           `bigquery:"last_ran_dns"`
	LastRanCertSans       time.Time           `bigquery:"last_ran_cert_sans"`
	LastRanSitemapParse   time.Time           `bigquery:"last_ran_sitemap_parse"`
	ARecords              []ARecordBQ         `bigquery:"a_records"`
	AAAARecords           []AAAARecordBQ      `bigquery:"aaaa_records"`
	MXRecords             []MXRecordBQ        `bigquery:"mx_records"`
	SOARecords            []SOARecordBQ       `bigquery:"soa_records"`
	Sitemaps              []SitemapBQ         `bigquery:"sitemaps"`
	WebRedirectDomains    []MatchedDomainBQ   `bigquery:"web_redirect_domains"`
	CertSANs              []MatchedDomainBQ   `bigquery:"cert_sans"`
	SitemapWebDomains     []MatchedDomainBQ   `bigquery:"sitemap_web_domains"`
	SitemapContactDomains []MatchedDomainBQ   `bigquery:"sitemap_contact_domains"`
}

func newDomainBQ(record *domains.Domain) DomainBQ {
	dbq := DomainBQ{
		CreatedAt:            record.CreatedAt,
		UpdatedAt:            record.UpdatedAt,
		DomainName:           record.DomainName,
		NonPublicDomain:      record.NonPublicDomain,
		Hostname:             bigquery.NullString{Valid: record.Hostname != "", StringVal: record.Hostname},
		Subdomain:            bigquery.NullString{Valid: record.Subdomain != "", StringVal: record.Subdomain},
		Suffix:               bigquery.NullString{Valid: record.Suffix != "", StringVal: record.Suffix},
		SuccessfulWebLanding: record.SuccessfulWebLanding,
		LastRanWebRedirect:   record.LastRanWebRedirect,
		LastRanDns:           record.LastRanDns,
		LastRanCertSans:      record.LastRanCertSans,
		LastRanSitemapParse:  record.LastRanSitemapParse,
	}
	var aRecords []ARecordBQ
	for _, a := range record.ARecords {
		aRecords = append(aRecords, newARecordBQ(a))
	}
	dbq.ARecords = aRecords

	var aaaaRecords []AAAARecordBQ
	for _, a := range record.AAAARecords {
		aaaaRecords = append(aaaaRecords, newAAAARecordBQ(a))
	}
	dbq.AAAARecords = aaaaRecords

	var mxRecords []MXRecordBQ
	for _, a := range record.MXRecords {
		mxRecords = append(mxRecords, newMXRecordBQ(a))
	}

	dbq.MXRecords = mxRecords

	var soaRecords []SOARecordBQ
	for _, a := range record.SOARecords {
		soaRecords = append(soaRecords, newSOARecordBQ(a))
	}
	dbq.SOARecords = soaRecords

	var sitemaps []SitemapBQ
	for _, a := range record.Sitemaps {
		sitemaps = append(sitemaps, newSitemapBQ(*a))
	}
	dbq.Sitemaps = sitemaps

	var webRedirectDomains []MatchedDomainBQ
	for _, a := range record.WebRedirectDomains {
		webRedirectDomains = append(webRedirectDomains, newMatchedDomainBQ(a.MatchedDomain))
	}
	dbq.WebRedirectDomains = webRedirectDomains

	var certSANs []MatchedDomainBQ
	for _, a := range record.CertSANs {
		certSANs = append(certSANs, newMatchedDomainBQ(a.MatchedDomain))
	}
	dbq.CertSANs = certSANs

	var sitemapWebDomains []MatchedDomainBQ
	for _, a := range record.SitemapWebDomains {
		sitemapWebDomains = append(sitemapWebDomains, newMatchedDomainBQ(a.MatchedDomain))
	}
	dbq.SitemapWebDomains = sitemapWebDomains

	var sitemapContactDomains []MatchedDomainBQ
	for _, a := range record.SitemapContactDomains {
		sitemapContactDomains = append(sitemapContactDomains, newMatchedDomainBQ(a.MatchedDomain))
	}
	dbq.SitemapContactDomains = sitemapContactDomains

	return dbq
}

func (a *DomainBQ) parse() *domains.Domain {
	d := &domains.Domain{
		CreatedAt:            a.CreatedAt,
		UpdatedAt:            a.UpdatedAt,
		DomainName:           a.DomainName,
		NonPublicDomain:      a.NonPublicDomain,
		Hostname:             a.Hostname.StringVal,
		Subdomain:            a.Subdomain.StringVal,
		Suffix:               a.Suffix.StringVal,
		SuccessfulWebLanding: a.SuccessfulWebLanding,
		LastRanWebRedirect:   a.LastRanWebRedirect,
		LastRanDns:           a.LastRanDns,
		LastRanCertSans:      a.LastRanCertSans,
		LastRanSitemapParse:  a.LastRanSitemapParse,
	}
	var aRecords []domains.ARecord
	for _, a := range a.ARecords {
		aRecords = append(aRecords, a.parse())
	}
	d.ARecords = aRecords

	var aaaaRecords []domains.AAAARecord
	for _, a := range a.AAAARecords {
		aaaaRecords = append(aaaaRecords, a.parse())
	}
	d.AAAARecords = aaaaRecords

	var mxRecords []domains.MXRecord
	for _, a := range a.MXRecords {
		mxRecords = append(mxRecords, a.parse())
	}
	d.MXRecords = mxRecords

	var soaRecords []domains.SOARecord
	for _, a := range a.SOARecords {
		soaRecords = append(soaRecords, a.parse())
	}
	d.SOARecords = soaRecords

	var sitemaps []*domains.Sitemap
	for _, a := range a.Sitemaps {
		sitemaps = append(sitemaps, a.parse())
	}
	d.Sitemaps = sitemaps

	var webRedirectDomains []domains.WebRedirectDomain
	for _, a := range a.WebRedirectDomains {
		webRedirectDomains = append(webRedirectDomains, domains.WebRedirectDomain{MatchedDomain: a.parse()})
	}
	d.WebRedirectDomains = webRedirectDomains

	var certSANs []domains.CertSansDomain
	for _, a := range a.CertSANs {
		certSANs = append(certSANs, domains.CertSansDomain{MatchedDomain: a.parse()})
	}
	d.CertSANs = certSANs

	var sitemapWebDomains []domains.SitemapWebDomain
	for _, a := range a.SitemapWebDomains {
		sitemapWebDomains = append(sitemapWebDomains, domains.SitemapWebDomain{MatchedDomain: a.parse()})
	}
	d.SitemapWebDomains = sitemapWebDomains

	var sitemapContactDomains []domains.SitemapContactDomain
	for _, a := range a.SitemapContactDomains {
		sitemapContactDomains = append(sitemapContactDomains, domains.SitemapContactDomain{MatchedDomain: a.parse()})
	}
	d.SitemapContactDomains = sitemapContactDomains

	return d
}

type ARecordBQ struct {
	CreatedAt time.Time `bigquery:"created_at"`
	UpdatedAt time.Time `bigquery:"updated_at"`
	IP        string    `bigquery:"ip"`
}

func newARecordBQ(record domains.ARecord) ARecordBQ {
	return ARecordBQ{
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		IP:        record.IP,
	}
}

func (a *ARecordBQ) parse() domains.ARecord {
	return domains.ARecord{
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		IP:        a.IP,
	}
}

type AAAARecordBQ struct {
	CreatedAt time.Time `bigquery:"created_at"`
	UpdatedAt time.Time `bigquery:"updated_at"`
	IPV6      string    `bigquery:"ip_v6"`
}

func newAAAARecordBQ(record domains.AAAARecord) AAAARecordBQ {
	return AAAARecordBQ{
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		IPV6:      record.IPV6,
	}
}

func (a *AAAARecordBQ) parse() domains.AAAARecord {
	return domains.AAAARecord{
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		IPV6:      a.IPV6,
	}
}

type MXRecordBQ struct {
	CreatedAt time.Time `bigquery:"created_at"`
	UpdatedAt time.Time `bigquery:"updated_at"`
	Mx        string    `bigquery:"mx"`
}

func newMXRecordBQ(record domains.MXRecord) MXRecordBQ {
	return MXRecordBQ{
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		Mx:        record.Mx,
	}
}

func (a *MXRecordBQ) parse() domains.MXRecord {
	return domains.MXRecord{
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		Mx:        a.Mx,
	}
}

type SOARecordBQ struct {
	CreatedAt time.Time           `bigquery:"created_at"`
	UpdatedAt time.Time           `bigquery:"updated_at"`
	NS        bigquery.NullString `bigquery:"ns"`
	MBox      bigquery.NullString `bigquery:"mbox"`
	Serial    bigquery.NullInt64  `bigquery:"serial"`
}

func newSOARecordBQ(record domains.SOARecord) SOARecordBQ {
	return SOARecordBQ{
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		NS:        bigquery.NullString{StringVal: record.NS, Valid: record.NS != ""},
		MBox:      bigquery.NullString{StringVal: record.MBox, Valid: record.MBox != ""},
		Serial:    bigquery.NullInt64{Int64: int64(record.Serial), Valid: record.Serial != 0},
	}
}

func (a *SOARecordBQ) parse() domains.SOARecord {
	return domains.SOARecord{
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		NS:        a.NS.StringVal,
		MBox:      a.MBox.StringVal,
		Serial:    uint32(a.Serial.Int64),
	}
}

type SitemapBQ struct {
	CreatedAt  time.Time           `bigquery:"created_at"`
	UpdatedAt  time.Time           `bigquery:"updated_at"`
	SitemapLoc bigquery.NullString `bigquery:"sitemap_loc"`
}

func newSitemapBQ(record domains.Sitemap) SitemapBQ {
	return SitemapBQ{
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
		SitemapLoc: bigquery.NullString{StringVal: record.SitemapLoc, Valid: record.SitemapLoc != ""},
	}
}

func (a *SitemapBQ) parse() *domains.Sitemap {
	return &domains.Sitemap{
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		SitemapLoc: a.SitemapLoc.StringVal,
	}
}

type MatchedDomainBQ struct {
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	DomainName string    `bigquery:"domain_name"`
}

func newMatchedDomainBQ(record domains.MatchedDomain) MatchedDomainBQ {
	return MatchedDomainBQ{
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
		DomainName: record.DomainName,
	}
}

func (a *MatchedDomainBQ) parse() domains.MatchedDomain {
	return domains.MatchedDomain{
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		DomainName: a.DomainName,
	}
}
