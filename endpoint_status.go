package main

type StatusResponse struct {
	WiFiStatus      WiFi    `json:"wifi_sta"`
	Cloud           Cloud   `json:"cloud"`
	MQTT            MQTT    `json:"mqtt"`
	Relays          []Relay `json:"relays"`
	Meters          []Meter `json:"meters"`
	Inputs          []Input `json:"inputs"`
	Temperature     float32 `json:"temperature"`
	Overtemperature bool    `json:"overtemperature"`
	Tmp             Tmp     `json:"tmp"`
	HasUpdate       bool    `json:"has_update"`
	Update          Update  `json:"update"`
	MACAddress      string  `json:"mac"`
	Serial          int     `json:"serial"`
	MemoryTotal     int     `json:"ram_total"`
	MemoryFree      int     `json:"ram_free"`
	FilesystemSize  int     `json:"fs_size"`
	FilesystemFree  int     `json:"fs_free"`
	Voltage         float32 `json:"voltage"`
	Uptime          int     `json:"uptime"`
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
