package composite

import "detector/internal/inspection/domain/inspector"

type ExtraInspectionInfo struct {
	Results map[string]inspector.InspectionResult
}

type InspectorConfig struct {
	Inspectors map[string]inspector.Inspector
}
