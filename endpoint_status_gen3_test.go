package main

import (
	"testing"
)

func TestStatusResponse_Gen3ShellyGetStatus(t *testing.T) {
	// Sample JSON mimicking Shelly 2PM Gen3 response to Shelly.GetStatus
	// Includes switch:0, switch:1, pm1:0, pm1:1
	jsonData := []byte(`
{
  "ble": {},
  "cloud": {"connected": false},
  "switch:0": {"id": 0, "output": true, "source": "WS_in", "temperature": {"tC": 45.0, "tF": 113.0}},
  "switch:1": {"id": 1, "output": false, "source": "WS_in", "temperature": {"tC": 46.0, "tF": 114.8}},
  "pm1:0": {"id": 0, "apower": 12.5, "voltage": 230.1, "current": 0.054},
  "pm1:1": {"id": 1, "apower": 0.0, "voltage": 230.1, "current": 0.0},
  "sys": {"mac": "GEN3MACADDR", "uptime": 12345},
  "wifi": {"ssid": "homelab", "rssi": -55}
}
`)

	var d device
	d.DisplayName = "Test Gen3 Device"

	// Directly call the parsing logic (assumed exported or available in same package)
	status, err := parseGen3StatusResponse(jsonData, d)
	if err != nil {
		t.Fatalf("Failed to parse Gen3 response: %v", err)
	}

	// Verify Relays (Switches)
	if len(status.Relays) != 2 {
		t.Errorf("Expected 2 relays, got %d", len(status.Relays))
	}
	if status.Relays[0].Ison != true {
		t.Errorf("Relay[0].Ison = %v; want true", status.Relays[0].Ison)
	}
	if status.Relays[1].Ison != false {
		t.Errorf("Relay[1].Ison = %v; want false", status.Relays[1].Ison)
	}

	// Verify Meters (PM1s)
	if len(status.Meters) != 2 {
		t.Errorf("Expected 2 meters, got %d", len(status.Meters))
	}
	if status.Meters[0].Power != 12.5 {
		t.Errorf("Meter[0].Power = %f; want 12.5", status.Meters[0].Power)
	}
	if status.Meters[1].Power != 0.0 {
		t.Errorf("Meter[1].Power = %f; want 0.0", status.Meters[1].Power)
	}
}
