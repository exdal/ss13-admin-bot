package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math"
	"net"
	"time"
)

const BYOND_PACKET_TYPE_GET_STEAM_STATUS = 0x36
const BYOND_PACKET_TYPE_HOST_MESSAGE = 0x3A
const BYOND_PACKET_TYPE_TOPIC = 0x83

const TOPIC_TYPE_NULL = 0x0
const TOPIC_TYPE_DECIMAL = 0x2A
const TOPIC_TYPE_STRING = 0x6

const BYOND_PACKET_HEADER_SIZE = 4

type ByondPacketHeader struct {
	PacketType uint16
	DataSize   uint16
}

var (
	byondConnection net.Conn
)

func byond_read_number(data []byte) int {
	return int(math.Float32frombits(binary.LittleEndian.Uint32(data)))
}

func byond_session_start(address string) bool {
	connection, err := net.DialTimeout("tcp", address, 15*time.Second)

	if err != nil {
		log.Fatalf("Failed to start BYOND session, error: %v\n", err)
		return false
	}

	byondConnection = connection

	return true
}

func byond_session_shutdown() {
	byondConnection.Close()
}

func byond_session_safe_connection() bool {
	for i := 0; i < 5; i++ {
		if byond_session_start(discordConfig.ByondServer) {
			return true
		}
	}

	panic("Failed to connect to server after 5 retries.")
}

func byond_session_read_safe(data []byte) bool {
	_, err := byondConnection.Read(data)

	if err != nil {
		if byond_session_safe_connection() {
			byondConnection.Read(data)
			return true
		}

		log.Fatalf("Failed to read session packet, error: %v", err)
		return false
	}

	return true
}

func byond_session_write_safe(data []byte) bool {
	_, err := byondConnection.Write(data)

	if err != nil {
		if byond_session_safe_connection() {
			byondConnection.Write(data)
			return true
		}

		log.Fatalf("Failed to write session packet, error: %v", err)
		return false
	}

	return true
}

func byond_write_packet_header(buffer *bytes.Buffer, packetType uint16, packetSize uint16) {
	binary.Write(buffer, binary.BigEndian, packetType)
	binary.Write(buffer, binary.BigEndian, packetSize)
}

func byond_read_packet_header(data []byte) ByondPacketHeader {
	return ByondPacketHeader{binary.BigEndian.Uint16(data[0:]), binary.BigEndian.Uint16(data[2:])}
}

func byond_on_topic(data string) (error, int, []byte) {
	packetData := new(bytes.Buffer)

	byond_write_packet_header(packetData, BYOND_PACKET_TYPE_TOPIC, uint16(len(data)+6))
	packetData.WriteString("\x00\x00\x00\x00\x00")
	packetData.WriteString(data)
	packetData.WriteString("\x00")

	if !byond_session_write_safe(packetData.Bytes()) {
		return nil, 0, nil
	}

	responseHeaderData := make([]byte, BYOND_PACKET_HEADER_SIZE)

	if !byond_session_read_safe(responseHeaderData) {
		return nil, 0, nil
	}

	responseHeader := byond_read_packet_header(responseHeaderData)
	responseData := make([]byte, responseHeader.DataSize)

	if !byond_session_read_safe(responseData) {
		return nil, 0, nil
	}

	responseDataType := int(responseData[0])
	if responseDataType == TOPIC_TYPE_NULL {
		return nil, responseDataType, nil
	}

	return nil, responseDataType, responseData[1:]
}
