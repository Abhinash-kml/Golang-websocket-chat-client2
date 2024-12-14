package main

import (
	"bufio"
	"encoding/json"
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
					textPayload := &models.TextMessage{}
					jsonPayload, _ := json.Marshal(message.Payload)
					json.Unmarshal(jsonPayload, textPayload)
					fmt.Printf("\033[38;2;255;0;0mSender: %s\nType: %s\nMessage: %+v\n\033[0m", message.Sender, message.Type, textPayload.Data)
				}
			case "binary":
				{
					binaryPayload := &models.BinaryMessage{}
					jsonPayload, _ := json.Marshal(message.Payload)
					json.Unmarshal(jsonPayload, binaryPayload)
					fmt.Printf("\033[38;2;0;255;0mSender: %s\nType: %s\nMessage: %+v\n\033[0m", message.Sender, message.Type, *binaryPayload)
				}
			case "bytearray":
				{
					bytePayload := &models.ByteArray{}
					jsonPayload, _ := json.Marshal(message.Payload)
					json.Unmarshal(jsonPayload, bytePayload)
					fmt.Printf("\033[38;2;0;0;255mSender: %s\nType: %s\nMessage: %+v\n\033[0m", message.Sender, message.Type, *bytePayload)
				}
			}
		}
	}()

	for scanner.Scan() {
		text := scanner.Text()
		message := models.Message{
			Sender: username,
		}

		if len(text) <= 0 {
			fmt.Println("next")
			continue
		}

		switch text[0] {
		case '0':
			{
				message.Type = "text"
				message.Payload = models.TextMessage{
					Data: text[2:],
				}
			}
		case '1':
			{
				message.Type = "binary"
				message.Payload = models.BinaryMessage{
					Data: text[2:],
				}
			}
		case '2':
			{
				message.Type = "bytearray"
				message.Payload = models.ByteArray{
					Data: []byte(text[2:]),
				}
			}
		}

		switch text[1] {
		case 'g':
			{
				message.Channel = "general"
			}
		case 'h':
			{
				message.Channel = "hindi"
			}
		case 'e':
			{
				message.Channel = "english"
			}
		case 'b':
			{
				message.Channel = "bakchodi"
			}
		}

		err := conn.WriteJSON(message)
		if err != nil {
			log.Fatal(err)
		}
	}
}
