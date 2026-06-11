package ping

import (
	"fmt"
	"time"

	"detector/internal/inspection/domain/inspector"
	routedomain "detector/internal/route/domain"

	probing "github.com/prometheus-community/pro-bing"
)

type Inspector struct {
	config InspectorConfig
}

func NewInspector(config InspectorConfig) *Inspector {
	return &Inspector{
		config: config,
	}
}

func (i *Inspector) Inspect(route routedomain.Route) (inspector.InspectionResult, error) {
	pinger, err := probing.NewPinger(route.URL.Hostname())
	if err != nil {
		return inspector.InspectionResult{}, fmt.Errorf("ping inspector: failed to create pinger: %w", err)
	}
	if i.config.PingCount != nil {
		pinger.Count = *i.config.PingCount
	}
	if i.config.Interval != nil {
		pinger.Interval = *i.config.Interval
	}
	if i.config.Timeout != nil {
		pinger.Timeout = *i.config.Timeout
	}

	start := time.Now()
	err = pinger.Run()
	if err != nil {
		return inspector.InspectionResult{}, fmt.Errorf("ping inspector: failed to run pinger: %w", err)
	}

	stats := pinger.Statistics()
	return inspector.InspectionResult{
		Status: i.determineStatus(stats),
		Start:  start,
		End:    time.Now(),
		Config: i.config,
		Extra: ExtraInspectionInfo{
			PacketLoss: stats.PacketLoss,
			MinRtt:     stats.MinRtt,
			AvgRtt:     stats.AvgRtt,
			MaxRtt:     stats.MaxRtt,
		},
	}, nil
}

func (i *Inspector) determineStatus(stats *probing.Statistics) inspector.InspectionStatus {
	if stats.PacketLoss <= 1-*i.config.Threshold {
		return inspector.StatusSuccess
	}
	return inspector.StatusError
}
