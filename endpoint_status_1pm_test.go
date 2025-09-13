package main

import (
	"encoding/json"
	"testing"
)

func TestStatusResponse_Shelly1PM(t *testing.T) {
	jsonData := []byte(`
{
  "wifi_sta": {"connected": true, "ssid": "test1pm", "ip": "192.168.1.10", "rssi": -55},
  "cloud": {"enabled": true, "connected": true},
  "mqtt": {"connected": true},
  "relays": [
    {"ison": true, "has_timer": false, "timer_started": 0, "timer_duration": 0, "timer_remaining": 0, "overpower": false, "source": "cloud"}
  ],
  "meters": [
    {"power": 23.4, "overpower": 0, "is_valid": true, "timestamp": 654321, "counters": [7.7,8.8,9.9], "total": 321}
  ],
  "inputs": [
    {"input": 0, "event": "btn_toggle", "event_cnt": 5}
  ],
  "temperature": 40.1,
  "overtemperature": false,
  "tmp": {"tC": 40.1, "tF": 104.2, "is_valid": true},
  "has_update": true,
  "update": {"status": "pending", "has_update": true, "new_version": "20210901-123456/v1.10.0@abcdef", "old_version": "20210801-123456/v1.9.0@123456"},
  "mac": "1PMABC654321",
  "serial": 1111,
  "ram_total": 40000,
  "ram_free": 25000,
  "fs_size": 150000,
  "fs_free": 90000,
  "voltage": 231.5,
  "uptime": 987654
}
`)
	var status StatusResponse
	if err := json.Unmarshal(jsonData, &status); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}
	if len(status.Relays) != 1 {
		t.Errorf("Expected 1 relay, got %d", len(status.Relays))
	}
	if len(status.Meters) != 1 {
		t.Errorf("Expected 1 meter, got %d", len(status.Meters))
	}
	if status.Relays[0].Ison != true {
		t.Errorf("Relay state incorrect: %+v", status.Relays[0])
	}
	if status.Meters[0].Power != 23.4 {
		t.Errorf("Meter power incorrect: %+v", status.Meters[0])
	}
	if status.Inputs[0].Event != "btn_toggle" {
		t.Errorf("Input event incorrect: %+v", status.Inputs[0])
	}
       if status.Temperature.Value != 40.1 {
	       t.Errorf("Temperature = %f; want 40.1", status.Temperature.Value)
       }
	if status.Voltage != 231.5 {
		t.Errorf("Voltage = %f; want 231.5", status.Voltage)
	}
}
