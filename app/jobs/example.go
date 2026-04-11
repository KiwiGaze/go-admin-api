package jobs

import (
	"fmt"
	"time"
)

// InitJob
// Add the defined structs to the map.
// The map key can be configured as the scheduled task invoke target.
func InitJob() {
	jobList = map[string]JobExec{
		"ExamplesOne": ExamplesOne{},
		// ...
	}
}

// ExamplesOne
// Newly added jobs must follow this format and implement the Exec function.
type ExamplesOne struct {
}

func (t ExamplesOne) Exec(arg interface{}) error {
	str := time.Now().Format(timeFormat) + " [INFO] JobCore ExamplesOne exec success"
	// TODO: Examples receives a string argument, so use arg.(string). Convert according to the corresponding type.
	switch arg.(type) {

	case string:
		if arg.(string) != "" {
			fmt.Println("string", arg.(string))
			fmt.Println(str, arg.(string))
		} else {
			fmt.Println("arg is nil")
			fmt.Println(str, "arg is nil")
		}
		break
	}

	return nil
}
