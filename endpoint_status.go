package main

type StatusResponse struct {
	Device          device
	WiFiStatus      WiFi          `json:"wifi_sta"`
	Cloud           Cloud         `json:"cloud"`
	MQTT            MQTT          `json:"mqtt"`
	Relays          []Relays      `json:"relays"`
	Meters          []OnePMMeters `json:"meters"`
	EMeters         []EMMeters    `json:"emeters"`
	Serial          int           `json:"serial"`
	HasUpdate       bool          `json:"has_update"`
	MACAddress      string        `json:"mac"`
	Temperature     float32       `json:"temperature"`
	Overtemperature bool          `json:"overtemperature"`
	MemoryTotal     int           `json:"ram_total"`
	MemoryFree      int           `json:"ram_free"`
	FilesystemSize  int           `json:"fs_size"`
	FilesystemFree  int           `json:"fs_free"`
	Voltage         float32       `json:"voltage"`
	Uptime          int           `json:"uptime"`
}

// Gen2:
type StatusResponseGen2 struct {
	Device  device
	WiFi    WiFi   `json:"wifi"`
	Cloud   Cloud  `json:"cloud"`
	MQTT    MQTT   `json:"mqtt"`
	Switch0 Switch `json:"switch:0"`
	Switch1 Switch `json:"switch:1"`
	Switch2 Switch `json:"switch:2"`
	Switch3 Switch `json:"switch:3"`
	System 	Sys    `json:"sys"`
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

type OnePMMeters struct {
	Connected bool      `json:"connected"`
	Power     float32   `json:"power"`
	Overpower float32   `json:"overpower"`
	Valid     bool      `json:"is_valid"`
	Timestamp int       `json:"timestamp"`
	Counters  []float32 `json:"counters"`
	Total     int       `json:"total"`
}

type OnePlusPMMeters struct {
	Power float32 `json:"apower"`
}

type EMMeters struct {
	Power    float32 `json:"power"`
	Reactive float32 `json:"reactive"`
	Voltage  float32 `json:"voltage"`
	Valid    bool    `json:"is_valid"`
	Total    float32 `json:"total"`
	TotalRet float32 `json:"total_returned"`
}

type Relays struct {
	State           bool   `json:"ison"`
	HasTimer        bool   `json:"has_timer"`
	TimerStarted    int    `json:"timer_started"`
	TimerDuration   int    `json:"timer_duration"`
	timer_remaining int    `json:"timer_remaining"`
	Overpower       bool   `json:"overpower"`
	Source          string `json:"source"`
}

// Gen2:
type Switch struct {
	Id      int     `json:"id"`
	Output  bool    `json:"output"`
	Power   float32 `json:"apower"`
	Voltage float32 `json:"voltage"`
	Current float32 `json:"current"`
	Temperature Temperature `json:"temperature"`
}

type Temperature struct {
	Celcius float32 `json:"tC"`
}

type Sys struct {
	Mac     		string  `json:"mac"`
	RestartRequired bool    `json:"restart_required"`
	Uptime			int     `json:"uptime"`
	Updates         interface{} `json:"available_updates"`
}