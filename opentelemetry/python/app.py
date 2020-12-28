from flask import Flask
import random
import time
from opentelemetry import metrics
from opentelemetry.exporter.cloud_monitoring import (
    CloudMonitoringMetricsExporter,
)
from opentelemetry.sdk.metrics import Counter, MeterProvider

# set up exporter
metrics.set_meter_provider(MeterProvider())
meter = metrics.get_meter(__name__)
metrics.get_meter_provider().start_pipeline(
    meter, CloudMonitoringMetricsExporter("stack-doctor"), 5
)

# set up metrics
# total requests
requests_counter = meter.create_counter(
    name="total_requests_python",
    description="number of requests",
    unit="1",
    value_type=int,
)
# errors
errors_counter = meter.create_counter(
    name="total_errors_python",
    description="number of errors",
    unit="1",
    value_type=int,
)
# latency
response_latency = meter.create_valuerecorder(
    name="response_latency_python",
    description="number of requests",
    unit="ms",
    value_type=int,
)

# define labels
metric_labels = {}

app = Flask(__name__)

@app.route('/')
def homePage():
    # start timer
    start_time = time.perf_counter()
    # count request
    requests_counter.add(1, metric_labels)
    # fail 10% of the time
    if random.randint(0, 100) > 90:
        errors_counter.add(1, metric_labels)
        latency = time.perf_counter() - start_time
        response_latency.record(latency)
        return("error!", 500)
    else:
        random_delay = random.randint(0,5000)/1000
        # delay for a bit to vary latency measurement
        time.sleep(random_delay)
        # record latency
        latency = time.perf_counter() - start_time
        response_latency.record(latency, metric_labels)
        return ("Responding in " + str(latency) + "ms")

if __name__ == '__main__':
    app.run()
