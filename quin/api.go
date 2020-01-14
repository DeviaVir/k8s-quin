package quin

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// RunServer handle traffic
func RunServer(port int) {
	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
