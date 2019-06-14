package handler

// Playground playground
type Playground struct {
	players map[string]*Player
}

var playground *Playground

func init() {
	playground = CreatePlayground()
}

// CreatePlayground createPlayground
func CreatePlayground() *Playground {
	playground := &Playground{
		players: make(map[string]*Player),
	}
	return playground
}

// AddPlayer add player
func (p *Playground) AddPlayer(player *Player) {
	p.players[player.ID] = player
}

// RemovePlayer remove player
func (p *Playground) RemovePlayer(id string) {
	if _, ok := p.players[id]; ok {
		delete(p.players, id)
	}
}

// GetPlayer get player
func (p *Playground) GetPlayer(id string) *Player {
	if player, ok := p.players[id]; ok {
		return player
	}
	return nil
}

// GetAllPlayer get all players
func (p *Playground) GetAllPlayer() []*Player {
	list := make([]*Player, 0)
	for _, player := range p.players {
		list = append(list, player)
	}
	return list
}
