package handler

import (
	"encoding/json"
	"math/rand"
	"server"
)

// Player player
type Player struct {
	ID      string          `json:"id"`
	XAxis   int             `json:"xaxis"`
	YAxis   int             `json:"yaxis"`
	Name    string          `json:"name"`
	Action  int             `json:"action"`
	Session *server.Session `json:"-"`
}

// CreatePlayer create nre player
func CreatePlayer(name string, s *server.Session) *Player {
	player := &Player{
		ID:      name,
		Name:    name,
		XAxis:   rand.Intn(10),
		YAxis:   rand.Intn(10),
		Action:  rand.Intn(5),
		Session: s,
	}
	return player
}

// ToJSON convert to JSON
func (p *Player) ToJSON() []byte {
	data, _ := json.Marshal(p)
	return data
}
