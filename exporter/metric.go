package exporter

import "github.com/prometheus/client_golang/prometheus"

// Metric informations
type MetricInfo struct {
	// Prometheus metric description
	Desc *prometheus.Desc
	// Prometheus metric type
	Type prometheus.ValueType
}
