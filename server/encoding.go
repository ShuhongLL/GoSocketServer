package server

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// Encode encode to binary
func Encode(msg *Message) ([]byte, error) {
	buffer := new(bytes.Buffer)

	err := binary.Write(buffer, binary.LittleEndian, msg.msgSize)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, msg.msgID)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, msg.data)
	if err != nil {
		return nil, err
	}

	err = binary.Write(buffer, binary.LittleEndian, msg.checkSum)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// Decode decode binary to Message
func Decode(data []byte) (*Message, error) {
	bufReader := bytes.NewReader(data)

	dataSize := len(data)

	// read message ID
	var msgID int32
	err := binary.Read(bufReader, binary.LittleEndian, &msgID)
	if err != nil {
		return nil, err
	}

	// read message data
	dataBufSize := dataSize - 8
	dataBuf := make([]byte, dataBufSize)
	err = binary.Read(bufReader, binary.LittleEndian, &dataBuf)
	if err != nil {
		return nil, err
	}

	// verify checkSum
	var checkSum uint32
	err = binary.Read(bufReader, binary.LittleEndian, &checkSum)
	if err != nil {
		return nil, err
	}

	message := &Message{}
	message.msgID = msgID
	message.msgSize = int32(dataSize)
	message.data = dataBuf
	message.checkSum = checkSum

	if message.Verify() {
		return message, nil
	}

	return nil, errors.New("Error in checkSum")
}
