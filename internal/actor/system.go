package actor

type PID struct{
	mailbox *Mailbox
}

func (p *PID) Tell(msg Message) {
	p.mailbox.Send(msg)
}

type actorCell struct{
	actor   Actor
	mailbox *Mailbox
}

type ActorSystem struct{}

func NewActorSystem() *ActorSystem {
	return &ActorSystem{}
}

func(s *ActorSystem) Spawn(a Actor , mailboxSize int)*PID{
	cell := &actorCell{
		actor:   a,
		mailbox: NewMailbox(mailboxSize),
	}
	go func() {
		for msg := range cell.mailbox.Channel() {
			cell.actor.Receive(msg)
		}
	}()
	
	return &PID{
		mailbox: cell.mailbox,
	}
}