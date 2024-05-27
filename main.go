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

	econ "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/econ"
	masterservers "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server"
)

func main() {
	configPath := flag.String("config-path", "./config.yml", "Teeworlds configuration YAML file path")
	port := flag.Uint("port", 8080, "Prometheus exporter port")
	endpoint := flag.String("endpoint", "/metrics", "Prometheus exporter HTTP endpoint")

	flag.Parse()

	// Get configuration as Golang struct
	c, err := config.ConfigFromFile(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	// Master server manager
	msm := masterservers.NewMasterServerManager()

	// Econ server manager
	em := econ.NewEconManager()

	// Process and parse the configuration
	if err := config.ProcessConfig(em, msm, *c); err != nil {
		log.Fatalln(err)
	}

	// Start refreshing the master servers
	msm.StartRefresh()

	// Register the events for metrics
	if err := em.RegisterEconEvents(); err != nil {
		log.Fatalln(err)
	}

	// Start handling events
	if err := em.StartHandle(); err != nil {
		log.Fatalln(err)
	}

	// Register the exporter
	exporter := exporter.NewExporter(msm, em)
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

	log.Fatalln(err)
}
