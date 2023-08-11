package tpp

// TUTO EXAMPLE
// type metrics struct {
// 	cpuTemp    prometheus.Gauge
// 	hdFailures *prometheus.CounterVec
// }

// func NewMetrics(reg prometheus.Registerer) *metrics {
// 	m := &metrics{
// 		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
// 			Name: "cpu_temperature_celsius",
// 			Help: "Current temperature of the CPU.",
// 		}),
// 		hdFailures: prometheus.NewCounterVec(
// 			prometheus.CounterOpts{
// 				Name: "hd_errors_total",
// 				Help: "Number of hard-disk errors.",
// 			},
// 			[]string{"device"},
// 		),
// 	}
// 	reg.MustRegister(m.cpuTemp)
// 	reg.MustRegister(m.hdFailures)
// 	return m
// }

// func incTemperature(m *metrics) {
// 	go func() {
// 		init := 15
// 		for {
// 			init = init + 1
// 			m.cpuTemp.Set(float64(init + 1))
// 			time.Sleep(2 * time.Second)
// 		}
// 	}()
// }
