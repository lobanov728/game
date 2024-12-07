package game

import (
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
)

type Pointable interface {
	Points() [][2]float64
	GetBox() Box
}

func NewRay(x, y, length, angle float64) Line {
	return Line{
		X1: x,
		Y1: y,
		X2: x + length*math.Cos(angle),
		Y2: y + length*math.Sin(angle),
	}
}

func LineHasIntersection(l1, l2 Line) bool {
	_, _, res := LineIntersection(l1, l2)

	return res
}

func LineIntersection(l1, l2 Line) (float64, float64, bool) {
	// https://en.wikipedia.org/wiki/Line%E2%80%93line_intersection#Given_two_points_on_each_line
	denom := (l1.X1-l1.X2)*(l2.Y1-l2.Y2) - (l1.Y1-l1.Y2)*(l2.X1-l2.X2)
	tNum := (l1.X1-l2.X1)*(l2.Y1-l2.Y2) - (l1.Y1-l2.Y1)*(l2.X1-l2.X2)
	uNum := (l1.Y1-l1.Y2)*(l1.X1-l2.X1) - (l1.X1-l1.X2)*(l1.Y1-l2.Y1)
	if denom == 0 {
		return 0, 0, false
	}

	t := tNum / denom

	if t > 1 || t < 0 {
		return 0, 0, false
	}

	u := uNum / denom
	if u > 1 || u < 0 {
		return 0, 0, false
	}
	x := l1.X1 + t*(l1.X2-l1.X1)
	y := l1.Y1 + t*(l1.Y2-l1.Y1)

	return x, y, true
}

func RayCasting(cx, cy, rayLength float64, objects []Pointable, playerBox Pointable) []Line {
	var rays []Line

	for _, obj := range objects {
		for _, p := range obj.Points() {
			l := Line{cx, cy, p[0], p[1]}
			angle := l.Angle()
			for _, angleOffset := range []float64{-0.0005, 0.0005} {
				points := [][2]float64{}
				ray := NewRay(cx, cy, rayLength, angle+angleOffset)
				for _, o := range objects {
					for _, line := range o.GetBox() {
						if px, py, ok := LineIntersection(ray, line); ok {
							points = append(points, [2]float64{px, py})
						}
					}
				}

				min := math.Inf(1)
				minI := -1
				for i, p := range points {
					d2 := (cx-p[0])*(cx-p[0]) + (cy-p[1])*(cy-p[1])
					if d2 < min {
						min = d2
						minI = i
					}
				}

				rays = append(rays, Line{cx, cy, points[minI][0], points[minI][1]})
			}
		}
	}

	for _, playerLine := range playerBox.GetBox() {
		points := [][2]float64{}

		for _, o := range objects {
			for _, line := range o.GetBox() {
				if px, py, ok := LineIntersection(playerLine, line); ok {
					points = append(points, [2]float64{px, py})
				}
			}
		}

		min := math.Inf(1)
		minI := -1
		for i, p := range points {
			d2 := (cx-p[0])*(cx-p[0]) + (cy-p[1])*(cy-p[1])
			if d2 < min {
				min = d2
				minI = i
			}
		}

		ray := Line{cx, cy, points[minI][0], points[minI][1]}
		for _, o := range objects {
			for _, line := range o.GetBox() {
				if px, py, ok := LineIntersection(ray, line); ok {
					points = append(points, [2]float64{px, py})
					ray = Line{cx, cy, px, py}
				}
			}
		}

		rays = append(rays, ray)
	}

	sort.Slice(rays, func(i, j int) bool {
		return rays[i].Angle() < rays[j].Angle()
	})

	return rays
}

func RayVertices(x1, y1, x2, y2, x3, y3 float64) []ebiten.Vertex {
	return []ebiten.Vertex{
		{
			DstX:   float32(x1),
			DstY:   float32(y1),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(x2),
			DstY:   float32(y2),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
		{
			DstX:   float32(x3),
			DstY:   float32(y3),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		},
	}
}
