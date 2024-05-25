package exporter

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/theobori/teeworlds-prometheus-exporter/internal/debug"
	masterservers "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server"
)

// Prometheus exporter collector
type Exporter struct {
	// Teeworlds master servers manager
	msm *masterservers.MasterServerManager
}

// Create a new exporter struct
func NewExporter(msm *masterservers.MasterServerManager) *Exporter {
	return &Exporter{
		msm: msm,
	}
}

// Collect the Teeworlds servers metrics
func (e *Exporter) collectServers(ch chan<- prometheus.Metric) {
	for metricInfo, f := range ServerMetrics {
		err := ServerMetric(metricInfo, e.msm, ch, f)
		if err != nil {
			debug.Debug(err.Error())
		}
	}
}

// Collect the Teeworlds master servers metrics
func (e *Exporter) collectMasterServers(ch chan<- prometheus.Metric) {
	// Getting master servers from the master servers manager
	masterServers := e.msm.MasterServers()

	for metricInfo, f := range MasterServerMetrics {
		err := MasterServerMetric(metricInfo, masterServers, ch, f)
		if err != nil {
			debug.Debug(err.Error())
		}
	}
}

// Send Prometheus metric description that represents the metrics attributes
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	// Teeworlds server
	for metricInfo := range ServerMetrics {
		ch <- metricInfo.Desc
	}

	// Teeworlds master server
	for metricInfo := range MasterServerMetrics {
		ch <- metricInfo.Desc
	}
}

// Collect implements required collect function for all promehteus exporters
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Teeworlds servers
	e.collectServers(ch)

	// Teeworlds master servers
	e.collectMasterServers(ch)
}
