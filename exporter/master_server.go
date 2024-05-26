package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	masterserver "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/master_server"
)

var (
	// Master server Prometheus labels
	MasterServerLabels = []string{
		"address",
		"protocol",
	}

	// Teeworlds master server metrics informations associated with function to scrape a metric
	MasterServerMetrics = map[*MetricInfo]func(masterServer *masterserver.MasterServer) float64{
		{
			Desc: prometheus.NewDesc("teeworlds_master_server_players", "Total number of players on a master server.", MasterServerLabels, nil),
			Type: prometheus.GaugeValue,
		}: func(masterServer *masterserver.MasterServer) float64 {
			s := 0

			// Assuming masterServer cannot be nil
			servers, _ := (*masterServer).Servers()

			for _, server := range servers {
				if server == nil {
					continue
				}

				s += len(server.Info.Clients)
			}

			return float64(s)
		},
		{
			Desc: prometheus.NewDesc("teeworlds_master_server_servers", "Total number of servers registered on a master server.", MasterServerLabels, nil),
			Type: prometheus.GaugeValue,
		}: func(masterServer *masterserver.MasterServer) float64 {
			// Assuming masterServer cannot be nil
			servers, _ := (*masterServer).Servers()

			return float64(len(servers))
		},
		{
			Desc: prometheus.NewDesc("teeworlds_master_server_request_duration_seconds", "Request duration when refreshing a master server. From client request to full data server response.", MasterServerLabels, nil),
			Type: prometheus.GaugeValue,
		}: func(masterServer *masterserver.MasterServer) float64 {
			// Assuming masterServer cannot be nil
			metrics := (*masterServer).Metrics()

			return float64(metrics.RequestTime)
		},
		{
			Desc: prometheus.NewDesc("teeworlds_master_server_request_total", "Total number of master server requests.", MasterServerLabels, prometheus.Labels{"state": "failed"}),
			Type: prometheus.CounterValue,
		}: func(masterServer *masterserver.MasterServer) float64 {
			// Assuming masterServer cannot be nil
			metrics := (*masterServer).Metrics()

			return float64(metrics.FailedRefreshCount)
		},
		{
			Desc: prometheus.NewDesc("teeworlds_master_server_request_total", "Total number of master server requests.", MasterServerLabels, prometheus.Labels{"state": "success"}),
			Type: prometheus.CounterValue,
		}: func(masterServer *masterserver.MasterServer) float64 {
			// Assuming masterServer cannot be nil
			metrics := (*masterServer).Metrics()

			return float64(metrics.SuccessRefreshCount)
		},
	}
)

// Send Teeworlds master servers Prometheus metric
func SendMasterServerMetrics(
	metricInfo *MetricInfo,
	masterServers []*masterserver.MasterServer,
	ch chan<- prometheus.Metric,
	f func(*masterserver.MasterServer) float64,
) error {
	if metricInfo == nil {
		return fmt.Errorf("missing metric info")
	}

	for _, masterServer := range masterServers {
		if masterServer == nil {
			continue
		}

		metadata := (*masterServer).Metadata()

		labelValues := []string{
			metadata.Address,
			metadata.Protocol,
		}

		metricValue := f(masterServer)

		ch <- prometheus.MustNewConstMetric(
			metricInfo.Desc,
			metricInfo.Type,
			metricValue,
			labelValues...,
		)
	}

	return nil
}
