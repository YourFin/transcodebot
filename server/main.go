package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/yourfin/TranscodeBot/common"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

//func handler(ww http.ResponseWriter, rr *http.Request) {
//	fmt.Fprintf(ww, "Hi, I love %s!", rr.URL.Path[1:])
//	conn, err := upgrader.Upgrade(ww, rr, nil)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	if err := conn.WriteMessage()
//}

func echo(ww http.ResponseWriter, rr *http.Request) {
	conn, err := upgrader.Upgrade(ww, rr, nil)
	if err != nil {
		log.Println("upgrade err:", err)
		return
	}
	fmt.Printf("Connected to: %s", rr)
	defer conn.Close()
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("write err: ", err)
			break
		}
		log.Printf("recv: %s\n", message)
		log.Printf("recv mt: %s\n", mt)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write err:", err)
		}
	}
}

func main() {
	fmt.Printf("%s\n", common.Computer{})
	http.HandleFunc("/ws", echo)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
