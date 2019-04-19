package metrics

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/AliyunContainerService/kube2ram/version"
)

const (
	namespace = "kube2ram"

	// RAMSuccessCode is the code used for metrics when an RAM request succeeds.
	RAMSuccessCode = "Success"
	// RAMUnknownFailCode is the code used for metrics when an RAM request fails with an error not reported by Alibaba Cloud.
	RAMUnknownFailCode = "UnknownError"
)

var (
	// RAMRequestSec tracks timing of RAM requests.
	RAMRequestSec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "ram",
			Name:      "request_duration_seconds",
			Help:      "Time taken to complete RAM request in seconds.",

			Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
		},
		[]string{
			// The HTTP status code Alibaba Cloud returned
			"code",
			// The arn of the RAM role being requested
			"role_arn",
		},
	)

	// RAMCacheHitCount tracks total number of RAM cache hits. Cache misses can be
	// calculated by looking at the total number of RAM requests for the same role_arn.
	RAMCacheHitCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "ram",
			Name:      "cache_hits_total",
			Help:      "Total number of RAM cache hits.",
		},
		[]string{
			// The arn of the RAM role being requested
			"role_arn",
		},
	)

	// HTTPRequestSec tracks timing of served HTTP requests.
	HTTPRequestSec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "Time taken for kube2ram to serve HTTP request.",

			Buckets: prometheus.ExponentialBuckets(0.001, 2, 14),
		},
		[]string{
			// The HTTP status code kube2ram returned
			"code",
			// The HTTP method being served
			"method",
			// The name of the handler being served
			"handler",
		},
	)

	// HealthcheckStatus reports the current healthcheck status of kube2ram.
	HealthcheckStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "healthcheck",
			Name:      "status",
			Help:      "The healthcheck status at scrape time. A value of 1 means it is passing, 0 means it is failing.",
		},
	)

	// Info reports various static information about the running kube2ram binary.
	Info = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "info",
			Help:      "Informational labels about the kube2ram process.",
		},
		[]string{
			// The version of kube2ram running
			"version",
			// The build date of the kube2ram version
			"build_date",
			// The commit hash of the kube2ram version
			"commit_hash",
		},
	)
)

func init() {
	prometheus.MustRegister(RAMRequestSec)
	prometheus.MustRegister(RAMCacheHitCount)
	prometheus.MustRegister(HTTPRequestSec)
	prometheus.MustRegister(HealthcheckStatus)
	prometheus.MustRegister(Info)

	for _, val := range []string{RAMSuccessCode, RAMUnknownFailCode} {
		RAMRequestSec.WithLabelValues(val, "")
	}
	Info.WithLabelValues(version.Version, version.BuildDate, version.GitCommit).Set(1)
}

// StartMetricsServer registers a prometheus /metrics handler and starts a HTTP server
// listening on the provided port to service it.
func StartMetricsServer(metricsPort string) {
	r := mux.NewRouter()
	r.Handle("/metrics", GetHandler())

	go func() {
		if err := http.ListenAndServe(":"+metricsPort, r); err != nil {
			log.Fatalf("Error creating metrics http server: %+v", err)
		}
	}()
}

// GetHandler creates a prometheus HTTP handler that will serve metrics.
func GetHandler() http.Handler {
	return promhttp.Handler()
}

// NewFunctionTimer creates a new timer for a generic function that can be observed to time the duration of the handler.
// The metric is labeled with the values produced by the lvsProducer to allow for late binding of label values.
// If provided, the timer value is stored in storeValue to allow callers access to the reported value.
func NewFunctionTimer(histVec *prometheus.HistogramVec, lvsFn func() []string, storeValue *float64) *prometheus.Timer {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		if storeValue != nil {
			*storeValue = v
		}
		histVec.WithLabelValues(lvsFn()...).Observe(v)
	}))
	return timer
}
