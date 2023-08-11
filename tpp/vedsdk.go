package tpp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

// Subject Alternative Name types
type Sans struct {
	Dns []string `json:"DNS"`
	Ip  []string `json:"IP"`
}

// Details of X509 certificates
type X509 struct {
	Cn         string    `json:"CN"`
	Sans       Sans      `json:"SANS"`
	Serial     string    `json:"Serial"`
	Thumbprint string    `json:"Thumbprint"`
	ValidFrom  time.Time `json:"ValidFrom"`
	ValidTo    time.Time `json:"ValidTo"`
}

// General certificates details
type Certificate struct {
	CreatedOn string `json:"CreatedOn"`
	Dn        string `json:"DN"`
	Guid      string `json:"Guid"`
	Name      string `json:"Name"`
	ParentDn  string `json:"ParentDn"`
	X509      X509   `json:"X509"`
}

type Certificates struct {
	Certificates []Certificate
	TotalCount   int
}

func (v *VenafiTpp) TotalCertificates(parameters string) Certificates {
	var certs Certificates
	req, err := http.NewRequest(
		"GET",
		string(v.Url+"/vedsdk/certificates"+parameters),
		nil,
	)
	CheckErr(err)
	req.Header.Add("Authorization", "Bearer "+v.AccessToken)
	client := &http.Client{}
	log.Printf("http_request; type=GET; endpoint=%s", string(v.Url+"/vedsdk/certificates"+parameters))
	resp, err := client.Do(req)
	CheckErr(err)

	defer resp.Body.Close()

	response, err := io.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(response, &certs)
	CheckErr(err)
	log.Printf("http_response; msg=Query %d certificate(s).", certs.TotalCount)
	return certs

}

func (v *VenafiTpp) QueryAll() {
	var certs Certificates
	var limit, offset int = 100, 0
	v.Certificates = nil

	log.Print("system; request; QueryAll() called.")
	for totalCerts := 100; totalCerts >= offset; {
		// var certsTmp Certificates
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/vedsdk/certificates?limit=%d&offset=%d", v.Url, limit, offset),
			nil,
		)
		log.Printf("http_request; type=GET; endpoint=%s", fmt.Sprintf("%s/vedsdk/certificates?limit=%d&offset=%d", v.Url, limit, offset))
		CheckErr(err)
		req.Header.Add("Authorization", "Bearer "+v.AccessToken)
		client := &http.Client{}
		resp, err := client.Do(req)
		CheckErr(err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		CheckErr(err)
		err = json.Unmarshal(body, &certs)
		CheckErr(err)
		v.Certificates = append(v.Certificates, certs.Certificates...)
		offset += limit
		totalCerts = certs.TotalCount
	}
	log.Printf("system; QueryAll result: %d certificates collected (%s)", len(v.Certificates), v.Url)
}

// Query all certificates in an error state.
func (v *VenafiTpp) QueryError() {
	var certs Certificates
	var limit, offset int = 100, 0

	log.Print("system; request; QueryError() called.")
	for totalCerts := 100; totalCerts >= offset; {
		// var certsTmp Certificates
		req, err := http.NewRequest(
			"GET",
			fmt.Sprintf("%s/vedsdk/certificates?Stage=200,500,800&limit=%d&offset=%d", v.Url, limit, offset),
			nil,
		)
		log.Printf("http_request; type=GET; endpoint=%s", fmt.Sprintf("%s/vedsdk/certificates?Stage=200,500,800&limit=%d&offset=%d", v.Url, limit, offset))
		CheckErr(err)
		req.Header.Add("Authorization", "Bearer "+v.AccessToken)
		client := &http.Client{}
		resp, err := client.Do(req)
		CheckErr(err)
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		CheckErr(err)
		err = json.Unmarshal(body, &certs)
		CheckErr(err)
		v.Certificates = append(v.Certificates, certs.Certificates...)
		offset += limit
		totalCerts = certs.TotalCount
	}
	log.Printf("system; QueryError result: %d certificates collected (%s)", len(v.Certificates), v.Url)
}

// Return number of valid certificates.
func (v VenafiTpp) GetValid() int {
	var valid []Certificate
	for _, details := range v.Certificates {
		if details.X509.ValidTo.After(time.Now()) {
			valid = append(valid, details)
		}
	}
	log.Printf("system; msg=%d certificates are still valid.", len(valid))
	return len(valid)
}

// Return the number of certificates expiring soon.
func (v VenafiTpp) ExpiringSoon(before int) int {
	var expire []Certificate
	for _, details := range v.Certificates {
		if details.X509.ValidTo.After(time.Now()) {
			if details.X509.ValidTo.Before(time.Now().AddDate(0, 0, before)) {
				expire = append(expire, details)
			}
		}
	}
	log.Printf("system; msg=%d certificates are expiring soon.", len(expire))
	return len(expire)
}

// Return the full list of policies.
func (v VenafiTpp) GetPolicies() []string {
	var policies []string
	for _, cert := range v.Certificates {
		dn_split := strings.Split(cert.Dn, "\\")
		dn := strings.Join(dn_split[0:len(dn_split)-1], "\\")

		// Get list of unique Dn
		if !IsIn(policies, dn) {
			policies = append(policies, dn)
		}
	}
	log.Printf("system; msg=%d unique policies detected.", len(policies))
	return policies
}

// Return a map with the number of certificate in each policies.
func (v VenafiTpp) GetCertPerPolicy(policies []string) map[string]int {
	counter := make(map[string]int)

	// Init map with policy names and counter 0
	for _, p := range policies {
		counter[p] = 0
	}

	for _, cert := range v.Certificates {
		p := strings.Split(cert.Dn, "\\")
		counter[(strings.Join(p[0:len(p)-1], "\\"))] += 1
	}
	return counter
}

// Return a map with the number of valid certificate in each policies.
func (v VenafiTpp) GetCertValidPerPolicy(policies []string) map[string]int {
	counter := make(map[string]int)

	// Init map with policy names and counter 0
	for _, p := range policies {
		counter[p] = 0
	}

	for _, cert := range v.Certificates {
		if cert.X509.ValidTo.After(time.Now().UTC()) {
			p := strings.Split(cert.Dn, "\\")
			counter[(strings.Join(p[0:len(p)-1], "\\"))] += 1
		}
	}
	return counter
}
