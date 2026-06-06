package storage

import "fmt"

func (s *StorageActor) getOrCreateHash(key string) (*HashValue, error) {

	entry, exists := s.getEntry(key)

	if !exists {

		hash := &HashValue{
			Data: make(map[string]string),
		}

		s.data[key] = Entry{
			Value: hash,
		}

		return hash, nil
	}

	hash, ok := entry.Value.(*HashValue)

	if !ok {

		return nil, fmt.Errorf(
			"WRONGTYPE operation against a key holding the wrong kind of value",
		)
	}

	return hash, nil
}