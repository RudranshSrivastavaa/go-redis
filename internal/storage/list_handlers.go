package storage

import "fmt"

func (s *StorageActor) getOrCreateList(key string) (*ListValue, error) {

	entry, exists := s.getEntry(key)

	if !exists {

		list := &ListValue{
			Data: []string{},
		}

		s.data[key] = Entry{
			Value: list,
		}

		return list, nil
	}

	list, ok := entry.Value.(*ListValue)

	if !ok {
		return nil, fmt.Errorf(
			"WRONGTYPE operation against a key holding the wrong kind of value",
		)
	}

	return list, nil
}