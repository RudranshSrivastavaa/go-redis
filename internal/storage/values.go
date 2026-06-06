package storage

type Value interface {
	Type() string
}

type StringValue struct {
	Data string
}

func (s *StringValue) Type() string {
	return "string"
}

type ListValue struct {
	Data []string
}

func (l *ListValue) Type() string {
	return "list"
}

type HashValue struct {
	Data map[string]string
}

func (h *HashValue) Type() string {
	return "hash"
}

type SetValue struct {
	Data map[string]struct{}
}

func (s *SetValue) Type() string {
	return "set"
}

type ZSetValue struct {
	Data map[string]float64
}

func (z *ZSetValue) Type() string {
	return "zset"
}