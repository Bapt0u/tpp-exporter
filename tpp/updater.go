package tpp

import (
	"log"
	"time"
)

func (m *VenafiMetrics) Updater(v *VenafiTpp, timer int) {
	ticker := time.NewTicker(time.Minute * time.Duration(timer))
	for ; true; <-ticker.C {
		defer ticker.Stop()

		v.QueryAll()
		m.UpdateTotalCert(v, timer)
		m.UpdateExpiringSoonCert(v, timer)
		m.UpdateCertPerPolicy(v, timer)
		m.UpdateValidCertPerPolicy(v, timer)
	}

}

// Update total of available certificate.
func (m *VenafiMetrics) UpdateTotalCert(v *VenafiTpp, timer int) {
	log.Printf("system; UpdateTotalCert() launched; timer=%d", timer)
	m.CertificatesTotal.Set(float64(len(v.Certificates)))
}

// Update the Gauge of expiring soon certificates.
func (m *VenafiMetrics) UpdateExpiringSoonCert(v *VenafiTpp, timer int) {
	log.Printf("system; UpdateExpiringSoonCert() launched; timer=%d", timer)
	m.CertificatesExpiringSoon.Set(float64(v.ExpiringSoon(v.ExpireSoon)))
}

// Update number of certificate per policy.
func (m *VenafiMetrics) UpdateCertPerPolicy(v *VenafiTpp, timer int) {
	log.Printf("system; UpdateCertPerPolicy() launched; timer=%d", timer)
	for k, e := range v.GetCertPerPolicy(v.GetPolicies()) {
		m.CertificatesPerPolicy.WithLabelValues(k).Set(float64(e))
	}
}

// Update number of valid certificate per policy.
func (m *VenafiMetrics) UpdateValidCertPerPolicy(v *VenafiTpp, timer int) {
	log.Printf("system; UpdateValidCertPerPolicy() launched; timer=%d", timer)
	for k, e := range v.GetCertValidPerPolicy(v.GetPolicies()) {
		m.CertificatesValidPerPolicy.WithLabelValues(k).Set(float64(e))
	}
}
