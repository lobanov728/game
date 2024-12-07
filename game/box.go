package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Box []Line

func BoxHasIntersection(b1, b2 Box) bool {
	for _, l1 := range b1 {
		for _, l2 := range b2 {
			if LineHasIntersection(l1, l2) {
				return true
			}
		}
	}

	return false
}

func GetBoxIntersection(b1, b2 Box) int {
	for number, l1 := range b1 {
		for _, l2 := range b2 {
			if LineHasIntersection(l1, l2) {
				return number
			}
		}
	}

	return -1
}

func NewRectBox(x, y, w, h float64) Box {
	return []Line{
		{x, y, x, y + h},
		{x, y + h, x + w, y + h},
		{x + w, y + h, x + w, y},
		{x + w, y, x, y},
	}
}

func NewCircleBox(x, y, r float64) Box {
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
