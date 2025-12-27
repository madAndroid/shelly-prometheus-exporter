package main

import (
	"encoding/json"
	"testing"
)

func TestStatusResponse_Gen3SwitchGetStatus(t *testing.T) {
	// Sample JSON mimicking Shelly 2PM Gen3 response to Switch.GetStatus (single channel)
	jsonData := []byte(`
{
  "id": 0,
  "source": "WS_in",
  "output": true,
  "apower": 12.5,
  "voltage": 230.1,
  "current": 0.054,
  "aenergy": {"total": 100.0},
  "temperature": {"tC": 45.0, "tF": 113.0}
}
`)

	var status StatusResponse
	if err := json.Unmarshal(jsonData, &status); err != nil {
		t.Fatalf("Failed to unmarshal Gen3 Switch.GetStatus: %v", err)
	}

	// Verify direct fields (Gen2 style)
	if status.Output != true {
		t.Errorf("output = %v; want true", status.Output)
	}
	if status.APower != 12.5 {
		t.Errorf("apower = %v; want 12.5", status.APower)
	}
	if status.Voltage != 230.1 {
		t.Errorf("voltage = %v; want 230.1", status.Voltage)
	}
	// Note: We don't check Relays/Meters lists here because getStatusResponseFromURL doesn't populate them for Gen2/Gen3.
	// The fetchDevices loop handles the metric generation based on .Output and .APower fields for Gen2/Gen3.
}
