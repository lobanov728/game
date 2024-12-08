package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	common "github.com/lobanov728/mud/example"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	router.GET("/"+common.WebsocketRoute, func() gin.HandlerFunc {
		return gin.HandlerFunc(func(context *gin.Context) {
			conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
			if err != nil {
				log.Fatalf("upgrade: %s", err)
			}

			go func() {
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
							// error log
							fmt.Printf("error: %v", err)
						}

						break
					}

					// log
					fmt.Println("get on server", string(message))
				}
			}()

			for {
				// I don't know why we sleep
				// time.Sleep(time.Second)

				conn.SetWriteDeadline(time.Now().Add(common.WriteWait))

				w, err := conn.NextWriter(websocket.TextMessage)
				if err != nil {
					return
				}

				// log
				w.Write([]byte("send from server"))

				if err := w.Close(); err != nil {
					return
				}
			}

		})
	}())

	router.Run(fmt.Sprintf(":%d", common.Port))
}
