package contest
//
//import (
//	"log"
//	"net/http"
//	"qoj/server/src/listener"
//	"strconv"
//
//	"github.com/gin-gonic/gin"
//	"github.com/gorilla/websocket"
//)
//
//type Message struct {
//	Type    string `json:"type"`
//	Message string `json:"message"`
//}
//
//var upgrader websocket.Upgrader
//var listenerList map[int]*listener.List
//
//func socketHandler(w http.ResponseWriter, r *http.Request) {
//	subscriptionList := make(map[int]int)
//	conn, err := upgrader.Upgrade(w, r, nil)
//
//	defer func(conn *websocket.Conn) {
//		// After connection closes, remove all subscriptions
//		for id := range subscriptionList {
//			listenerList[id].Unsubscribe(conn)
//		}
//	}(conn)
//
//	if err != nil {
//		log.Fatalln(err)
//		return
//	}
//
//	for {
//		message := Message{}
//		if err := conn.ReadJSON(&message); err != nil {
//			log.Println(err)
//			return
//		}
//		switch message.Type {
//		case "subscribe":
//			submissionId, _ := strconv.ParseInt(message.Message, 10, 16)
//			if _, ok := subscriptionList[int(submissionId)]; ok || listenerList[int(submissionId)] == nil {
//				// Already subscribed
//			} else {
//				// Subscribe
//				subscriptionList[int(submissionId)] = 1
//
//				// Add listener
//				listenerList[int(submissionId)].Subscribe(conn)
//			}
//		case "unsubscribe":
//			submissionId, _ := strconv.ParseInt(message.Message, 10, 16)
//			if _, ok := subscriptionList[int(submissionId)]; ok && listenerList[int(submissionId)] != nil {
//				// Unsubscribe
//				delete(subscriptionList, int(submissionId))
//
//				// Remove listener
//				listenerList[int(submissionId)].Unsubscribe(conn)
//			} else {
//				// Haven't subscribed, ignore
//			}
//		}
//	}
//}
//
//func InitialiseSubmissionSocket(app *gin.Engine) {
//	upgrader.CheckOrigin = func(r *http.Request) bool {
//		return true
//	}
//
//	judges = make(map[int]chan interface{})
//	listenerList = make(map[int]*listener.List)
//
//	app.GET("/ws", func(ctx *gin.Context) {
//		socketHandler(ctx.Writer, ctx.Request)
//	})
//}
