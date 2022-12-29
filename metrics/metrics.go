package metrics

import "github.com/prometheus/client_golang/prometheus"

type metricDetails struct {
	desc       *prometheus.Desc
	metricType prometheus.ValueType
}

type metricInstance struct {
	details metricDetails
	value   float64
}
