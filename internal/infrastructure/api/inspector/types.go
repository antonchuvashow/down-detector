package apiinspector

import (
	"encoding/json"
	"net/http"
	"time"
)

type DurationMs struct {
	time.Duration
}

func (d *DurationMs) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.Milliseconds())
}

func (d *DurationMs) UnmarshalJSON(data []byte) error {
	var res int
	if err := json.Unmarshal(data, &res); err != nil {
		return err
	}

	*d = DurationMs{time.Duration(res) * time.Millisecond}
	return nil
}

type Inspector struct {
	Type   string `json:"type" binding:"required"`
	Config Config `json:"config" binding:"required"`
}

func (i *Inspector) UnmarshalJSON(data []byte) error {
	var aux struct {
		Type   string          `json:"type"`
		Config json.RawMessage `json:"config"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	i.Type = aux.Type

	var config any
	switch aux.Type {
	case "ping":
		config = &PingConfig{}
	case "http":
		config = &HttpConfig{}
	case "composite":
		config = &CompositeConfig{}
	default:
		return ErrInspectorUnknown{aux.Type}
	}

	if err := json.Unmarshal(aux.Config, config); err != nil {
		return err
	}
	i.Config = config

	return nil
}

type Config any

type PingConfig struct {
	PingCount  *int        `json:"ping_count"`
	IntervalMs *DurationMs `json:"interval_ms"`
	TimeoutMs  *DurationMs `json:"timeout_ms"`
	Threshold  *float64    `json:"threshold"`
}

type HttpConfig struct {
	TimeoutMs     *DurationMs `json:"timeout_ms"`
	ExpectedCodes []int       `json:"expected_codes"`
	Method        *string     `json:"method"`
	Header        http.Header `json:"header"`
}

type CompositeConfig struct {
	Inspectors map[string]Inspector `json:"inspectors"`
}
