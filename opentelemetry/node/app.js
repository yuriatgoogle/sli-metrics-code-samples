"use strict";

// imports
const express = require ('express')
const { MeterProvider } = require('@opentelemetry/metrics');
const { PrometheusExporter } = require('@opentelemetry/exporter-prometheus');

function sleep (n) {
    Atomics.wait(new Int32Array(new SharedArrayBuffer(4)), 0, 0, n);
}

// set up prometheus exporter
const prometheusPort = 8081;
const app = express();
const exporter = new PrometheusExporter(
  {
    startServer: true,
    port: prometheusPort
  },
  () => {
    console.log("prometheus scrape endpoint: http://localhost:"
      + prometheusPort 
      + "/metrics");
  }
);
const meter = new MeterProvider({
    exporter,
    interval: 2000,
}).getMeter('example-prometheus');

var measuredLatency = 0;

// define metrics with descriptions
const requestCount = meter.createCounter("request_count", {
  description: "Counts total number of requests"
});
const errorCount = meter.createCounter("error_count", {
    description: "Counts total number of errors"
});
meter.createValueObserver('response_latency', {
    description: 'Records latency of response',
  }, async (observerResult) => { // this callback is called once per each interval
    await new Promise((resolve) => {
      setTimeout(()=> {resolve()}, 50);
    });
    observerResult.observe(measuredLatency, { testLabel: 'test'});
  });

// set metric values on request
app.get('/', (req, res) => {
    // start latency timer
    const requestReceived = new Date().getTime();
    console.log('request made');
    // increment total requests counter
    requestCount.add(1);
    // return an error 10% of the time
    if ((Math.floor(Math.random() * 100)) > 90) {
        // increment error counter
        errorCount.add(1);
        // return error code
        res.status(500).send("error!")
    } 
    else {
        // delay for a bit
        sleep(Math.floor(Math.random()*1000));
        measuredLatency = new Date().getTime() - requestReceived;
        res.status(200).send("did not delay - success in " + measuredLatency + " ms");
    }
})

app.listen(8080, () => console.log(`Example app listening on port 8080!`))