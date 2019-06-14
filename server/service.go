package server

import (
	"context"
	"errors"
	"log"
	"net"
	"sync"
	"time"
)

// SocketService socket
type SocketService struct {
	OnMessage    func(*Session, *Message)
	OnConnect    func(*Session)
	OnDisconnect func(*Session, error)
	sessions     *sync.Map
	hbInterval   time.Duration
	hbTimeout    time.Duration
	laddr        string
	status       int
	listener     net.Listener
	stopCh       chan error
}

// CreateSocket Create new socket
func CreateSocket(laddr string) (*SocketService, error) {
	l, err := net.Listen("tcp", laddr)

	if err != nil {
		return nil, err
	}

	s := &SocketService{
		sessions:   &sync.Map{},
		stopCh:     make(chan error),
		hbInterval: 0 * time.Second,
		hbTimeout:  0 * time.Second,
		laddr:      laddr,
		status:     STInited,
		listener:   l,
	}

	return s, nil
}

// CreateMessageHandler Register Message Handler
// trigger socket onMessage event
func (s *SocketService) CreateMessageHandler(handler func(*Session, *Message)) {
	s.OnMessage = handler
}

// CreateConnectHandler Register Connection Handler
// trigger socket onConnect event
func (s *SocketService) CreateConnectHandler(handler func(*Session)) {
	s.OnConnect = handler
}

// CreateDisconnectHandler Register disconnect Handler
// trigger socket onDisconnect event
func (s *SocketService) CreateDisconnectHandler(handler func(*Session, error)) {
	s.OnDisconnect = handler
}

// Serv Start socker service
func (s *SocketService) Serv() {
	s.status = STRunning
	ctx, cancel := context.WithCancel(context.Background())

	defer func() {
		s.status = STStop
		cancel()
		s.listener.Close()
	}()

	go s.acceptHandler(ctx)

	for {
		select {
		case <-s.stopCh:
			return
		}
	}
}

// acceptHandler:
// start a new go routine listening to the connection
// only stop the service when receive stopCh signal
func (s *SocketService) acceptHandler(ctx context.Context) {
	for {
		c, err := s.listener.Accept()
		if err != nil {
			s.stopCh <- err
			return
		}

		go s.ConnectHandler(ctx, c)
	}
}

// ConnectHandler connect handlers
// after a connection is built, create a new go routine to create conn and session
// save the session into the SocketService.sessions
// after conn is created, it will start read and write
func (s *SocketService) ConnectHandler(ctx context.Context, c net.Conn) {
	conn := CreateConn(c, s.hbInterval, s.hbTimeout)
	session := CreateSession(conn)
	s.sessions.Store(session.GetSessionID(), session)

	connContext, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()
		conn.Close()
		s.sessions.Delete(session.GetSessionID)
	}()

	go conn.ReadCoroutine(connContext)
	go conn.WriteCoroutine(connContext)

	if s.OnConnect != nil {
		s.OnConnect(session)
	}

	for {
		select {
		case err := <-conn.done:
			if s.OnDisconnect != nil {
				s.OnDisconnect(session, err)
			}
			return
		case msg := <-conn.messageCh:
			if s.OnMessage != nil {
				s.OnMessage(session, msg)
			}
		}
	}
}

// GetStatus get socket service status
func (s *SocketService) GetStatus() int {
	return s.status
}

// Stop stop socket service
func (s *SocketService) Stop(msg string) {
	s.stopCh <- errors.New(msg)
}

// SetHeartbeat set heart beats
func (s *SocketService) SetHeartbeat(hbInterval time.Duration, hbTimeout time.Duration) error {
	if s.status == STRunning {
		return errors.New("Service is running, Cannot set heartbeats")
	}

	s.hbInterval = hbInterval
	s.hbTimeout = hbTimeout

	return nil
}

// GetConnectionCount get the number of the connections
func (s *SocketService) GetConnectionCount() int {
	var count int
	s.sessions.Range(func(k, v interface{}) bool {
		count++
		return true
	})
	return count
}

// Unicast unicast with session ID
func (s *SocketService) Unicast(sid string, msg *Message) {
	v, ok := s.sessions.Load(sid)
	if ok {
		session := v.(*Session)
		err := session.GetConn().SendMsg(msg)
		if err != nil {
			return
		}
	}
}

// Broadcast send message to all connections
func (s *SocketService) Broadcast(msg *Message) {
	s.sessions.Range(func(k, v interface{}) bool {
		s := v.(*Session)
		err := s.GetConn().SendMsg(msg)
		if err != nil {
			log.Println(err)
		}
		return true
	})
}
