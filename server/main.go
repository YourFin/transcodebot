// Copyright © 2018 Patrick Nuckolls <nuckollsp at gmail>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package server

import (
	"fmt"
	"log"
	"net/http"
	"github.com/yourfin/transcodebot/common"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"html/template"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

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
