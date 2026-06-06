package storage

import (
	"time"

	"go-redis/internal/actor"
)

type TimerActor struct {
	storagePID *actor.PID
}

func NewTimerActor(storagePID *actor.PID) *TimerActor {
	return &TimerActor{
		storagePID: storagePID,
	}
}

func (t *TimerActor) Start() {

	ticker := time.NewTicker(
		1 * time.Second,
	)

	for range ticker.C {

		t.storagePID.Tell(
			ExpireSweepMsg{},
		)
	}
}