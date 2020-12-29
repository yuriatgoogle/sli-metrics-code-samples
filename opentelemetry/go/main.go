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
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	mexporter "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
	"go.opencensus.io/tag"

	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric/controller/push"
	"go.opentelemetry.io/otel/sdk/resource"

	"os"
)

var (
	env          = os.Getenv("ENV")
	projectID    = os.Getenv("PROJECT_ID")
	latencyValue = 0.0
)

func main() {

	// Initialization. In order to pass the credentials to the exporter,
	// prepare credential file following the instruction described in this doc.
	// https://pkg.go.dev/golang.org/x/oauth2/google?tab=doc#FindDefaultCredentials
	opts := []mexporter.Option{
		mexporter.WithProjectID(projectID),
	}

	// NOTE: In current implementation of exporter, this resource is ignored
	resOpt := push.WithResource(resource.NewWithAttributes(
		label.String("instance_id", "abc123"),
		label.String("application", "example-app"),
	))
	pusher, err := mexporter.InstallNewPipeline(opts, resOpt)
	if err != nil {
		log.Fatalf("Failed to establish pipeline: %v", err)
	}
	defer pusher.Stop()

	meter := pusher.MeterProvider().Meter("cloudmonitoring/example")
	// Register request counter
	totalRequestsCounter := metric.Must(meter).NewInt64Counter("total_requests")
	// Register error counter
	errorsCounter := metric.Must(meter).NewInt64Counter("failed_requests")
	// Register latency observer
	olabels := []label.KeyValue{}
	callback := func(_ context.Context, result metric.Float64ObserverResult) {
		v := latencyValue
		result.Observe(v, olabels...)
	}
	metric.Must(meter).NewFloat64ValueObserver("response_latency", callback)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := tag.New(context.Background())
		requestReceived := time.Now()
		// count the request
		totalRequestsCounter.Add(ctx, 1)

		// fail 10% of the time
		if rand.Intn(100) > 90 {
			// count the failed request
			errorsCounter.Add(ctx, 1)
			fmt.Fprintf(w, "error!")
			// record latency for failure
			latencyValue = float64(time.Since(requestReceived).Seconds())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		} else {
			delay := time.Duration(rand.Intn(1000)) * time.Millisecond
			time.Sleep(delay)
			// record latency for success
			latencyValue = float64(time.Since(requestReceived).Seconds())
			fmt.Fprintf(w, "Responded after "+delay.String())
		}
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
