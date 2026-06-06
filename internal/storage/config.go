package storage

type EvictionPolicy string

const (
	NoEviction EvictionPolicy = "noeviction"
	AllKeysLRU EvictionPolicy = "allkeys-lru"
)

type Config struct {
	MaxKeys int
	Policy  EvictionPolicy
}


