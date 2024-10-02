package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/temoto/robotstxt"
	"github.com/weppos/publicsuffix-go/publicsuffix"
)

type Domain struct {
	DomainName            string `json:"domainName,omitempty" gorm:"primaryKey"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
	NonPublicDomain       bool                   `json:"nonPublicDomain,omitempty"`
	Hostname              string                 `json:"hostname,omitempty"`
	Subdomain             string                 `json:"subdomain,omitempty"`
	Suffix                string                 `json:"suffix,omitempty"`
	SuccessfulWebLanding  bool                   `json:"successfulWebLanding,omitempty"`
	WebRedirectURLFinal   string                 `json:"webRedirectURLFinal,omitempty"`
	LastRanWebRedirect    time.Time              `json:"lastRanWebRedirect,omitempty"`
	LastRandDNS           time.Time              `json:"lastRanDNS,omitempty"`
	LastRanCertSANs       time.Time              `json:"lastRanCertSANs,omitempty"`
	LastRanSitemapParse   time.Time              `json:"lastRanSitemapParse,omitempty"`
	WebRedirectDomains    []WebRedirect          `json:"landedWebHost,omitempty" gorm:"foreignKey:MatchedDomainName;references:DomainName"`
	CertSANs              []CertSAN              `json:"certSANs,omitempty" gorm:"foreignKey:MatchedDomainName;references:DomainName"`
	ARecords              []ARecord              `gorm:"foreignKey:DomainName;references:DomainName"`
	AAAARecords           []AAAARecord           `gorm:"foreignKey:DomainName;references:DomainName"`
	MXRecords             []MXRecord             `gorm:"foreignKey:DomainName;references:DomainName"`
	SOARecords            []SOARecord            `gorm:"foreignKey:DomainName;references:DomainName"`
	Sitemaps              []*Sitemap             `json:"sitemaps,omitempty" gorm:"foreignKey:DomainName;references:DomainName"`
	SitemapWebDomains     []SitemapWebDomain     `json:"sitemapWebDomains,omitempty" gorm:"foreignKey:MatchedDomainName;references:DomainName"`
	SitemapContactDomains []SitemapContactDomain `json:"sitemapContactDomains,omitempty" gorm:"foreignKey:MatchedDomainName;references:DomainName"`

	sitemapURLs  []string `json:"-" gorm:"-"`
	contactPages []string `json:"-" gorm:"-"`

	*robotstxt.RobotsData `json:"-" gorm:"-:all"`
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
	d.DomainName = fmt.Sprintf("%s.%s", dom.SLD, dom.TLD)
	d.Hostname = dom.SLD
	d.Subdomain = dom.TRD
	d.Suffix = dom.TLD
	d.NonPublicDomain = false
	return nil
}

func NewDomain(domain_name string) (*Domain, error) {
	d := &Domain{DomainName: domain_name}
	err := d.parseDomain()
	return d, err
}
