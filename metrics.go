package main

import "github.com/prometheus/client_golang/prometheus"

type metrics struct {
	ScanCnt *prometheus.CounterVec
	ScanDur *prometheus.SummaryVec
}

func NewMetrics() *metrics {
	metrics := &metrics{
		ScanDur: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Namespace: "http",
				Subsystem: "client",
				Name:      "requests_duration",
				Help:      "The virusscan latencies in seconds.",
			},
			[]string{"phase"},
		),
	}

	prometheus.MustRegister(metrics.ScanDur)
	return metrics
}
