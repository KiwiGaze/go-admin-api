package storage

const PrefixKey = "__host"

type Messager interface {
	SetID(string)
	SetStream(string)
	SetValues(map[string]interface{})
	GetID() string
	GetStream() string
	GetValues() map[string]interface{}
	GetPrefix() string
	SetPrefix(string)
	SetErrorCount(count int)
	GetErrorCount() int
}

type ConsumerFunc func(Messager) error

type AdapterQueue interface {
	String() string
	Append(message Messager) error
	Register(name string, f ConsumerFunc)
	Run()
	Shutdown()
}
