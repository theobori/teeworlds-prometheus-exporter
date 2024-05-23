package udp

import (
	"testing"
)

func TestMasterServerRefresh(t *testing.T) {
	ms := NewDefaultMasterServer()

	if err := ms.Connect(); err != nil {
		t.Error(err)
	}

	if err := ms.Refresh(); err != nil {
		t.Error(err)
	}
}

func TestMasterServerServersInfo(t *testing.T) {
	ms := NewDefaultMasterServer()

	if err := ms.Connect(); err != nil {
		t.Error(err)
	}

	if err := ms.Refresh(); err != nil {
		t.Error(err)
	}

	serversInfo := ms.ServersInfo()

	if len(serversInfo) == 0 {
		t.Errorf("no servers")
	}
}
