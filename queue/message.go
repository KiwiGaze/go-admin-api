package queue

import (
	"sync"

	"go-admin-api/storage"
)

type Message struct {
	ID         string
	Stream     string
	Values     map[string]interface{}
	ErrorCount int
	mux        sync.RWMutex
}

func (m *Message) SetID(id string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.ID = id
}

func (m *Message) SetStream(stream string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.Stream = stream
}

func (m *Message) SetValues(values map[string]interface{}) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.Values = values
}

func (m *Message) GetID() string {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.ID
}

func (m *Message) GetStream() string {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.Stream
}

func (m *Message) GetValues() map[string]interface{} {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.Values
}

func (m *Message) GetPrefix() string {
	m.mux.RLock()
	defer m.mux.RUnlock()
	if m.Values == nil {
		return ""
	}
	prefix, ok := m.Values[storage.PrefixKey].(string)
	if !ok {
		return ""
	}
	return prefix
}

func (m *Message) SetPrefix(prefix string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	if m.Values == nil {
		m.Values = make(map[string]interface{})
	}
	m.Values[storage.PrefixKey] = prefix
}

func (m *Message) SetErrorCount(count int) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.ErrorCount = count
}

func (m *Message) GetErrorCount() int {
	m.mux.RLock()
	defer m.mux.RUnlock()
	return m.ErrorCount
}
