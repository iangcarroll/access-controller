package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	tcpServerPort            = flag.String("port", "11020", "Port for the TCP server to listen on")
	reader                   = flag.String("reader", "idemia", "Reader to receive connections from (only 'idemia' is supported)")
	homeAssistantAccessToken = flag.String("token", "", "Access token for Home Assistant")
	homeAssistantEntityID    = flag.String("entity", "", "Lock entity ID for Home Assistant")
	homeAssistantBaseURL     = flag.String("base-url", "http://homeassistant.local:8123", "Base URL for Home Assistant")
	allowedUsers             = flag.String("allowed-users", "123", "Comma-separated list of allowed user IDs")
	allowedUsersMap          = make(map[string]bool)
)

func idemiaReader() {
	ln, err := net.Listen("tcp", ":"+*tcpServerPort)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Printf("Listening on port %s...\n", *tcpServerPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleReaderConnection(conn)
	}
}

func main() {
	// Define command-line flags
	flag.Parse()

	if *homeAssistantBaseURL == "" || *homeAssistantAccessToken == "" || *homeAssistantEntityID == "" || *allowedUsers == "" {
		flag.Usage()
		log.Fatal("Access token and entity ID must be provided")
	}

	// Populate allowed users map
	for _, user := range strings.Split(*allowedUsers, ",") {
		allowedUsersMap[user] = true
	}

	switch *reader {
	case "idemia":
		idemiaReader()
	default:
		log.Fatalf("Unsupported reader: %s", *reader)
	}
}
