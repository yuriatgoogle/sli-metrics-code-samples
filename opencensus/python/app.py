from flask import Flask
import random
import time

from opencensus.stats import aggregation as aggregation_module
from opencensus.stats import measure as measure_module
from opencensus.stats import stats as stats_module
from opencensus.stats import view as view_module
from opencensus.tags import tag_key as tag_key_module
from opencensus.tags import tag_map as tag_map_module
from opencensus.tags import tag_value as tag_value_module

from opencensus.ext.prometheus import stats_exporter as prometheus

# set up measures
m_request_count = measure_module.MeasureInt("python_request_count", "total requests", "requests")
m_failed_request_count = measure_module.MeasureInt("python_failed_request_count", "failed requests", "requests")
m_response_latency = measure_module.MeasureFloat("python_response_latency", "response latency", "s")

# set up stats recorder
stats_recorder = stats_module.stats.stats_recorder

# set up views
latency_view = view_module.View("python_response_latency", "The distribution of the latencies",
    [],
    m_response_latency,
    aggregation_module.DistributionAggregation([0, 1000, 2000, 3000, 4000, 5000, 6000, 7000, 8000, 9000, 10000]))

request_count_view = view_module.View("python_request_count", "total requests",
    [],
    m_request_count,
    aggregation_module.CountAggregation())

failed_request_count_vew = view_module.View("python_failed_request_count", "failed requests",
    [],m_failed_request_count,
    aggregation_module.CountAggregation())

# register views
def registerAllViews(view_manager):
    view_manager.register_view(latency_view)
    view_manager.register_view(request_count_view)
    view_manager.register_view(failed_request_count_vew)

# set up exporter
def setupOpenCensusAndPrometheusExporter():
    stats = stats_module.stats
    view_manager = stats.view_manager
    exporter = prometheus.new_stats_exporter(prometheus.Options(namespace="oc_python", port=8080))
    view_manager.register_exporter(exporter)
    registerAllViews(view_manager)


app = Flask(__name__)

@app.route('/')
def homePage():
    # start timer
    start_time = time.perf_counter()
    mmap = stats_recorder.new_measurement_map()
    # count request
    mmap.measure_int_put(m_request_count, 1)
    # fail 10% of the time
    if random.randint(0, 100) > 90:
        mmap.measure_int_put(m_failed_request_count, 1)
        response_latency = time.perf_counter() - start_time
        mmap.measure_float_put(m_response_latency, response_latency)
        tmap = tag_map_module.TagMap()
        mmap.record(tmap)
        return("error!", 500)
    else:
        random_delay = random.randint(0,5000)/1000
        # delay for a bit to vary latency measurement
        time.sleep(random_delay)
        # record latency
        response_latency = time.perf_counter() - start_time
        mmap.measure_float_put(m_response_latency, response_latency)
        tmap = tag_map_module.TagMap()
        mmap.record(tmap)
        return ("home page")

if __name__ == '__main__':
    setupOpenCensusAndPrometheusExporter()
    app.run()