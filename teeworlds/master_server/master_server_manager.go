package masterserver

import (
	"fmt"
	"time"

	"github.com/theobori/teeworlds-prometheus-exporter/internal/debug"
	masterserver "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/master_server"
)

// Master server entry used to register a master server
type MasterServerManagerEntry struct {
	// Master server implementation
	MasterServer masterserver.MasterServer
	// Cooldown between each refresh
	RefreshCooldown uint
	// Refresh state
	IsRefreshing bool
}

// Create a new master server entry
func NewMasterServerManagerEntry(
	masterServer masterserver.MasterServer,
	refreshcooldown uint,
) *MasterServerManagerEntry {
	return &MasterServerManagerEntry{
		MasterServer:    masterServer,
		RefreshCooldown: refreshcooldown,
		IsRefreshing:    false,
	}
}

// Map used to manage master servers
type MasterServersMap map[masterserver.MasterServerMetadata]MasterServerManagerEntry

// Master server manager
type MasterServerManager struct {
	masterServers MasterServersMap
}

// Create a new master server manager
func NewMasterServerManager() *MasterServerManager {
	return &MasterServerManager{
		masterServers: make(MasterServersMap),
	}
}

// Register a master server
func (msm *MasterServerManager) Register(entry MasterServerManagerEntry) error {
	masterServer := entry.MasterServer

	if masterServer == nil {
		return fmt.Errorf("masterServer is nil")
	}

	msm.masterServers[masterServer.Metadata()] = entry

	return nil
}

// Delete a master server
func (msm *MasterServerManager) Delete(masterServerMetadata masterserver.MasterServerMetadata) {
	delete(msm.masterServers, masterServerMetadata)
}

// Return a Slice of pointers on master server
func (msm *MasterServerManager) MasterServers() []*masterserver.MasterServer {
	var masterServers []*masterserver.MasterServer

	for _, entry := range msm.masterServers {
		masterServers = append(masterServers, &entry.MasterServer)
	}

	return masterServers
}

// Goroutine that start refreshing a master server
func startRefresh(entry *MasterServerManagerEntry, errorCh chan error) {
	if entry == nil {
		errorCh <- fmt.Errorf("missing entry")
		return
	}

	masterServer := (*entry).MasterServer

	if masterServer == nil {
		errorCh <- fmt.Errorf("missing masterServer")
		return
	}

	errorCh <- nil

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

// Control loop that start refreshing every master server
func (msm *MasterServerManager) StartRefresh() {
	errorCh := make(chan error)

	for _, entry := range msm.masterServers {
		if entry.IsRefreshing {
			continue
		}

		go startRefresh(&entry, errorCh)

		err := <-errorCh
		if err != nil {
			debug.Debug(err.Error())
		}
	}
}
