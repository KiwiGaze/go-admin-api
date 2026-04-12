package runtime

import "go-admin-api/storage"

func NewQueue(prefix string, queue storage.AdapterQueue) storage.AdapterQueue {
	return &Queue{
		prefix: prefix,
		queue:  queue,
	}
}

type Queue struct {
	prefix string
	queue  storage.AdapterQueue
}

func (e *Queue) String() string { return e.queue.String() }

func (e *Queue) Register(name string, f storage.ConsumerFunc) {
	e.queue.Register(name, f)
}

func (e *Queue) Append(message storage.Messager) error {
	values := message.GetValues()
	if values == nil {
		values = make(map[string]interface{})
	}
	values[storage.PrefixKey] = e.prefix
	return e.queue.Append(message)
}

func (e *Queue) Run() { e.queue.Run() }

func (e *Queue) Shutdown() {
	if e.queue != nil {
		e.queue.Shutdown()
	}
}
