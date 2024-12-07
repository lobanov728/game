package game

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
