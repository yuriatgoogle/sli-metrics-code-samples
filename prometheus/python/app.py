from flask import Flask
from prometheus_client import Summary, Gauge, Counter, generate_latest, REGISTRY, Histogram
import random
import time


app = Flask(__name__)
PYTHON_REQUESTS_COUNTER=Counter('python_requests', 'total requests')
PYTHON_FAILED_REQUESTS_COUNTER=Counter('python_failed_requests', 'failed requests')
PYTHON_LATENCIES_HISTOGRAM=Histogram('python_request_latency', 'request latency by path')


@app.route('/')
@PYTHON_LATENCIES_HISTOGRAM.time()
def homePage():
    # count request
    PYTHON_REQUESTS_COUNTER.inc()
    # fail 10% of the time
    if random.randint(0, 100) > 90:
        PYTHON_FAILED_REQUESTS_COUNTER.inc()
        return("error!", 500)
    else:
        random_delay = random.randint(0,5000)/1000
        # delay for a bit to vary latency measurement
        time.sleep(random_delay)
        return ("home page")

@app.route('/metrics', methods=['GET'])
def stats():
    return generate_latest(REGISTRY), 200

if __name__ == '__main__':
    app.run(debug=True,host='0.0.0.0', port=8080)