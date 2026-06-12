package ping

import (
	"fmt"
	"time"

	"detector/internal/inspector/domain"
	"detector/internal/route/domain"

	"github.com/prometheus-community/pro-bing"
)

type Inspector struct {
	config InspectorConfig
}

func NewInspector(config InspectorConfig) *Inspector {
	return &Inspector{
		config: config,
	}
}

func (i *Inspector) Inspect(route route.Route) (inspector.Result, error) {
	pinger, err := probing.NewPinger(route.URL.Hostname())
	if err != nil {
		return inspector.Result{}, fmt.Errorf("ping inspector: failed to create pinger: %w", err)
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
		return inspector.Result{}, fmt.Errorf("ping inspector: failed to run pinger: %w", err)
	}

	stats := pinger.Statistics()
	return inspector.Result{
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

func (i *Inspector) determineStatus(stats *probing.Statistics) inspector.ResultStatus {
	if stats.PacketLoss <= 1-*i.config.Threshold {
		return inspector.ResultStatusSuccess
	}
	return inspector.ResultStatusError
}
