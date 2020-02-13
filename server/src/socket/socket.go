package socket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Listener struct {
	EventName string
	Callback  func(conn *websocket.Conn, data string)
}

type Socket struct {
	upgrader     websocket.Upgrader
	url          string
	listenerList []Listener
}

func NewSocket(url string) Socket {
	upgrader := websocket.Upgrader{}
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}
	return Socket{
		upgrader:     upgrader,
		url:          url,
		listenerList: make([]Listener, 0),
	}
}

func (socket *Socket) On(event string, callback func(conn *websocket.Conn, data string)) {
	socket.listenerList = append(socket.listenerList, Listener{
		EventName: event,
		Callback:  callback,
	})
}

func (socket *Socket) Deploy(app *gin.Engine) {
	app.GET(socket.url, func(ctx *gin.Context) {
		conn, err := socket.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		defer func(conn *websocket.Conn) {
			if conn == nil {
				return
			}
			for _, eventListener := range socket.listenerList {
				if eventListener.EventName == "destroy" {
					eventListener.Callback(conn, "")
				}
			}
		}(conn)

		if err != nil {
			log.Fatalln(err)
			return
		}

		for {
			message := Message{}
			if err := conn.ReadJSON(&message); err != nil {
				log.Println(err)
				return
			}
			for _, eventListener := range socket.listenerList {
				if eventListener.EventName == message.Type {
					eventListener.Callback(conn, message.Message)
				}
			}
		}
	})
}