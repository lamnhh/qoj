package contest

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"qoj/server/src/listener"
	"qoj/server/src/socket"
	"strconv"
	"sync"
)

var listenerList map[int]*listener.List
var listLock sync.Mutex

func SendResult(res map[string]interface{}) {
	contestId := res["contestId"].(int)
	if listenerList[contestId] == nil {
		return
	}
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
		if ok != nil {
			return
		}
		contestId := int(contestId64)

		if _, subscribed := subscriptionList[contestId]; subscribed {
			// Already subscribed, ignore
			return
		}
		subscriptionList[contestId] = 1

		// Add listener
		listLock.Lock()
		if listenerList[contestId] == nil {
			listenerList[contestId] = &listener.List{}
		}
		listLock.Unlock()
		listenerList[contestId].Subscribe(conn)
	})
	server.On("unsubscribe", func(conn *websocket.Conn, data string) {
		contestId64, ok := strconv.ParseInt(data, 10, 16)
		if ok != nil {
			return
		}
		contestId := int(contestId64)

		if _, subscribed := subscriptionList[contestId]; subscribed == false {
			// Haven't subscribed, ignore
			return
		}
		delete(subscriptionList, contestId)

		// Remove listener
		listLock.Lock()
		if listenerList[contestId] != nil {
			listenerList[contestId].Unsubscribe(conn)
		}
		listLock.Unlock()
	})
	server.On("destroy", func(conn *websocket.Conn, _ string) {
		// After connection closes, remove all subscriptions
		listLock.Lock()
		for id := range subscriptionList {
			if listenerList[id] != nil {
				listenerList[id].Unsubscribe(conn)
			}
		}
		listLock.Unlock()
	})

	server.Deploy(app)
}