package game

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
