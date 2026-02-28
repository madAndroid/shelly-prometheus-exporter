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
	Generation   int               // Optional: explicit generation (2, 3)
	ChannelNames map[string]string // Optional: channel index to friendly name
}

// getStatusURLs returns a list of status endpoints to poll for a device.
// For Gen2/Gen3 (Plus/Pro/Gen3) devices, returns /rpc/Switch.GetStatus?id=N for each channel.
// For Gen1, returns /status.
func (d device) getStatusURLs() []string {
	var urls []string

	// Gen2/Gen3 detection
	isGen2 := false
	if d.Generation >= 2 {
		isGen2 = true
	}
	t := strings.ToLower(d.Type)
	if t != "" {
		if strings.Contains(t, "gen3") || strings.Contains(t, "plus") || strings.Contains(t, "pro") {
			isGen2 = true
		}
		if t == "2pm" || t == "2pmplus" || t == "1pmplus" {
			isGen2 = true
		}
	}
	if isGen2 {
		// Assume 2 channels for 2PM/2PM Plus/Gen3, 1 for 1PM Plus
		numChannels := 2
		// Adjust simply for single channel devices if known
		if t == "1pmplus" || t == "1pm" || t == "1" || t == "1gen3" || t == "1pmgen3" {
			// This might be too broad if user has "1pm" Gen1 but sets generation 2 (unlikely)
			// But sticking to type names. "1pmplus" is definitely 1 channel.
			// "2pmgen3" will fall through to 2 channels.
			if strings.Contains(t, "1pm") || strings.Contains(t, "1gen3") {
				// Assuming 1 channel for 1PM variants if Gen2/3
				// But explicit check is safer?
				// Existing code only checked "1pmplus".
				// Let's keep it safe. 2PM is the target.
				if t == "1pmplus" {
					numChannels = 1
				}
			}
		}
		// If explicit "1pmplus" check was enough before, sticking close to that logic:
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
