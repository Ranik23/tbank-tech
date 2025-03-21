package metrics

import "github.com/prometheus/client_golang/prometheus"


var (
	TotalRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",
			Help: "Total number of gRPC requests",
		},
	)
	ErrorRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "grpc_requests_errors_total",
			Help: "Total number of failed gRPC requests",
		},
	)
	RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",
			Help:    "Duration of gRPC requests",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func init() {
	prometheus.MustRegister(TotalRequests, ErrorRequests, RequestDuration)
}