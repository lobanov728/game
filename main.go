package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"log"
	"time"

	"github.com/lobanov728/mud/game"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
)

var world game.World

var (
	hitCircleImage *ebiten.Image
	tilesImage     *ebiten.Image
)

var (
	//go:embed sprites/floor_1.png
	Tiles_png []byte

	//go:embed sprites/hit/circle02.png
	Circle_2_png []byte
)

func main() {
	img, _, err := image.Decode(bytes.NewReader(Tiles_png))
	if err != nil {
		log.Fatal(err)
	}
	tilesImage = ebiten.NewImageFromImage(img)

	hitImage, _, err := image.Decode(bytes.NewReader(Circle_2_png))
	if err != nil {
		log.Fatal(err)
	}
	hitCircleImage = ebiten.NewImageFromImage(hitImage)

	connection, _, _ := websocket.DefaultDialer.Dial("ws://127.0.0.1:3000/ws", nil)
	defer connection.Close()

	go func() {
		for {
			var event game.Event
			err := connection.ReadJSON(&event)
			if err != nil {
				fmt.Println("err", err)
			}
			world.HandleEvent(&event)
		}
	}()

	// sleep to block goroutine for init connection to the world
	time.Sleep(time.Millisecond * 10)
	g := game.NewGame(
		connection,
		world,
		tilesImage,
		hitCircleImage,
	)

	ebiten.SetWindowSize(g.GetWindowSize())
	ebiten.SetWindowTitle(" Hello world ")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
