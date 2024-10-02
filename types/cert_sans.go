package types

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"
)

type CertSAN struct {
	CreatedAt         time.Time `bigquery:"created_at"`
	UpdatedAt         time.Time `bigquery:"updated_at"`
	DomainName        string    `json:"reqDomainName,omitempty" bigquery:"domain_name"`
	MatchedDomainName string    `json:"certSanDomainName,omitempty" bigquery:"matched_domain_name"`
	Domain            Domain    `json:"certSanDomain,omitempty" gorm:"foreignKey:MatchedDomainName" bigquery:"-"`
}

func (d *Domain) GetCertSANs() error {
	d.LastRanCertSans = time.Now()
	proxyURL, err := url.Parse(os.Getenv("HTTPS_PROXY"))
	if err != nil {
		panic(err)
	}
	dom := d.DomainName
	// Create a DialTLS function that uses the proxy
	dialTLS := func(network, addr string) (net.Conn, error) {
		// Connect to the proxy
		proxyConn, err := net.Dial("tcp", proxyURL.Host)
		if err != nil {
			return nil, err
		}

		// Send CONNECT request to the proxy
		connectReq := fmt.Sprintf("CONNECT %s:443 HTTP/1.1\r\nHost: %s\r\n\r\n", dom, dom)
		_, err = proxyConn.Write([]byte(connectReq))
		if err != nil {
			return nil, err
		}

		// Read proxy response (should be 200 Connection established)
		// (You might need more robust handling here)
		var buf [1024]byte
		_, err = proxyConn.Read(buf[:])
		if err != nil {
			return nil, err
		}

		// Now establish the TLS connection over the proxy connection
		tlsConn := tls.Client(
			proxyConn, &tls.Config{
				ServerName:         dom,
				InsecureSkipVerify: true, // Use with caution!
			},
		)
		err = tlsConn.Handshake()
		if err != nil {
			return nil, err
		}

		return tlsConn, nil
	}

	// Use the dialTLS function to connect
	conn, err := dialTLS("tcp", dom+":443")
	if err != nil {
		return err
	}
	defer conn.Close()
	tlsConn, ok := conn.(*tls.Conn)
	if !ok {
		return fmt.Errorf("failed to cast connection to tls.Conn")
	}
	cert := tlsConn.ConnectionState().PeerCertificates[0]
	domsFound := make(map[string]bool)
	for _, san := range cert.DNSNames {
		dm, err := NewDomain(san)
		if err != nil {
			log.Println("Error parsing domain: ", err)
		}
		if dm.DomainName == dom {
			continue
		}
		if _, exists := domsFound[dm.DomainName]; !exists {
			domsFound[dm.DomainName] = true
			certSAN := CertSAN{DomainName: dom, Domain: *dm}
			d.CertSANs = append(d.CertSANs, certSAN)
		}
	}
	return nil
}
