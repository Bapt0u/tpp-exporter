package tpp

import (
	"github.com/prometheus/client_golang/prometheus"
)

type VenafiMetrics struct {
	CertificatesTotal          prometheus.Gauge
	CertificatesExpiringSoon   prometheus.Gauge
	CertificatesPerPolicy      *prometheus.GaugeVec
	CertificatesValidPerPolicy *prometheus.GaugeVec
	// CertificatesError          prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *VenafiMetrics {
	m := &VenafiMetrics{
		CertificatesTotal: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "venafi_certificate_total",
			Help: "Total of available certificates.",
		}),
		// CertificatesError: prometheus.NewGauge(prometheus.GaugeOpts{
		// 	Name: "venafi_certificate_error",
		// 	Help: "Number of certificate in error state (code 200, 500, 800)",
		// }),
		CertificatesExpiringSoon: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "venafi_certificate_expiring_soon",
			Help: "Number of certificate expiring soon.",
		}),
		CertificatesPerPolicy: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "venafi_certificate_per_policy",
			Help: "Number of certificate in each policy.",
		}, []string{"policy"}),
		CertificatesValidPerPolicy: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "venafi_certificate_valid_per_policy",
			Help: "Number of valid certificate in each policy.",
		}, []string{"policy"}),
	}

	// reg.MustRegister(m.CertificatesError)
	reg.MustRegister(m.CertificatesTotal)
	reg.MustRegister(m.CertificatesExpiringSoon)
	reg.MustRegister(m.CertificatesPerPolicy)
	reg.MustRegister(m.CertificatesValidPerPolicy)

	return m
}
