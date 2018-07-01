package main

import (
	//"fmt"
	"log"
	"net/url"
	"os"
	"time"
	"os/signal"
	//"github.com/yourfin/TranscodeBot/common"
	"github.com/gorilla/websocket"
)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/"}
	log.Printf("Connecting to %s...", u.String())
	connection, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial: ", err)
	}
	defer connection.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := connection.ReadMessage()
			if err != nil {
				log.Printf("read: %s\n", err)
			}
			log.Printf("recv: %s\n", message)
		}
	}()
	go func() {
		for {
			connection.WriteMessage(websocket.TextMessage, []byte("potato"))
		}
	}()
	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("Interrupt sent")

			err := connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close err: ", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
