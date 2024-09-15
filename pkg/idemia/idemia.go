package idemia

import (
	"net"
)

type idemiaMessage struct {
	ID     uint8  // 1 byte
	Length uint16 // 2 bytes
	Data   []byte // variable length
}

const (
	IdemiaMessageUserControlSuccessful = 0x00

	IdemiaAccessGranted = 0x00
	IdemiaAccessDenied  = 0xFF
)

func ReadMessage(conn net.Conn) (*idemiaMessage, error) {
	// Create a buffer to hold the first 3 bytes (ID and Length).
	header := make([]byte, 3)

	// Read the first 3 bytes.
	_, err := conn.Read(header)
	if err != nil {
		return nil, err
	}

	// Parse the header to extract the ID and Length.
	msg := &idemiaMessage{
		ID:     header[0],                                // First byte is the ID.
		Length: uint16(header[1]) | uint16(header[2])<<8, // Next two bytes are the Length.
	}

	// Now read the remaining Data based on the Length.
	msg.Data = make([]byte, msg.Length)
	_, err = conn.Read(msg.Data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func SendMessage(conn net.Conn, id uint8, data []byte) error {
	// Create the message header.
	length := uint16(len(data))
	header := []byte{id, byte(length & 0xFF), byte(length >> 8)}

	// Write the header and data to the connection.
	_, err := conn.Write(append(header, data...))
	return err
}

func SendBasicApproval(conn net.Conn) error {
	return SendMessage(conn, 0x50, []byte{IdemiaAccessGranted})
}

func SendBasicDenial(conn net.Conn) error {
	return SendMessage(conn, 0x50, []byte{IdemiaAccessDenied})
}
