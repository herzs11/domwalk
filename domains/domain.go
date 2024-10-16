package domains

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"cloud.google.com/go/bigquery"
	"github.com/temoto/robotstxt"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Domain struct {
	DomainName            string                 `json:"domainName,omitempty" bigquery:"domain_name"`
	CreatedAt             time.Time              `json:"createdAt,omitempty" bigquery:"created_at"`
	UpdatedAt             time.Time              `json:"updatedAt,omitempty" bigquery:"updated_at"`
	NonPublicDomain       bool                   `json:"nonPublicDomain,omitempty" bigquery:"non_public_domain"`
	Hostname              string                 `json:"hostname,omitempty" bigquery:"hostname"`
	Subdomain             string                 `json:"subdomain,omitempty" bigquery:"subdomain"`
	Suffix                string                 `json:"suffix,omitempty" bigquery:"suffix"`
	SuccessfulWebLanding  bool                   `json:"successfulWebLanding,omitempty" bigquery:"successful_web_landing"`
	WebRedirectURLFinal   string                 `json:"webRedirectURLFinal,omitempty" bigquery:"web_redirect_url_final"`
	LastRanWebRedirect    time.Time              `json:"lastRanWebRedirect,omitempty" bigquery:"last_ran_web_redirect"`
	LastRanDns            time.Time              `json:"lastRanDNS,omitempty" bigquery:"last_ran_dns"`
	LastRanCertSans       time.Time              `json:"lastRanCertSANs,omitempty" bigquery:"last_ran_cert_sans"`
	LastRanSitemapParse   time.Time              `json:"lastRanSitemapParse,omitempty" bigquery:"last_ran_sitemap_parse"`
	ARecords              []ARecord              `json:"aRecords" bigquery:"-"`
	AAAARecords           []AAAARecord           `json:"aaaaRecords" bigquery:"-"`
	MXRecords             []MXRecord             `json:"mxRecords" bigquery:"-"`
	SOARecords            []SOARecord            `json:"soaRecords" bigquery:"-"`
	Sitemaps              []*Sitemap             `json:"-" bigquery:"-"`
	WebRedirectDomains    []WebRedirectDomain    `json:"webRedirectDomains" bigquery:"-"`
	CertSANs              []CertSansDomain       `json:"certSANs" bigquery:"-"`
	SitemapWebDomains     []SitemapWebDomain     `json:"sitemapWebDomains" bigquery:"-"`
	SitemapContactDomains []SitemapContactDomain `json:"sitemapContactDomains" bigquery:"-"`

	sitemapURLs  []string
	contactPages []string

	*robotstxt.RobotsData
}

func (d *Domain) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "domain_name"}},
			DoUpdates: clause.Assignments(
				map[string]interface{}{
					"updated_at":             gorm.Expr("MAX(updated_at, excluded.updated_at)"),
					"last_ran_web_redirect":  gorm.Expr("MAX(last_ran_web_redirect, excluded.last_ran_web_redirect)"),
					"last_ran_dns":           gorm.Expr("MAX(last_ran_dns, excluded.last_ran_dns)"),
					"last_ran_cert_sans":     gorm.Expr("MAX(last_ran_cert_sans, excluded.last_ran_cert_sans)"),
					"last_ran_sitemap_parse": gorm.Expr("MAX(last_ran_sitemap_parse, excluded.last_ran_sitemap_parse)"),
				},
			),
		},
	)
	return nil
}

func (d *Domain) ToBQ() DomainBQ {
	return DomainBQ{
		DomainName:           d.DomainName,
		CreatedAt:            d.CreatedAt,
		UpdatedAt:            d.UpdatedAt,
		NonPublicDomain:      d.NonPublicDomain,
		Hostname:             bigquery.NullString{d.Hostname, d.Hostname != ""},
		Subdomain:            bigquery.NullString{d.Subdomain, d.Subdomain != ""},
		Suffix:               bigquery.NullString{d.Suffix, d.Suffix != ""},
		SuccessfulWebLanding: d.SuccessfulWebLanding,
		WebRedirectURLFinal:  bigquery.NullString{d.WebRedirectURLFinal, d.WebRedirectURLFinal != ""},
		LastRanWebRedirect:   d.LastRanWebRedirect,
		LastRanDns:           d.LastRanDns,
		LastRanCertSans:      d.LastRanCertSans,
		LastRanSitemapParse:  d.LastRanSitemapParse,
	}
}

func (d *Domain) parseDomain() error {
	dom, err := publicsuffix.ParseFromListWithOptions(
		publicsuffix.DefaultList, d.DomainName, &publicsuffix.FindOptions{IgnorePrivate: true},
	)
	if err != nil {
		return err
	}
	if dom == nil {
		d.NonPublicDomain = true
		return errors.New("Unable to parse domain from public suffix list")
	}
	d.DomainName = fmt.Sprintf("%s.%s", strings.ToLower(dom.SLD), strings.ToLower(dom.TLD))
	d.Hostname = strings.ToLower(dom.SLD)
	d.Subdomain = strings.ToLower(dom.TRD)
	d.Suffix = strings.ToLower(dom.TLD)
	d.NonPublicDomain = false
	return nil
}

func NewDomain(domain_name string) (*Domain, error) {
	dn := strings.TrimSpace(domain_name)
	d := &Domain{DomainName: dn}
	err := d.parseDomain()
	return d, err
}
