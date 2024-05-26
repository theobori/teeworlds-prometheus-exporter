package config

import (
	"fmt"

	twecon "github.com/theobori/teeworlds-econ"
	econ "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/econ"
	masterservers "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server"
	mhttp "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/http"
	masterserver "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/master_server"
	mudp "github.com/theobori/teeworlds-prometheus-exporter/teeworlds/master_server/udp"
)

// Function prototype
type getConfigMasterServerFunc func(m *MasterServer) (masterserver.MasterServer, error)

var (
	// Master server configuration error
	ErrMasterServerConfig = fmt.Errorf("missing master server configuration")

	// Master server configuration protocol
	MasterServerConfigProtocol = map[string]getConfigMasterServerFunc{
		"http": processMasterServerHTTP,
		"udp":  processMasterServerUDP,
	}
)

// Return a Teeworlds HTTP master server controller from the configuration
func processMasterServerHTTP(m *MasterServer) (masterserver.MasterServer, error) {
	if m == nil {
		return nil, ErrMasterServerConfig
	}

	masterServer := mhttp.NewMasterServer(m.URL)

	return masterServer, nil
}

// Return a Teeworlds UDP master server controller from the configuration
func processMasterServerUDP(m *MasterServer) (masterserver.MasterServer, error) {
	if m == nil {
		return nil, ErrMasterServerConfig
	}

	masterServer := mudp.NewMasterServer(m.Host, m.Port)

	err := masterServer.Connect()
	if err != nil {
		return nil, err
	}

	return masterServer, nil
}

func processMasterServer(
	msm *masterservers.MasterServerManager,
	masterServerConfig MasterServer,
) error {
	f, found := MasterServerConfigProtocol[masterServerConfig.Protocol]
	if !found {
		return fmt.Errorf("invalid master server protocol")
	}

	masterServer, err := f(&masterServerConfig)
	if err != nil {
		return err
	}

	entry := masterservers.NewMasterServerManagerEntry(
		masterServer,
		masterServerConfig.RefreshCooldown,
	)

	err = msm.Register(*entry)
	if err != nil {
		return err
	}

	return nil
}

// Process the Teeworlds master servers depending of its configuration (protocol),
// it registers them on the manager `msm`
func processMasterServers(
	msm *masterservers.MasterServerManager,
	masterServerConfigs []MasterServer,
) error {
	if msm == nil {
		return ErrMasterServerConfig
	}

	for _, masterServerConfig := range masterServerConfigs {
		err := processMasterServer(msm, masterServerConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func processEconServer(em *econ.EconManager, econConfig EconServer) error {
	c := twecon.EconConfig{
		Host:     econConfig.Host,
		Port:     econConfig.Port,
		Password: econConfig.Password,
	}

	e := twecon.NewEcon(&c)

	if err := e.Connect(); err != nil {
		return err
	}

	if r, err := e.Authenticate(); err != nil || !r.State {
		return fmt.Errorf("error: %v, response: %v", err, r)
	}

	if err := em.Register(e); err != nil {
		return err
	}

	return nil
}

func processEconServers(em *econ.EconManager, econConfigs []EconServer) error {
	for _, econConfig := range econConfigs {
		err := processEconServer(em, econConfig)
		if err != nil {
			return err
		}
	}

	return nil
}

func ProcessConfig(em *econ.EconManager, msm *masterservers.MasterServerManager, c Config) error {
	if err := processEconServers(em, c.Servers.Econ); err != nil {
		return err
	}

	if err := processMasterServers(msm, c.Servers.Master); err != nil {
		return err
	}

	return nil
}
