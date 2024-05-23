package http

import (
	"testing"
)

func TestMasterServerRefresh(t *testing.T) {
	ms := NewDefaultMasterServer()

	if err := ms.RefreshWithoutContext(); err != nil {
		t.Error(err)
	}
}

func TestMasterServerServer(t *testing.T) {
	ms := NewDefaultMasterServer()

	if err := ms.RefreshWithoutContext(); err != nil {
		t.Error(err)
	}

	_, err := ms.Server("", 0)
	if err == nil {
		t.Error(err)
	}
}

func TestMasterServerServers(t *testing.T) {
	ms := NewDefaultMasterServer()

	if err := ms.RefreshWithoutContext(); err != nil {
		t.Error(err)
	}

	servers, err := ms.Servers()
	if err != nil {
		t.Error(err)
	}

	if len(servers) == 0 {
		t.Errorf("no servers found")
	}
}
