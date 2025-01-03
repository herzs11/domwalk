package bq

import (
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/herzs11/go-doms/domain"
)

type WhoisContactBQ struct {
	Name         bigquery.NullString `bigquery:"name"`
	Organization bigquery.NullString `bigquery:"organization"`
	Street1      bigquery.NullString `bigquery:"street1"`
	Street2      bigquery.NullString `bigquery:"street2"`
	Street3      bigquery.NullString `bigquery:"street3"`
	Street4      bigquery.NullString `bigquery:"street4"`
	City         bigquery.NullString `bigquery:"city"`
	State        bigquery.NullString `bigquery:"state"`
	PostalCode   bigquery.NullString `bigquery:"postal_code"`
	Country      bigquery.NullString `bigquery:"country"`
	CountryCode  bigquery.NullString `bigquery:"country_code"`
	Email        bigquery.NullString `bigquery:"email"`
	Telephone    bigquery.NullString `bigquery:"telephone"`
	TelephoneExt bigquery.NullString `bigquery:"telephone_ext"`
	Fax          bigquery.NullString `bigquery:"fax"`
	FaxExt       bigquery.NullString `bigquery:"fax_ext"`
}

type WhoisDataBQ struct {
	DomainName            string              `bigquery:"domain_name"`
	CreatedDate           time.Time           `bigquery:"created_date"`
	UpdatedDate           time.Time           `bigquery:"updated_date"`
	RegistrarName         string              `bigquery:"registrar_name"`
	RegistrarIANAId       bigquery.NullString `bigquery:"registrar_ianaid"`
	Status                bigquery.NullString `bigquery:"status"`
	Registrant            *WhoisContactBQ     `bigquery:"registrant,nullable"`
	AdministrativeContact *WhoisContactBQ     `bigquery:"administrative_contact,nullable"`
	TechnicalContact      *WhoisContactBQ     `bigquery:"technical_contact,nullable"`
	BillingContact        *WhoisContactBQ     `bigquery:"billing_contact,nullable"`
	ZoneContact           *WhoisContactBQ     `bigquery:"zone_contact,nullable"`
	Header                bigquery.NullString `bigquery:"header"`
	Footer                bigquery.NullString `bigquery:"footer"`
	EstimatedDomainAge    bigquery.NullInt64  `bigquery:"estimated_domain_age"`
	Ips                   []string            `bigquery:"ips"`
	LastUpdated           time.Time           `bigquery:"last_updated"`
}

type DomainBQ struct {
	CreatedAt             time.Time           `bigquery:"created_at"`
	UpdatedAt             time.Time           `bigquery:"updated_at"`
	DomainName            string              `bigquery:"domain_name"`
	NonPublicDomain       bool                `bigquery:"non_public_domain"`
	Hostname              bigquery.NullString `bigquery:"hostname"`
	Subdomain             bigquery.NullString `bigquery:"subdomain"`
	Suffix                bigquery.NullString `bigquery:"suffix"`
	SuccessfulWebLanding  bool                `bigquery:"successful_web_landing"`
	CertOrgNames          []string            `bigquery:"cert_org_names"`
	WebRedirectURLFinal   bigquery.NullString `bigquery:"web_redirect_url_final"`
	LastRanWebRedirect    time.Time           `bigquery:"last_ran_web_redirect"`
	LastRanDns            time.Time           `bigquery:"last_ran_dns"`
	LastRanCertSans       time.Time           `bigquery:"last_ran_cert_sans"`
	LastRanSitemapParse   time.Time           `bigquery:"last_ran_sitemap_parse"`
	LastRanWhois          time.Time           `bigquery:"last_ran_whois"`
	LastRanReverseWhois   time.Time           `bigquery:"last_ran_reverse_whois"`
	ARecords              []ARecordBQ         `bigquery:"a_records"`
	AAAARecords           []AAAARecordBQ      `bigquery:"aaaa_records"`
	MXRecords             []MXRecordBQ        `bigquery:"mx_records"`
	SOARecords            []SOARecordBQ       `bigquery:"soa_records"`
	Sitemaps              []SitemapBQ         `bigquery:"sitemaps"`
	*WhoisDataBQ          `bigquery:"whois_data,nullable"`
	WebRedirectDomains    []MatchedDomainBQ `bigquery:"web_redirect_domains"`
	CertSANs              []MatchedDomainBQ `bigquery:"cert_sans"`
	SitemapWebDomains     []MatchedDomainBQ `bigquery:"sitemap_web_domains"`
	SitemapContactDomains []MatchedDomainBQ `bigquery:"sitemap_contact_domains"`
	ReverseWhoisDomains   []MatchedDomainBQ `bigquery:"reverse_whois_domains"`
}

func newWhoisDataBQ(record *domain.WhoisData) *WhoisDataBQ {
	wdb := WhoisDataBQ{
		DomainName:    record.DomainName,
		CreatedDate:   record.CreatedDate,
		UpdatedDate:   record.UpdatedDate,
		RegistrarName: record.RegistrarName,
		RegistrarIANAId: bigquery.NullString{
			StringVal: record.RegistrarIANAID,
			Valid:     record.RegistrarIANAID != "",
		},
		Status: bigquery.NullString{
			StringVal: record.Status,
			Valid:     record.Status != "",
		},
		Registrant:            newWhoisContactBQ(&record.Registrant),
		AdministrativeContact: newWhoisContactBQ(&record.AdministrativeContact),
		TechnicalContact:      newWhoisContactBQ(&record.TechnicalContact),
		BillingContact:        newWhoisContactBQ(&record.BillingContact),
		ZoneContact:           newWhoisContactBQ(&record.ZoneContact),
		Header: bigquery.NullString{
			StringVal: record.Header,
			Valid:     record.Header != "",
		},
		Footer: bigquery.NullString{
			StringVal: record.Footer,
			Valid:     record.Footer != "",
		},
		EstimatedDomainAge: bigquery.NullInt64{
			Int64: int64(record.EstimatedDomainAge),
			Valid: record.EstimatedDomainAge != 0,
		},
		Ips:         record.Ips,
		LastUpdated: record.LastUpdated,
	}
	return &wdb
}

func newWhoisContactBQ(record *domain.WhoisContact) *WhoisContactBQ {
	wic := &WhoisContactBQ{
		Name: bigquery.NullString{
			StringVal: record.Name,
			Valid:     record.Name != "",
		},
		Organization: bigquery.NullString{
			StringVal: record.Organization,
			Valid:     record.Organization != "",
		},
		Street1: bigquery.NullString{
			StringVal: record.Street1,
			Valid:     record.Street1 != "",
		},
		Street2: bigquery.NullString{
			StringVal: record.Street2,
			Valid:     record.Street2 != "",
		},
		Street3: bigquery.NullString{
			StringVal: record.Street3,
			Valid:     record.Street3 != "",
		},
		Street4: bigquery.NullString{
			StringVal: record.Street4,
			Valid:     record.Street4 != "",
		},
		City: bigquery.NullString{
			StringVal: record.City,
			Valid:     record.City != "",
		},
		State: bigquery.NullString{
			StringVal: record.State,
			Valid:     record.State != "",
		},
		PostalCode: bigquery.NullString{
			StringVal: record.PostalCode,
			Valid:     record.PostalCode != "",
		},
		Country: bigquery.NullString{
			StringVal: record.Country,
			Valid:     record.Country != "",
		},
		CountryCode: bigquery.NullString{
			StringVal: record.CountryCode,
			Valid:     record.CountryCode != "",
		},
		Email: bigquery.NullString{
			StringVal: record.Email,
			Valid:     record.Email != "",
		},
		Telephone: bigquery.NullString{
			StringVal: record.Telephone,
			Valid:     record.Telephone != "",
		},
		TelephoneExt: bigquery.NullString{
			StringVal: record.TelephoneExt,
			Valid:     record.TelephoneExt != "",
		},
		Fax: bigquery.NullString{
			StringVal: record.Fax,
			Valid:     record.Fax != "",
		},
		FaxExt: bigquery.NullString{
			StringVal: record.FaxExt,
			Valid:     record.FaxExt != "",
		},
	}
	return wic
}

func newDomainBQ(record *domain.Domain) DomainBQ {
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
		LastRanWhois:         record.LastRanWhois,
		LastRanReverseWhois:  record.LastRanReverseWhois,
		WhoisDataBQ:          newWhoisDataBQ(record.Whois),
		CertOrgNames:         record.CertOrgNames,
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
		webRedirectDomains = append(webRedirectDomains, newMatchedDomainBQ(a))
	}
	dbq.WebRedirectDomains = webRedirectDomains

	var certSANs []MatchedDomainBQ
	for _, a := range record.CertSANs {
		certSANs = append(certSANs, newMatchedDomainBQ(a))
	}
	dbq.CertSANs = certSANs

	var sitemapWebDomains []MatchedDomainBQ
	for _, a := range record.SitemapWebDomains {
		sitemapWebDomains = append(sitemapWebDomains, newMatchedDomainBQ(a))
	}
	dbq.SitemapWebDomains = sitemapWebDomains

	var sitemapContactDomains []MatchedDomainBQ
	for _, a := range record.SitemapContactDomains {
		sitemapContactDomains = append(sitemapContactDomains, newMatchedDomainBQ(a))
	}
	dbq.SitemapContactDomains = sitemapContactDomains

	var reverseWhoisDomains []MatchedDomainBQ
	for _, a := range record.ReverseWhoisDomains {
		reverseWhoisDomains = append(reverseWhoisDomains, newMatchedDomainBQ(a))
	}
	dbq.ReverseWhoisDomains = reverseWhoisDomains

	return dbq
}

func (a *DomainBQ) parse() *domain.Domain {
	d := &domain.Domain{
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
		LastRanWhois:         a.LastRanWhois,
		LastRanReverseWhois:  a.LastRanReverseWhois,
		Whois:                a.WhoisDataBQ.parse(),
		CertOrgNames:         a.CertOrgNames,
	}
	var aRecords []domain.ARecord
	for _, a := range a.ARecords {
		aRecords = append(aRecords, a.parse())
	}
	d.ARecords = aRecords

	var aaaaRecords []domain.AAAARecord
	for _, a := range a.AAAARecords {
		aaaaRecords = append(aaaaRecords, a.parse())
	}
	d.AAAARecords = aaaaRecords

	var mxRecords []domain.MXRecord
	for _, a := range a.MXRecords {
		mxRecords = append(mxRecords, a.parse())
	}
	d.MXRecords = mxRecords

	var soaRecords []domain.SOARecord
	for _, a := range a.SOARecords {
		soaRecords = append(soaRecords, a.parse())
	}
	d.SOARecords = soaRecords

	var sitemaps []*domain.Sitemap
	for _, a := range a.Sitemaps {
		sitemaps = append(sitemaps, a.parse())
	}
	d.Sitemaps = sitemaps

	var webRedirectDomains []domain.MatchedDomain
	for _, a := range a.WebRedirectDomains {
		webRedirectDomains = append(webRedirectDomains, a.parse())
	}
	d.WebRedirectDomains = webRedirectDomains

	var certSANs []domain.MatchedDomain
	for _, a := range a.CertSANs {
		certSANs = append(certSANs, a.parse())
	}
	d.CertSANs = certSANs

	var sitemapWebDomains []domain.MatchedDomain
	for _, a := range a.SitemapWebDomains {
		sitemapWebDomains = append(sitemapWebDomains, a.parse())
	}
	d.SitemapWebDomains = sitemapWebDomains

	var sitemapContactDomains []domain.MatchedDomain
	for _, a := range a.SitemapContactDomains {
		sitemapContactDomains = append(sitemapContactDomains, a.parse())
	}
	d.SitemapContactDomains = sitemapContactDomains

	var reverseWhoisDomains []domain.MatchedDomain
	for _, a := range a.ReverseWhoisDomains {
		reverseWhoisDomains = append(reverseWhoisDomains, a.parse())
	}
	d.ReverseWhoisDomains = reverseWhoisDomains

	return d
}

func (w *WhoisContactBQ) parse() *domain.WhoisContact {
	return &domain.WhoisContact{
		Name:         w.Name.String(),
		Organization: w.Organization.String(),
		Street1:      w.Street1.String(),
		Street2:      w.Street2.String(),
		Street3:      w.Street3.String(),
		Street4:      w.Street4.String(),
		City:         w.City.String(),
		State:        w.State.String(),
		PostalCode:   w.PostalCode.String(),
		Country:      w.Country.String(),
		CountryCode:  w.CountryCode.String(),
		Email:        w.Email.String(),
		Telephone:    w.Telephone.String(),
		TelephoneExt: w.TelephoneExt.String(),
		Fax:          w.Fax.String(),
		FaxExt:       w.FaxExt.String(),
	}
}

func (w *WhoisDataBQ) parse() *domain.WhoisData {
	return &domain.WhoisData{
		DomainName:            w.DomainName,
		CreatedDate:           w.CreatedDate,
		UpdatedDate:           w.UpdatedDate,
		RegistrarName:         w.RegistrarName,
		RegistrarIANAID:       w.RegistrarIANAId.String(),
		Status:                w.Status.String(),
		Registrant:            *w.Registrant.parse(),
		AdministrativeContact: *w.AdministrativeContact.parse(),
		TechnicalContact:      *w.TechnicalContact.parse(),
		BillingContact:        *w.BillingContact.parse(),
		ZoneContact:           *w.ZoneContact.parse(),
		Header:                w.Header.String(),
		Footer:                w.Footer.String(),
		EstimatedDomainAge:    int(w.EstimatedDomainAge.Int64),
		Ips:                   w.Ips,
		LastUpdated:           w.LastUpdated,
	}
}

type ARecordBQ struct {
	CreatedAt time.Time `bigquery:"created_at"`
	UpdatedAt time.Time `bigquery:"updated_at"`
	IP        string    `bigquery:"ip"`
}

func newARecordBQ(record domain.ARecord) ARecordBQ {
	return ARecordBQ{
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		IP:        record.IP,
	}
}

func (a *ARecordBQ) parse() domain.ARecord {
	return domain.ARecord{
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

func newAAAARecordBQ(record domain.AAAARecord) AAAARecordBQ {
	return AAAARecordBQ{
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		IPV6:      record.IPV6,
	}
}

func (a *AAAARecordBQ) parse() domain.AAAARecord {
	return domain.AAAARecord{
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

func newMXRecordBQ(record domain.MXRecord) MXRecordBQ {
	return MXRecordBQ{
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		Mx:        record.Mx,
	}
}

func (a *MXRecordBQ) parse() domain.MXRecord {
	return domain.MXRecord{
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

func newSOARecordBQ(record domain.SOARecord) SOARecordBQ {
	return SOARecordBQ{
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
		NS:        bigquery.NullString{StringVal: record.NS, Valid: record.NS != ""},
		MBox:      bigquery.NullString{StringVal: record.MBox, Valid: record.MBox != ""},
		Serial:    bigquery.NullInt64{Int64: int64(record.Serial), Valid: record.Serial != 0},
	}
}

func (a *SOARecordBQ) parse() domain.SOARecord {
	return domain.SOARecord{
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

func newSitemapBQ(record domain.Sitemap) SitemapBQ {
	return SitemapBQ{
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
		SitemapLoc: bigquery.NullString{StringVal: record.SitemapLoc, Valid: record.SitemapLoc != ""},
	}
}

func (a *SitemapBQ) parse() *domain.Sitemap {
	return &domain.Sitemap{
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

func newMatchedDomainBQ(record domain.MatchedDomain) MatchedDomainBQ {
	return MatchedDomainBQ{
		CreatedAt:  record.CreatedAt,
		UpdatedAt:  record.UpdatedAt,
		DomainName: record.DomainName,
	}
}

func (a *MatchedDomainBQ) parse() domain.MatchedDomain {
	return domain.MatchedDomain{
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
		DomainName: a.DomainName,
	}
}
