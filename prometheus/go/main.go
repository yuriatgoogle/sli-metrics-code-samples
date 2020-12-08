package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// set up metrics
var (
	requestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_request_count",
		Help: "total request count",
	})
	failedRequestCount = promauto.NewCounter(prometheus.CounterOpts{
		Name: "go_failed_request_count",
		Help: "failed request count",
	})
	responseLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "go_response_latency",
		Help: "response latencies",
	})
)

func main() {
	log.Printf("main function")
	http.HandleFunc("/", handle)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	requestCount.Inc()
	// start timer
	// fail
	fmt.Fprintf(w, "Hello")
}
