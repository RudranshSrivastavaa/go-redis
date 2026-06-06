package actor

type Mailbox struct{
	ch chan Message
}

func NewMailbox(buffer int) *Mailbox {
	return &Mailbox{
		ch: make(chan Message, buffer),
	}
}

func(m *Mailbox) Send(msg Message){
	m.ch <- msg;
}

func (m *Mailbox) Channel() <-chan Message {
	return m.ch
}