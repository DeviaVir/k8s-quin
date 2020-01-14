package quin

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	pings           *prometheus.HistogramVec
	k8sConnFailures *prometheus.CounterVec
)

// RegisterMetrics initiate prometheus exported labels
func RegisterMetrics() {
	log.Println("registering prometheus metrics")
	pings = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "ping_seconds",
		Help: "Ping RTT in seconds",
		// In ms 1, 5 , 10, 20, 35, 50, 75, 100, 150, 200, 500, 1000, 2000, 5000
		Buckets: []float64{0.001, 0.005, 0.01, 0.02, 0.035, 0.05, 0.075, 0.1, 0.15, 0.2, 0.5, 1, 2, 5},
	}, []string{"hostname"})
	prometheus.MustRegister(pings)

	k8sConnFailures = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "k8s_conn_failed_count",
		Help: "number of failed attempts to communicate with k8s cluster",
	}, []string{"source"})
	prometheus.MustRegister(k8sConnFailures)
}
