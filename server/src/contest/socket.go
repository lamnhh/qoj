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

func createListenerList(contestId int) {
	listLock.Lock()
	if _, exists := listenerList[contestId]; exists == false {
		listenerList[contestId] = &listener.List{}
	}
	listLock.Unlock()
}

func InitialiseSocket(app *gin.Engine) {
	listenerList = make(map[int]*listener.List)
	subscriptionList := make(map[*websocket.Conn]map[int]int)

	server := socket.NewSocket("/ws/contest")
	server.On("start", func(conn *websocket.Conn, s string) {
		subscriptionList[conn] = make(map[int]int)
	})
	server.On("subscribe", func(conn *websocket.Conn, data string) {
		contestId64, ok := strconv.ParseInt(data, 10, 16)
		if ok != nil {
			return
		}
		contestId := int(contestId64)

		if _, subscribed := subscriptionList[conn][contestId]; subscribed {
			// Already subscribed, ignore
			return
		}

		subscriptionList[conn][contestId] = 1
		createListenerList(contestId)
		listenerList[contestId].Subscribe(conn)
	})
	server.On("unsubscribe", func(conn *websocket.Conn, data string) {
		contestId64, ok := strconv.ParseInt(data, 10, 16)
		if ok != nil {
			return
		}
		contestId := int(contestId64)

		if listenerList[contestId] == nil {
			return
		}

		if _, subscribed := subscriptionList[conn][contestId]; subscribed == false {
			// Has not subscribed, ignore
			return
		}
		delete(subscriptionList[conn], contestId)

		// Remove listener
		if listenerList[contestId] != nil {
			listenerList[contestId].Unsubscribe(conn)
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