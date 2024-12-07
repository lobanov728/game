package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 1) / 10

	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/ws", func() gin.HandlerFunc {
		return gin.HandlerFunc(func(c *gin.Context) {
			conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
			if err != nil {
				log.Println(err)
				return
			}

			go func() {
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
							fmt.Printf("error: %v", err)
						}

						break
					}
					fmt.Println("get on server", string(message))
				}
			}()

			for {
				time.Sleep(time.Second)
				conn.SetWriteDeadline(time.Now().Add(writeWait))

				w, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				w.Write([]byte("send from server"))

				if err := w.Close(); err != nil {
					return
				}
			}

		})
	}())

	r.Run(":3000")
}
