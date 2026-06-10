package ping

import "time"

type PingExtraInspectionInfo struct {
	PacketLoss float64
	MinRtt     time.Duration
	AvgRtt     time.Duration
	MaxRtt     time.Duration
}

type PingInspectorConfig struct {
	PingCount *int
	Interval  *time.Duration
	Timeout   *time.Duration
	Threshold *float64
}

func NewInspectorConfig() *PingInspectorConfig {
	return &PingInspectorConfig{
		PingCount: new(5),
		Interval:  new(time.Second),
		Timeout:   new(5 * time.Second),
		Threshold: new(1.0),
	}
}
