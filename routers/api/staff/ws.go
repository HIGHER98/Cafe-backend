package staff

import (
	"cafe/models"
	"cafe/pkg/e"
	"cafe/pkg/logging"
	"cafe/pkg/util"
	"encoding/json"
	"net/http"
	"strconv"
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
	logging.Debug("In wsWriter for connection: ", conn.RemoteAddr())
	conn.WriteMessage(websocket.TextMessage, respReq("auth"))
	for {
		select {
		case o := <-order:
			logging.Debug("Updating order channel with o: ", o, " to client: ", conn.RemoteAddr())
			if auth {
				conn.WriteMessage(websocket.TextMessage, respSuccess(o))
			}
		case <-time.After(5 * 60 * time.Second):
			auth = false
			conn.WriteMessage(websocket.TextMessage, respReq("auth"))
		}
	}
}

func wsReader(conn *websocket.Conn) {
	logging.Debug("In wsReader for ", conn.RemoteAddr())
	for {
		t, p, err := conn.ReadMessage()
		if err != nil {
			logging.Error("Failed to read websocket connection: ", err)
			break
		}
		if t != websocket.TextMessage {
			logging.Debug("Not a websocket text message. Instead type: ", t, " for client: ", conn.RemoteAddr())
			conn.WriteMessage(websocket.TextMessage, respErr(e.BAD_REQUEST))
			continue
		}
		var r Request
		err = json.Unmarshal(p, &r)
		if err != nil {
			logging.Error("Failed to read request: ", err)
			continue
		}

		if !auth && strings.Compare(r.Req, "auth") != 0 {
			conn.WriteMessage(websocket.TextMessage, respReq("auth"))
			continue
		}

		switch r.Req {
		case "auth":
			token, ok := r.Data["token"].(string)
			if !ok {
				logging.Error("Failed to get token")
				break
			}
			_, err := util.ParseToken(token)
			if err != nil {
				if err = conn.WriteMessage(websocket.TextMessage, respErr(e.UNAUTHORIZED)); err != nil {
					logging.Error("Failed to write message: ", err)
				}
				if err = conn.Close(); err != nil {
					logging.Error("Failed to close websocket connection: ", err)
				}
			}
			auth = true
			conn.WriteMessage(websocket.TextMessage, respSuccess(nil))
			logging.Debug(conn.RemoteAddr(), " has been successfully authenticated")

		case "update":
			//Ensuring it's a valid update
			if r.Data["status"] == 1 {
				if err := conn.WriteMessage(websocket.TextMessage, respErr(e.BAD_REQUEST)); err != nil {
					logging.Error("Failed to write websocket: ", err)
					break
				}
			}
			logging.Info(r.Data)
			id := int(r.Data["id"].(float64))
			//Usual typecasting was not working on status for some reason
			status, err := strconv.Atoi(r.Data["status"].(string))
			if err != nil {
				logging.Error("Failed to case status to string: ", err)
				conn.WriteMessage(websocket.TextMessage, respErr(e.BAD_REQUEST))
				continue
			}
			user := int(r.Data["user"].(float64))
			reply := Order{Id: id, Status: status}
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
			err = models.AddPurchaseActivity(id, status, user)
			if err != nil {
				logging.Error("Failed to add purchase activity: ", err)
			}
			go func(orderChan chan Order, order Order) {
				logging.Debug("Updating channel in anonymous function: ", orderChan, ", Order:", order, " for connection: ", conn.RemoteAddr())
				UpdateOrder(orderChan, order)
			}(O, reply)

		case "ping":
			conn.WriteMessage(websocket.TextMessage, []byte("pong"))

		default:
			if err = conn.Close(); err != nil {
				logging.Error("Failed to close websocket connection: ", err)
			}
		}

	}

}
