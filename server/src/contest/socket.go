package contest

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"qoj/server/src/listener"
	"qoj/server/src/socket"
	"strconv"
)

var listenerList map[int]*listener.List

func SendResult(res map[string]interface{}) {
	contestId := res["contestId"].(int)
	connList := listenerList[contestId].GetSubscriptionList()
	for _, conn := range connList {
		_ = conn.WriteJSON(res)
	}
}

func initialiseContestSocket(app *gin.Engine) {
	listenerList = make(map[int]*listener.List)
	subscriptionList := make(map[int]int)

	server := socket.NewSocket("/ws/contest")
	server.On("subscribe", func(conn *websocket.Conn, data string) {
		contestId64, ok := strconv.ParseInt(data, 10, 16)
		if ok == nil {
			return
		}
		contestId := int(contestId64)

		if _, subscribed := subscriptionList[contestId]; subscribed {
			// Already subscribed, ignore
			return
		}
		subscriptionList[contestId] = 1

		// Add listener
		listenerList[contestId].Subscribe(conn)
	})
	server.On("unsubscribe", func(conn *websocket.Conn, data string) {
		contestId64, ok := strconv.ParseInt(data, 10, 16)
		if ok == nil {
			return
		}
		contestId := int(contestId64)

		if _, subscribed := subscriptionList[contestId]; subscribed == false {
			// Haven't subscribed, ignore
			return
		}
		delete(subscriptionList, contestId)

		// Remove listener
		listenerList[contestId].Unsubscribe(conn)
	})
	server.On("destroy", func(conn *websocket.Conn, _ string) {
		// After connection closes, remove all subscriptions
		for id := range subscriptionList {
			listenerList[id].Unsubscribe(conn)
		}
	})

	server.Deploy(app)
}