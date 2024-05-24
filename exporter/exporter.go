package exporter

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/theobori/teeworlds-prometheus-exporter/internal/debug"
	masterservers "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server"
	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/server"
)

type MetricInfo struct {
	Desc *prometheus.Desc
	Type prometheus.ValueType
}

var (
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

	ServerMetrics = map[*MetricInfo]func(server *server.Server) float64{
		{
			Desc: prometheus.NewDesc("teeworlds_server_player_total", "Total number of player in a Teeworlds server", ServerLabels, nil),
			Type: prometheus.GaugeValue,
		}: func(server *server.Server) float64 {
			return float64(len(server.Info.Clients))
		},
	}
)

func ServerMetric(
	metricInfo *MetricInfo,
	ms *masterservers.MasterServers,
	ch chan<- prometheus.Metric,
	f func(server *server.Server) float64,
) error {
	if ms == nil || metricInfo == nil {
		return fmt.Errorf("missing master servers and metric info")
	}

	servers, err := ms.Servers()
	if err != nil {
		return err
	}

	var passworded string

	for metadata, value := range servers {
		for _, server := range value {
			if len(server.Addresses) == 0 {
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

type Exporter struct {
	ms *masterservers.MasterServers
}

func NewExporter(ms *masterservers.MasterServers) *Exporter {
	return &Exporter{
		ms: ms,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	// Teeworlds server
	for metricInfo := range ServerMetrics {
		ch <- metricInfo.Desc
	}
}

// Collect implements required collect function for all promehteus exporters
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	// Teeworlds server
	for metricInfo, f := range ServerMetrics {
		err := ServerMetric(metricInfo, e.ms, ch, f)
		if err != nil {
			debug.Debug("%v", err)
		}
	}
}
