package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type memoryDB struct {
	items map[string]any
	mu    sync.RWMutex
}

func newMemoryDB() memoryDB {
	f, err := os.Open("db.json")
	if err != nil {
		return memoryDB{items: map[string]any{}}
	}

	items := map[string]any{}
	if err := json.NewDecoder(f).Decode(&items); err != nil {
		fmt.Println("could not decode db.json:", err.Error())
	}
	return memoryDB{items: items}
}

func (m *memoryDB) set(key, value string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.items[key] = value
}

func (m *memoryDB) get(key string) (any, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	value, found := m.items[key]
	return value, found
}

func (m *memoryDB) keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, 0, len(m.items))
	for k := range m.items {
		keys = append(keys, k)
	}
	return keys
}

func (m *memoryDB) delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.items, key)
}

func (m *memoryDB) save() {
	f, err := os.Create("db.json")
	if err != nil {
		fmt.Println("could not create file", err.Error())
	}
	if err := json.NewEncoder(f).Encode(m.items); err != nil {
		fmt.Println("could not encode", err.Error())
	}
	fmt.Println("successfully saved db to file")
}
