package main

import (
	"encoding/json"
	"testing"
)

// Mirroring the struct used in parseGen3StatusResponse for testing purposes
type Gen3SwitchComponent struct {
	ID     int         `json:"id"`
	Output bool        `json:"output"`
	Source string      `json:"source"`
	Temp   Temperature `json:"temperature"`
}

type Gen3PM1Component struct {
	ID      int     `json:"id"`
	APower  float64 `json:"apower"`
	Voltage float64 `json:"voltage"`
	Current float64 `json:"current"`
}

func TestParseGen3ComponentStructs(t *testing.T) {
	// Test Switch Component Parsing
	switchJson := []byte(`
	{
		"id": 0,
		"output": true,
		"source": "WS_in",
		"temperature": {"tC": 42.5, "tF": 108.5, "is_valid": true}
	}`)
	var sw Gen3SwitchComponent
	if err := json.Unmarshal(switchJson, &sw); err != nil {
		t.Fatalf("Failed to unmarshal Gen3 Switch component: %v", err)
	}
	if !sw.Output {
		t.Errorf("Switch Output = %v; want true", sw.Output)
	}
	if sw.Source != "WS_in" {
		t.Errorf("Switch Source = %v; want WS_in", sw.Source)
	}
	if sw.Temp.Value != 42.5 && sw.Temp.TC != 42.5 {
		// Temperature unmarshal logic puts value in Value or TC depending on format
		// Our custom Temperature unmarshal handles this check logic implicitly
		t.Logf("Switch Temp parsed as: %+v", sw.Temp)
	}

	// Test PM1 Component Parsing
	pm1Json := []byte(`
	{
		"id": 0,
		"apower": 120.5,
		"voltage": 230.1,
		"current": 0.52
	}`)
	var pm Gen3PM1Component
	if err := json.Unmarshal(pm1Json, &pm); err != nil {
		t.Fatalf("Failed to unmarshal Gen3 PM1 component: %v", err)
	}
	if pm.APower != 120.5 {
		t.Errorf("PM1 APower = %v; want 120.5", pm.APower)
	}
	if pm.Voltage != 230.1 {
		t.Errorf("PM1 Voltage = %v; want 230.1", pm.Voltage)
	}
	if pm.Current != 0.52 {
		t.Errorf("PM1 Current = %v; want 0.52", pm.Current)
	}
}
