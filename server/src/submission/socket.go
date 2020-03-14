package submission

import (
	"qoj/server/src/listener"
	"qoj/server/src/socket"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var listenerList map[int]*listener.List

func InitialiseSocket(app *gin.Engine) {
	judges = make(map[int]chan interface{})

	// listenerList[submissionId] is a list of connection, waiting to receive results from this submission
	listenerList = make(map[int]*listener.List)

	// subscriptionList[conn] is a list of submission IDs that this connection has subscribed to
	subscriptionList := make(map[*websocket.Conn]map[int]int)

	server := socket.NewSocket("/ws/status")
	server.On("start", func(conn *websocket.Conn, s string) {
		subscriptionList[conn] = make(map[int]int)
	})
	server.On("subscribe", func(conn *websocket.Conn, message string) {
		submissionId64, err := strconv.ParseInt(message, 10, 16)
		if err != nil {
			return
		}
		submissionId := int(submissionId64)

		if listenerList[submissionId] == nil {
			// Submission is judged completely, no need to subscribe
			return
		}

		// Subscribe if hasn't
		if listenerList[submissionId].HasSubscribed(conn) == false {
			listenerList[submissionId].Subscribe(conn)
			subscriptionList[conn][submissionId] = 1
		}
	})
	server.On("unsubscribe", func(conn *websocket.Conn, message string) {
		submissionId64, err := strconv.ParseInt(message, 10, 16)
		if err != nil {
			return
		}
		submissionId := int(submissionId64)

		if listenerList[submissionId] == nil {
			// Submission is judged completely, unsubscription is not allowed
			return
		}

		// Unsubscribe
		if listenerList[submissionId].HasSubscribed(conn) {
			listenerList[submissionId].Unsubscribe(conn)
			delete(subscriptionList[conn], submissionId)
		}
	})
	server.On("destroy", func(conn *websocket.Conn, _ string) {
		// After connection closes, remove all subscriptions
		for id := range subscriptionList[conn] {
			listenerList[id].Unsubscribe(conn)
		}
		delete(subscriptionList, conn)
	})

	server.Deploy(app)
}