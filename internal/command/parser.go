package command

import (
	"fmt"
	"strings"
)

const (
	PING CommandType = "PING"
	ECHO CommandType = "ECHO"
)

func FromRESP(parts []string) (Command, error) {

	if len(parts) == 0 {
		return Command{}, fmt.Errorf("empty command")
	}

	return Command{
		Type: CommandType(
			strings.ToUpper(parts[0]),
		),
		Args: parts[1:],
	}, nil
}