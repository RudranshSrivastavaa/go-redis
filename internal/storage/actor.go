package storage

import (
	"fmt"
	"time"
	"go-redis/internal/actor"
	"go-redis/internal/command"
)


type StorageActor struct {
    data map[string]Entry

    config Config

    lruPool []LRUCandidate
}


func NewStorageActor(cfg Config) *StorageActor {
	return &StorageActor{
		data:    make(map[string]Entry),
		config:  cfg,
		lruPool: make([]LRUCandidate, 0, LRUPoolSize),
	}
}

func (s *StorageActor) touch(key string,entry Entry) {
	entry.LastAccess = time.Now().UnixNano()
	s.data[key] = entry
}

func (s *StorageActor) Receive(msg actor.Message) {

	// switch m := msg.(type) {

	// case SetMsg:

	// now := time.Now()

	// s.data[m.Key] = Entry{
	// 	Value: &StringValue{
	// 		Data: m.Value,
	// 	},
	// 	CreatedAt: now,
	// 	UpdatedAt: now,
	// }

	// m.ReplyCh <- actor.Response{
	// 	Value: "OK",
	// }


	// case GetMsg:

	// entry, exists := s.data[m.Key]

	// if !exists {
	// 	m.ReplyCh <- actor.Response{
	// 		Err: fmt.Errorf("key not found"),
	// 	}
	// 	return
	// }

	// strValue, ok := entry.Value.(*StringValue)

	// if !ok {
	// 	m.ReplyCh <- actor.Response{
	// 		Err: fmt.Errorf("wrong type"),
	// 	}
	// 	return
	// }

	// m.ReplyCh <- actor.Response{
	// 	Value: strValue.Data,
	// }
	
	// case DelMsg:

	// 	_, exists := s.data[m.Key]

	// 	if exists {
	// 		delete(s.data, m.Key)

	// 		m.ReplyCh <- actor.Response{
	// 			Value: 1,
	// 		}
	// 	} else {
	// 		m.ReplyCh <- actor.Response{
	// 			Value: 0,
	// 		}
	// 	}

	// case ExistsMsg:

	// 	_, exists := s.data[m.Key]

	// 	if exists {
	// 		m.ReplyCh <- actor.Response{
	// 			Value: 1,
	// 		}
	// 	} else {
	// 		m.ReplyCh <- actor.Response{
	// 			Value: 0,
	// 		}
	// 	}
	// }

	switch m := msg.(type) {

    case CommandMsg:

    s.executeCommand(
        m.Command,
        m.ReplyCh,
    )
	
	case ExpireSweepMsg:

	s.sweepExpired()
}
}

func (s *StorageActor) executeCommand(
	cmd command.Command,
	replyCh chan command.Response,
) {

	switch cmd.Type {

	case command.SET:
		s.handleSet(cmd, replyCh)

	case command.GET:
		s.handleGet(cmd, replyCh)

	case command.DEL:
		s.handleDel(cmd, replyCh)

	case command.EXISTS:
		s.handleExists(cmd, replyCh)
	
	case command.LPUSH:
	    s.handleLPush(cmd, replyCh)

    case command.RPUSH:
	    s.handleRPush(cmd, replyCh)

    case command.LPOP:
	    s.handleLPop(cmd, replyCh)

    case command.RPOP:
	    s.handleRPop(cmd, replyCh)

	case command.LLEN:
	    s.handleLLen(cmd, replyCh)

	case command.HSET:
	    s.handleHSet(cmd, replyCh)

    case command.HGET:
	    s.handleHGet(cmd, replyCh)

    case command.HDEL:
	    s.handleHDel(cmd, replyCh)

	case command.HEXISTS:
	    s.handleHExists(cmd, replyCh)

    case command.HLEN:
	    s.handleHLen(cmd, replyCh)

	case command.EXPIRE:
	    s.handleExpire(cmd, replyCh)

    case command.TTL:
	    s.handleTTL(cmd, replyCh)

    case command.PERSIST:
	    s.handlePersist(cmd, replyCh)

	default:

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"unknown command: %s",
				cmd.Type,
			),
		}
	}
}

func (s *StorageActor) handleSet(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 2 {

		replyCh <- command.Response{
			Err: fmt.Errorf("SET requires 2 args"),
		}

		return
	}

	key := cmd.Args[0]
	value := cmd.Args[1]

	now := time.Now()

    s.data[key] = Entry{
    Value: &StringValue{
        Data: value,
    },
    CreatedAt: now,
    UpdatedAt: now,
    LastAccess: now.UnixNano(),
}
	s.evictIfNeeded()
	replyCh <- command.Response{
		Value: "OK",
	}
}

func (s *StorageActor) handleGet(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf("GET requires 1 arg"),
		}

		return
	}

	key := cmd.Args[0]

	entry, exists := s.getEntry(key)

	s.touch(key,entry);
	if !exists {

		replyCh <- command.Response{
			Err: fmt.Errorf("key not found"),
		}

		return
	}

	str, ok := entry.Value.(*StringValue)

	if !ok {

		replyCh <- command.Response{
			Err: fmt.Errorf("wrong type"),
		}

		return
	}

	replyCh <- command.Response{
		Value: str.Data,
	}
}


func (s *StorageActor) handleDel(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf("DEL requires 1 arg"),
		}

		return
	}

	key := cmd.Args[0]

	_, exists := s.data[key]

	if !exists {

		replyCh <- command.Response{
			Value: 0, //Here our redis wil return number of keys removed
		}

		return
	}

	delete(s.data, key)

	replyCh <- command.Response{
		Value: 1,
	}
}

func (s *StorageActor) handleExists(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf("EXISTS requires 1 arg"),
		}

		return
	}

	key := cmd.Args[0]

	_, exists := s.getEntry(key)

	if exists {

		replyCh <- command.Response{
			Value: 1,
		}

		return
	}

	replyCh <- command.Response{
		Value: 0,
	}
}


func (s *StorageActor) handleLPush(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 2 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"LPUSH requires key and value",
			),
		}
		return
	}

	key := cmd.Args[0]
	value := cmd.Args[1]

	list, err := s.getOrCreateList(key)

	if err != nil {
		replyCh <- command.Response{
			Err: err,
		}
		return
	}

	list.Data = append(
		[]string{value},
		list.Data...,
	)
	s.evictIfNeeded()
	replyCh <- command.Response{
		Value: len(list.Data),
	}
}

func (s *StorageActor) handleRPush(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 2 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"RPUSH requires key and value",
			),
		}
		return
	}

	key := cmd.Args[0]
	value := cmd.Args[1]

	list, err := s.getOrCreateList(key)

	if err != nil {
		replyCh <- command.Response{
			Err: err,
		}
		return
	}

	list.Data = append(
		list.Data,
		value,
	)
	s.evictIfNeeded()
	replyCh <- command.Response{
		Value: len(list.Data),
	}
}

func (s *StorageActor) handleLLen(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"LLEN requires key",
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
	s.touch(key, entry)

	list, ok := entry.Value.(*ListValue)

	if !ok {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"WRONGTYPE operation against a key holding the wrong kind of value",
			),
		}

		return
	}

	replyCh <- command.Response{
		Value: len(list.Data),
	}
}

func (s *StorageActor) handleLPop(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"LPOP requires key",
			),
		}
		return
	}

	key := cmd.Args[0]

	entry, exists := s.getEntry(key)

	if !exists {
		replyCh <- command.Response{
			Value: nil,
		}
		return
	}
	s.touch(key, entry)
	list, ok := entry.Value.(*ListValue)

	if !ok {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"WRONGTYPE operation against a key holding the wrong kind of value",
			),
		}
		return
	}

	if len(list.Data) == 0 {

		replyCh <- command.Response{
			Value: nil,
		}
		return
	}

	value := list.Data[0]

	list.Data = list.Data[1:]

	now := time.Now()

    entry.LastAccess = now.UnixNano()
    entry.UpdatedAt = now

    s.data[key] = entry

	replyCh <- command.Response{
		Value: value,
	}
}

func (s *StorageActor) handleRPop(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"RPOP requires key",
			),
		}
		return
	}

	key := cmd.Args[0]

	entry, exists := s.getEntry(key)
	

	if !exists {

		replyCh <- command.Response{
			Value: nil,
		}
		return
	}

	s.touch(key, entry)
	list, ok := entry.Value.(*ListValue)

	if !ok {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"WRONGTYPE operation against a key holding the wrong kind of value",
			),
		}
		return
	}

	if len(list.Data) == 0 {

		replyCh <- command.Response{
			Value: nil,
		}
		return
	}

	last := len(list.Data) - 1

	value := list.Data[last]

	list.Data = list.Data[:last]
	
	now := time.Now()

    entry.LastAccess = now.UnixNano()
    entry.UpdatedAt = now

    s.data[key] = entry

	replyCh <- command.Response{
		Value: value,
	}
}

func (s *StorageActor) handleHSet(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 3 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"HSET requires key field value",
			),
		}
		return
	}

	key := cmd.Args[0]
	field := cmd.Args[1]
	value := cmd.Args[2]

	hash, err := s.getOrCreateHash(key)

	if err != nil {

		replyCh <- command.Response{
			Err: err,
		}

		return
	}

	_, exists := hash.Data[field]

	hash.Data[field] = value

	s.evictIfNeeded()

	if exists {

		replyCh <- command.Response{
			Value: 0,
		}

		return
	}

	replyCh <- command.Response{
		Value: 1,
	}
}

func (s *StorageActor) handleHGet(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 2 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"HGET requires key field",
			),
		}
		return
	}

	key := cmd.Args[0]
	field := cmd.Args[1]

	entry, exists := s.getEntry(key)


	if !exists {

		replyCh <- command.Response{
			Value: nil,
		}
		return
	}
	s.touch(key, entry)

	hash, ok := entry.Value.(*HashValue)

	if !ok {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"WRONGTYPE operation against a key holding the wrong kind of value",
			),
		}

		return
	}

	value, exists := hash.Data[field]

	if !exists {

		replyCh <- command.Response{
			Value: nil,
		}

		return
	}

	replyCh <- command.Response{
		Value: value,
	}
}

func (s *StorageActor) handleHDel(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 2 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"HDEL requires key field",
			),
		}

		return
	}

	key := cmd.Args[0]
	field := cmd.Args[1]

	entry, exists := s.getEntry(key)

	if !exists {

		replyCh <- command.Response{
			Value: 0,
		}

		return
	}

	hash, ok := entry.Value.(*HashValue)

	if !ok {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"WRONGTYPE operation against a key holding the wrong kind of value",
			),
		}

		return
	}

	_, exists = hash.Data[field]

	if !exists {

		replyCh <- command.Response{
			Value: 0,
		}

		return
	}

	delete(hash.Data, field)

	replyCh <- command.Response{
		Value: 1,
	}
}

func (s *StorageActor) handleHExists(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 2 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"HEXISTS requires key field",
			),
		}
		return
	}

	key := cmd.Args[0]
	field := cmd.Args[1]

	entry, exists := s.getEntry(key)


	if !exists {

		replyCh <- command.Response{
			Value: 0,
		}
		return
	}
	s.touch(key, entry)

	hash, ok := entry.Value.(*HashValue)

	if !ok {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"WRONGTYPE operation against a key holding the wrong kind of value",
			),
		}
		return
	}

	_, exists = hash.Data[field]

	if exists {

		replyCh <- command.Response{
			Value: 1,
		}
	} else {

		replyCh <- command.Response{
			Value: 0,
		}
	}
}

func (s *StorageActor) handleHLen(cmd command.Command,replyCh chan command.Response) {

	if len(cmd.Args) != 1 {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"HLEN requires key",
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

	s.touch(key, entry)

	hash, ok := entry.Value.(*HashValue)

	if !ok {

		replyCh <- command.Response{
			Err: fmt.Errorf(
				"WRONGTYPE operation against a key holding the wrong kind of value",
			),
		}
		return
	}

	replyCh <- command.Response{
		Value: len(hash.Data),
	}
}