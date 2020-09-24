package websocket

import (
	"cafe/models"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"cafe/routers/api/staff"
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

//TODO Setup websocket to send down the pipe when channel is updated. Channel will be updated when order status is changed.
//Authenticate on establishing websocket connection - can be trusted after that; verify that claim more than Stackoverflow person did.

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func Wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("Failed to set websocket upgrade: %v", err)
		return
	}

	go writer(conn, staff.Order)
	for {
		t, _, err := conn.ReadMessage()
		if err != nil {
			logging.Error("Failed to read websocket connection: ", err)
			break
		}
		var reply staff.Purchase
		err = conn.ReadJSON(&reply)
		if err != nil {
			logging.Error(err)
			conn.WriteMessage(t, []byte(e.GetMsg(e.ERROR)))
			break
		}

		err = models.UpdatePurchaseStatus(reply.Id, reply.Status)
		if err != nil {
			logging.Error(err)
			conn.WriteMessage(t, []byte(e.GetMsg(e.ERROR)))
		}
		go staff.UpdateOrder(staff.Order, reply)
		logging.Info("Success!")
	}
}

func writer(conn *websocket.Conn, order chan staff.Purchase) {
	for {
		select {
		case o := <-order:
			logging.Info("Order has been updated")
			orderBytes, err := json.Marshal(&o)
			if err != nil {
				logging.Error("Failed to marshal json: ", err)
			}
			conn.WriteMessage(websocket.TextMessage, orderBytes)
		}
	}

}
