package masterserver

import (
	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/server"
)

type MasterServerMetadata struct {
	Protocol string
	Address  string
}

type MasterServer interface {
	Servers() ([]*server.Server, error)
	Refresh() error
	Metadata() MasterServerMetadata
}
