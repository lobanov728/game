package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/lobanov728/mud/game"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var frame uint64
var x, y float64
var world game.World

type Game struct {
	conn *websocket.Conn
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.conn.WriteJSON(game.Event{
			Type: game.PlayerEventMove,
			Data: game.PlayerMove{
				UnitID:    world.MyID,
				Direction: game.DirectionUp,
			},
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.conn.WriteJSON(game.Event{
			Type: game.PlayerEventMove,
			Data: game.PlayerMove{
				UnitID:    world.MyID,
				Direction: game.DirectionDown,
			},
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.conn.WriteJSON(game.Event{
			Type: game.PlayerEventMove,
			Data: game.PlayerMove{
				UnitID:    world.MyID,
				Direction: game.DirectionRight,
			},
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.conn.WriteJSON(game.Event{
			Type: game.PlayerEventMove,
			Data: game.PlayerMove{
				UnitID:    world.MyID,
				Direction: game.DirectionLeft,
			},
		})
	}
	if world.Units[world.MyID].Action == game.ActionRun {
		g.conn.WriteJSON(game.Event{
			Type: game.PlayerEventIdle,
			Data: game.PlayerMove{
				UnitID: world.MyID,
			},
		})
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	frame++

	// img, _, _ := ebitenutil.NewImageFromFile("sprites/background.png")
	// screen.DrawImage(img, nil)
	for _, tile := range world.Tiles {

		img, _, _ := ebitenutil.NewImageFromFile(
			fmt.Sprintf("sprites/%s.png", tile.SpriteName),
		)
		drawOptions := &ebiten.DrawImageOptions{}
		drawOptions.GeoM.Translate(tile.X, tile.Y)

		screen.DrawImage(img, drawOptions)
	}

	for _, obj := range world.Objects {
		if obj.SpriteName != "" {
			img, _, _ := ebitenutil.NewImageFromFile(
				fmt.Sprintf("sprites/%s.png", obj.SpriteName),
			)
			drawOptions := &ebiten.DrawImageOptions{}
			drawOptions.GeoM.Translate(obj.X, obj.Y)

			screen.DrawImage(img, drawOptions)
		}

		for i, l := range obj.Box {
			vector.StrokeLine(
				screen,
				float32(l.X1),
				float32(l.Y1),
				float32(l.X2),
				float32(l.Y2),
				1,
				color.RGBA{uint8(40 * i), 0, 0, 255},
				true,
			)
		}
	}

	var rays []game.Line
	var playerX, playerY float64

	for _, unit := range world.Units {
		img, _, _ := ebitenutil.NewImageFromFile(
			fmt.Sprintf("sprites/%s_%s_anim_f%d.png", unit.SpriteName, unit.Action, (frame/10)%4),
		)
		drawOptions := &ebiten.DrawImageOptions{}
		drawOptions.GeoM.Translate(unit.X, unit.Y)

		screen.DrawImage(img, drawOptions)

		if unit.ID == world.MyID {
			playerX, playerY = unit.X, unit.Y
			rays = game.RayCasting(unit.X+8, unit.Y+16, 100, world.Objects)
		}
	}

	triangleImage := ebiten.NewImage(320, 240)
	triangleImage.Fill(color.White)
	opt := &ebiten.DrawTrianglesOptions{}
	// opt.Address = ebiten.AddressRepeat
	// opt.Blend = ebiten.BlendSourceOut
	for i, line := range rays {
		if i+1 == len(rays) {
			// break
		}
		nextLine := rays[(i+1)%len(rays)]

		v := game.RayVertices(
			playerX+8, playerY+16,
			nextLine.X2, nextLine.Y2,
			line.X2, line.Y2,
		)
		screen.DrawTriangles(v, []uint16{0, 1, 2}, triangleImage, opt)
	}

	for i, l := range rays {
		fmt.Println("i x2, y2 angle", i, l.X2, l.Y2, l.Angle())
		vector.StrokeLine(
			screen,
			float32(l.X1),
			float32(l.Y1),
			float32(l.X2),
			float32(l.Y2),
			1,
			color.RGBA{uint8(40 * i), 0, 0, 255},
			true,
		)
	}

	shadowImage := ebiten.NewImage(100, 100)

	shadowImage.Fill(color.Black)

	circle := ebiten.NewImage(100, 100)
	circle.Fill(color.White)

	opt = &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.Blend = ebiten.BlendSourceOut
	vector.DrawFilledCircle(shadowImage, 50, 50, 50, color.White, true)
	// shadowImage.DrawTriangles([]ebiten.Vertex{
	// 	{
	// 		DstX:   30,
	// 		DstY:   30,
	// 		SrcX:   0,
	// 		SrcY:   0,
	// 		ColorR: 1,
	// 		ColorG: 1,
	// 		ColorB: 1,
	// 		ColorA: 1,
	// 	},
	// 	{
	// 		DstX:   30,
	// 		DstY:   50,
	// 		SrcX:   0,
	// 		SrcY:   0,
	// 		ColorR: 1,
	// 		ColorG: 1,
	// 		ColorB: 1,
	// 		ColorA: 1,
	// 	},
	// 	{
	// 		DstX:   90,
	// 		DstY:   90,
	// 		SrcX:   0,
	// 		SrcY:   0,
	// 		ColorR: 1,
	// 		ColorG: 1,
	// 		ColorB: 1,
	// 		ColorA: 1,
	// 	},
	// 	{
	// 		DstX:   50,
	// 		DstY:   90,
	// 		SrcX:   0,
	// 		SrcY:   0,
	// 		ColorR: 1,
	// 		ColorG: 1,
	// 		ColorB: 1,
	// 		ColorA: 1,
	// 	},
	// 	{
	// 		DstX:   5,
	// 		DstY:   90,
	// 		SrcX:   0,
	// 		SrcY:   0,
	// 		ColorR: 1,
	// 		ColorG: 1,
	// 		ColorB: 1,
	// 		ColorA: 1,
	// 	},
	// 	{
	// 		DstX:   5,
	// 		DstY:   5,
	// 		SrcX:   0,
	// 		SrcY:   0,
	// 		ColorR: 1,
	// 		ColorG: 1,
	// 		ColorB: 1,
	// 		ColorA: 1,
	// 	},
	// }, []uint16{
	// 	0,
	// 	1,
	// 	2,
	// 	3,
	// 	4,
	// 	5,
	// },
	// 	circle,
	// 	opt)

	op := &ebiten.DrawImageOptions{}
	op.ColorScale.ScaleAlpha(0.7)
	screen.DrawImage(shadowImage, op)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	c, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:3000/ws", nil)

	go func(c *websocket.Conn) {
		defer c.Close()

		for {
			var event game.Event
			c.ReadJSON(&event)
			world.HandleEvent(&event)
		}
	}(c)

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle(" Hello world ")
	if err := ebiten.RunGame(&Game{conn: c}); err != nil {
		log.Fatal(err)
	}
}
