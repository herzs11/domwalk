package domains

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"
)

type CertSansDomain struct {
	MatchedDomain `json:"certSanDomain,omitempty"`
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
	now := time.Now()
	domsFound := make(map[string]CertSansDomain)
	for _, df := range d.CertSANs {
		domsFound[df.DomainName] = df
	}
	var cs []CertSansDomain
	for _, san := range cert.DNSNames {
		dm, err := NewDomain(san)
		if err != nil {
			log.Println("Error parsing domain: ", err)
			continue
		}
		if dm.DomainName == dom {
			continue
		}
		if _, exists := domsFound[dm.DomainName]; !exists {
			certSAN := CertSansDomain{MatchedDomain{DomainName: dm.DomainName}}
			domsFound[dm.DomainName] = certSAN
			cs = append(cs, certSAN)
		} else {
			certSAN := domsFound[dm.DomainName]
			certSAN.UpdatedAt = now
			cs = append(cs, certSAN)
		}
	}
	d.CertSANs = cs
	return nil
}
