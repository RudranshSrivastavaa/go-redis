package network

import (
	"fmt"
	"go-redis/internal/actor"
	"go-redis/internal/command"
	"go-redis/internal/pubsub"
	"go-redis/internal/resp"
	"go-redis/internal/storage"
	"net"
)

type ConnActor struct {
	conn   net.Conn
	writer *resp.Writer
}

func NewConnActor(conn net.Conn) *ConnActor {

	return &ConnActor{
		conn:   conn,
		writer: resp.NewWriter(conn),
	}
}

// Required to satisfy Actor interface

func (c *ConnActor) Receive(msg actor.Message) {

	fmt.Println("receive called")

	switch m := msg.(type) {

	case pubsub.DeliverMsg:

		fmt.Println(
			"delivering:",
			m.Message,
		)
		c.writer.WriteArray(
			[]string{
				"message",
				m.Channel,
				m.Message,
			},
		)
	}
}

func (c *ConnActor) Start(storagePID *actor.PID, pubsubPID *actor.PID, connPID *actor.PID) {
	defer c.conn.Close()

	parser := resp.NewParser(c.conn)
	c.writer = resp.NewWriter(c.conn)

	for {
		// Parse RESP
		parts, err := parser.Parse()
		if err != nil {
			return
		}

		// Convert to Command

		cmd, err := command.FromRESP(parts)
		if err != nil {

			c.writer.WriteError(
				err.Error(),
			)

			continue
		}

		// Local Commands

		switch cmd.Type {

		case command.PING:

			c.writer.WriteSimpleString("PONG")

			continue

		case command.ECHO:

			if len(cmd.Args) != 1 {

				c.writer.WriteError("ECHO requires argument")

				continue
			}

			c.writer.WriteBulkString(cmd.Args[0])

			continue

		case command.SUBSCRIBE:

			if len(cmd.Args) != 1 {

				c.writer.WriteError("SUBSCRIBE requires channel")

				continue
			}

			pubsubPID.Tell(
				pubsub.SubscribeMsg{
					Channel: cmd.Args[0],
					ConnPID: connPID,
				},
			)

			c.writer.WriteSimpleString(
				"SUBSCRIBED",
			)

			continue

		case command.PUBLISH:

			if len(cmd.Args) != 2 {

				c.writer.WriteError(
					"PUBLISH requires channel message",
				)

				continue
			}

			pubsubPID.Tell(
				pubsub.PublishMsg{
					Channel: cmd.Args[0],
					Message: cmd.Args[1],
				},
			)

			c.writer.WriteInteger(1)

			continue
		}

		// Storage Commands

		replyCh := make(
			chan command.Response,
			1,
		)

		storagePID.Tell(
			storage.CommandMsg{
				Command: cmd,
				ReplyCh: replyCh,
			},
		)

		response := <-replyCh

		// Error

		if response.Err != nil {

			c.writer.WriteError(response.Err.Error())

			continue
		}

		// Serialize Response

		switch v := response.Value.(type) {

		case nil:
			c.writer.WriteNullBulkString()
		case string:
			if v == "OK" {
				c.writer.WriteSimpleString(v)

			} else {
				c.writer.WriteBulkString(v)
			}

		case int:
			c.writer.WriteInteger(int64(v))

		case int64:
			c.writer.WriteInteger(v)

		default:
			c.writer.WriteError("unsupported response type")
		}
	}
}
