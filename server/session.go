package server

import uuid "github.com/satori/go.uuid"

// Session struct
type Session struct {
	sID      string
	uID      string
	conn     *Conn
	settings map[string]interface{}
}

// CreateSession create a new session
func CreateSession(conn *Conn) *Session {
	id := uuid.NewV4()
	session := &Session{
		sID:      id.String(),
		uID:      "",
		conn:     conn,
		settings: make(map[string]interface{}),
	}

	return session
}

// GetSessionID get the current session id
func (s *Session) GetSessionID() string {
	return s.sID
}

// BindUserID bind a userID to the current session ID
func (s *Session) BindUserID(uid string) {
	s.uID = uid
}

// GetUserID get userID
func (s *Session) GetUserID() string {
	return s.uID
}

// GetConn get connection from the current session
func (s *Session) GetConn() *Conn {
	return s.conn
}

// SetConn set the connectio
func (s *Session) SetConn(conn *Conn) {
	s.conn = conn
}

// GetSetting get k v from the settings
func (s *Session) GetSetting(key string) interface{} {
	if v, ok := s.settings[key]; ok {
		return v
	}
	return nil
}

// SetSetting set settings
func (s *Session) SetSetting(k string, v interface{}) {
	s.settings[k] = v
}
