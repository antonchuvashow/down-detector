package composite

import (
	"fmt"
	"time"

	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"
)

type CompositeInspector struct {
	config CompositeInspectorConfig
}

func (c *CompositeInspector) Inspect(route routedomain.Route) (inspector.InspectionResult, error) {
	result := inspector.InspectionResult{
		Start:  time.Now(),
		Config: c.config,
		Status: inspector.StatusSuccess,
	}
	allResults := make(map[string]inspector.InspectionResult)

	for name, insp := range c.config.Inspectors {
		partialResult, err := insp.Inspect(route)
		if err != nil {
			return inspector.InspectionResult{}, fmt.Errorf("error in %s inspector: %w", name, err)
		}
		if partialResult.Status == inspector.StatusError {
			result.Status = inspector.StatusError
		}
		allResults[name] = partialResult
	}
	result.Extra = allResults

	return result, nil
}

func NewCompositeInspector(config CompositeInspectorConfig) *CompositeInspector {
	return &CompositeInspector{config: config}
}
