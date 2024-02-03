package game

import (
	"fmt"
	"testing"
	"time"
)

func TestActionCooldown(t *testing.T) {

	hitAction := NewAction(ActionHit, time.Millisecond*500)
	moveAction := NewAction(PlayerEventMove, time.Millisecond*70)
	act := []*Action{hitAction, moveAction}

	readyForAction := make(chan *Action)

	for _, action := range act {
		action := action
		go func() {
			for {
				readyForAction <- action.Subcribe()
			}
		}()
	}

	go func() {
		for {
			readyAction := <-readyForAction
			readyAction.Fire()
			fmt.Println(readyAction.GetName())
		}
	}()

	time.Sleep(time.Second * 2)

	t.Fail()

	// for i := 0; i < 10; i++ {
	// 	select {
	// 	case <-hitAction.Subcribe():
	// 		hitAction.Fire()
	// 	case <-moveAction.Subcribe():
	// 		moveAction.Fire()
	// 	default:
	// 		time.Sleep(time.Millisecond * 200)
	// 	}

	// }

	// for i := 0; i < 10; i++ {
	// 	select {
	// 	case <-hitAction.Subcribe():
	// 		hitAction.Fire()
	// 	case <-moveAction.Subcribe():
	// 		moveAction.Fire()
	// 	default:
	// 		time.Sleep(time.Millisecond * 200)
	// 	}

	// }
}
