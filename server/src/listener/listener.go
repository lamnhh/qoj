package listener

import (
	"github.com/gorilla/websocket"
	"sync"
)

type List struct {
	list map[*websocket.Conn]int
	mux  sync.Mutex
}

func (l *List) Subscribe(conn *websocket.Conn) {
	l.mux.Lock()
	if l.list == nil {
		l.list = make(map[*websocket.Conn]int)
	}
	l.list[conn] = 1
	l.mux.Unlock()
}

func (l *List) HasSubscribed(conn *websocket.Conn) bool {
	l.mux.Lock()
	var ans bool
	if _, ok := l.list[conn]; ok {
		ans = true
	} else {
		ans = false
	}
	l.mux.Unlock()
	return ans
}

func (l *List) Unsubscribe(conn *websocket.Conn) {
	l.mux.Lock()
	delete(l.list, conn)
	l.mux.Unlock()
}

func (l *List) GetSubscriptionList() []*websocket.Conn {
	list := make([]*websocket.Conn, 0)
	l.mux.Lock()
	for k := range l.list {
		list = append(list, k)
	}
	l.mux.Unlock()
	return list
}
