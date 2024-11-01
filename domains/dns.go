package domains

import (
	"errors"
	"log"
	"time"

	"github.com/miekg/dns"
)

type AAAARecord struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	IPV6      string    `json:"ip_v6"`
}

type ARecord struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	IP        string    `json:"ip"`
}

type SOARecord struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	NS        string    `json:"ns,omitempty"`
	MBox      string    `json:"MBox,omitempty"`
	Serial    uint32    `json:"serial,omitempty"`
}

type MXRecord struct {
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
	Mx        string    `json:"mx,omitempty"`
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
	foundMx := make(map[string]MXRecord)
	for _, m := range d.MXRecords {
		foundMx[m.Mx] = m
	}
	now := time.Now()
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.MX); ok {
			if m, ok := foundMx[a.Mx]; ok {
				m.UpdatedAt = now
				foundMx[a.Mx] = m
			} else {
				r := MXRecord{
					CreatedAt: now,
					UpdatedAt: now,
					Mx:        a.Mx,
				}
				foundMx[a.Mx] = r
			}
		}
	}
	var mxs []MXRecord
	for _, mx := range foundMx {
		mxs = append(mxs, mx)
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
	foundA := make(map[string]ARecord)
	for _, a := range d.ARecords {
		foundA[a.IP] = a
	}
	now := time.Now()
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.A); ok {
			if ip, ok := foundA[a.A.String()]; ok {
				ip.UpdatedAt = now
				foundA[a.A.String()] = ip
				continue
			}
			r := ARecord{IP: a.A.String(), CreatedAt: now, UpdatedAt: now}
			foundA[a.A.String()] = r
		}
	}
	var ips []ARecord
	for _, ip := range foundA {
		ips = append(ips, ip)
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
	foundA := make(map[string]AAAARecord)
	for _, a := range d.AAAARecords {
		foundA[a.IPV6] = a
	}
	now := time.Now()
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.AAAA); ok {
			if ip, ok := foundA[a.AAAA.String()]; ok {
				ip.UpdatedAt = now
				foundA[a.AAAA.String()] = ip
				continue
			}
			r := AAAARecord{IPV6: a.AAAA.String(), CreatedAt: now, UpdatedAt: now}
			foundA[a.AAAA.String()] = r
		}
	}
	var ips []AAAARecord
	for _, ip := range foundA {
		ips = append(ips, ip)
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
	now := time.Now()
	var soas []SOARecord
	for _, ans := range r.Answer {
		if a, ok := ans.(*dns.SOA); ok {
			soas = append(soas, SOARecord{CreatedAt: now, UpdatedAt: now, NS: a.Ns, MBox: a.Mbox, Serial: a.Serial})
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
