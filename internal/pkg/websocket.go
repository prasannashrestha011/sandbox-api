package pkg

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func userCmdHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		userCmd := string(msg)
	}
}
