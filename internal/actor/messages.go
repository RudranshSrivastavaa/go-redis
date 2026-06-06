package actor

type Message interface{}

type Response struct {
	Value any
	Err   error
}