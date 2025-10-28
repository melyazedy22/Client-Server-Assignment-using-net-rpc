package main

import (
	"errors"
	"log"
	"net"
	"net/rpc"
	"sync"
	"time"
)

// Message must be exported for net/rpc
type Message struct {
	User string
	Text string
	Time time.Time
}

// Chat service
type Chat struct {
	mu   sync.Mutex
	msgs []Message
}

// SendMessage appends the incoming message and returns the entire history
func (c *Chat) SendMessage(msg Message, reply *[]Message) error {
	if msg.Text == "" {
		return errors.New("empty message")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	// add timestamp server-side to ensure consistency
	msg.Time = time.Now()
	c.msgs = append(c.msgs, msg)

	// return a copy of history
	history := make([]Message, len(c.msgs))
	copy(history, c.msgs)
	*reply = history
	return nil
}

// GetHistory returns the full history without adding a new message
func (c *Chat) GetHistory(dummy struct{}, reply *[]Message) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	history := make([]Message, len(c.msgs))
	copy(history, c.msgs)
	*reply = history
	return nil
}

func main() {
	chat := &Chat{
		msgs: make([]Message, 0, 100),
	}
	err := rpc.RegisterName("Chat", chat)
	if err != nil {
		log.Fatalf("rpc register: %v", err)
	}

	ln, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	log.Println("Chat RPC server listening on :1234")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		// Serve each connection in a goroutine
		go rpc.ServeConn(conn)
	}
}
