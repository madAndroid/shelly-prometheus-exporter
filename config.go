package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type configuration struct {
	Port           int
	ScrapeInterval time.Duration
	RequestTimeout time.Duration
	Devices        []device
}

type device struct {
	DisplayName  string
	Username     string
	Password     string
	IPAddress    string
	MACAddress   string
	Type         string
	ChannelNames map[string]string // Optional: channel index to friendly name
}

// getStatusURLs returns a list of status endpoints to poll for a device.
// For Gen2 (Plus/Pro) devices, returns /rpc/Switch.GetStatus?id=N for each channel.
// For Gen1, returns /status.
func (d device) getStatusURLs() []string {
	var urls []string
	// Gen3 detection
	if strings.Contains(strings.ToLower(d.Type), "gen3") {
		if d.IPAddress != "" {
			urls = append(urls, fmt.Sprintf("http://%s/rpc/Shelly.GetStatus", d.IPAddress))
		}
		return urls
	}

	// Gen2 detection: type contains "plus" or "pro" or is "2pm"/"2pmplus" (case-insensitive)
	isGen2 := false
	t := strings.ToLower(d.Type)
	if t != "" {
		if t == "2pm" || t == "2pmplus" || t == "plus" || t == "pro" || t == "1pmplus" {
			isGen2 = true
		}
		if strings.HasSuffix(t, "plus") || strings.HasSuffix(t, "pro") {
			isGen2 = true
		}
	}
	if isGen2 {
		// Assume 2 channels for 2PM/2PM Plus, 1 for 1PM Plus
		numChannels := 2
		if t == "1pmplus" {
			numChannels = 1
		}
		for i := 0; i < numChannels; i++ {
			if d.IPAddress != "" {
				urls = append(urls, fmt.Sprintf("http://%s/rpc/Switch.GetStatus?id=%d", d.IPAddress, i))
			}
		}
	} else {
		if d.IPAddress != "" {
			urls = append(urls, fmt.Sprintf("http://%s/status", d.IPAddress))
		}
		if d.MACAddress != "" {
			urls = append(urls, fmt.Sprintf("http://shelly%s-%s/status", d.Type, d.MACAddress))
		}
	}
	return urls
}

func getConfig() configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error in config file: %s", err))
	}

	var config configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to unmarshal config into struct, %v", err)
	}

	return config
}
