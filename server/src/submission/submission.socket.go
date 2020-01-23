package submission

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader websocket.Upgrader

type Message struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func distributeResult(judge chan interface{}, conn *websocket.Conn) {
	for {
		select {
		case res := <-judge:
			_ = conn.WriteJSON(res.(map[string]interface{}))
		}
	}
}

func socketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
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
			judge := judges[int(submissionId)]
			go distributeResult(judge, conn)
		}
	}
}

func InitialiseSubmissionSocket(app *gin.Engine) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	judges = make(map[int]chan interface{})
	app.GET("/ws", func(ctx *gin.Context) {
		socketHandler(ctx.Writer, ctx.Request)
	})
}

// TIMEOUT CPU 0.51 MEM 18612 MAXMEM 18612 STALE 0 MAXMEM_RSS 2500