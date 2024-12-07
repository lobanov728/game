package game

import (
	"math"
)

type Unit struct {
	ID                  UnitID    `json:"id"`
	X                   float64   `json:"x"`
	Y                   float64   `json:"y"`
	SpriteName          string    `json:"sprite_name"`
	Action              EventName `json:"action"`
	Frame               int       `json:"frame"`
	HorizontalDirection int       `json:"direction"`
	ActionVector        *Line     `json:"action_vector"`
	Box                 Box       `json:"box"`
	TriggerBox          Box       `json:"trigger_box"`
	HitPoints           int       `json:"hit_points"`
	Actions             map[EventName]*Action
}

func (u *Unit) GetBox() Box {
	return u.Box
}

func (u *Unit) CouldMakeAction(name EventName) bool {
	return false
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

type Line struct {
	X1 float64 `json:"x1"`
	Y1 float64 `json:"y1"`
	X2 float64 `json:"x2"`
	Y2 float64 `json:"y2"`
}

func (l *Line) Angle() float64 {
	return math.Atan2(l.Y2-l.Y1, l.X2-l.X1)
}

func (l *Line) Length() float64 {
	return math.Sqrt(math.Pow(l.X2-l.X1, 2) + math.Pow(l.Y2-l.Y1, 2))
}

type Tile struct {
	ID         UnitID    `json:"id"`
	X          float64   `json:"x"`
	Y          float64   `json:"y"`
	SpriteName string    `json:"sprite_name"`
	Action     EventName `json:"action"`
	Frame      int       `json:"frame"`
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
