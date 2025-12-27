package main

import (
	"encoding/json"
	"testing"
)

func TestStatusResponse_Gen2SwitchGetStatus(t *testing.T) {
	jsonData := []byte(`
{
  "id": 1,
  "source": "WS_in",
  "output": true,
  "apower": 64.6,
  "voltage": 229.6,
  "freq": 49.9,
  "current": 0.419,
  "pf": 0.67,
  "aenergy": {"total": 108948.948, "by_minute": [1167.192,1140.855,1149.166], "minute_ts": 1757783700},
  "temperature": {"tC": 38.06, "tF": 100.52, "is_valid": true},
  "errors": [],
  "update": {"status": "idle", "has_update": false, "new_version": "20230913-114008/v1.14.0-gcb84623", "old_version": "20230913-114008/v1.14.0-gcb84623"},
  "ram_total": 49672,
  "ram_free": 37316,
  "fs_size": 233681,
  "fs_free": 117719,
  "uptime": 6423051
}
`)
	var status map[string]interface{}
	if err := json.Unmarshal(jsonData, &status); err != nil {
		t.Fatalf("Failed to unmarshal Gen2 Switch.GetStatus: %v", err)
	}
	if status["output"] != true {
		t.Errorf("output = %v; want true", status["output"])
	}
	if status["apower"].(float64) != 64.6 {
		t.Errorf("apower = %v; want 64.6", status["apower"])
	}
	if status["voltage"].(float64) != 229.6 {
		t.Errorf("voltage = %v; want 229.6", status["voltage"])
	}
	if status["uptime"].(float64) != 6423051 {
		t.Errorf("uptime = %v; want 6423051", status["uptime"])
	}
}
