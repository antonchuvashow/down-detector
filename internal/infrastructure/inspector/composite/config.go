package composite

import "detector/internal/inspection/domain/inspector"

type CompositeInspectorConfig struct {
	Inspectors map[string]inspector.Inspector
}
