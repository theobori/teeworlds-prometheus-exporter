package masterserver

import (
	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/server"
)

// Master server metadata
type MasterServerMetadata struct {
	// Master server protocol
	Protocol string
	// Master server address
	Address string
}

// Master server metrics
type MasterServerMetrics struct {
	// Master server success refresh count
	SuccessRefreshCount uint
	// Master server failed refresh count
	FailedRefreshCount uint
	// Master server request time in seconds
	RequestTime uint
}

type MasterServer interface {
	Servers() ([]*server.Server, error)
	Refresh() error
	Metadata() MasterServerMetadata
	Metrics() MasterServerMetrics
}
