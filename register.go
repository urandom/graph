package graph

import (
	"encoding/json"
	"sync"
)

type LinkerJSONConstructor func(opts json.RawMessage) (Linker, error)

var (
	operationsMu sync.Mutex
	operations   = make(map[string]LinkerJSONConstructor)
)

func RegisterLinker(name string, constructor LinkerJSONConstructor) {
	operationsMu.Lock()
	defer operationsMu.Unlock()

	if constructor == nil {
		panic("drawgl: Register operation constructor is nil")
	}

	if _, dup := operations[name]; dup {
		panic("drawgl: Register called twice for constructor " + name)
	}

	operations[name] = constructor
}
