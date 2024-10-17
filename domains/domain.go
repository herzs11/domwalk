package domains

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/temoto/robotstxt"
	"github.com/weppos/publicsuffix-go/publicsuffix"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Domain struct {
	DomainName            string                 `json:"domainName,omitempty"`
	CreatedAt             time.Time              `json:"createdAt,omitempty"`
	UpdatedAt             time.Time              `json:"updatedAt,omitempty"`
	NonPublicDomain       bool                   `json:"nonPublicDomain,omitempty"`
	Hostname              string                 `json:"hostname,omitempty"`
	Subdomain             string                 `json:"subdomain,omitempty"`
	Suffix                string                 `json:"suffix,omitempty"`
	SuccessfulWebLanding  bool                   `json:"successfulWebLanding,omitempty"`
	WebRedirectURLFinal   string                 `json:"webRedirectURLFinal,omitempty"`
	LastRanWebRedirect    time.Time              `json:"lastRanWebRedirect,omitempty"`
	LastRanDns            time.Time              `json:"lastRanDNS,omitempty"`
	LastRanCertSans       time.Time              `json:"lastRanCertSANs,omitempty"`
	LastRanSitemapParse   time.Time              `json:"lastRanSitemapParse,omitempty"`
	ARecords              []ARecord              `json:"aRecords"`
	AAAARecords           []AAAARecord           `json:"aaaaRecords"`
	MXRecords             []MXRecord             `json:"mxRecords"`
	SOARecords            []SOARecord            `json:"soaRecords"`
	Sitemaps              []*Sitemap             `json:"sitemaps"`
	WebRedirectDomains    []WebRedirectDomain    `json:"webRedirectDomains"`
	CertSANs              []CertSansDomain       `json:"certSANs"`
	SitemapWebDomains     []SitemapWebDomain     `json:"sitemapWebDomains"`
	SitemapContactDomains []SitemapContactDomain `json:"sitemapContactDomains"`

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
