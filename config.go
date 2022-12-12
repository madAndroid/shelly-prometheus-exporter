package main

import (
	"fmt"
	"log"
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
	DisplayName string
	Username    string
	Password    string
	IPAddress   string
	MACAddress  string
	Type        string
	Generation  int
}

func (d device) getStatusURL() string {
	var statusURL string
	if d.IPAddress != "" {
		statusURL = fmt.Sprintf("http://%s/status", d.IPAddress)
	}
	if d.MACAddress != "" {
		statusURL = fmt.Sprintf("http://shelly%s-%s/status", d.Type, d.MACAddress)
	}
	if d.Generation == 2 {
		statusURL = fmt.Sprintf("http://%s/rpc/Shelly.GetStatus", d.IPAddress)
	}
	return statusURL
}

func getConfig() configuration {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(fmt.Errorf("fatal error in config file: %s \n", err))
	}

	var config configuration
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to unmarshal config into struct, %v", err)
	}

	return config
}
