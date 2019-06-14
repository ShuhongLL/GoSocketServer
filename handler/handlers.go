package handler

import (
	"encoding/json"
	"log"
	"server"
)

// HandleConnect handle onConnect event
func HandleConnect(s *server.Session) {
	log.Println(s.GetConn().GetName() + " connected")
}

// HandleMesssage handle onMessage event
func HandleMesssage(s *server.Session, msg *server.Message) {
	msgID := msg.GetID()

	switch msgID {
	case RequestJoin:
		name := s.GetConn().GetName()
		player := CreatePlayer(name, s)

		mp := make(map[string]interface{})
		mp["self"] = player
		mp["list"] = playground.GetAllPlayer()

		data, _ := json.Marshal(mp)
		response := server.CreateMessage(RequestJoin, data)
		s.GetConn().SendMsg(response)

		for _, p := range playground.GetAllPlayer() {
			response := server.CreateMessage(BroadcastJoin, player.ToJSON())
			p.Session.GetConn().SendMsg(response)
		}

		playground.AddPlayer(player)
		s.BindUserID(player.ID)

		break
	case RequestMove:
		var f interface{}
		err := json.Unmarshal(msg.GetData(), &f)
		if err != nil {
			return
		}
		mp := f.(map[string]interface{})
		x := mp["xaxis"]
		y := mp["yaxis"]
		playerID := mp["id"].(string)

		player := playground.GetPlayer(playerID)
		player.XAxis = int(x.(float64))
		player.YAxis = int(y.(float64))

		players := playground.GetAllPlayer()

		for _, p := range players {
			message := server.CreateMessage(BroadcastMove, player.ToJSON())
			p.Session.GetConn().SendMsg(message)
		}
		break
	}

}

// HandleDisconnect handle onDisconnect event
func HandleDisconnect(s *server.Session, err error) {
	log.Println(s.GetConn().GetName() + "disconnected.")
	id := s.GetUserID()
	quitPlayer := playground.GetPlayer(id)
	if quitPlayer == nil {
		return
	}

	playground.RemovePlayer(id)
	for _, p := range playground.GetAllPlayer() {
		message := server.CreateMessage(BroadcastLeave, quitPlayer.ToJSON())
		p.Session.GetConn().SendMsg(message)
	}
}
