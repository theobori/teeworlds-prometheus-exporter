package masterserver

import (
	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/server"
)

type MasterServer interface {
	Id() string
	Servers() ([]*server.Server, error)
	Refresh() error
	Kind() string
}
