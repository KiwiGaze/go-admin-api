package runtime

import (
	"sync"

	"go-admin-api/queue"
	"go-admin-api/storage"
)

type Runtime interface {
	GetMemoryQueue(prefix string) storage.AdapterQueue
	GetStreamMessage(id, stream string, value map[string]interface{}) (storage.Messager, error)
}

type Application struct {
	memoryQueue storage.AdapterQueue
	mux         sync.RWMutex
}

func NewConfig() *Application {
	return &Application{
		memoryQueue: queue.NewMemory(10000),
	}
}

func (e *Application) GetMemoryQueue(prefix string) storage.AdapterQueue {
	return NewQueue(prefix, e.memoryQueue)
}

func (e *Application) GetStreamMessage(id, stream string, value map[string]interface{}) (storage.Messager, error) {
	message := &queue.Message{}
	message.SetID(id)
	message.SetStream(stream)
	message.SetValues(value)
	return message, nil
}
