package game

import (
	"fmt"
	"image"
	"image/color"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var frame uint64
var x, y float64

type writer interface {
	WriteJSON(v interface{}) error
}

type Game struct {
	writer         writer
	world          World
	tilesImage     *ebiten.Image
	hitCircleImage *ebiten.Image
}

func NewGame(
	writer writer,
	world World,
	tilesImage *ebiten.Image,
	hitCircleImage *ebiten.Image,
) *Game {
	return &Game{
		writer:         writer,
		world:          world,
		tilesImage:     tilesImage,
		hitCircleImage: hitCircleImage,
	}
}

var totalDown = 0

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.writer.WriteJSON(Event{
			ID:   uuid.New().String(),
			Type: PlayerEventMove,
			Data: UnitMove{
				UnitID:    g.world.MyID,
				Direction: DirectionUp,
			},
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		totalDown++
		fmt.Println("totalDown", totalDown)
		g.writer.WriteJSON(Event{
			ID:   uuid.New().String(),
			Type: PlayerEventMove,
			Data: UnitMove{
				UnitID:    g.world.MyID,
				Direction: DirectionDown,
			},
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.writer.WriteJSON(Event{
			ID:   uuid.New().String(),
			Type: PlayerEventMove,
			Data: UnitMove{
				UnitID:    g.world.MyID,
				Direction: DirectionRight,
			},
		})
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.writer.WriteJSON(Event{
			ID:   uuid.New().String(),
			Type: PlayerEventMove,
			Data: UnitMove{
				UnitID:    g.world.MyID,
				Direction: DirectionLeft,
			},
		})
	}
	fmt.Println(g.world.Units == nil)
	if g.world.Units[g.world.MyID].Action == ActionIdle {
		g.writer.WriteJSON(Event{
			ID:   uuid.New().String(),
			Type: PlayerEventIdle,
			Data: UnitMove{
				UnitID: g.world.MyID,
			},
		})
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	shadowImage := ebiten.NewImage(screenWidth, screenHeight)
	shadowImage.Fill(color.Black)
	unitImage := ebiten.NewImage(screenWidth, screenHeight)
	unitSightImage := ebiten.NewImage(screenWidth, screenHeight)
	triangleImage := ebiten.NewImage(screenWidth, screenHeight)
	smallWhiteImage := ebiten.NewImage(screenWidth, screenHeight)
	smallWhiteImage1 := ebiten.NewImage(screenWidth, screenHeight)
	smallWhiteImage1.Fill(color.White)
	smallWhiteImage.Fill(color.White)

	frame++

	for x := 0; x < 20; x++ {
		for y := 0; y < 15; y++ {
			drawOptions := &ebiten.DrawImageOptions{}
			drawOptions.GeoM.Translate(float64(x*16), float64(y*16))

			screen.DrawImage(g.tilesImage.SubImage(image.Rect(0, 0, 16, 16)).(*ebiten.Image), drawOptions)
		}
	}

	for _, tile := range g.world.Tiles {
		img, _, _ := ebitenutil.NewImageFromFile(
			fmt.Sprintf("sprites/%s.png", tile.SpriteName),
		)
		drawOptions := &ebiten.DrawImageOptions{}
		drawOptions.GeoM.Translate(tile.X, tile.Y)

		screen.DrawImage(img, drawOptions)
	}

	for _, obj := range g.world.Objects {
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

	var objects []Pointable
	for _, obj := range g.world.Objects {
		//if id == "box" {
		objects = append(objects, obj)
		//}
	}

	var rays []Line
	var player *Unit
	var playerHitPoints int
	var playerX, playerY float64

	for _, unit := range g.world.Units {
		img, _, _ := ebitenutil.NewImageFromFile(
			fmt.Sprintf("sprites/%s_%s_anim_f%d.png", unit.SpriteName, unit.Action, (frame/10)%4),
		)
		drawOptions := &ebiten.DrawImageOptions{}
		drawOptions.GeoM.Translate(unit.X-8, unit.Y-16)

		unitImage.DrawImage(img, drawOptions)

		if unit.ID == g.world.MyID {
			playerX, playerY = unit.X, unit.Y
			playerHitPoints = unit.HitPoints

			unit.Box = NewCircleBox(unit.X, unit.Y, 800)
			player = unit
			objects = append(objects, unit)
		} else {
			// fmt.Println(unit)

		}
		if unit.TriggerBox != nil {
			// box around active objects (players and mobs)
			// for _, l := range unit.TriggerBox {
			// 	vector.StrokeLine(screen,
			// 		float32(l.X1), float32(l.Y1),
			// 		float32(l.X2), float32(l.Y2),
			// 		1,
			// 		color.RGBA{255, 0, 255, 255},
			// 		true,
			// 	)
			// }
		}
		// vector.StrokeRect(screen, float32(unit.TriggerBox), float32(unit.Y), 10, 1, color.White, 1 true)

		if unit.ActionVector != nil {
			// vector where mod is look
			fmt.Println("unit.ActionVector", unit.ActionVector)
			// vector.StrokeLine(screen,
			// 	float32(unit.ActionVector.X1), float32(unit.ActionVector.Y1),
			// 	float32(unit.ActionVector.X2), float32(unit.ActionVector.Y2),
			// 	1,
			// 	color.White,
			// 	true,
			// )
		}
	}

	hitOp := &ebiten.DrawImageOptions{}
	hitOp.GeoM.Translate(-float64(32)/2, -float64(32)/2)
	hitOp.GeoM.Translate(playerX, playerY)
	hitIndex := (frame / 10) % 4
	hitSx, hitSy := 64*hitIndex+0, 64
	screen.DrawImage(
		g.hitCircleImage.SubImage(
			image.Rect(int(hitSx), hitSy, int(hitSx)+64, hitSy+64),
		).(*ebiten.Image),
		hitOp,
	)

	// playerBox := &Unit{
	// 	ID:    "",
	// 	X:     playerX,
	// 	Y:     playerY,
	// 	Sight: NewCircleBox(playerX, playerY, 800),
	// }

	for i, l := range player.Box {
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
	rays = RayCasting(playerX, playerY, 1000, objects, player)

	// opt := &ebiten.DrawTrianglesOptions{}
	// opt.Address = ebiten.AddressRepeat
	// opt.Blend = ebiten.BlendDestinationIn
	for i, line := range rays {
		nextLine := rays[(i+1)%len(rays)]

		v := RayVertices(
			playerX, playerY,
			nextLine.X2, nextLine.Y2,
			line.X2, line.Y2,
		)
		unitSightImage.DrawTriangles(v, []uint16{0, 1, 2}, smallWhiteImage1, nil)
		triangleImage.DrawTriangles(v, []uint16{0, 1, 2}, smallWhiteImage, nil)
	}

	offscreen := ebiten.NewImage(screenWidth, screenHeight)

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

	unitScreen := ebiten.NewImage(screenWidth, screenHeight)
	unitImageOpt := &ebiten.DrawImageOptions{}
	unitImageOpt.Blend = ebiten.BlendDestinationIn
	unitImage.DrawImage(unitSightImage, unitImageOpt)

	unitScreen.DrawImage(unitImage, nil)

	offscreen.DrawImage(shadowImage, shadowImageOpt)

	screen.DrawImage(unitScreen, nil)
	screen.DrawImage(offscreen, nil)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %0.2f", ebiten.ActualTPS()), 51, 51)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Hit points: %d", playerHitPoints), 51, 31)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) GetWindowSize() (int, int) {
	return screenWidth * 2, screenHeight * 2
}
