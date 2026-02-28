package main

import (
	"encoding/json"
	"testing"
)

func TestStatusResponse_ShellyEM(t *testing.T) {
	emJson := []byte(`{
  "wifi_sta": {"connected": true, "ssid": "IOTDevonDale", "ip": "192.168.66.167", "rssi": -69},
  "cloud": {"enabled": true, "connected": true},
  "mqtt": {"connected": false},
  "relays": [{"ison": true, "has_timer": false, "timer_started": 0, "timer_duration": 0, "timer_remaining": 0, "overpower": false, "is_valid": true, "source": "input"}],
  "emeters": [
    {"power": 852.88, "reactive": -736.96, "pf": -0.76, "voltage": 231.35, "is_valid": true, "total": 20372011.8, "total_returned": 2022.2},
    {"power": 21.79, "reactive": 22.44, "pf": 0.70, "voltage": 231.35, "is_valid": true, "total": 3893409.9, "total_returned": 0.3}
  ],
  "update": {"status": "idle", "has_update": false, "new_version": "20230913-114150/v1.14.0-gcb84623", "old_version": "20230913-114150/v1.14.0-gcb84623", "beta_version": "20231107-164916/v1.14.1-rc1-g0617c15"},
  "ram_total": 51064,
  "ram_free": 35752,
  "fs_size": 233681,
  "fs_free": 156875,
  "uptime": 1202943
}`)

	var status StatusResponse
	if err := json.Unmarshal(emJson, &status); err != nil {
		t.Fatalf("Failed to unmarshal Shelly EM JSON: %v", err)
	}
	if len(status.EMeters) != 2 {
		t.Errorf("Expected 2 emeters, got %d", len(status.EMeters))
	}
	if status.EMeters[0].Power != 852.88 {
		t.Errorf("EMeter 0 power incorrect: got %f, want 852.88", status.EMeters[0].Power)
	}
	if status.EMeters[1].Power != 21.79 {
		t.Errorf("EMeter 1 power incorrect: got %f, want 21.79", status.EMeters[1].Power)
	}
	if status.EMeters[0].TotalReturned != 2022.2 {
		t.Errorf("EMeter 0 total_returned incorrect: got %f, want 2022.2", status.EMeters[0].TotalReturned)
	}
	if status.EMeters[1].TotalReturned != 0.3 {
		t.Errorf("EMeter 1 total_returned incorrect: got %f, want 0.3", status.EMeters[1].TotalReturned)
	}
}
