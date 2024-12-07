package game

const StepSize = 1.5
const GameSpeedMs = 100

const (
	DirectionUp = iota + 1
	DirectionRight
	DirectionDown
	DirectionLeft
	DirectionVector
)

const (
	PlayerEventConnect EventName = "connect"
	PlayerEventMove    EventName = "move"
	PlayerEventIdle    EventName = "idle"
	PlayerEventInit    EventName = "init"

	ActionRun  EventName = "run"
	ActionHit  EventName = "hit"
	ActionIdle EventName = "idle"
)
