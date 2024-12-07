package game

import (
	"fmt"
	"time"
)

type Action struct {
	name      EventName
	cooldown  time.Duration
	fired     bool
	firedTime time.Time

	timer *time.Timer
	reset chan struct{}
}

func NewAction(n EventName, c time.Duration) *Action {
	timer := time.NewTimer(0)
	<-timer.C

	a := &Action{
		name:     n,
		cooldown: c,
		timer:    timer,
		reset:    make(chan struct{}),
	}

	go func() {
		// a.reset <- struct{}{}
	}()

	return a
}

func (a *Action) IsReady() time.Duration {
	if !a.fired {
		return 0
	}

	return time.Since(a.firedTime)
}

func (a *Action) GetName() EventName {
	return a.name
}

func (a *Action) Fire() {
	fmt.Println("Fire", a.name)
	a.fired = true
	a.firedTime = time.Now()
	a.timer.Reset(a.cooldown)

	go func() {
		<-a.timer.C
		// a.reset <- struct{}{}
		a.fired = false
	}()
}

func (a *Action) Subcribe() *Action {
	<-a.reset
	return a
}
