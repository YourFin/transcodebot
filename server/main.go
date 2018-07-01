package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/yourfin/TranscodeBot/common"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"html/template"
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

func rootHandler(ww http.ResponseWriter, rr *http.Request) {
	files, err := ioutil.ReadDir("./clients/")
	if err != nil {
		log.Fatal("find files err: ", err)
	}
	tmpl, err := template.ParseFiles("root.html")
	if err != nil {
		log.Fatal("template err: ", err)
	}
		tmpl.Execute(ww, files)
}

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
	fs := http.FileServer(http.Dir("clients"))
	http.Handle("/clients/", http.StripPrefix("/clients", fs))
	http.HandleFunc("/ws", echo)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
