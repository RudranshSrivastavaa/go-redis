package pubsub

import ("go-redis/internal/actor"
"fmt"
)

type PubSubActor struct {
	subscribers map[string][]*actor.PID
}

func NewPubSubActor() *PubSubActor {
	return &PubSubActor{
		subscribers: make(map[string][]*actor.PID),
	}
}

func (p *PubSubActor) Receive(msg actor.Message) {
	switch m := msg.(type) {

	case SubscribeMsg:

		p.subscribers[m.Channel] =
			append(
				p.subscribers[m.Channel],
				m.ConnPID,
			)

	case PublishMsg:

    subs := p.subscribers[m.Channel]

    fmt.Println(
        "subscribers:",
        len(subs),
    )

    for _, pid := range subs {

        fmt.Println("sending")

        pid.Tell(
            DeliverMsg{
                Channel: m.Channel,
                Message: m.Message,
            },
        )
    }
}
}