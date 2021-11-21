package caches

type Status struct {
	Count     int64 `json:"count"`
	KeySize   int64 `json:"KeySize"`
	ValueSize int64 `json:"ValueSize"`
}

func NewStauts() *Status {
	return &Status{
		Count:     0, // record cache nums
		KeySize:   0, // record keys size
		ValueSize: 0, // record values size
	}
}

func (s *Status) addEntry(key string, value []byte) {
	s.Count++
	s.KeySize += int64(len(key))
	s.ValueSize += int64(len(value))
}

func (s *Status) subEntry(key string, value []byte) {
	s.Count--
	s.KeySize -= int64(len(key))
	s.ValueSize -= int64(len(value))
}
func (s *Status) entrySize() int64 {
	return s.ValueSize + s.KeySize
}
