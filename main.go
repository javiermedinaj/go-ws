package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New connection", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed")
			} else {
				fmt.Println("Error reading", err)
			}
			delete(s.conns, ws)
			return

		}
		msg := string(buf[:n])
		s.broadcast([]byte(msg))
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Error writing", err)
			}
		}(ws)
	}
}

func (s *Server) HandleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("New connection", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("order book data -> %s", time.Now().String())
		ws.Write([]byte(payload))
		time.Sleep(1 * time.Second)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/ws-orderbook", websocket.Handler(server.HandleWSOrderbook))
	http.ListenAndServe(":3030", nil)
}
