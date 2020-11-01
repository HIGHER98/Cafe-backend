package staff

import (
	"cafe/models"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"cafe/pkg/util"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
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
	fmt.Println("Have response written")
	fmt.Println(b)
	return b
}

var wsupgrader = websocket.Upgrader{
	//Allow any origin
	CheckOrigin:     func(r *http.Request) bool { return true },
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var auth bool = false

func Wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		logging.Error("Failed to set websocket upgrade: %v", err)
		return
	}
	logging.Info("Connected websocket to remote client on: ", conn.RemoteAddr())
	go wsWriter(conn, O)
	go wsReader(conn)
}

func wsWriter(conn *websocket.Conn, order chan Order) {
	conn.WriteMessage(websocket.TextMessage, respReq("auth"))
	for {
		select {
		case o := <-order:
			if auth {
				logging.Info("Order has been updated")
				fmt.Println("Updating order")
				fmt.Println(o)
				conn.WriteMessage(websocket.TextMessage, respSuccess(o))
				fmt.Println("Written")
			}
		case <-time.After(5 * 60 * time.Second):
			auth = false
			conn.WriteMessage(websocket.TextMessage, respReq("auth"))
		}
	}
}

func wsReader(conn *websocket.Conn) {
	for {
		t, p, err := conn.ReadMessage()
		if err != nil {
			logging.Error("Failed to read websocket connection: ", err)
			break
		}
		if t != websocket.TextMessage {
			conn.WriteMessage(websocket.TextMessage, respErr(e.BAD_REQUEST))
			continue
		}
		var r Request
		err = json.Unmarshal(p, &r)
		if err != nil {
			logging.Error("Failed to read request: ", err)
			continue
		}
		fmt.Println(time.Now().Format(time.ANSIC), "Got a message: ", r.Req, "\t", r.Data)

		if !auth && strings.Compare(r.Req, "auth") != 0 {
			conn.WriteMessage(websocket.TextMessage, respReq("auth"))
			continue
		}

		switch r.Req {
		case "auth":
			logging.Info("Checking auth")
			token := r.Data["token"].(string)
			_, err := util.ParseToken(token)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, respErr(e.UNAUTHORIZED))
				if err = conn.Close(); err != nil {
					logging.Error("Failed to close websocket connection: ", err)
				}
			}
			auth = true
			logging.Info("Successfully authenticated client's websocket connection")
			conn.WriteMessage(websocket.TextMessage, respSuccess(nil))

		case "update":
			if r.Data["status"] == 1 {
				conn.WriteMessage(websocket.TextMessage, respErr(e.BAD_REQUEST))
				break
			}
			var reply = Order{Id: int(r.Data["id"].(float64)), Status: int(r.Data["status"].(float64))}
			err = models.UpdatePurchaseStatus(reply.Id, reply.Status)
			if err != nil {
				switch err {
				case gorm.ErrRecordNotFound:
					conn.WriteMessage(websocket.TextMessage, respErr(e.ID_NOT_FOUND))
				default:
					logging.Error(err)
					conn.WriteMessage(websocket.TextMessage, respErr(e.ERROR))
				}
				continue
			}
			//If not setting status to 'Processing payment'
			if reply.Status != 1 {
				go UpdateOrder(O, reply)
				logging.Info("Successfully updated order through websockets")
			}

		case "ping":
			conn.WriteMessage(websocket.TextMessage, []byte("pong"))

		default:
			if err = conn.Close(); err != nil {
				logging.Error("Failed to close websocket connection: ", err)
			}
		}

	}

}
