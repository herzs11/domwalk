package domains

import (
	"testing"
	"time"
)

func TestDomains(t *testing.T) {
	return
	doms := []string{
		"levistrauss.com",
		"google.com",
		"levi.com",
		"ukg.com",
	}

	for _, dom := range doms {
		d, err := NewDomain(dom)
		// skip domains that are already in the table
		if err != nil {
			t.Errorf("Error parsing domain %s: %s\n", dom, err)
		}
		d.GetDNSRecords()
		d.GetRedirectDomains()
		d.GetCertSANs()
		d.GetDomainsFromSitemap()
		time.Sleep(1 * time.Second)
		d2, err := NewDomain("levistrauss.com")
		if err != nil {
			t.Errorf("Error parsing domain %s: %s\n", d2.DomainName, err)
		}
		d2.GetCertSANs()
	}

}
