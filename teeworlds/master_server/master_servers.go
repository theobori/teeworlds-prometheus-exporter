package masterserver

import (
	"fmt"
	"time"

	"github.com/theobori/teeworlds-prometheus-exporter/internal/debug"
	masterserver "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/master_server"
	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/server"
)

type MasterServersEntry struct {
	MasterServer    masterserver.MasterServer
	RefreshCooldown uint
	IsRefreshing    bool
}

func NewMasterServersEntry(masterServer masterserver.MasterServer, refreshcooldown uint) *MasterServersEntry {
	return &MasterServersEntry{
		MasterServer:    masterServer,
		RefreshCooldown: refreshcooldown,
		IsRefreshing:    false,
	}
}

type ServersMap map[masterserver.MasterServerMetadata][]*server.Server
type MasterServersMap map[masterserver.MasterServerMetadata]MasterServersEntry

type MasterServers struct {
	masterServers MasterServersMap
}

func NewMasterServers() *MasterServers {
	return &MasterServers{
		masterServers: make(MasterServersMap),
	}
}

func (ms *MasterServers) Register(entry MasterServersEntry) error {
	masterServer := entry.MasterServer

	if masterServer == nil {
		return fmt.Errorf("masterServer is nil")
	}

	ms.masterServers[masterServer.Metadata()] = entry

	return nil
}

func (ms *MasterServers) Delete(masterServerMetadata masterserver.MasterServerMetadata) {
	delete(ms.masterServers, masterServerMetadata)
}

func (ms *MasterServers) Servers() (ServersMap, error) {
	servers := make(ServersMap, len(ms.masterServers))

	for metadata, entry := range ms.masterServers {
		masterServer := entry.MasterServer

		if masterServer == nil {
			return nil, fmt.Errorf("masterserver is nil")
		}

		s, err := masterServer.Servers()
		if err != nil {
			return nil, err
		}

		servers[metadata] = s
	}

	return servers, nil
}

func startRefresh(entry *MasterServersEntry) error {
	if entry == nil {
		return fmt.Errorf("missing entry")
	}

	masterServer := (*entry).MasterServer

	if masterServer == nil {
		return fmt.Errorf("missing masterServer")
	}

	metadata := masterServer.Metadata()
	duration := time.Duration((*entry).RefreshCooldown) * time.Second

	(*entry).IsRefreshing = true

	for {
		err := masterServer.Refresh()
		if err != nil {
			debug.Debug(
				"could not refresh %s with protocol %s",
				metadata.Address,
				metadata.Protocol,
			)
		}

		time.Sleep(duration)
	}
}

func (ms *MasterServers) StartRefresh() {
	for _, entry := range ms.masterServers {
		if entry.IsRefreshing {
			continue
		}

		go startRefresh(&entry)
	}
}

// Failed during refresh
// Success during refresh
