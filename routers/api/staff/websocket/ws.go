package websocket

import (
	"cafe/models"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"encoding/json"
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
		pur, err := models.GetTodaysOrders()
		if err != nil {
			logging.Error(err)
			conn.WriteMessage(t, []byte(e.GetMsg(e.ERROR)))
		}
		purBytes, err := json.Marshal(&pur)
		if err != nil {
			logging.Error(err)
			conn.WriteMessage(t, []byte(e.GetMsg(e.ERROR)))
		}
		conn.WriteMessage(t, purBytes)
	}
}
