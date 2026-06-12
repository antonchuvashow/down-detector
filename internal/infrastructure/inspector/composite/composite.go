package composite

import (
	"fmt"
	"time"

	"detector/internal/inspector/domain"
	"detector/internal/route/domain"
)

type Inspector struct {
	config InspectorConfig
}

func (c *Inspector) Inspect(route route.Route) (inspector.Result, error) {
	result := inspector.Result{
		Start:  time.Now(),
		Config: c.config,
		Status: inspector.ResultStatusSuccess,
	}
	allResults := ExtraInspectionInfo{make(map[string]inspector.Result)}

	for name, insp := range c.config.Inspectors {
		partialResult, err := insp.Inspect(route)
		if err != nil {
			return inspector.Result{}, fmt.Errorf("error in %s inspector: %w", name, err)
		}
		if partialResult.Status == inspector.ResultStatusError {
			result.Status = inspector.ResultStatusError
		}
		allResults.Results[name] = partialResult
	}
	result.Extra = allResults

	return result, nil
}

func NewInspector(config InspectorConfig) *Inspector {
	return &Inspector{config: config}
}
