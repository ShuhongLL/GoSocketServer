package server

import (
	"fmt"
	"net"
	"testing"
	"time"
)

// TestSocket mock
func TestSocket(t *testing.T) {
	host := "127.0.0.1:9543"

	ss, err := CreateSocket(host)
	if err != nil {
		return
	}

	//set Heartbeat
	ss.SetHeartbeat(5*time.Second, 30*time.Second)

	ss.CreateMessageHandler(HandleMessage)
	ss.CreateConnectHandler(HandleConnect)
	ss.CreateDisconnectHandler(HandleDisconnect)

	go Client()

	timer := time.NewTimer(time.Second * 1)
	go func() {
		<-timer.C
		ss.Stop("stopping socket...")
	}()

	ss.Serv()
}

// HandleMessage mock
func HandleMessage(s *Session, msg *Message) {
	fmt.Println("receive msgID:", msg)
	fmt.Println("receive data:", string(msg.GetData()))
}

// HandleConnect mock
func HandleConnect(s *Session) {
	fmt.Println(s.GetConn().GetName() + " connected.")
}

// HandleDisconnect mock
func HandleDisconnect(s *Session, err error) {
	fmt.Println(s.GetConn().GetName() + " disconnected.")
}

// Client mock
func Client() {
	host := "127.0.0.1:9543"
	tcpAddr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return
	}

	msg := CreateMessage(1, []byte("Hi Go"))
	data, err := Encode(msg)
	if err != nil {
		return
	}
	conn.Write(data)
}
