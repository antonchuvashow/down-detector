package ping

import "time"

type ExtraInspectionInfo struct {
	PacketLoss float64
	MinRtt     time.Duration
	AvgRtt     time.Duration
	MaxRtt     time.Duration
}

type InspectorConfig struct {
	PingCount *int
	Interval  *time.Duration
	Timeout   *time.Duration
	Threshold *float64
}

func NewInspectorConfig() *InspectorConfig {
	return &InspectorConfig{
		PingCount: new(5),
		Interval:  new(time.Second),
		Timeout:   new(5 * time.Second),
		Threshold: new(1.0),
	}
}
