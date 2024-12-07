package game

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

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
