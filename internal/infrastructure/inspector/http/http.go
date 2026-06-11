package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"
)

type Inspector struct {
	config InspectorConfig
}

func NewInspector(config InspectorConfig) *Inspector {
	// TODO: Add validation of config
	return &Inspector{config: config}
}

func (h *Inspector) Inspect(route routedomain.Route) (inspector.InspectionResult, error) {
	client := http.Client{Timeout: *h.config.Timeout}
	req, err := http.NewRequest(*h.config.Method, route.URL.String(), nil)
	if err != nil {
		return inspector.InspectionResult{}, fmt.Errorf("http inspector: unable to create http request: %w", err)
	}

	req.Header = h.config.Header
	start := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		err, _ := errors.AsType[*url.Error](err)
		if !err.Timeout() {
			return inspector.InspectionResult{}, fmt.Errorf("http inspector: unable to send  inspector: %w", err)
		}
		res := inspector.InspectionResult{
			Status: inspector.StatusError,
			Start:  start,
			End:    time.Now(),
			Config: h.config,
			Extra: ExtraInspectionInfo{
				IsTimeout:  true,
				StatusCode: -1,
			},
		}
		return res, nil
	}
	defer resp.Body.Close()

	if _, ok := h.config.ExpectedCodes[resp.StatusCode]; !ok {
		res := inspector.InspectionResult{
			Status: inspector.StatusError,
			Start:  start,
			End:    time.Now(),
			Config: h.config,
			Extra: ExtraInspectionInfo{
				IsTimeout:  false,
				StatusCode: resp.StatusCode,
			},
		}
		return res, nil
	}

	res := inspector.InspectionResult{
		Status: inspector.StatusSuccess,
		Start:  start,
		End:    time.Now(),
		Config: h.config,
		Extra: ExtraInspectionInfo{
			IsTimeout:  false,
			StatusCode: resp.StatusCode,
		},
	}
	return res, nil
}
