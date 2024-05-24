package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/master_server"
	twserver "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/server"
)

var (
	// Default HTTP master server urls
	MasterServerHTTPUrls = []string{
		"https://master1.ddnet.tw/ddnet/15/servers.json",
		"https://master2.ddnet.tw/ddnet/15/servers.json",
	}

	// Master server kind
	MasterServerKind = "http"
)

// HTTP master server controller
type MasterServerHTTP struct {
	// Master server url
	url string
	// Represents the Teeworlds servers
	servers []*twserver.Server
	// HTTP client used to perform every requests
	httpClient *http.Client
	mu         sync.Mutex
}

// Creates a new MasterServerHTTP struct
func NewMasterServer(url string) *MasterServerHTTP {
	return &MasterServerHTTP{
		url:        url,
		servers:    []*twserver.Server{},
		httpClient: &httpDefaultClient,
	}
}

// Create a new default MasterServerHTTP struct
func NewDefaultMasterServer() *MasterServerHTTP {
	return NewMasterServer(MasterServerHTTPUrls[0])
}

// Get the master server HTTP(s) url
func (ms *MasterServerHTTP) Url() string {
	return ms.url
}

// Get the master server metadata
func (ms *MasterServerHTTP) Metadata() masterserver.MasterServerMetadata {
	return masterserver.MasterServerMetadata{
		Protocol: "http",
		Address:  ms.Url(),
	}
}

// Set a HTTP client
func (ms *MasterServerHTTP) SetHTTPClient(httpClient *http.Client) {
	ms.httpClient = httpClient
}

// Get Teeworlds servers
func (ms *MasterServerHTTP) Servers() ([]*twserver.Server, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.servers, nil
}

// Refresh the Teeworlds servers with a context stored within the struct,
// this method has to be called at least one time
// if you want to get data.
func (ms *MasterServerHTTP) RefreshWithContext(ctx context.Context) error {
	var serversField twserver.Servers
	var servers []*twserver.Server

	err := HTTPGetJson(ctx, ms.httpClient, ms.url, &serversField)
	if err != nil {
		return err
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	for _, server := range serversField.Servers {
		servers = append(servers, &server)
	}

	ms.servers = servers

	return nil
}

// Refresh the Teeworlds servers with a default context
func (ms *MasterServerHTTP) RefreshWithoutContext() error {
	return ms.RefreshWithContext(context.Background())
}

// Refresh the Teeworlds server
func (ms *MasterServerHTTP) Refresh() error {
	return ms.RefreshWithoutContext()
}

// Get server informations
func (ms *MasterServerHTTP) Server(host string, port uint16) (*twserver.Server, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	targetAddress := fmt.Sprintf("%s:%d", host, port)

	for _, server := range ms.servers {
		for _, address := range server.Addresses {
			if strings.HasSuffix(address, targetAddress) {
				return server, nil
			}
		}
	}

	return nil, fmt.Errorf("the server %s is not registered", targetAddress)
}
