package storage

import "go-redis/internal/command"

type CommandMsg struct {
	Command command.Command
	ReplyCh chan command.Response
}

type ExpireSweepMsg struct{}