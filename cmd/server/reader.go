package main

import (
	"log"
	"net"

	"github.com/iangcarroll/access-controller/pkg/idemia"
)

func handleReaderConnection(conn net.Conn) {
	defer conn.Close()

	message, err := idemia.ReadMessage(conn)
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	log.Printf("Received message ID: %.2x\n", message.ID)

	switch message.ID {
	case idemia.IdemiaMessageUserControlSuccessful:
		log.Printf("User ID: %s\n", message.Data)

		if allowedUsersMap[string(message.Data)] {
			log.Println("User is allowed, sending approval.")

			if err := idemia.SendBasicApproval(conn); err != nil {
				log.Println("error sending approval:", err)
			}

			if err := unlockSmartLock(*homeAssistantBaseURL, *homeAssistantAccessToken, *homeAssistantEntityID); err != nil {
				log.Println("error unlocking smart lock:", err)
			} else {
				log.Println("Smart lock unlocked successfully.")
			}
		} else {
			log.Println("User is not allowed, sending denial. Add this ID to -allowed-users flag to grant access.")

			if err := idemia.SendBasicDenial(conn); err != nil {
				log.Println("error sending denial:", err)
			}
		}
	default:
		log.Println("unknown message type from the IDEMIA reader, ignoring.")
	}
}
