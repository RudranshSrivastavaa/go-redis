package storage

import "time"

func (s *StorageActor) isExpired(entry Entry) bool {
	if entry.ExpiresAt == nil {
		return false
	}

	return time.Now().After(
		*entry.ExpiresAt,
	)
}

func (s *StorageActor) getEntry(key string) (Entry, bool) {

	entry, exists := s.data[key]

	if !exists {
		return Entry{}, false
	}

	if s.isExpired(entry) {

		delete(s.data, key)

		return Entry{}, false
	}

	return entry, true
}

func (s *StorageActor) sweepExpired() {

	now := time.Now()

	for key, entry := range s.data {

		if entry.ExpiresAt == nil {
			continue
		}

		if now.After(*entry.ExpiresAt) {
			delete(s.data, key)
		}
	}
}