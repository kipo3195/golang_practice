package java

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

func Connect(wg *sync.WaitGroup) {
	defer wg.Done()

	var conns []*websocket.Conn
	conns = make([]*websocket.Conn, 0, 3000)

	for i := 0; i < 3000; i++ {

		id := fmt.Sprintf("id%d", i)
		url := "ws://172.16.10.114/notificator/ws/connect?user_id=" + id

		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Println("connection failed:", err)
			continue
		}

		conns = append(conns, conn)
	}

	log.Println("connections:", len(conns))

	select {}
}
