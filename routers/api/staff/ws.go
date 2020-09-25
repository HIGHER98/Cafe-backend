package staff

import (
	"cafe/models"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"cafe/pkg/util"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
)

type Request struct {
	Req  string                 `json:"req"`
	Data map[string]interface{} `json:"data"`
}

type Response struct {
	//Server may request from client, eg. If requesting auth to verify token still valid
	Req  string      `json:"req"`
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

//Error response to the client
func respErr(code int) []byte {
	b, err := json.Marshal(Response{Code: code, Data: e.GetMsg(code)})
	if err != nil {
		logging.Error("Failed to marshal response: ", err)
		return nil
	}
	return b
}

//Successful response to the client
func respSuccess(d interface{}) []byte {
	b, err := json.Marshal(Response{Code: e.SUCCESS, Data: d})
	if err != nil {
		logging.Error("Failed to marshal response: ", err)
		return nil
	}
	return b
}

//Make a request to the client
func respReq(req string) []byte {
	b, err := json.Marshal(Response{Req: req})
	if err != nil {
		logging.Error("Failed to marshal response: ", err)
		return nil
	}
	return b
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var a bool = false

func Wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("Failed to set websocket upgrade: %v", err)
		return
	}

	go wsWriter(conn, O)
	for {
		t, _, err := conn.ReadMessage()
		if err != nil {
			logging.Error("Failed to read websocket connection: ", err)
			break
		}

		var r Request
		err = conn.ReadJSON(&r)
		if err != nil {
			logging.Error("Failed to read request: ", err)
			continue
		}
		if !a && strings.Compare(r.Req, "auth") != 0 {
			conn.WriteMessage(websocket.TextMessage, respReq("auth"))
			/*if err = conn.Close(); err != nil {
				logging.Error("Failed to close websocket connection: ", err)
				break
			}*/
			continue
		}
		switch r.Req {
		case "auth":
			token := r.Data["token"].(string)
			_, err := util.ParseToken(token)
			if err != nil {
				conn.WriteMessage(t, respErr(e.UNAUTHORIZED))
				continue
			}
			a = true
			conn.WriteMessage(t, respSuccess(nil))
		case "update":
			var reply = Order{Id: int(r.Data["id"].(float64)), Status: int(r.Data["status"].(float64))}
			err = models.UpdatePurchaseStatus(reply.Id, reply.Status)
			if err != nil {
				switch err {
				case gorm.ErrRecordNotFound:
					conn.WriteMessage(t, respErr(e.ID_NOT_FOUND))
				default:
					logging.Error(err)
					conn.WriteMessage(t, respErr(e.ERROR))
				}
				continue
			}
			if reply.Status == 1 || reply.Status == 2 {
				go UpdateOrder(O, reply)
			}
			logging.Info("Successfully updated order through websockets")
		case "ping":
			conn.WriteMessage(websocket.TextMessage, []byte("pong"))
		default:

		}
	}
}

func wsWriter(conn *websocket.Conn, order chan Order) {
	for {
		select {
		case o := <-order:
			logging.Info("Order has been updated")
			conn.WriteMessage(websocket.TextMessage, respSuccess(o))
		case <-time.After(5 * 60 * time.Second):
			a = false
			conn.WriteMessage(websocket.TextMessage, respReq("auth"))
		}
	}

}
