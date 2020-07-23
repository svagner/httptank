package httptank

import "github.com/prometheus/client_golang/prometheus"

var (
	queries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "tank",
			Name:      "queries_count",
			Help:      "Count of queries",
		},
		[]string{"code"})

	queriesErrors = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "tank",
			Name:      "queries_errors",
			Help:      "Queries Errors",
		})

	queriesLatency = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "tank",
			Name:       "query_latency",
			Help:       "Query latency",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
		[]string{"stage"})
)
