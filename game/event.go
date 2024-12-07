package game

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID   string      `json:"id"`
	Type EventName   `json:"type"`
	Data interface{} `json:"data"`
}

type PlayerConnect struct {
	Unit
}

type UnitMove struct {
	UnitID    UnitID  `json:"unit_id"`
	Direction int     `json:"direction"`
	Angle     float64 `json:"angle"`
}

type Hit struct {
	ToUnitID   UnitID `json:"to_unit_id"`
	FromUnitID UnitID `json:"from_unit_id"`
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

type EventName string

const StepSize = 1.5
const GameSpeedMs = 100

const (
	PlayerEventConnect EventName = "connect"
	PlayerEventMove    EventName = "move"
	PlayerEventIdle    EventName = "idle"
	PlayerEventInit    EventName = "init"

	ActionRun  EventName = "run"
	ActionHit  EventName = "hit"
	ActionIdle EventName = "idle"
)

const (
	DirectionUp = iota + 1
	DirectionRight
	DirectionDown
	DirectionLeft
	DirectionVector
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
		var ev UnitMove
		json.Unmarshal(str, &ev)

		unit := world.Units[ev.UnitID]
		unit.Action = ActionRun

		var newX, newY float64
		newX = unit.X
		newY = unit.Y

		switch ev.Direction {
		case DirectionUp:
			newY -= StepSize
		case DirectionDown:
			newY += StepSize
		case DirectionRight:
			newX += StepSize
			unit.HorizontalDirection = ev.Direction
		case DirectionLeft:
			newX -= StepSize
			unit.HorizontalDirection = ev.Direction
		case DirectionVector:
			newX += StepSize * math.Cos(ev.Angle)
			newY += StepSize * math.Sin(ev.Angle)
			unit.ActionVector = &Line{
				X1: unit.X,
				Y1: unit.Y,
				X2: unit.X + 150*math.Cos(ev.Angle),
				Y2: unit.Y + 150*math.Sin(ev.Angle),
			}
			fmt.Println("ev.Angle", ev.Angle)
		}
		handleMove := true

		if unit.TriggerBox != nil {
			newBox := NewRectBox(newX, newY, unit.TriggerBox[1].X2-unit.TriggerBox[1].X1, unit.TriggerBox[0].Y2-unit.TriggerBox[0].Y1)
			for _, o := range world.Objects {
				if BoxHasIntersection(newBox, o.Box) {
					handleMove = false
				}
			}

			if handleMove {
				unit.TriggerBox = newBox
			}
		}
		fmt.Println("newX, newY", newX, newY)

		if handleMove {
			unit.X = newX
			unit.Y = newY
		}

	case PlayerEventIdle:
		str, _ := json.Marshal(event.Data)
		var ev UnitMove
		json.Unmarshal(str, &ev)

		unit := world.Units[ev.UnitID]
		unit.Action = ActionIdle
	case ActionHit:
		str, _ := json.Marshal(event.Data)
		var ev Hit
		json.Unmarshal(str, &ev)

		world.Units[ev.ToUnitID].HitPoints = world.Units[ev.ToUnitID].HitPoints - 1
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
		X:          24,
		Y:          24,
		SpriteName: skins[rnd.Intn(len(skins))],
		Action:     ActionIdle,
		Frame:      rnd.Intn(4),
		HitPoints:  10,
		TriggerBox: NewRectBox(24, 24, 16, 16),
	}
	world.Units[UnitID(id)] = unit

	return unit
}
