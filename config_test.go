package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	"github.com/spf13/viper"
)

var yamlContent = []byte(`port: 9123
requestTimeout: 5s
scrapeInterval: 60s
devices:
  - IPAddress: ""
    MACAddress: "ABC12345"
    displayName: "livingRoomShutter"
    type: "switch25"
    username: "some-user"
    password: "pass123"
  - IPAddress: ""
    MACAddress: "123DEF"
    displayName: "kitchenShutter"
    type: "switch25"
    username: "another-user"
    password: "secure"
  - IPAddress: "192.168.88.39"
    displayName: "FamilyRoom2PM"
    type: "2pm"
    channelNames:
      "0": "LeftBlind"
      "1": "RightBlind"
`)

func TestReadConfig(t *testing.T) {

	viper.SetConfigType("yaml")
	viper.ReadConfig(bytes.NewBuffer(yamlContent))

	var givenConfig configuration
	err := viper.Unmarshal(&givenConfig)
	if err != nil {
		t.Errorf("could not unmarshal status endpoint content: %v", err)
	}

	expectedConfig := configuration{
		Port:           9123,
		RequestTimeout: time.Second * time.Duration(5),
		ScrapeInterval: time.Second * time.Duration(60),
		Devices: []device{
			{
				IPAddress:   "",
				MACAddress:  "ABC12345",
				DisplayName: "livingRoomShutter",
				Type:        "switch25",
				Username:    "some-user",
				Password:    "pass123",
			},
			{
				IPAddress:   "",
				MACAddress:  "123DEF",
				DisplayName: "kitchenShutter",
				Type:        "switch25",
				Username:    "another-user",
				Password:    "secure",
			},
			{
				IPAddress:   "192.168.88.39",
				DisplayName: "FamilyRoom2PM",
				Type:        "2pm",
				ChannelNames: map[string]string{
					"0": "LeftBlind",
					"1": "RightBlind",
				},
			},
		},
	}

	if !reflect.DeepEqual(givenConfig, expectedConfig) {
		t.Error("given config file does not match expected config")
	}
	// Additional assertion: check channelNames parsed
	found := false
	for _, d := range givenConfig.Devices {
		if d.DisplayName == "FamilyRoom2PM" {
			if d.ChannelNames == nil || d.ChannelNames["0"] != "LeftBlind" || d.ChannelNames["1"] != "RightBlind" {
				t.Errorf("channelNames not parsed correctly: %+v", d.ChannelNames)
			}
			found = true
		}
	}
	if !found {
		t.Error("FamilyRoom2PM device with channelNames not found in parsed config")
	}

	if urls := givenConfig.Devices[0].getStatusURLs(); len(urls) == 0 || urls[0] != "http://shellyswitch25-ABC12345/status" {
		t.Errorf("getStatusURLs()[0] = %v; want [http://shellyswitch25-ABC12345/status]", urls)
	}

	if urls := givenConfig.Devices[1].getStatusURLs(); len(urls) == 0 || urls[0] != "http://shellyswitch25-123DEF/status" {
		t.Errorf("getStatusURLs()[0] = %v; want [http://shellyswitch25-123DEF/status]", urls)
	}
}
