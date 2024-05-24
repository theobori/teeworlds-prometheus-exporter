package udp

import (
	"fmt"
	"sync"

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
	// Mutex to protect `servers`
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

// Set the servers
func (ms *MasterServerUDP) setServers(servers []*browser.ServerInfo) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.servers = servers
}

// Update the teeworlds servers informations
func (ms *MasterServerUDP) Refresh() error {
	if ms.client == nil {
		return fmt.Errorf("missing client")
	}

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

	// Update the servers
	ms.setServers(serversInfo)

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
