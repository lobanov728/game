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
			targetUnits := make([]*game.Unit, 0, 100)
			for id := range world.Units {
				if id != "mob" {
					targetUnits = append(targetUnits, world.Units[id])
				}
			}
			for id := range world.Units {
				if id == "mob" {
					for i := range targetUnits {
						targetUnit := targetUnits[i]
						mob := world.Units[id]
						targetVector := game.Line{
							X1: mob.X,
							Y1: mob.Y,
							X2: targetUnit.X,
							Y2: targetUnit.Y,
						}
						fmt.Println("target", targetUnit.X, targetUnit.Y, targetUnit.ID)
						fmt.Println("mob", mob.X, mob.Y)
						difX := mob.X - targetUnit.X
						difY := mob.Y - targetUnit.Y

						for _, action := range mob.Actions {
							fmt.Println(action.GetName(), action.IsReady())
							if action.IsReady() == 0 {
								var ev game.Event
								if action.GetName() == game.ActionHit && math.Abs(difX) <= 9 && math.Abs(difY) <= 9 {
									ev = game.Event{
										Type: game.ActionHit,
										Data: game.Hit{
											ToUnitID:   targetUnit.ID,
											FromUnitID: id,
										},
									}
								} else if action.GetName() == game.PlayerEventMove {
									var direction int
									fmt.Println("difX", difX)
									fmt.Println("difY", difY)
									type interactionData struct {
										x, y float64
										line game.Line
									}
									interactions := make([]interactionData, 0)

									newX := mob.X + game.MobStepSize*math.Cos(targetVector.Angle())
									newY := mob.Y + game.MobStepSize*math.Sin(targetVector.Angle())
									newBox := game.NewRectBox(newX, newY, mob.TriggerBox[1].X2-mob.TriggerBox[1].X1, mob.TriggerBox[0].Y2-mob.TriggerBox[0].Y1)

									for _, obj := range world.Objects {
										lineIntersectionNumber := game.GetBoxIntersection(newBox, obj.Box)
										if lineIntersectionNumber != -1 {

										}
										for _, l := range obj.Box {
											x, y, hasInteraction := game.LineIntersection(targetVector, l)
											if hasInteraction {
												interactions = append(interactions, interactionData{
													x:    x,
													y:    y,
													line: l,
												})
											}
										}
									}
									if len(interactions) > 0 {
										closestInteraction := interactions[0]
										for _, item := range interactions {
											closestLine := game.Line{
												X1: mob.X,
												Y1: mob.Y,
												X2: closestInteraction.x,
												Y2: closestInteraction.y,
											}
											currLine := game.Line{
												X1: mob.X,
												Y1: mob.Y,
												X2: item.x,
												Y2: item.y,
											}
											if closestLine.Length() > currLine.Length() {
												closestInteraction = item
											}
										}

										closestVertexLine := game.Line{
											X1: targetVector.X2,
											Y1: targetVector.Y2,
											X2: closestInteraction.line.X1,
											Y2: closestInteraction.line.Y1,
										}

										secondVertex := game.Line{
											X1: targetVector.X2,
											Y1: targetVector.Y2,
											X2: closestInteraction.line.X2,
											Y2: closestInteraction.line.Y2,
										}
										fmt.Println("closestInteraction", closestInteraction.line)
										fmt.Println("closestVertexLine", closestVertexLine)
										fmt.Println("secondVertex", secondVertex)
										if closestVertexLine.Length() > secondVertex.Length() {
											closestVertexLine = secondVertex
										}

										targetVector = game.Line{
											X1: mob.X,
											Y1: mob.Y,
											X2: closestVertexLine.X2,
											Y2: closestVertexLine.Y2,
										}
										// fmt.Println("targetVector", targetVector, targetVector.Length())
										// targetVector.X2 -= 20 * math.Cos(targetVector.Angle())
										// targetVector.Y2 -= 20 * math.Sin(targetVector.Angle())
									}
									if mob.TriggerBox != nil {

									}
									fmt.Println("targetVector", targetVector, targetVector.Length())

									direction = game.DirectionVector

									if direction != 0 {

										ev = game.Event{
											Type: game.PlayerEventMove,
											Data: game.UnitMove{
												UnitID:    id,
												Direction: direction,
												Angle:     targetVector.Angle(),
											},
										}

									}
								}

								if ev.Type != "" {
									message, _ := json.Marshal(ev)
									hub.broadcast <- message

									world.HandleEvent(&ev)
									action.Fire()
								}
							}
						}
						time.Sleep(time.Millisecond * game.GameSpeedMs)
					}
				}
			}
		}
	}()

	go client.writePump()
	go client.readPump(world)
}
