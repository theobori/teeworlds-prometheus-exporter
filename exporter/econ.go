package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	econ "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/econ"
)

var (
	// Econ Prometheus labels
	EconLabels = []string{
		"address",
		"port",
		"event",
	}

	// Econ Prometheus metric
	EconMetric = MetricInfo{
		Desc: prometheus.NewDesc("teeworlds_econ_event_total", "Total number of received econ events.", EconLabels, nil),
		Type: prometheus.CounterValue,
	}
)

// Send Teeworlds econ servers Prometheus metric
func SendEconServerMetrics(
	metadata econ.EconMananagerKey,
	econMetrics econ.EconMetrics,
	ch chan<- prometheus.Metric,
) error {
	for metricName, metricValue := range econMetrics {
		labelValues := []string{
			metadata.Host,
			fmt.Sprintf("%d", metadata.Port),
			metricName,
		}

		ch <- prometheus.MustNewConstMetric(
			EconMetric.Desc,
			EconMetric.Type,
			float64(metricValue),
			labelValues...,
		)
	}

	return nil
}
