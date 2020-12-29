/*
Copyright 2020 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
