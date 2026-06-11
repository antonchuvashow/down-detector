package composite

import (
	"fmt"
	"time"

	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"
)

type Inspector struct {
	config InspectorConfig
}

func (c *Inspector) Inspect(route routedomain.Route) (inspector.InspectionResult, error) {
	result := inspector.InspectionResult{
		Start:  time.Now(),
		Config: c.config,
		Status: inspector.InspectionStatusSuccess,
	}
	allResults := ExtraInspectionInfo{make(map[string]inspector.InspectionResult)}

	for name, insp := range c.config.Inspectors {
		partialResult, err := insp.Inspect(route)
		if err != nil {
			return inspector.InspectionResult{}, fmt.Errorf("error in %s inspector: %w", name, err)
		}
		if partialResult.Status == inspector.InspectionStatusError {
			result.Status = inspector.InspectionStatusError
		}
		allResults.Results[name] = partialResult
	}
	result.Extra = allResults

	return result, nil
}

func NewCompositeInspector(config InspectorConfig) *Inspector {
	return &Inspector{config: config}
}
