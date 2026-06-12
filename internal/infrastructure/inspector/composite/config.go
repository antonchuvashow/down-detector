package composite

import (
	"detector/internal/inspector/domain"
)

type ExtraInspectionInfo struct {
	Results map[string]inspector.Result
}

type InspectorConfig struct {
	Inspectors map[string]inspector.Inspector
}
