package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/theobori/teeworlds-prometheus-exporter/exporter"
	"github.com/theobori/teeworlds-prometheus-exporter/internal/config"

	masterservers "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server"
	mhttp "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/http"
	masterserver "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/master_server"
	mudp "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/udp"
)

type getConfigMasterServerFunc func(m *config.MasterServer) (masterserver.MasterServer, error)

var (
	ErrMasterServerConfig = fmt.Errorf("missing master server configuration")

	MasterServerConfigProtocol = map[string]getConfigMasterServerFunc{
		"http": getConfigMasterServerHTTP,
		"udp": getConfigMasterServerUDP,
	}
)

func getConfigMasterServerHTTP(m *config.MasterServer) (masterserver.MasterServer, error) {
	if m == nil {
		return nil, ErrMasterServerConfig
	}

	masterServer := mhttp.NewMasterServer(m.URL)

	return masterServer, nil
}

func getConfigMasterServerUDP(m *config.MasterServer) (masterserver.MasterServer, error) {
	if m == nil {
		return nil, ErrMasterServerConfig
	}

	masterServer := mudp.NewMasterServer(m.Host, m.Port)

	err := masterServer.Connect()
	if err != nil {
		return nil, err
	}

	return masterServer, nil
}

func processConfigMasterServer(ms *masterservers.MasterServers, cm []config.MasterServer) error {
	if ms == nil {
		return ErrMasterServerConfig
	}

	for _, m := range cm {
		f, found := MasterServerConfigProtocol[m.Protocol]
		if !found {
			return fmt.Errorf("invalid master server protocol")
		}

		masterServer, err := f(&m)
		if err != nil {
			return err
		}

		entry := masterservers.NewMasterServersEntry(
			masterServer,
			m.RefreshCooldown,
		)

		err = ms.Register(*entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	configPath := flag.String("config-path", "./config.yaml", "Configuration YAML file path")
	port := flag.Uint("port", 8080, "Prometheus exporter port")
	endpoint := flag.String("endpoint", "/metrics", "Prometheus exporter HTTP endpoint")

	flag.Parse()

	// Get configuration as Golang struct
	c, err := config.ConfigFromFile(*configPath)
	if err != nil {
		log.Println(err)
		return
	}

	ms := masterservers.NewMasterServers()

	// Process and parse the configuration
	err = processConfigMasterServer(ms, c.Servers.Master)
	if err != nil {
		log.Println(err)
		return
	}

	// Start refreshing the servers
	ms.StartRefresh()

	// Register the exporter
	exporter := exporter.NewExporter(ms)
	prometheus.MustRegister(exporter)

	http.Handle(*endpoint, promhttp.Handler())

	log.Printf("Listening at endpoint %s on port %d", *endpoint, *port)
	
	pattern := fmt.Sprintf(":%d", *port)
	log.Fatal(http.ListenAndServe(pattern, nil))
}
