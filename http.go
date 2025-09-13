package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// getStatusResponseFromURL fetches and parses a status response from a given URL for a device.
func getStatusResponseFromURL(config configuration, d device, url string) (*StatusResponse, error) {
	httpClient := &http.Client{Timeout: config.RequestTimeout}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	if d.Username != "" && d.Password != "" {
		request.SetBasicAuth(d.Username, d.Password)
	}

	debug := os.Getenv("DEBUG") != ""
	response, err := httpClient.Do(request)
	if err != nil {
		if debug {
			log.Printf("[DEBUG] Device: %s (%s) HTTP request error: %v\n", d.DisplayName, d.IPAddress, err)
		}
		return nil, fmt.Errorf("error while doing the request for device '%s': %v", d.DisplayName, err)
	}
	if debug {
		log.Printf("[DEBUG] Device: %s (%s) HTTP %d\nHeaders: %v\n", d.DisplayName, d.IPAddress, response.StatusCode, response.Header)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes := make([]byte, 1024)
		n, _ := response.Body.Read(bodyBytes)
		body := string(bodyBytes[:n])
		return nil, fmt.Errorf("device '%s' returned HTTP %d: %s", d.DisplayName, response.StatusCode, body)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body for device '%s': %v", d.DisplayName, err)
	}

	// Log the raw response body for debugging
	if debug {
		log.Printf("[DEBUG] Device: %s (%s) raw response:\n%s\n", d.DisplayName, d.IPAddress, string(bodyBytes))
	}

	statusResponse := new(StatusResponse)
	err = json.Unmarshal(bodyBytes, statusResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON for device '%s': %v\nRaw body: %s", d.DisplayName, err, string(bodyBytes))
	}

	return statusResponse, nil
}

func urlHasOutputField(status StatusResponse) bool {
	// In Go, bool fields default to false, so we can't distinguish missing from false.
	// But for Gen2, if APower is present, Output is always present, so just return true if APower is present.
	return true
}

func bool2float64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func fetchDevices(config configuration) {
	debug := os.Getenv("DEBUG") != ""
	for _, device := range config.Devices {
		if debug {
			log.Printf("[DEBUG] Polling device: %s (%s)", device.DisplayName, device.IPAddress)
		}
		urls := device.getStatusURLs()
		for idx, url := range urls {
			if debug {
				log.Printf("[DEBUG] Polling URL: %s", url)
			}
			// For Gen2, treat each channel as a separate device metric
			statusResponse, err := getStatusResponseFromURL(config, device, url)
			if err != nil {
				labels := map[string]string{
					"name":    device.DisplayName,
					"address": device.IPAddress,
					"type":    device.Type,
				}
				log.Printf("[ERROR] Device: %s (%s)\n%v\n", device.DisplayName, device.IPAddress, err)
				errorCounter.With(labels).Inc()
				continue
			}

			// Per-device/channel metrics
			labels := map[string]string{
				"name":    fmt.Sprintf("%s-Channel-%d", device.DisplayName, idx),
				"address": fmt.Sprintf("%s-Channel-%d", device.IPAddress, idx),
				"type":    device.Type,
			}
			// Use .Value if set (Gen1), else .TC (Gen2)
			temp := float64(0)
			if statusResponse.Temperature.Valid {
				if statusResponse.Temperature.Value != 0 {
					temp = float64(statusResponse.Temperature.Value)
				} else if statusResponse.Temperature.TC != 0 {
					temp = float64(statusResponse.Temperature.TC)
				}
			}
			temperatureGauge.With(labels).Set(temp)
			isOvertemperatureGauge.With(labels).Set(bool2float64(statusResponse.Overtemperature))
			voltageGauge.With(labels).Set(float64(statusResponse.Voltage))
			uptimeGauge.With(labels).Set(float64(statusResponse.Uptime))
			isUpdateAvailableGauge.With(labels).Set(bool2float64(statusResponse.HasUpdate))

			// For single-relay devices, omit -Relay-0 suffix
			if len(statusResponse.Relays) == 1 {
				relayLabels := map[string]string{
					"name":    device.DisplayName,
					"address": device.IPAddress,
					"type":    device.Type,
				}
				relayStateGauge.With(relayLabels).Set(bool2float64(statusResponse.Relays[0].Ison))
			} else {
				for i, relay := range statusResponse.Relays {
					relayLabels := map[string]string{
						"name":    fmt.Sprintf("%s-Relay-%d", device.DisplayName, i),
						"address": fmt.Sprintf("%s-Relay-%d", device.IPAddress, i),
						"type":    device.Type,
					}
					relayStateGauge.With(relayLabels).Set(bool2float64(relay.Ison))
				}
			}

			// Gen2: if no relays but APower present, emit relay state for each channel (2PM etc)
			if len(statusResponse.Relays) == 0 && (statusResponse.APower != 0 || urlHasOutputField(*statusResponse)) {
				// Omit -Channel-0 for single-channel 1pmPlus/1pmplus
				// Gen2/Plus: use Output field for relay state if present, else fallback to APower > 0
				relayState := bool2float64(statusResponse.Output)
				if len(device.getStatusURLs()) == 1 && (device.Type == "1pmplus" || device.Type == "1pmPlus") {
					relayLabels := map[string]string{
						"name":    device.DisplayName,
						"address": device.IPAddress,
						"type":    device.Type,
					}
					relayStateGauge.With(relayLabels).Set(relayState)
				} else if len(device.getStatusURLs()) == 1 {
					relayLabels := map[string]string{
						"name":    device.DisplayName,
						"address": device.IPAddress,
						"type":    device.Type,
					}
					relayStateGauge.With(relayLabels).Set(relayState)
				} else {
					relayLabels := map[string]string{
						"name":    fmt.Sprintf("%s-Channel-%d", device.DisplayName, idx),
						"address": fmt.Sprintf("%s-Channel-%d", device.IPAddress, idx),
						"type":    device.Type,
					}
					relayStateGauge.With(relayLabels).Set(relayState)
				}
			}

			// Emit power metrics (shelly_power) for all devices
			// Gen1: emit per-meter power metrics, but for single-meter devices (e.g., 1PM, 1PMPlus), omit the -Meter-0 suffix
			if len(statusResponse.Meters) == 1 && (device.Type == "1pm" || device.Type == "1PM" || device.Type == "1pmplus" || device.Type == "1pmPlus") {
				meterLabels := map[string]string{
					"name":    device.DisplayName,
					"address": device.IPAddress,
					"type":    device.Type,
				}
				powerGauge.With(meterLabels).Set(float64(statusResponse.Meters[0].Power))
			} else if len(statusResponse.Meters) > 1 {
				for i, meter := range statusResponse.Meters {
					meterLabels := map[string]string{
						"name":    fmt.Sprintf("%s-Meter-%d", device.DisplayName, i),
						"address": fmt.Sprintf("%s-Meter-%d", device.IPAddress, i),
						"type":    device.Type,
					}
					powerGauge.With(meterLabels).Set(float64(meter.Power))
				}
			}
			// Gen2: emit per-channel power metrics using APower if present and no meters
			if len(statusResponse.Meters) == 0 && statusResponse.APower != 0 {
				var meterLabels map[string]string
				// Omit -Channel-0 for single-channel 1pmPlus/1pmplus
				if len(device.getStatusURLs()) == 1 && (device.Type == "1pmplus" || device.Type == "1pmPlus") {
					meterLabels = map[string]string{
						"name":    device.DisplayName,
						"address": device.IPAddress,
						"type":    device.Type,
					}
				} else if len(device.getStatusURLs()) == 1 {
					meterLabels = map[string]string{
						"name":    device.DisplayName,
						"address": device.IPAddress,
						"type":    device.Type,
					}
				} else {
					meterLabels = map[string]string{
						"name":    fmt.Sprintf("%s-Channel-%d", device.DisplayName, idx),
						"address": fmt.Sprintf("%s-Channel-%d", device.IPAddress, idx),
						"type":    device.Type,
					}
				}
				powerGauge.With(meterLabels).Set(statusResponse.APower)
			}
		}
	}
}
