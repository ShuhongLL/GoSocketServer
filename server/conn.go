package server

import (
	"bytes"
	"context"
	"encoding/binary"
	"io"
	"net"
	"time"
)

// Conn net.Conn
type Conn struct {
	sid        string
	rawConn    net.Conn
	sendCh     chan []byte
	done       chan error
	hbTimer    *time.Timer
	name       string
	messageCh  chan *Message
	hbInterval time.Duration
	hbTimeout  time.Duration
}

// CreateConn create new conn
func CreateConn(c net.Conn, hbInterval time.Duration, hbTimeout time.Duration) *Conn {
	conn := &Conn{
		rawConn:    c,
		sendCh:     make(chan []byte, 100),
		done:       make(chan error),
		messageCh:  make(chan *Message, 100),
		hbInterval: hbInterval,
		hbTimeout:  hbTimeout,
	}

	conn.name = c.RemoteAddr().String()
	conn.hbTimer = time.NewTimer(conn.hbInterval)

	if conn.hbInterval == 0 {
		conn.hbTimer.Stop()
	}

	return conn
}

//GetName get conn name
func (c *Conn) GetName() string {
	return c.name
}

// Close close connection
func (c *Conn) Close() {
	c.hbTimer.Stop()
	c.rawConn.Close()
}

// SendMsg send message
func (c *Conn) SendMsg(msg *Message) error {
	pkg, err := Encode(msg)

	if err != nil {
		return err
	}

	c.sendCh <- pkg
	return nil
}

// WriteCoroutine write coroutine
func (c *Conn) WriteCoroutine(ctx context.Context) {
	hbData := make([]byte, 0)

	for {
		select {
		case <-ctx.Done():
			return
		case pkt := <-c.sendCh:
			if pkt == nil {
				continue
			}

			if _, err := c.rawConn.Write(pkt); err != nil {
				c.done <- err
			}
		case <-c.hbTimer.C:
			hbMessage := CreateMessage(MsgHeartbeat, hbData)
			c.SendMsg(hbMessage)
		}
	}
}

// ReadCoroutine read coroutine
func (c *Conn) ReadCoroutine(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			// set time out
			if c.hbInterval > 0 {
				err := c.rawConn.SetReadDeadline(time.Now().Add(c.hbTimeout))
				if err != nil {
					c.done <- err
					continue
				}
			}

			buf := make([]byte, 4)
			_, err := io.ReadFull(c.rawConn, buf)
			if err != nil {
				c.done <- err
				continue
			}

			bufReader := bytes.NewReader(buf)

			var dataSize int32
			rErr := binary.Read(bufReader, binary.LittleEndian, &dataSize)
			if rErr != nil {
				c.done <- rErr
				continue
			}

			// read data
			databuf := make([]byte, dataSize)

			_, derr := io.ReadFull(c.rawConn, databuf)
			if err != nil {
				c.done <- derr
				continue
			}

			// decode
			msg, err := Decode(databuf)
			if err != nil {
				c.done <- err
				continue
			}

			// set heartbeats timer
			if c.hbInterval > 0 {
				c.hbTimer.Reset(c.hbInterval)
			}

			if msg.GetID() == MsgHeartbeat {
				continue
			}

			c.messageCh <- msg
		}
	}
}
