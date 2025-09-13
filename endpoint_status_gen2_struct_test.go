
import (
	package main

	import (
		"encoding/json"
		"testing"
	)

	type Gen2SwitchStatus struct {
		ID          int         `json:"id"`
		Output      bool        `json:"output"`
		APower      float64     `json:"apower"`
		Voltage     float64     `json:"voltage"`
		Uptime      int         `json:"uptime"`
		Temperature Temperature `json:"temperature"`
	}
	"testing"
)

	ID          int         `json:"id"`
	Output      bool        `json:"output"`
	APower      float64     `json:"apower"`
	Voltage     float64     `json:"voltage"`
	Uptime      int         `json:"uptime"`
	Temperature Temperature `json:"temperature"`
}

func TestParseGen2SwitchGetStatusStruct(t *testing.T) {
	jsonData := []byte(`
{
  "id": 1,
  "output": true,
  "apower": 64.6,
  "voltage": 229.6,
  "uptime": 6423051,
  "temperature": {"tC": 38.06, "tF": 100.52, "is_valid": true}
}
`)
	var status Gen2SwitchStatus
	if err := json.Unmarshal(jsonData, &status); err != nil {
		t.Fatalf("Failed to unmarshal Gen2 Switch.GetStatus: %v", err)
	}
	if !status.Output {
		t.Errorf("output = %v; want true", status.Output)
	}
	if status.APower != 64.6 {
		t.Errorf("apower = %v; want 64.6", status.APower)
	}
	if status.Voltage != 229.6 {
		t.Errorf("voltage = %v; want 229.6", status.Voltage)
	}
	if status.Uptime != 6423051 {
		t.Errorf("uptime = %v; want 6423051", status.Uptime)
	}
       if !status.Temperature.Valid || float32(status.Temperature.TC) != 38.06 {
	       t.Errorf("temperature = %+v; want valid and tC=38.06", status.Temperature)
       }
}
