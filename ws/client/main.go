package main

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	c, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:3000/ws", nil)
	defer c.Close()
	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				fmt.Println("err", err)
			}
			fmt.Println("client get msg", string(msg))

			time.Sleep(time.Second)
		}
	}()
	for {
		time.Sleep(time.Second)

		c.WriteMessage(websocket.TextMessage, []byte("send from client"))
	}
}
