package main

import (
	"fmt"

	"github.com/gorilla/websocket"

	common "github.com/lobanov728/mud/example"
)

func main() {
	c, _, _ := websocket.DefaultDialer.Dial(fmt.Sprintf("ws://127.0.0.1:%d/%s", common.Port, common.WebsocketRoute), nil)
	defer c.Close()

	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("err", err)
			}

			fmt.Println("client get msg", string(msg))

			// I don't why we sleep
			// time.Sleep(time.Second)
		}
	}()

	for {
		// I don't why we sleep
		// time.Sleep(time.Second)

		c.WriteMessage(websocket.TextMessage, []byte("send from client"))
	}
}
