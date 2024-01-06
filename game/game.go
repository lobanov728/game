package game

import (
	"encoding/json"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Unit struct {
	ID                  UnitID  `json:"id"`
	X                   float64 `json:"x"`
	Y                   float64 `json:"y"`
	SpriteName          string  `json:"sprite_name"`
	Action              string  `json:"action"`
	Frame               int     `json:"frame"`
	HorizontalDirection int     `json:"direction"`
	Box                 []Line  `json:"line"`
}

func (u *Unit) Points() [][2]float64 {
	var points [][2]float64

	for _, box := range u.Box {
		points = append(points, [2]float64{box.X2, box.Y2})
	}
	p := [2]float64{u.Box[0].X1, u.Box[0].Y1}

	if p[0] != points[len(points)-1][0] && p[1] != points[len(points)-1][1] {
		points = append(points, [2]float64{u.Box[0].X1, u.Box[0].Y1})
	}

	return points
}

func (u *Unit) GetBox() []Line {
	return u.Box
}

func NewRectBox(x, y, w, h float64) []Line {
	return []Line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}

func NewCircleBox(x, y, r float64) []Line {
	var path vector.Path
	path.Arc(float32(x), float32(y), float32(r), 0, 2*math.Pi, vector.Clockwise)
	vs, _ := path.AppendVerticesAndIndicesForFilling(nil, nil)

	var res []Line
	for i := 0; i < len(vs)-1; i++ {
		nextLine := vs[i+1]
		res = append(res, Line{float64(vs[i].DstX), float64(vs[i].DstY), float64(nextLine.DstX), float64(nextLine.DstY)})
	}

	return res
}

type Line struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

func (l *Line) Angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

type Tile struct {
	ID         UnitID  `json:"id"`
	X          float64 `json:"x"`
	Y          float64 `json:"y"`
	SpriteName string  `json:"sprite_name"`
	Action     string  `json:"action"`
	Frame      int     `json:"frame"`
}

type UnitID string

type Units map[UnitID]*Unit
type Objects map[UnitID]*Unit
type Tiles map[UnitID]*Tile

type World struct {
	MyID     UnitID  `json:"_"`
	IsServer bool    `json:"_"`
	Units    Units   `json:"units"`
	Objects  Objects `json:"objects"`
	Tiles    Tiles   `json:"tiles"`
}

type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type PlayerConnect struct {
	Unit
}

type PlayerMove struct {
	UnitID    UnitID `json:"unit_id"`
	Direction int    `json:"direction"`
}

type PlayerIdle struct {
	UnitID UnitID `json:"unit_id"`
}

type PlayerInit struct {
	PlayerID UnitID  `json:"player_ud"`
	Units    Units   `json:"units"`
	Tiles    Tiles   `json:"tiles"`
	Objects  Objects `json:"objects"`
}

const (
	PlayerEventConnect = "connect"
	PlayerEventMove    = "move"
	PlayerEventIdle    = "idle"
	PlayerEventInit    = "init"

	ActionRun  = "run"
	ActionIdle = "idle"
)

const (
	DirectionUp = iota
	DirectionRight
	DirectionDown
	DirectionLeft
)

func (world *World) HandleEvent(event *Event) {

	switch event.Type {
	case PlayerEventConnect:
		str, _ := json.Marshal(event.Data)
		var ev PlayerConnect
		json.Unmarshal(str, &ev)

		world.Units[ev.ID] = &ev.Unit

	case PlayerEventInit:
		str, _ := json.Marshal(event.Data)
		var ev PlayerInit
		json.Unmarshal(str, &ev)

		if !world.IsServer {
			world.MyID = ev.PlayerID
			world.Units = ev.Units
			world.Tiles = ev.Tiles
			world.Objects = ev.Objects
		}
	case PlayerEventMove:
		str, _ := json.Marshal(event.Data)
		var ev PlayerMove
		json.Unmarshal(str, &ev)

		unit := world.Units[ev.UnitID]
		unit.Action = ActionRun

		switch ev.Direction {
		case DirectionUp:
			unit.Y -= 1.5
		case DirectionDown:
			unit.Y += 1.5
		case DirectionRight:
			unit.X += 1.5
			unit.HorizontalDirection = ev.Direction
		case DirectionLeft:
			unit.X -= 1.5
			unit.HorizontalDirection = ev.Direction
		}
	case PlayerEventIdle:
		str, _ := json.Marshal(event.Data)
		var ev PlayerMove
		json.Unmarshal(str, &ev)

		unit := world.Units[ev.UnitID]
		unit.Action = ActionIdle
	}
}

func (world *World) AddPlayer() *Unit {
	skins := []string{
		"elf_f", "elf_m",
		"knight_f", "knight_m",
		"lizard_f", "lizard_m",
	}

	id := uuid.New().String()
	rnd := rand.New(rand.NewSource(time.Now().Unix()))
	unit := &Unit{
		ID:         UnitID(id),
		X:          50,
		Y:          50,
		SpriteName: skins[rnd.Intn(len(skins))],
		Action:     ActionIdle,
		Frame:      rnd.Intn(4),
	}
	world.Units[UnitID(id)] = unit

	return unit
}
