package actor

type Actor interface{
	Receive(msg Message)
}