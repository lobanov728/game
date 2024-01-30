package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/lobanov728/mud/game"

	"github.com/gorilla/websocket"
)

const (
	writeWait = 10 * time.Second

	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 1) / 10

	maxMessageSize = 512
)

var (
	newLine = []byte("\n")
	space   = []byte(" ")
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	id   string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

func (c *Client) readPump(world *game.World) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.conn.SetPongHandler(func(string) error {
	// 	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	// 	return nil
	// })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			break
		}

		c.hub.broadcast <- message
		event := &game.Event{}
		err = json.Unmarshal(message, &event)
		world.HandleEvent(event)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			fmt.Println("<-c.send", string(message))
			// c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newLine)
			// 	w.Write(<-c.send)
			// }

			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func initPlayerConnection(hub *Hub, world *game.World, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	client.hub.register <- client

	player := world.AddPlayer()
	conn.WriteJSON(game.Event{
		Type: game.PlayerEventInit,
		Data: game.PlayerInit{
			PlayerID: player.ID,
			Units:    world.Units,
			Tiles:    world.Tiles,
			Objects:  world.Objects,
		},
	})

	message, _ := json.Marshal(game.Event{
		Type: game.PlayerEventConnect,
		Data: game.PlayerConnect{
			Unit: *world.Units[player.ID],
		},
	})
	hub.broadcast <- message

	go func() {
		for {
			targetUnit := &game.Unit{}
			for id := range world.Units {
				if id != "mob" {
					targetUnit = world.Units[id]
				}
			}
			if targetUnit.ID != "" {
				for id := range world.Units {
					if id == "mob" {
						fmt.Println("target", targetUnit.X, targetUnit.Y, targetUnit.ID)
						fmt.Println("mob", world.Units[id].X, world.Units[id].Y)
						difX := world.Units[id].X - targetUnit.X
						difY := world.Units[id].Y - targetUnit.Y
						var direction int
						fmt.Println("x", difX)
						fmt.Println("y", difY)
						var ev game.Event

						if math.Abs(difX) <= 9 && math.Abs(difY) <= 9 {
							ev = game.Event{
								Type: game.ActionHit,
								Data: game.Hit{
									ToUnitID:   targetUnit.ID,
									FromUnitID: id,
								},
							}
						} else {
							if math.Abs(difX) >= math.Abs(difY) {
								if difX >= 9 {
									direction = game.DirectionLeft
								} else {
									direction = game.DirectionRight
								}
							} else {
								if difY >= 9 {
									direction = game.DirectionUp
								} else {
									direction = game.DirectionDown
								}
							}
							if direction != 0 {
								ev = game.Event{
									Type: game.PlayerEventMove,
									Data: game.PlayerMove{
										UnitID:    id,
										Direction: direction,
									},
								}
							}
						}

						if ev.Type != "" {
							message, _ := json.Marshal(ev)
							hub.broadcast <- message

							world.HandleEvent(&ev)
						}
					}
				}
				time.Sleep(time.Millisecond * 70)
			}
		}
	}()

	go client.writePump()
	go client.readPump(world)
}
