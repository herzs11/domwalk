package domains

import (
	"crypto/tls"
	"log"
	"net"
	"time"
)

type CertSansDomain struct {
	MatchedDomain
}

func (d *Domain) GetCertSANs() error {
	d.LastRanCertSans = time.Now()
	dom := d.DomainName
	dialer := net.Dialer{Timeout: 3 * time.Second}
	// Use the dialTLS function to connect
	conn, err := dialer.Dial(
		"tcp", dom+":443",
	)
	if err != nil {
		return err
	}
	defer conn.Close()
	tlsConn := tls.Client(
		conn, &tls.Config{
			ServerName:         dom,
			InsecureSkipVerify: true,
		},
	)
	err = tlsConn.Handshake()
	if err != nil {
		return err
	}
	defer tlsConn.Close()
	cert := tlsConn.ConnectionState().PeerCertificates[0]
	now := time.Now()
	domsFound := make(map[string]CertSansDomain)
	for _, df := range d.CertSANs {
		domsFound[df.DomainName] = df
	}
	for _, san := range cert.DNSNames {
		dm, err := NewDomain(san)
		if err != nil {
			log.Println("Error parsing domain: ", err)
			continue
		}
		if dm.DomainName == dom {
			continue
		}
		if c, exists := domsFound[dm.DomainName]; !exists {
			certSAN := CertSansDomain{MatchedDomain{CreatedAt: now, UpdatedAt: now, DomainName: dm.DomainName}}
			domsFound[dm.DomainName] = certSAN
		} else {
			c.UpdatedAt = now
			domsFound[dm.DomainName] = c
		}
	}
	var cs []CertSansDomain
	for _, c := range domsFound {
		cs = append(cs, c)
	}
	d.CertSANs = cs
	return nil
}
