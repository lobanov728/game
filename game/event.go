package game

type Event struct {
	ID   string      `json:"id"`
	Type EventName   `json:"type"`
	Data interface{} `json:"data"`
}

type EventName string

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
