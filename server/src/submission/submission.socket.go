package submission

import (
	"qoj/server/src/listener"
	"qoj/server/src/socket"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var listenerList map[int]*listener.List

func InitialiseSubmissionSocket(app *gin.Engine) {
	judges = make(map[int]chan interface{})
	listenerList = make(map[int]*listener.List)

	subscriptionList := make(map[int]int)

	server := socket.NewSocket("/ws/status")
	server.On("subscribe", func(conn *websocket.Conn, message string) {
		submissionId, _ := strconv.ParseInt(message, 10, 16)
		if _, ok := subscriptionList[int(submissionId)]; ok || listenerList[int(submissionId)] == nil {
			// Already subscribed
		} else {
			// Subscribe
			subscriptionList[int(submissionId)] = 1

			// Add listener
			listenerList[int(submissionId)].Subscribe(conn)
		}
	})
	server.On("unsubscribe", func(conn *websocket.Conn, message string) {
		submissionId, _ := strconv.ParseInt(message, 10, 16)
		if _, ok := subscriptionList[int(submissionId)]; ok && listenerList[int(submissionId)] != nil {
			// Unsubscribe
			delete(subscriptionList, int(submissionId))

			// Remove listener
			listenerList[int(submissionId)].Unsubscribe(conn)
		} else {
			// Haven't subscribed, ignore
		}

	})
	server.On("destroy", func(conn *websocket.Conn, _ string) {
		// After connection closes, remove all subscriptions
		for id := range subscriptionList {
			listenerList[id].Unsubscribe(conn)
		}
	})

	server.Deploy(app)
}