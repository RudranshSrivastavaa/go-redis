package storage

import (
	"math/rand"
	"sort"
)

type LRUCandidate struct {
	Key        string
	LastAccess int64
}

const (
	LRUSamples  = 5
	LRUPoolSize = 16
)

func (s *StorageActor) sampleKeys(count int) []string {

	keys := make([]string, 0, len(s.data))

	for key := range s.data {
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		return nil
	}

	if count > len(keys) {
		count = len(keys)
	}

	result := make([]string, 0, count)

	for i := 0; i < count; i++ {

		idx := rand.Intn(len(keys))

		result = append(
			result,
			keys[idx],
		)
	}

	return result
}

func (s *StorageActor) refreshLRUPool() {

	samples := s.sampleKeys(
		LRUSamples,
	)
	for _, key := range samples {

		entry, exists := s.getEntry(key)

		if !exists {
			continue
		}

		if s.candidateExists(key) {
		continue
	    }

		s.lruPool = append(
			s.lruPool,
			LRUCandidate{
				Key:        key,
				LastAccess: entry.LastAccess,
			},
		)
	}
	sort.Slice(
		s.lruPool,
		func(i, j int) bool {

			return s.lruPool[i].LastAccess <
				s.lruPool[j].LastAccess
		},
	)
	if len(s.lruPool) > LRUPoolSize {

		s.lruPool = s.lruPool[:LRUPoolSize]
	}
}

func (s *StorageActor) prunePool() {

	filtered := s.lruPool[:0]

	for _, candidate := range s.lruPool {

		if _, exists := s.data[candidate.Key]; exists {

			filtered = append(
				filtered,
				candidate,
			)
		}
	}

	s.lruPool = filtered
}

func (s *StorageActor) evictLRU() {
	s.prunePool()
	s.refreshLRUPool()
	if len(s.lruPool) == 0 {
		return
	}
	victim := s.lruPool[0]
	if _, exists := s.data[victim.Key]; !exists {

		s.lruPool = s.lruPool[1:]

		return
	}
	delete(
		s.data,
		victim.Key,
	)
	s.lruPool = s.lruPool[1:]
}

func (s *StorageActor) evictIfNeeded() {

	for len(s.data) > s.config.MaxKeys {

		switch s.config.Policy {

		case AllKeysLRU:
			s.evictLRU()

		case NoEviction:
			return
		}
	}
}

func (s *StorageActor) candidateExists(
	key string,
) bool {

	for _, c := range s.lruPool {

		if c.Key == key {
			return true
		}
	}

	return false
}
