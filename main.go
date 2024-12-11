package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"v/models"

	"github.com/gorilla/websocket"
)

func main() {
	var username string

	fmt.Println("Enter your username.")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		username = scanner.Text()
	}
	fmt.Printf("\n")

	serverurl := "ws://localhost:8000/ws"
	conn, _, err := websocket.DefaultDialer.Dial(serverurl, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	go func() {
		for {
			message := &models.Message{}
			err := conn.ReadJSON(message)
			if err != nil {
				fmt.Println(err)
				return
			}

			switch message.Type {
			case "text":
				{
					fmt.Printf("Sender: %s\nType: %s\nMessage: %s\n", message.Sender, message.Type, message.Payload)
				}
			case "binary":
				{
					fmt.Printf("Sender: %s\nType: %s\nMessage: %s\n", message.Sender, message.Type, message.Payload)
				}
			case "bytearray":
				{
					fmt.Printf("Sender: %s\nType: %s\nMessage: %s\n", message.Sender, message.Type, message.Payload)
				}
			}
		}
	}()

	for scanner.Scan() {
		text := scanner.Text()
		message := models.Message{
			Sender: username,
			Type:   "text",
			Payload: models.TextMessage{
				Data: text,
			},
		}

		err := conn.WriteJSON(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}
