package domains

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/weppos/publicsuffix-go/publicsuffix"
)

type WebRedirectDomain struct {
	MatchedDomain
}

func (d *Domain) GetRedirectDomains() error {
	d.LastRanWebRedirect = time.Now()
	hosts := make(map[string]bool)
	finalURL := fmt.Sprintf("https://%s", d.DomainName)
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // Maximum amount of time to wait for a dial to complete
			KeepAlive: 3 * time.Second, // Keep-alive period for an active network connection
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second, // Maximum amount of time to wait for a TLS handshake
		ResponseHeaderTimeout: 5 * time.Second, // Maximum amount of time to wait for a server's response headers
		ExpectContinueTimeout: 1 * time.Second, // Maximum amount of time to wait for a 100-continue response from the server
	}
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			dom, err := publicsuffix.ParseFromListWithOptions(
				publicsuffix.DefaultList, req.URL.Hostname(), &publicsuffix.FindOptions{IgnorePrivate: false},
			)
			if err == nil && dom != nil {
				if dn := fmt.Sprintf("%s.%s", dom.SLD, dom.TLD); dn != d.DomainName {
					hosts[fmt.Sprintf("%s.%s", dom.SLD, dom.TLD)] = true
				}
			}
			finalURL = req.URL.String()
			return nil
		},
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	// Make the initial request
	resp, err := client.Get(fmt.Sprintf("http://%s", d.DomainName))
	if err != nil {
		d.SuccessfulWebLanding = false
		d.WebRedirectDomains = []WebRedirectDomain{}
		return fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	d.WebRedirectURLFinal = finalURL
	if len(hosts) == 0 {
		d.SuccessfulWebLanding = true
		d.WebRedirectDomains = []WebRedirectDomain{}
		return nil
	}
	now := time.Now()
	wrs := []WebRedirectDomain{}
	for host := range hosts {
		rdom, err := NewDomain(host)
		if err != nil {
			log.Println(err)
			continue
		}
		wr := WebRedirectDomain{MatchedDomain{CreatedAt: now, UpdatedAt: now, DomainName: rdom.DomainName}}
		wrs = append(wrs, wr)
	}
	d.WebRedirectDomains = wrs
	return nil
}
