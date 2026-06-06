package pubsub

import "go-redis/internal/actor"

type DeliverMsg struct {
	Channel string
	Message string
}

type SubscribeMsg struct {
	Channel string
	ConnPID *actor.PID
}

type UnsubscribeMsg struct {
	Channel string
	ConnPID *actor.PID
}

type PublishMsg struct {
	Channel string
	Message string
}