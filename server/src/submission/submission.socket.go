package submission

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

var upgrader websocket.Upgrader
var listenerList map[int]*ListenerList

func socketHandler(w http.ResponseWriter, r *http.Request) {
	subscriptionList := make(map[int]int)
	conn, err := upgrader.Upgrade(w, r, nil)

	defer func() {
		// After connection closes, remove all subscriptions
		for id := range subscriptionList {
			listenerList[id].Unsubscribe(conn)
		}
	}()

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
		switch message.Type {
		case "subscribe":
			submissionId, _ := strconv.ParseInt(message.Message, 10, 16)
			if _, ok := subscriptionList[int(submissionId)]; ok {
				// Already subscribed
			} else {
				// Subscribe
				subscriptionList[int(submissionId)] = 1

				// Add listener
				listenerList[int(submissionId)].Subscribe(conn)
			}
		}
	}
}

func InitialiseSubmissionSocket(app *gin.Engine) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	judges = make(map[int]chan interface{})
	listenerList = make(map[int]*ListenerList)

	app.GET("/ws", func(ctx *gin.Context) {
		socketHandler(ctx.Writer, ctx.Request)
	})
}

// TIMEOUT CPU 0.51 MEM 18612 MAXMEM 18612 STALE 0 MAXMEM_RSS 2500