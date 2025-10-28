package main

import (
	"bufio"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strings"
	"time"
)

// Message struct must match server
type Message struct {
	User string
	Text string
	Time time.Time
}

func dialServer() (*rpc.Client, error) {
	return rpc.Dial("tcp", "127.0.0.1:1234")
}

func printHistory(history []Message) {
	fmt.Println("----- chat history -----")
	for _, m := range history {
		// format time nicely
		t := m.Time.Format("2006-01-02 15:04:05")
		fmt.Printf("[%s] %s: %s\n", t, m.User, m.Text)
	}
	fmt.Println("------------------------")
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter your name: ")
	nameRaw, _ := reader.ReadString('\n')
	name := strings.TrimSpace(nameRaw)
	if name == "" {
		name = "anonymous"
	}

	// try to connect once
	client, err := dialServer()
	if err != nil {
		log.Printf("couldn't connect to server on start: %v", err)
		client = nil
	} else {
		log.Println("connected to server")
	}

	for {
		fmt.Print("> ")
		textRaw, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("read error: %v", err)
			continue
		}
		text := strings.TrimSpace(textRaw)
		if text == "" {
			continue
		}
		if text == "exit" {
			fmt.Println("Goodbye.")
			return
		}

		msg := Message{
			User: name,
			Text: text,
			Time: time.Now(),
		}

		if client == nil {
			var dErr error
			for i := 0; i < 3; i++ {
				client, dErr = dialServer()
				if dErr == nil {
					log.Println("reconnected to server")
					break
				}
				log.Printf("reconnect attempt %d failed: %v", i+1, dErr)
			}
			if client == nil {
				log.Println("could not connect to server. message not sent.")
				continue
			}
		}

		var history []Message
		callErr := client.Call("Chat.SendMessage", msg, &history)
		if callErr != nil {
			log.Printf("RPC error: %v", callErr)
			client.Close()
			client = nil
			continue
		}

		printHistory(history)
	}
}
