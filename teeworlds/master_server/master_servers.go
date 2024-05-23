package masterserver

import (
	"fmt"

	masterserver "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/master_server"
)

type MasterServers struct {
	masterServers map[string]*masterserver.MasterServer
}

func (ms *MasterServers) Register(masterServer *masterserver.MasterServer) error {
	if masterServer == nil {
		return fmt.Errorf("masterServer is nil")
	}

	ms.masterServers[(*masterServer).Id()] = masterServer

	return nil
}

func (ms *MasterServers) Delete(s string) {
	delete(ms.masterServers, s)
}
