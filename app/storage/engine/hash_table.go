package engine

import (
	"errors"
	"sync"
)

type HashTable struct {
	mx   sync.RWMutex
	data map[string]string
}

func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
	}
}

func (s *HashTable) Set(key, value string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.data[key] = value
}

func (s *HashTable) Get(key string) string {
	s.mx.RLock()
	defer s.mx.RUnlock()
	if data, ok := s.data[key]; ok {
		return data
	} else {
		return errors.New("not found").Error()
	}
}

func (s *HashTable) Delete(key string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	delete(s.data, key)
}
