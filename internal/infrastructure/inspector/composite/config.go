package composite

import "detector/internal/inspection/domain/inspector"

type InspectorConfig struct {
	Inspectors map[string]inspector.Inspector
}
