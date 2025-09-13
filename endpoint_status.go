package main

import (
	"encoding/json"
)

type StatusResponse struct {
	WiFiStatus      WiFi        `json:"wifi_sta"`
	Cloud           Cloud       `json:"cloud"`
	MQTT            MQTT        `json:"mqtt"`
	Relays          []Relay     `json:"relays"`
	Meters          []Meter     `json:"meters"`
	Inputs          []Input     `json:"inputs"`
	Temperature     Temperature `json:"temperature"`
	Overtemperature bool        `json:"overtemperature"`
	Tmp             Tmp         `json:"tmp"`
	HasUpdate       bool        `json:"has_update"`
	Update          Update      `json:"update"`
	MACAddress      string      `json:"mac"`
	Serial          int         `json:"serial"`
	MemoryTotal     int         `json:"ram_total"`
	MemoryFree      int         `json:"ram_free"`
	FilesystemSize  int         `json:"fs_size"`
	FilesystemFree  int         `json:"fs_free"`
	Voltage         float32     `json:"voltage"`
	Uptime          int         `json:"uptime"`
}

type Temperature struct {
	Value float32
	TC    float32
	TF    float32
	Valid bool
}

func (t *Temperature) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as float
	var f float32
	if err := json.Unmarshal(data, &f); err == nil {
		t.Value = f
		t.Valid = true
		return nil
	}
	// Try to unmarshal as object
	var obj struct {
		TC float32 `json:"tC"`
		TF float32 `json:"tF"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		t.TC = obj.TC
		t.TF = obj.TF
		t.Valid = true
		return nil
	}
	t.Valid = false
	return nil
}

type WiFi struct {
	Connected bool    `json:"connected"`
	SSID      string  `json:"ssid"`
	IP        string  `json:"ip"`
	RSSI      float32 `json:"rssi"`
}

type Cloud struct {
	Enabled   bool `json:"enabled"`
	Connected bool `json:"connected"`
}

type MQTT struct {
	Connected bool `json:"connected"`
}

type Meter struct {
	Power     float32   `json:"power"`
	Overpower float32   `json:"overpower"`
	IsValid   bool      `json:"is_valid"`
	Timestamp int       `json:"timestamp"`
	Counters  []float32 `json:"counters"`
	Total     int       `json:"total"`
}

type EMMeters struct {
	Power    float32 `json:"power"`
	Reactive float32 `json:"reactive"`
	Voltage  float32 `json:"voltage"`
	Valid    bool    `json:"is_valid"`
	Total    float32 `json:"total"`
	TotalRet float32 `json:"total_returned"`
}

type Relay struct {
	Ison           bool   `json:"ison"`
	HasTimer       bool   `json:"has_timer"`
	TimerStarted   int    `json:"timer_started"`
	TimerDuration  int    `json:"timer_duration"`
	TimerRemaining int    `json:"timer_remaining"`
	Overpower      bool   `json:"overpower"`
	Source         string `json:"source"`
}

type Input struct {
	Input    int    `json:"input"`
	Event    string `json:"event"`
	EventCnt int    `json:"event_cnt"`
}

type Tmp struct {
	TC      float32 `json:"tC"`
	TF      float32 `json:"tF"`
	IsValid bool    `json:"is_valid"`
}

type Update struct {
	Status     string `json:"status"`
	HasUpdate  bool   `json:"has_update"`
	NewVersion string `json:"new_version"`
	OldVersion string `json:"old_version"`
}
