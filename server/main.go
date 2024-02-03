package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lobanov728/mud/game"
)

func main() {
	world := &game.World{
		IsServer: true,
		Units:    game.Units{},
		Tiles:    game.Tiles{},
		Objects:  game.Objects{},
	}

	rnd := rand.New(rand.NewSource(time.Now().Unix()))

	skins := []string{
		"big_demon", "big_zombie",
		"goblin", "ice_zombie",
		"imp", "muddy",
		"ogre", "necromancer",
		"orc_shaman", "orc_warrior",
		"masked_orc", "orc_warrior",
	}

	for i := 0; i < 1; i++ {
		id := game.UnitID(uuid.New().String())

		actionMap := map[game.EventName]*game.Action{
			game.ActionHit:       game.NewAction(game.ActionHit, time.Second),
			game.PlayerEventMove: game.NewAction(game.PlayerEventMove, time.Millisecond*50),
		}
		world.Units["mob"] = &game.Unit{
			ID:         id,
			X:          float64(16 * i),
			Y:          16,
			SpriteName: skins[rnd.Intn(len(skins))],
			Action:     game.ActionIdle,
			Frame:      0,
			Actions:    actionMap,
		}
	}

	for x := 0; x < 20; x++ {
		for y := 0; y < 15; y++ {
			id := game.UnitID(uuid.New().String())
			world.Tiles[id] = &game.Tile{
				ID:         id,
				X:          float64(x * 16),
				Y:          float64(y * 16),
				SpriteName: fmt.Sprintf("floor_%d", rnd.Intn(7)+1),
				Action:     game.ActionIdle,
				Frame:      0,
			}
		}
	}

	world.Objects["box"] = &game.Unit{
		ID:         "box",
		X:          10,
		Y:          10,
		SpriteName: "",
		Box:        game.NewRectBox(10, 10, 300, 220),
	}

	world.Objects["door"] = &game.Unit{
		ID:         "door",
		X:          50,
		Y:          50,
		SpriteName: "doors_all",
		Box:        game.NewRectBox(50, 50, 64, 35),
	}

	world.Objects["door1"] = &game.Unit{
		ID:         "door1",
		X:          200,
		Y:          150,
		SpriteName: "doors_all",
		Box:        game.NewRectBox(200, 150, 64, 35),
	}

	hub := newHub()
	go hub.run()

	r := gin.New()
	gin.SetMode(gin.ReleaseMode)
	r.GET("/ws", func(hub *Hub, world *game.World) gin.HandlerFunc {
		return gin.HandlerFunc(func(c *gin.Context) {
			initPlayerConnection(hub, world, c.Writer, c.Request)
		})
	}(hub, world))

	r.Run(":3000")
}
