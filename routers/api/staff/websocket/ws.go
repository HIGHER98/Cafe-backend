package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %v", err)
		return
	}

	for {
		t, _, err := conn.ReadMessage()
		if err != nil {
			break
		}

		conn.WriteMessage(t, []byte("pong"))
	}
}
