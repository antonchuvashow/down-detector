package apiinspector

import "fmt"

type ErrInspectorUnknown struct {
	Type string
}

func (e ErrInspectorUnknown) Error() string {
	return fmt.Sprintf("unknown Inspector type: %s", e.Type)
}
