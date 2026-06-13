package apiinspector

import (
	"fmt"
	"time"

	"detector/internal/infrastructure/inspector/composite"
	"detector/internal/infrastructure/inspector/http"
	"detector/internal/infrastructure/inspector/ping"
	inspector "detector/internal/inspector/domain"
)

func InspectorToDomain(model Inspector) (inspector.Inspector, error) {
	var instance inspector.Inspector
	switch model.Type {
	case "ping":
		cfg, ok := model.Config.(*PingConfig)
		if !ok {
			return nil, fmt.Errorf("api inspector: invalid ping config")
		}
		instance = pingConfigToDomainInspector(*cfg)
	case "http":
		cfg, ok := model.Config.(*HttpConfig)
		if !ok {
			return nil, fmt.Errorf("api inspector: invalid http config")
		}
		instance = httpConfigToDomainInspector(*cfg)
	case "composite":
		cfg, ok := model.Config.(*CompositeConfig)
		if !ok {
			return nil, fmt.Errorf("api inspector: invalid composite config")
		}

		var err error
		instance, err = compositeConfigToDomainInspector(*cfg)
		if err != nil {
			return nil, fmt.Errorf("api inspector: composite mapping failed: %w", err)
		}
	}

	return instance, nil
}

func InspectorFromDomain(domainInspector inspector.Inspector) (Inspector, error) {
	switch insp := domainInspector.(type) {
	case *ping.Inspector:
		return pingInspectorFromDomain(insp), nil

	case *http.Inspector:
		return httpInspectorFromDomain(insp), nil

	case *composite.Inspector:
		return compositeInspectorFromDomain(insp)

	default:
		return Inspector{}, fmt.Errorf("api inspector: unknown inspector type: %T", domainInspector)
	}
}

func pingConfigToDomainInspector(model PingConfig) *ping.Inspector {
	cfg := ping.InspectorConfig{
		PingCount: model.PingCount,
		Interval:  toDuration(model.IntervalMs),
		Timeout:   toDuration(model.TimeoutMs),
		Threshold: model.Threshold,
	}
	return ping.NewInspector(cfg)
}

func pingInspectorFromDomain(insp *ping.Inspector) Inspector {
	cfg := insp.Config()

	return Inspector{
		Type: "ping",
		Config: &PingConfig{
			PingCount:  cfg.PingCount,
			IntervalMs: fromDuration(cfg.Interval),
			TimeoutMs:  fromDuration(cfg.Timeout),
			Threshold:  cfg.Threshold,
		},
	}
}

func httpConfigToDomainInspector(model HttpConfig) *http.Inspector {
	expectedCodes := make(map[int]struct{})
	for _, code := range model.ExpectedCodes {
		expectedCodes[code] = struct{}{}
	}
	if model.ExpectedCodes == nil {
		expectedCodes = nil
	}

	cfg := http.InspectorConfig{
		Timeout:       toDuration(model.TimeoutMs),
		ExpectedCodes: expectedCodes,
		Method:        model.Method,
		Header:        model.Header,
	}
	return http.NewInspector(cfg)
}

func httpInspectorFromDomain(insp *http.Inspector) Inspector {
	cfg := insp.Config()

	var expectedCodes []int
	if cfg.ExpectedCodes != nil {
		expectedCodes = make([]int, 0, len(cfg.ExpectedCodes))
		for code := range cfg.ExpectedCodes {
			expectedCodes = append(expectedCodes, code)
		}
	}

	return Inspector{
		Type: "http",
		Config: &HttpConfig{
			TimeoutMs:     fromDuration(cfg.Timeout),
			ExpectedCodes: expectedCodes,
			Method:        cfg.Method,
			Header:        cfg.Header,
		},
	}
}

func compositeConfigToDomainInspector(model CompositeConfig) (*composite.Inspector, error) {
	inspectors := make(map[string]inspector.Inspector)
	for name, instance := range model.Inspectors {
		var err error
		inspectors[name], err = InspectorToDomain(instance)
		if err != nil {
			return nil, err
		}
	}
	cfg := composite.InspectorConfig{
		Inspectors: inspectors,
	}
	return composite.NewInspector(cfg), nil
}

func compositeInspectorFromDomain(insp *composite.Inspector) (Inspector, error) {
	cfg := insp.Config()

	apiInspectors := make(map[string]Inspector, len(cfg.Inspectors))
	for name, domainInsp := range cfg.Inspectors {
		apiInsp, err := InspectorFromDomain(domainInsp)
		if err != nil {
			return Inspector{}, fmt.Errorf("failed to convert nested inspector %s: %w", name, err)
		}
		apiInspectors[name] = apiInsp
	}

	return Inspector{
		Type: "composite",
		Config: &CompositeConfig{
			Inspectors: apiInspectors,
		},
	}, nil
}

func toDuration(ms *DurationMs) *time.Duration {
	if ms == nil {
		return nil
	}
	return &ms.Duration
}

func fromDuration(d *time.Duration) *DurationMs {
	if d == nil {
		return nil
	}
	return &DurationMs{Duration: *d}
}
