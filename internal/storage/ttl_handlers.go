package storage

import(
	"fmt"
	"go-redis/internal/command"
	"strconv"
	"time"
)

func (s *StorageActor) handleExpire(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 2 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"EXPIRE requires key seconds",
			),
		}

		return
	}

	key := cmd.Args[0]

	seconds, err := strconv.Atoi(
		cmd.Args[1],
	)

	if err != nil {

		replyCh <- command.Response{
			Err: err,
		}

		return
	}

	entry, exists := s.getEntry(key)

	if !exists {

		replyCh <- command.Response{
			Value: 0,
		}

		return
	}

	expiry := time.Now().
		Add(time.Duration(seconds) * time.Second)

	entry.ExpiresAt = &expiry

	s.data[key] = entry

	replyCh <- command.Response{
		Value: 1,
	}
}

func (s *StorageActor) handleTTL(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"TTL requires key",
			),
		}

		return
	}

	key := cmd.Args[0]

	entry, exists := s.getEntry(key)

	if !exists {

		replyCh <- command.Response{
			Value: -2,
		}

		return
	}

	if entry.ExpiresAt == nil {

		replyCh <- command.Response{
			Value: -1,
		}

		return
	}

	remaining := int(
		time.Until(*entry.ExpiresAt).
			Seconds(),
	)

	replyCh <- command.Response{
		Value: remaining,
	}
}

func (s *StorageActor) handlePersist(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"PERSIST requires key",
			),
		}

		return
	}

	key := cmd.Args[0]

	entry, exists := s.getEntry(key)

	if !exists {

		replyCh <- command.Response{
			Value: 0,
		}

		return
	}

	if entry.ExpiresAt == nil {

		replyCh <- command.Response{
			Value: 0,
		}

		return
	}

	entry.ExpiresAt = nil

	s.data[key] = entry

	replyCh <- command.Response{
		Value: 1,
	}
}