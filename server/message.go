package server

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/adler32"
)

// Message struct
type Message struct {
	msgSize  int32
	msgID    int32
	data     []byte
	checkSum uint32
}

//CreateMessage create new message
func CreateMessage(msgID int32, data []byte) *Message {
	msg := &Message{
		msgSize: int32(len(data)) + 8,
		msgID:   msgID,
		data:    data,
	}

	msg.checkSum = msg.GetCheckSum()
	return msg
}

// GetCheckSum calculate the sum of checks
func (msg *Message) GetCheckSum() uint32 {
	if msg == nil {
		return 0
	}

	data := new(bytes.Buffer)

	iderr := binary.Write(data, binary.LittleEndian, msg.msgID)
	if iderr != nil {
		return 0
	}

	dataerr := binary.Write(data, binary.LittleEndian, msg.data)
	if dataerr != nil {
		return 0
	}

	checkSum := adler32.Checksum(data.Bytes())
	return checkSum
}

//GetData get data
func (msg *Message) GetData() []byte {
	return msg.data
}

// GetID get ID
func (msg *Message) GetID() int32 {
	return msg.msgID
}

// Verify verify checkSum
func (msg *Message) Verify() bool {
	return msg.checkSum == msg.GetCheckSum()
}

func (msg *Message) toString() string {
	return fmt.Sprintf("Size=%d ID=%d DataLen=%d Checksum=%d", msg.msgSize, msg.GetID(), len(msg.GetData()), msg.checkSum)
}
