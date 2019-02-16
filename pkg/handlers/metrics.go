package handlers

import "github.com/prometheus/client_golang/prometheus"

var Counter = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_requests_counter",
	Help: "total requests",
}, []string{"messageID"})

var Timings = prometheus.NewSummaryVec(prometheus.SummaryOpts{
	Name: "http_requests_duration",
	Help: "requests duration",
}, []string{"messageID"})

func RegisterMetrics() {
	prometheus.MustRegister(Counter, Timings)
}
