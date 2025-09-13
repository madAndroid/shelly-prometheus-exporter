package main

import (
	"encoding/json"
	"testing"
)

func TestStatusResponse_MultiChannel2PM(t *testing.T) {
	jsonData := []byte(`
{
  "wifi_sta": {"connected": true, "ssid": "test", "ip": "192.168.1.2", "rssi": -60},
  "cloud": {"enabled": false, "connected": false},
  "mqtt": {"connected": false},
  "relays": [
    {"ison": true, "has_timer": false, "timer_started": 0, "timer_duration": 0, "timer_remaining": 0, "overpower": false, "source": "cloud"},
    {"ison": false, "has_timer": false, "timer_started": 0, "timer_duration": 0, "timer_remaining": 0, "overpower": false, "source": "cloud"}
  ],
  "meters": [
    {"power": 12.3, "overpower": 0, "is_valid": true, "timestamp": 123456, "counters": [1.1,2.2,3.3], "total": 100},
    {"power": 45.6, "overpower": 0, "is_valid": true, "timestamp": 123457, "counters": [4.4,5.5,6.6], "total": 200}
  ],
  "inputs": [
    {"input": 0, "event": "btn_down", "event_cnt": 2},
    {"input": 1, "event": "btn_up", "event_cnt": 3}
  ],
  "temperature": 55.5,
  "overtemperature": false,
  "tmp": {"tC": 55.5, "tF": 131.9, "is_valid": true},
  "has_update": false,
  "update": {"status": "none", "has_update": false, "new_version": "", "old_version": ""},
  "mac": "2PMABC123456",
  "serial": 2222,
  "ram_total": 50000,
  "ram_free": 35000,
  "fs_size": 250000,
  "fs_free": 150000,
  "voltage": 230.0,
  "uptime": 1234567
}
`)
	var status StatusResponse
	if err := json.Unmarshal(jsonData, &status); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if len(status.Relays) != 2 {
		t.Errorf("Expected 2 relays, got %d", len(status.Relays))
	}
	if len(status.Meters) != 2 {
		t.Errorf("Expected 2 meters, got %d", len(status.Meters))
	}
	if status.Relays[0].Ison != true || status.Relays[1].Ison != false {
		t.Errorf("Relay states incorrect: %+v", status.Relays)
	}
	if status.Meters[0].Power != 12.3 || status.Meters[1].Power != 45.6 {
		t.Errorf("Meter powers incorrect: %+v", status.Meters)
	}
	if status.Inputs[0].Event != "btn_down" || status.Inputs[1].Event != "btn_up" {
		t.Errorf("Input events incorrect: %+v", status.Inputs)
	}
	if status.Temperature.Value != 55.5 {
		t.Errorf("Temperature = %f; want 55.5", status.Temperature.Value)
	}
	if status.Voltage != 230.0 {
		t.Errorf("Voltage = %f; want 230.0", status.Voltage)
	}
}
