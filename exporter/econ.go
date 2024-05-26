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
	}

	// Teeworlds master server metrics informations associated with function to scrape a metric
	EconMetrics = map[*MetricInfo]func(econMetrics econ.EconMetrics) float64{
		{
			Desc: prometheus.NewDesc("teeworlds_econ_event_total", "Total number of received econ events.", EconLabels, prometheus.Labels{"event": "message"}),
			Type: prometheus.CounterValue,
		}: func(econMetrics econ.EconMetrics) float64 {
			return float64(econMetrics.MessagesTotal)
		},
		{
			Desc: prometheus.NewDesc("teeworlds_econ_event_total", "Total number of received econ events.", EconLabels, prometheus.Labels{"event": "kill"}),
			Type: prometheus.CounterValue,
		}: func(econMetrics econ.EconMetrics) float64 {
			return float64(econMetrics.KillsTotal)
		},
		{
			Desc: prometheus.NewDesc("teeworlds_econ_event_total", "Total number of received econ events.", EconLabels, prometheus.Labels{"event": "captured_flag"}),
			Type: prometheus.CounterValue,
		}: func(econMetrics econ.EconMetrics) float64 {
			return float64(econMetrics.CapturedFlagsTotal)
		},
		{
			Desc: prometheus.NewDesc("teeworlds_econ_event_total", "Total number of received econ events.", EconLabels, prometheus.Labels{"event": "vote"}),
			Type: prometheus.CounterValue,
		}: func(econMetrics econ.EconMetrics) float64 {
			return float64(econMetrics.VotesTotal)
		},
	}
)

// Send Teeworlds econ servers Prometheus metric
func SendEconServerMetrics(
	metricInfo *MetricInfo,
	em *econ.EconManager,
	ch chan<- prometheus.Metric,
	f func(econMetrics econ.EconMetrics) float64,
) error {
	if metricInfo == nil || em == nil {
		return fmt.Errorf("missing metric info or econ manager")
	}

	serversMetrics := em.EconServersMetrics()

	for key, metrics := range serversMetrics {
		labelValues := []string{
			key.Host,
			fmt.Sprintf("%d", key.Port),
		}

		metricValue := f(metrics)

		ch <- prometheus.MustNewConstMetric(
			metricInfo.Desc,
			metricInfo.Type,
			metricValue,
			labelValues...,
		)
	}

	return nil
}
