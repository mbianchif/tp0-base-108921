package common

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

// Message kind
const (
	KIND_BATCH   = 0
	KIND_CONFIRM = 1
)

const DELIMITER = ","
const BET_DELIMITER = ";"
const BATCH_SIZE_SIZE = 4
const ID_SIZE = 1
const MESSAGE_KIND_SIZE = 1
const DNI_COUNT_SIZE = 4
const DNI_SIZE = 4

type Bet struct {
	Agency    string
	Name      string
	Surname   string
	Id        string
	Birthdate string
	Number    string
}

func (m Bet) Encode() []byte {
	fields := []string{
		m.Agency,
		m.Name,
		m.Surname,
		m.Id,
		m.Birthdate,
		m.Number,
	}

	return []byte(strings.Join(fields, DELIMITER))
}

type BetSockStream struct {
	conn net.Conn
	id   int
}

func BetSockConnect(address string, id string) (*BetSockStream, error) {
	parsedID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil, err
	}

	// Share id with the other end
	writer := bufio.NewWriter(conn)
	idBytes := []byte{byte(parsedID)}
	writer.Write(idBytes)

	if err := writer.Flush(); err != nil {
		conn.Close()
		return nil, err
	}

	return &BetSockStream{conn, parsedID}, nil
}

func (s BetSockStream) PeerAddr() net.Addr {
	return s.conn.RemoteAddr()
}

func (s *BetSockStream) SendBets(bets []Bet) error {
	writer := bufio.NewWriter(s.conn)

	// Write message kind
	writer.Write([]byte{KIND_BATCH})

	// Encode batch
	betsEncoded := make([][]byte, 0)
	for _, bet := range bets {
		betsEncoded = append(betsEncoded, bet.Encode())
	}
	batchBytes := bytes.Join(betsEncoded, []byte(BET_DELIMITER))

	// Write batch size and data
	batchSize := len(batchBytes)
	batchSizeBytes := make([]byte, BATCH_SIZE_SIZE)
	binary.BigEndian.PutUint32(batchSizeBytes, uint32(batchSize))
	writer.Write(batchSizeBytes)
	writer.Write(batchBytes)

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("couldn't send message: %v", err)
	}

	return nil
}

func (s *BetSockStream) Confirm() error {
	writer := bufio.NewWriter(s.conn)

	// Write message kind
	writer.Write([]byte{KIND_CONFIRM})

	if err := writer.Flush(); err != nil {
		return fmt.Errorf("couldn't send confirmation")
	}

	return nil
}

func (s *BetSockStream) RecvWinners() ([]int, error) {
	dniCountBytes := make([]byte, DNI_COUNT_SIZE)

	if n, err := io.ReadFull(s.conn, dniCountBytes); err != nil {
		return nil, fmt.Errorf("couldn't recv winner quantity, err: %v, read %v out of %v bytes", err, n, DNI_COUNT_SIZE)
	}

	count := binary.BigEndian.Uint32(dniCountBytes)
	dnisBytes := make([]byte, DNI_SIZE*count)
	if _, err := io.ReadFull(s.conn, dnisBytes); err != nil {
		return nil, fmt.Errorf("couldn't recv winners, err: %v", err)
	}

	dnis := make([]int, 0, count)
	for i := 0; i < len(dnisBytes); i += DNI_SIZE {
		dni := binary.BigEndian.Uint32(dnisBytes[i : i+DNI_SIZE])
		dnis = append(dnis, int(dni))
	}

	return dnis, nil
}

func (s *BetSockStream) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}
