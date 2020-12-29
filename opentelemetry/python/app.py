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
# waiting for https://github.com/GoogleCloudPlatform/opentelemetry-operations-python/issues/67 
# to be resolved to use valuerecorder

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
        # no way to record a histogram at present - see above
        return("error!", 500)
    else:
        random_delay = random.randint(0,5000)/1000
        # delay for a bit to vary latency measurement
        time.sleep(random_delay)
        # record latency
        latency = time.perf_counter() - start_time
        # no way to record a histogram at present - see above
        return ("Responding in " + str(latency) + "ms")

if __name__ == '__main__':
    app.run()
