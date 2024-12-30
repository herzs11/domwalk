package domains

import (
	`fmt`
	`testing`
)

func TestDomainSAN(t *testing.T) {
	d, err := NewDomain("johnmuirhealth.com")
	if err != nil {
		t.Log(err)
	}
	d.GetCertSANs()
	fmt.Println(d.OrgName)
}
