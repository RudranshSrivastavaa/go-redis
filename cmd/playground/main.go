package main

import (
	"fmt"
	"log"

	"go-redis/internal/actor"
	"go-redis/internal/network"
	"go-redis/internal/pubsub"
	"go-redis/internal/storage"
)

type PrintMessage struct {
	Text string
}

type PrinterActor struct{}

func (p *PrinterActor) Receive(msg actor.Message) {

	switch m := msg.(type) {

	case PrintMessage:
		fmt.Println("PrinterActor received:", m.Text)

	default:
		fmt.Println("unknown message")
	}
}

type PingMsg struct {
	Text    string
	ReplyCh chan actor.Response
}

type PingActor struct{}

func (p *PingActor) Receive(msg actor.Message) {

	switch m := msg.(type) {

	case PingMsg:

		m.ReplyCh <- actor.Response{
			Value: "PONG: " + m.Text,
			Err:   nil,
		}
	}
}

func main() {

	// system := actor.NewActorSystem()

	// storagePID := system.Spawn(
	// 	storage.NewStorageActor(),
	// 	100,
	// )

	// replyCh := make(chan command.Response)

	// //------------------------------------------------
	// // SET
	// //------------------------------------------------

	// storagePID.Tell(
	// 	storage.CommandMsg{
	// 		Command: command.Command{
	// 			Type: command.SET,
	// 			Args: []string{
	// 				"name",
	// 				"rudransh",
	// 			},
	// 		},
	// 		ReplyCh: replyCh,
	// 	},
	// )

	// fmt.Println(<-replyCh)

	// //------------------------------------------------
	// // GET
	// //------------------------------------------------

	// storagePID.Tell(
	// 	storage.CommandMsg{
	// 		Command: command.Command{
	// 			Type: command.GET,
	// 			Args: []string{
	// 				"name",
	// 			},
	// 		},
	// 		ReplyCh: replyCh,
	// 	},
	// )
	// fmt.Println(<-replyCh)

	// //------------------------------------------------
	// // EXISTS
	// //------------------------------------------------

	// storagePID.Tell(
	// 	storage.CommandMsg{
	// 		Command: command.Command{
	// 			Type: command.EXISTS,
	// 			Args: []string{
	// 				"name",
	// 			},
	// 		},
	// 		ReplyCh: replyCh,
	// 	},
	// )

	// fmt.Println(<-replyCh)

	// //------------------------------------------------
	// // DEL
	// //------------------------------------------------

	// storagePID.Tell(
	// 	storage.CommandMsg{
	// 		Command: command.Command{
	// 			Type: command.DEL,
	// 			Args: []string{
	// 				"name",
	// 			},
	// 		},
	// 		ReplyCh: replyCh,
	// 	},
	// )

	// fmt.Println(<-replyCh)

	//}

	system := actor.NewActorSystem()

	storagePID := system.Spawn(
		storage.NewStorageActor(
			storage.Config{
				MaxKeys: 1000,
				Policy:  storage.AllKeysLRU,
			},
		),
		1000,
	)
	pubsubPID := system.Spawn(
		pubsub.NewPubSubActor(),
		1000,
	)
	server := network.NewServer(":6379")

	log.Println("server listening on :6379")

	if err := server.Start(
		system,
		storagePID,
		pubsubPID,
	); err != nil {

		log.Fatal(err)
	}
	timer := storage.NewTimerActor(storagePID)
	go timer.Start()
}
