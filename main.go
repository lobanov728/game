package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
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

var (
	tilesImage *ebiten.Image
)

var (
	//go:embed sprites/floor_1.png
	Tiles_png []byte
)

func init() {
	// Decode an image from the image file's byte slice.
	img, _, err := image.Decode(bytes.NewReader(Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)
}

func (g *Game) Draw(screen *ebiten.Image) {
	shadowImage := ebiten.NewImage(320, 240)
	shadowImage.Fill(color.Black)
	unitImage := ebiten.NewImage(320, 240)
	unitSightImage := ebiten.NewImage(320, 240)
	triangleImage := ebiten.NewImage(320, 240)
	smallWhiteImage := ebiten.NewImage(320, 240)
	smallWhiteImage1 := ebiten.NewImage(320, 240)
	smallWhiteImage1.Fill(color.White)
	smallWhiteImage.Fill(color.White)

	frame++

	// for x := 0; x < 40; x++ {
	// 	for y := 0; y < 30; y++ {
	// 		drawOptions := &ebiten.DrawImageOptions{}
	// 		drawOptions.GeoM.Translate(float64(x*16), float64(y*16))

	// 		screen.DrawImage(tilesImage.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), drawOptions)
	// 	}
	// }

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

		unitImage.DrawImage(img, drawOptions)

		if unit.ID == world.MyID {
			playerX, playerY = unit.X+8, unit.Y+16
		}
	}

	var objects []game.Pointable
	for _, obj := range world.Objects {
		objects = append(objects, obj)
	}

	playerBox := &game.Unit{
		ID:  "",
		X:   playerX,
		Y:   playerY,
		Box: game.NewCircleBox(playerX, playerY, 50),
	}

	objects = append(objects, playerBox)
	for i, l := range playerBox.Box {
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
	rays = game.RayCasting(playerX, playerY, 1000, objects, playerBox)

	// opt := &ebiten.DrawTrianglesOptions{}
	// opt.Address = ebiten.AddressRepeat
	// opt.Blend = ebiten.BlendDestinationIn
	for i, line := range rays {
		nextLine := rays[(i+1)%len(rays)]

		v := game.RayVertices(
			playerX, playerY,
			nextLine.X2, nextLine.Y2,
			line.X2, line.Y2,
		)
		unitSightImage.DrawTriangles(v, []uint16{0, 1, 2}, smallWhiteImage1, nil)
		triangleImage.DrawTriangles(v, []uint16{0, 1, 2}, smallWhiteImage, nil)
	}

	offscreen := ebiten.NewImage(320, 240)

	shadowImageOpt := &ebiten.DrawImageOptions{}
	shadowImageOpt.ColorScale.ScaleAlpha(0.5)
	triangleImageOpt := &ebiten.DrawImageOptions{}
	triangleImageOpt.Blend = ebiten.BlendDestinationOut
	shadowImage.DrawImage(triangleImage, triangleImageOpt)

	//
	// triangleImage.DrawImage(unitImage, unitImageOpt)

	// triangleImageOpt := &ebiten.DrawImageOptions{}
	// triangleImageOpt.Blend = ebiten.BlendSourceOver
	// shadowImage.DrawImage(triangleImage, triangleImageOpt)
	// shadowImage.DrawImage(triangleImage, triangleImageOpt)

	unitScreen := ebiten.NewImage(320, 240)
	unitImageOpt := &ebiten.DrawImageOptions{}
	unitImageOpt.Blend = ebiten.BlendDestinationIn
	unitImage.DrawImage(unitSightImage, unitImageOpt)

	unitScreen.DrawImage(unitImage, nil)

	offscreen.DrawImage(shadowImage, shadowImageOpt)

	screen.DrawImage(unitScreen, nil)
	screen.DrawImage(offscreen, nil)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
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
		// if err := ebiten.RunGame(&Game{conn:c nil}); err != nil {
		log.Fatal(err)
	}
}
