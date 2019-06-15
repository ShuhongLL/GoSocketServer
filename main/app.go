package main

import (
	"handler"
	"log"
	"server"
	"time"
)

func main() {
	host := "127.0.0.1:9543"

	ss, err := server.CreateSocket(host)
	if err != nil {
		log.Println(err)
		return
	}

	ss.SetHeartbeat(5*time.Second, 30*time.Second)
	ss.CreateMessageHandler(handler.HandleMesssage)
	ss.CreateConnectHandler(handler.HandleConnect)
	ss.CreateDisconnectHandler(handler.HandleDisconnect)

	log.Println("Running on " + host)
	ss.Serv()
}
