package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

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
	requestReceived := time.Now()
	requestCount.Inc()

	// fail 10% of the time
	if rand.Intn(100) > 90 {
		failedRequestCount.Inc()
		fmt.Fprintf(w, "error!")
		responseLatency.Observe(time.Since(requestReceived).Seconds())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		responseLatency.Observe(time.Since(requestReceived).Seconds())
		fmt.Fprintf(w, "Hello")
	}
}
