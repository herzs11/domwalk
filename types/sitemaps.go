package types

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/temoto/robotstxt"
)

type URL struct {
	Loc string `xml:"loc"`
}

type URLSet struct {
	URLs []URL `xml:"url"`
}

type SitemapIndex struct {
	Sitemaps []URL `xml:"sitemap"`
}

type Sitemap struct {
	CreatedAt  time.Time `bigquery:"created_at"`
	UpdatedAt  time.Time `bigquery:"updated_at"`
	DomainName string    `json:"domainName,omitempty" bigquery:"domain_name"`
	Sitemap    string    `json:"sitemap,omitempty" bigquery:"sitemap"`
}

func (d *Domain) GetDomainsFromSitemap() error {
	if !d.SuccessfulWebLanding {
		return fmt.Errorf("Domain has not successfully landed on the web")
	}
	d.LastRanSitemapParse = time.Now()
	err := d.getRobotstxt()
	if err != nil {
		return fmt.Errorf("Error fetching robots.txt: %v", err)
	}
	d.getURLsFromSitemaps()
	d.GetWebDomainsFromSitemap()
	err = d.GetContactDomainsFromSitemap()
	if err != nil {
		return fmt.Errorf("Error fetching contact domains: %v", err)
	}
	return nil
}

func (d *Domain) getRobotstxt() error {
	url_raw := d.WebRedirectURLFinal
	url_parsed, err := url.Parse(url_raw)
	if err != nil {
		return err
	}
	host_root := url_parsed.Scheme + "://" + url_parsed.Host
	if !strings.HasPrefix(host_root, "http") {
		host_root = "http://" + url_parsed.Host
	}
	resp, err := http.Get(host_root + "/robots.txt")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	robots, err := robotstxt.FromResponse(resp)
	if err != nil {
		return err
	}
	d.RobotsData = robots
	if len(robots.Sitemaps) > 11 {
		robots.Sitemaps = robots.Sitemaps[:11]
	}
	for _, sitemap := range robots.Sitemaps {
		d.Sitemaps = append(d.Sitemaps, &Sitemap{DomainName: d.DomainName, Sitemap: sitemap})
	}

	return nil
}

func (s *Sitemap) readSitemap() (URLSet, []*Sitemap, error) {
	// add proxy
	proxyURL, err := url.Parse("http://localhost:9000")
	if err != nil {
		fmt.Printf("Error parsing proxy URL: %s\n", err)
		return URLSet{}, nil, err
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // Maximum amount of time to wait for a dial to complete
			KeepAlive: 3 * time.Second, // Keep-alive period for an active network connection
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second, // Maximum amount of time to wait for a TLS handshake
		ResponseHeaderTimeout: 5 * time.Second, // Maximum amount of time to wait for a server's response headers
		ExpectContinueTimeout: 1 * time.Second, // Maximum amount of time to wait for a 100-continue response from the server
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // Overall timeout for the request
	}
	resp, err := client.Get(s.Sitemap)
	if err != nil {
		fmt.Printf("Error performing GET request: %s\n", err)
		return URLSet{}, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return URLSet{}, nil, fmt.Errorf("error fetching sitemap: received status code %d", resp.StatusCode)
	}

	// Read the sitemap
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return URLSet{}, nil, fmt.Errorf("error reading sitemap: %v", err)
	}

	var (
		sitemapIndex   SitemapIndex
		urlSet         URLSet
		subSiteMapURLs URLSet
		sms            []*Sitemap
	)

	err = xml.Unmarshal(body, &sitemapIndex)
	if err == nil && len(sitemapIndex.Sitemaps) > 0 {
		// It's a sitemap index, process each sub-sitemap
		for _, sitemap := range sitemapIndex.Sitemaps {
			if !strings.Contains(sitemap.Loc, ".xml") || sitemap.Loc == s.Sitemap {
				continue
			}
			s := &Sitemap{DomainName: s.DomainName, Sitemap: sitemap.Loc}
			sms = append(sms, s)
			urls, sitemaps, err := s.readSitemap()
			sms = append(sms, sitemaps...)
			if err != nil {
				log.Println(err)
				continue
			}
			subSiteMapURLs.URLs = append(subSiteMapURLs.URLs, urls.URLs...)
		}
	}

	err = xml.Unmarshal(body, &urlSet)
	if err != nil {
		return subSiteMapURLs, sms, fmt.Errorf("error parsing sitemap: %v", err)
	}
	urlSet.URLs = append(urlSet.URLs, subSiteMapURLs.URLs...)

	return urlSet, sms, nil
}

func (d *Domain) getURLsFromSitemaps() {
	var smsParsed = make(map[string]bool)
	var urlsFound = make(map[string]bool)
	for _, sitemap := range d.Sitemaps {
		if _, exists := smsParsed[sitemap.Sitemap]; exists {
			continue
		}
		smsParsed[sitemap.Sitemap] = true
		urls, sitemaps, err := sitemap.readSitemap()
		if err != nil {
			log.Println(err)
		}
		if len(urls.URLs) == 0 {
			continue
		}
		for _, url := range urls.URLs {
			if _, exists := urlsFound[url.Loc]; !exists {
				urlsFound[url.Loc] = true
				d.sitemapURLs = append(d.sitemapURLs, url.Loc)
			}
		}
		for _, s := range sitemaps {
			if _, exists := smsParsed[s.Sitemap]; !exists {
				d.Sitemaps = append(d.Sitemaps, s)
				smsParsed[s.Sitemap] = true
			}
		}
		if len(d.Sitemaps) > 15 {
			return
		}
		if len(d.sitemapURLs) > 10000 {
			return
		}
	}
}

func (d *Domain) GetWebDomainsFromSitemap() {
	var domsFound = make(map[string]bool)
	for _, u := range d.sitemapURLs {
		up, err := url.Parse(u)
		if err != nil {
			log.Println(err)
			continue
		}
		dom, err := NewDomain(up.Host)
		if err != nil {
			log.Println(err)
		}
		if _, exists := domsFound[dom.DomainName]; !exists {
			domsFound[dom.DomainName] = true
			sd := MatchedDomain{}
			d.SitemapWebDomains = append(d.SitemapWebDomains, sd)
		}
	}
}

func (d *Domain) GetContactDomainsFromSitemap() error {
	d.getContactPagesFromSitemap()
	if len(d.contactPages) == 0 {
		return fmt.Errorf("No contact pages found in sitemap")
	}

	proxyURL, err := url.Parse("http://localhost:9000")
	if err != nil {
		log.Printf("Error parsing proxy URL: %s\n", err)
		return err
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // Maximum amount of time to wait for a dial to complete
			KeepAlive: 3 * time.Second, // Keep-alive period for an active network connection
		}).DialContext,
		TLSHandshakeTimeout:   5 * time.Second, // Maximum amount of time to wait for a TLS handshake
		ResponseHeaderTimeout: 5 * time.Second, // Maximum amount of time to wait for a server's response headers
		ExpectContinueTimeout: 1 * time.Second, // Maximum amount of time to wait for a 100-continue response from the server
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second, // Overall timeout for the request
	}

	var domsFound = make(map[string]bool)
	for _, url := range d.contactPages {
		resp, err := client.Get(url)
		if err != nil {
			log.Printf("Error performing GET request: %s\n", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Error fetching contact page: received status code %d\n", resp.StatusCode)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error reading contact page: %v\n", err)
			continue
		}

		// Find email addresses in the body
		emailRegex := regexp.MustCompile(`(?:^|\s|,|;)([a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.(com|net|org|edu|gov|info|biz|io|co\.uk|co|us|de|ca)(?:\s|,|;|$))`)
		emails := emailRegex.FindAllString(string(body), -1)
		for _, email := range emails {
			dom, err := NewDomain(strings.Split(email, "@")[1])
			if err != nil {
				log.Println(err)
			}
			if _, exists := domsFound[dom.DomainName]; !exists {
				domsFound[dom.DomainName] = true
				sd := SitemapContactDomain{SitemapDomain{DomainName: d.DomainName, Domain: *dom}}
				d.SitemapContactDomains = append(d.SitemapContactDomains, sd)
			}
		}
	}
	return nil

}

func (d *Domain) getContactPagesFromSitemap() {
	for _, url := range d.sitemapURLs {
		if strings.Contains(url, "contact") {
			if allow := d.RobotsData.TestAgent(url, "*"); allow {
				d.contactPages = append(d.contactPages, url)
			}
		}
	}
	if len(d.contactPages) > 10 {
		d.contactPages = d.contactPages[:10]
	}
}
