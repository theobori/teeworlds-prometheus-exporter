package server

import (
	"fmt"

	"github.com/jxsl13/twapi/browser"
	twclient "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/client"
)

type Servers struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	Addresses []string   `json:"addresses"`
	Location  string     `json:"location"`
	Info      ServerInfo `json:"info"`
}

type ServerInfo struct {
	MaxClients      int               `json:"max_clients"`
	MaxPlayers      int               `json:"max_players"`
	Passworded      bool              `json:"passworded"`
	GameType        string            `json:"game_type"`
	Name            string            `json:"name"`
	Map             ServerMap         `json:"map"`
	Version         string            `json:"version"`
	ClientScoreKind string            `json:"client_score_kind"`
	Clients         []twclient.Client `json:"clients"`
}

type ServerMap struct {
	Name   string `json:"name"`
	SHA256 string `json:"sha256"`
	Size   int    `json:"size"`
}

// Get a `*Server` based on the teeworlds UDP master server fields
func FromUDPFields(other *browser.ServerInfo) (*Server, error) {
	var server Server
	var clients []twclient.Client

	if other == nil {
		return nil, fmt.Errorf("nil server")
	}

	// Converting clients
	for _, playerInfo := range other.Players {
		client, err := twclient.FromUDPFields(&playerInfo)
		if err != nil {
			return nil, err
		}

		clients = append(clients, *client)
	}

	server.Info.Clients = clients

	// Converting map
	m := ServerMap{
		Name:   other.Map,
		SHA256: "",
		Size:   0,
	}

	// Check for the first (least) server flag bit
	passworded := (other.ServerFlags & 1) == 1

	// Converting the server informations
	serverInfo := ServerInfo{
		MaxClients: other.MaxClients,
		MaxPlayers: other.MaxPlayers,
		Passworded: passworded,
		GameType:   other.GameType,
		Name:       other.Name,
		Map:        m,
		Version:    other.Version,
		Clients:    clients,
	}

	server.Info = serverInfo
	server.Addresses = []string{other.Address}
	server.Location = ""

	return &server, nil
}
