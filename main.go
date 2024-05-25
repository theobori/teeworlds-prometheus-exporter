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

// Function prototype
type getConfigMasterServerFunc func(m *config.MasterServer) (masterserver.MasterServer, error)

var (
	// Master server configuration error
	ErrMasterServerConfig = fmt.Errorf("missing master server configuration")

	// Master server configuration protocol
	MasterServerConfigProtocol = map[string]getConfigMasterServerFunc{
		"http": getConfigMasterServerHTTP,
		"udp":  getConfigMasterServerUDP,
	}
)

// Get a Teeworlds HTTP master server controller
func getConfigMasterServerHTTP(m *config.MasterServer) (masterserver.MasterServer, error) {
	if m == nil {
		return nil, ErrMasterServerConfig
	}

	masterServer := mhttp.NewMasterServer(m.URL)

	return masterServer, nil
}

// Get a Teeworlds UDP master server controller
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

// Process the Teeworlds master server depending of its configuration (protocol) 
func processConfigMasterServer(msm *masterservers.MasterServerManager, cm []config.MasterServer) error {
	if msm == nil {
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

		err = msm.Register(*entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	configPath := flag.String("config-path", "./config.yml", "Configuration YAML file path")
	port := flag.Uint("port", 8080, "Prometheus exporter port")
	endpoint := flag.String("endpoint", "/metrics", "Prometheus exporter HTTP endpoint")

	flag.Parse()

	// Get configuration as Golang struct
	c, err := config.ConfigFromFile(*configPath)
	if err != nil {
		log.Println(err)
		return
	}

	msm := masterservers.NewMasterServers()

	// Process and parse the configuration
	err = processConfigMasterServer(msm, c.Servers.Master)
	if err != nil {
		log.Println(err)
		return
	}

	// Start refreshing the servers
	msm.StartRefresh()

	// Register the exporter
	exporter := exporter.NewExporter(msm)
	prometheus.MustRegister(exporter)

	http.Handle(*endpoint, promhttp.Handler())
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write(
				[]byte(`<html>
             <head><title>Teeworlds Exporter</title></head>
             <body>
             <h1>Teeworlds Exporter</h1>
             <p><a href='` + *endpoint + `'>Metrics</a></p>
             </body>
             </html>`),
			)
		},
	)

	log.Printf("Exposing metrics via HTTP at endpoint %s on port %d", *endpoint, *port)

	pattern := fmt.Sprintf(":%d", *port)
	err = http.ListenAndServe(pattern, nil)

	log.Fatal(err)
}
