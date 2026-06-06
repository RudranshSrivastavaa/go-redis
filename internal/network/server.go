package network

import (
	"net"

	"go-redis/internal/actor"
)

type Server struct {
	address string
}

func NewServer(address string) *Server {

	return &Server{
		address: address,
	}
}

func (s *Server) Start(system *actor.ActorSystem,storagePID *actor.PID,pubsubPID *actor.PID) error {

	listener, err := net.Listen(
		"tcp",
		s.address,
	)

	if err != nil {
		return err
	}

	defer listener.Close()

	for {

	conn, err := listener.Accept()

	if err != nil {
		continue
	}

	connActor := NewConnActor(conn)

	connPID := system.Spawn(
		connActor,
		100,
	)

	go connActor.Start(
		storagePID,
		pubsubPID,
		connPID,
	)
}
}