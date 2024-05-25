package udp

import (
	"fmt"
	"sync"
	"time"

	"github.com/jxsl13/twapi/browser"
	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/master_server"
	twserver "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/server"
)

// UDP master server controller
type MasterServerUDP struct {
	// Master server host
	host string
	// Master server port
	port uint16
	// Client that manage the UDP connection
	client *browser.Client
	// Teeworlds servers informations
	servers []*browser.ServerInfo
	// Master server metrics
	metrics masterserver.MasterServerMetrics
	// Mutex to protect `servers` and `metrics`
	mu sync.Mutex
}

// Create a new MasterServerUDP struct
func NewMasterServer(host string, port uint16) *MasterServerUDP {
	return &MasterServerUDP{
		host:   host,
		port:   port,
		client: nil,
	}
}

// Create a new default MasterServerUDP struct
func NewDefaultMasterServer() *MasterServerUDP {
	return NewMasterServer("master1.teeworlds.com", 8283)
}

// Get the master server address
func (ms *MasterServerUDP) Address() string {
	return fmt.Sprintf("%s:%d", ms.host, ms.port)
}

// Get the master server metadata
func (ms *MasterServerUDP) Metadata() masterserver.MasterServerMetadata {
	return masterserver.MasterServerMetadata{
		Protocol: "udp",
		Address:  ms.Address(),
	}
}

// Connect to the master server
func (ms *MasterServerUDP) Connect() error {
	addr := ms.Address()

	client, err := browser.NewClient(addr)
	if err != nil {
		return err
	}

	ms.client = client

	return nil
}

// Close the connection with the master server
func (ms *MasterServerUDP) Disconnect() error {
	if ms.client == nil {
		return fmt.Errorf("missing client")
	}

	return ms.client.Close()
}

// Update the teeworlds servers informations
func (ms *MasterServerUDP) refresh() error {
	if ms.client == nil {
		return fmt.Errorf("missing client")
	}

	// Starting before we get the server addresses
	start := time.Now()

	// Get the registered teeworlds server addresses
	addresses, err := ms.client.GetServerAddresses()
	if err != nil {
		return err
	}

	// Get the teeworlds servers informations
	serversInfo, err := browser.GetServerInfosOf(addresses)
	if err != nil {
		return err
	}

	// Get the elapsed time
	elapsed := time.Since(start).Seconds()

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.servers = serversInfo
	ms.metrics.RequestTime = uint(elapsed)

	return nil
}

// Update the teeworlds servers informations
func (ms *MasterServerUDP) Refresh() error {
	err := ms.refresh()

	ms.mu.Lock()
	defer ms.mu.Unlock()

	if err != nil {
		ms.metrics.FailedRefreshCount++

		return err
	}

	ms.metrics.SuccessRefreshCount++

	return nil
}

// Get the Teeworlds servers informations
func (ms *MasterServerUDP) Servers() ([]*twserver.Server, error) {
	var servers []*twserver.Server

	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, serverInfo := range ms.servers {
		server, err := twserver.FromUDPFields(serverInfo)
		if err != nil {
			return nil, err
		}

		servers = append(servers, server)
	}

	return servers, nil
}

// Get the Teeworlds servers informations with its original format
func (ms *MasterServerUDP) ServersInfo() []*browser.ServerInfo {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.servers
}

// Get the master server metrics
func (ms *MasterServerUDP) Metrics() masterserver.MasterServerMetrics {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.metrics
}
