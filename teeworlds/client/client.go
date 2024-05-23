package client

import (
	"fmt"

	"github.com/jxsl13/twapi/browser"
)

type Client struct {
	Name     string     `json:"name"`
	Clan     string     `json:"clan"`
	Country  int        `json:"country"`
	Score    int        `json:"score"`
	IsPlayer bool       `json:"is_player"`
	Skin     ClientSkin `json:"skin"`
	Afk      bool       `json:"afk"`
	Team     int        `json:"team"`
}

type ClientSkin struct {
	Name string `json:"name"`
}

// Get a `*Client` based on the teeworlds UDP master server fields
func FromUDPFields(other *browser.PlayerInfo) (*Client, error) {
	if other == nil {
		return nil, fmt.Errorf("invalid argument")
	}

	client := Client{
		Name:     other.Name,
		Clan:     other.Clan,
		Country:  other.Country,
		Score:    other.Score,
		IsPlayer: other.Type == 0,
		Skin:     ClientSkin{Name: ""},
		Afk:      false,
		Team:     0,
	}

	return &client, nil
}
