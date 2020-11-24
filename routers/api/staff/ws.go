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

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// connections.
type hub struct {
	// Registered connections.
	connections map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan []byte

	// Register requests from the connections.
	register chan *connection

	// Unregister requests from connections.
	unregister chan *connection
}

var h = hub{
	broadcast:   make(chan []byte),
	register:    make(chan *connection),
	unregister:  make(chan *connection),
	connections: make(map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {
				delete(h.connections, c)
				close(c.send)
			}
		case m := <-h.broadcast:
			logging.Debug("Going to broadcast...")
			for c := range h.connections {

				select {
				case c.send <- m:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

type Order struct {
	Id     int `json:"id"`
	Status int `json:"status"`
}

var O chan Order

func InitOrders() {
	O = make(chan Order)
	go h.run()
}

//Signals websocket connections of an update to an order
func UpdateOrder(o chan Order, p Order) {
	logging.Debug("An order is being updated: Channel: ", o, " Order: ", p)
	O <- p
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
	c := &connection{send: make(chan []byte, 256), ws: conn}
	h.register <- c
	go c.wsWriter(O)
	c.wsReader()
}

func (c *connection) wsWriter(order chan Order) {
	logging.Debug("In wsWriter for connection: ", c.ws.RemoteAddr())
	c.ws.WriteMessage(websocket.TextMessage, respReq("auth"))
	for {
		select {
		case o := <-order:
			logging.Debug("Updating order channel with o: ", o, " to client: ", c.ws.RemoteAddr())
			h.broadcast <- respSuccess(o)
		case message, ok := <-c.send:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.ws.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}

		case <-time.After(5 * 60 * time.Second):
			auth = false
			c.ws.WriteMessage(websocket.TextMessage, respReq("auth"))
		}
	}
}

func (c *connection) wsReader() {
	logging.Debug("In wsReader for ", c.ws.RemoteAddr())
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()
	for {
		t, p, err := c.ws.ReadMessage()
		if err != nil {
			logging.Error("Failed to read websocket connection: ", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				logging.Error("Unexpected websocket closure: ", err)
			}
			break
		}
		if t != websocket.TextMessage {
			logging.Debug("Not a websocket text message. Instead type: ", t, " for client: ", c.ws.RemoteAddr())
			c.ws.WriteMessage(websocket.TextMessage, respErr(e.BAD_REQUEST))
			continue
		}
		var r Request
		err = json.Unmarshal(p, &r)
		if err != nil {
			logging.Error("Failed to read request: ", err)
			continue
		}

		if !auth && strings.Compare(r.Req, "auth") != 0 {
			c.ws.WriteMessage(websocket.TextMessage, respReq("auth"))
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
				if err = c.ws.WriteMessage(websocket.TextMessage, respErr(e.UNAUTHORIZED)); err != nil {
					logging.Error("Failed to write message: ", err)
				}
				if err = c.ws.Close(); err != nil {
					logging.Error("Failed to close websocket connection: ", err)
				}
			}
			auth = true
			c.ws.WriteMessage(websocket.TextMessage, respSuccess(nil))
			logging.Debug(c.ws.RemoteAddr(), " has been successfully authenticated")

		case "update":
			//Ensuring it's a valid update
			if r.Data["status"] == 1 {
				if err := c.ws.WriteMessage(websocket.TextMessage, respErr(e.BAD_REQUEST)); err != nil {
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
				c.ws.WriteMessage(websocket.TextMessage, respErr(e.BAD_REQUEST))
				continue
			}
			user := int(r.Data["user"].(float64))
			reply := Order{Id: id, Status: status}
			err = models.UpdatePurchaseStatus(reply.Id, reply.Status)
			if err != nil {
				switch err {
				case gorm.ErrRecordNotFound:
					c.ws.WriteMessage(websocket.TextMessage, respErr(e.ID_NOT_FOUND))
				default:
					logging.Error(err)
					c.ws.WriteMessage(websocket.TextMessage, respErr(e.ERROR))
				}
				continue
			}
			err = models.AddPurchaseActivity(id, status, user)
			if err != nil {
				logging.Error("Failed to add purchase activity: ", err)
			}
			go func(orderChan chan Order, order Order) {
				logging.Debug("Updating channel in anonymous function: ", orderChan, ", Order:", order, " for connection: ", c.ws.RemoteAddr())
				UpdateOrder(orderChan, order)
			}(O, reply)

		case "ping":
			c.ws.WriteMessage(websocket.TextMessage, []byte("pong"))

		default:
			if err = c.ws.Close(); err != nil {
				logging.Error("Failed to close websocket connection: ", err)
			}
		}

	}

}
