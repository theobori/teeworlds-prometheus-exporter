package exporter

import (
	"github.com/prometheus/client_golang/prometheus"

	"github.com/theobori/teeworlds-prometheus-exporter/internal/debug"
	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/econ"
	masterservers "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server"
)

// Prometheus exporter collector
type Exporter struct {
	// Teeworlds master servers manager
	msm *masterservers.MasterServerManager
	// Teeworlds econ servers manager
	em *econ.EconManager
}

// Create a new exporter struct
func NewExporter(msm *masterservers.MasterServerManager, em *econ.EconManager) *Exporter {
	return &Exporter{
		msm: msm,
		em:  em,
	}
}

// Collect the Teeworlds servers metrics
func (e *Exporter) collectServers(ch chan<- prometheus.Metric) {
	for metricInfo, f := range ServerMetrics {
		err := SendServerMetrics(metricInfo, e.msm, ch, f)
		if err != nil {
			debug.Debug(err.Error())
		}
	}
}

// Collect the Teeworlds master servers metrics
func (e *Exporter) collectMasterServers(ch chan<- prometheus.Metric) {
	masterServers := e.msm.MasterServers()

	for metricInfo, f := range MasterServerMetrics {
		err := SendMasterServerMetrics(metricInfo, masterServers, ch, f)
		if err != nil {
			debug.Debug(err.Error())
		}
	}
}

// Collect the Teeworlds econ servers metrics
func (e *Exporter) collectEconServers(ch chan<- prometheus.Metric) {
	econServersMetrics := e.em.EconServersMetrics()

	for metadata, metrics := range econServersMetrics {
		err := SendEconServerMetrics(metadata, metrics, ch)
		if err != nil {
			debug.Debug(err.Error())
		}
	}
}

// Send Prometheus metric description that represents the metrics attributes
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	// Teeworlds server metrics
	for metricInfo := range ServerMetrics {
		ch <- metricInfo.Desc
	}

	// Teeworlds master server metrics
	for metricInfo := range MasterServerMetrics {
		ch <- metricInfo.Desc
	}

	// Teeworlds econ server metric
	ch <- EconMetric.Desc
}

// Collect implements required collect function for all promehteus exporters
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Teeworlds servers
	e.collectServers(ch)

	// Teeworlds master servers
	e.collectMasterServers(ch)

	// Teeworlds econ servers
	e.collectEconServers(ch)
}
