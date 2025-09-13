package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getStatusResponseFromDevice(config configuration, d device) (*StatusResponse, error) {
	httpClient := &http.Client{Timeout: config.RequestTimeout}

	request, err := http.NewRequest("GET", d.getStatusURL(), nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %v", err)
	}

	if d.Username != "" && d.Password != "" {
		request.SetBasicAuth(d.Username, d.Password)
	}

	response, err := httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error while doing the request for device '%s': %v", d.DisplayName, err)
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

	statusResponse := new(StatusResponse)
	err = json.Unmarshal(bodyBytes, statusResponse)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON for device '%s': %v\nRaw body: %s", d.DisplayName, err, string(bodyBytes))
	}

	return statusResponse, nil
}

func bool2float64(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func fetchDevices(config configuration) {
	for _, device := range config.Devices {
		statusResponse, err := getStatusResponseFromDevice(config, device)
		if err != nil {
			labels := map[string]string{
				"name":    device.DisplayName,
				"address": device.IPAddress,
				"type":    device.Type,
			}
			fmt.Printf("[ERROR] Device: %s (%s)\n%s\n", device.DisplayName, device.IPAddress, err)
			errorCounter.With(labels).Inc()
			continue
		}

		// Per-device metrics
		labels := map[string]string{
			"name":    device.DisplayName,
			"address": device.IPAddress,
			"type":    device.Type,
		}
		temperatureGauge.With(labels).Set(float64(statusResponse.Temperature))
		isOvertemperatureGauge.With(labels).Set(bool2float64(statusResponse.Overtemperature))
		voltageGauge.With(labels).Set(float64(statusResponse.Voltage))
		uptimeGauge.With(labels).Set(float64(statusResponse.Uptime))
		isUpdateAvailableGauge.With(labels).Set(bool2float64(statusResponse.HasUpdate))

		// Per-relay metrics (Shelly 2PM has 2 relays)
		for i, relay := range statusResponse.Relays {
			relayLabels := map[string]string{
				"name":    fmt.Sprintf("%s-Relay-%d", device.DisplayName, i),
				"address": fmt.Sprintf("%s-Relay-%d", device.IPAddress, i),
				"type":    device.Type,
			}
			relayStateGauge.With(relayLabels).Set(bool2float64(relay.Ison))
		}

		// Per-meter metrics (Shelly 2PM has 2 meters)
		for i, meter := range statusResponse.Meters {
			meterLabels := map[string]string{
				"name":    fmt.Sprintf("%s-Meter-%d", device.DisplayName, i),
				"address": fmt.Sprintf("%s-Meter-%d", device.IPAddress, i),
				"type":    device.Type,
			}
			powerGauge.With(meterLabels).Set(float64(meter.Power))
		}
		// Add more metrics as needed for new fields (inputs, tmp, etc.)
	}
}
