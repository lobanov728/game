package game

import (
	"encoding/json"
	"fmt"
	"math"
)

func (world *World) HandleEvent(event *Event) {
	fmt.Println("event.Type", event.Type)
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
			newY -= PlayerStepSize
		case DirectionDown:
			newY += PlayerStepSize
		case DirectionRight:
			newX += PlayerStepSize
			unit.HorizontalDirection = ev.Direction
		case DirectionLeft:
			newX -= PlayerStepSize
			unit.HorizontalDirection = ev.Direction
		case DirectionVector:
			newX += MobStepSize * math.Cos(ev.Angle)
			newY += MobStepSize * math.Sin(ev.Angle)
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
