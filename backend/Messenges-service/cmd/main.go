package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				log.Println("client closed the conn")
				return
			}

			log.Fatal("error to read conn: ", err)
		}

		if string(data) == "test" {
			fmt.Printf("response: %s\n", data)
		}
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func messageChanel(w http.ResponseWriter, r *http.Request) {
	mc, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("error to create conn ", err)
	}

	err = mc.WriteMessage(1, []byte("Server is listening\n"))
	if err != nil {
		log.Fatal("error to write message", err)
	}

	reader(mc)
}

func main() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/message", messageChanel)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
