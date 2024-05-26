package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	masterservers "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server"
	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/server"
)

var (
	// Teeworlds server Prometheus labels
	ServerLabels = []string{
		"name",
		"address",
		"gametype",
		"max_players",
		"password",
		"map",
		"version",
		"master_server_protocol",
		"master_server_address",
	}

	// Teeworlds master server metrics informations associated with function to scrape a metric
	ServerMetrics = map[*MetricInfo]func(server *server.Server) float64{
		{
			Desc: prometheus.NewDesc("teeworlds_server_players", "Total number of players in a Teeworlds server", ServerLabels, nil),
			Type: prometheus.GaugeValue,
		}: func(server *server.Server) float64 {
			// Assuming server cannot be nil
			return float64(len(server.Info.Clients))
		},
	}
)

// Send Teeworlds servers Prometheus metric
func SendServerMetrics(
	metricInfo *MetricInfo,
	msm *masterservers.MasterServerManager,
	ch chan<- prometheus.Metric,
	f func(server *server.Server) float64,
) error {
	if msm == nil || metricInfo == nil {
		return fmt.Errorf("missing master servers and metric info")
	}

	masterServers := msm.MasterServers()

	var passworded string

	for _, masterServer := range masterServers {
		if masterServer == nil {
			continue
		}

		metadata := (*masterServer).Metadata()

		servers, err := (*masterServer).Servers()
		if err != nil {
			continue
		}

		for _, server := range servers {
			if server == nil || len(server.Addresses) == 0 {
				continue
			}

			if server.Info.Passworded {
				passworded = "true"
			} else {
				passworded = "false"
			}

			labelValues := []string{
				server.Info.Name,
				server.Addresses[0],
				server.Info.GameType,
				fmt.Sprintf("%d", server.Info.MaxPlayers),
				passworded,
				server.Info.Map.Name,
				server.Info.Version,
				metadata.Protocol,
				metadata.Address,
			}

			metricValue := f(server)

			ch <- prometheus.MustNewConstMetric(
				metricInfo.Desc,
				metricInfo.Type,
				metricValue,
				labelValues...,
			)

		}
	}

	return nil
}
