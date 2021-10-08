package main

type StatusResponse struct {
	Device          device
	WiFiStatus      WiFi     `json:"wifi_sta"`
	Cloud           Cloud    `json:"cloud"`
	MQTT            MQTT     `json:"mqtt"`
	Meters          []Meters `json:"meters"`
	Relays          []Relays `json:"relays"`
	Serial          int      `json:"serial"`
	HasUpdate       bool     `json:"has_update"`
	MACAddress      string   `json:"mac"`
	Temperature     float32  `json:"temperature"`
	Overtemperature bool     `json:"overtemperature"`
	MemoryTotal     int      `json:"ram_total"`
	MemoryFree      int      `json:"ram_free"`
	FilesystemSize  int      `json:"fs_size"`
	FilesystemFree  int      `json:"fs_free"`
	Voltage         float32  `json:"voltage"`
	Uptime          int      `json:"uptime"`
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

type Meters struct {
	Connected bool      `json:"connected"`
	Power     float32   `json:"power"`
	Overpower float32   `json:"overpower"`
	Valid     bool      `json:"is_valid"`
	Timestamp int       `json:"timestamp"`
	Counters  []float32 `json:"counters"`
	Total     int       `json:"total"`
}

type Relays struct {
	State           bool    `json:"ison"`
	HasTimer        bool    `json:"has_timer"`
	TimerStarted    int     `json:"timer_started"`
	TimerDuration   int     `json:"timer_duration"`
	timer_remaining int     `json:"timer_remaining"`
	Overpower       float32 `json:"overpower"`
	Source          string  `json:"source"`
}
