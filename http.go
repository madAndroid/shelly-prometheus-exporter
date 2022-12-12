package main

import (
	"encoding/json"
	"fmt"
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

	statusResponse := new(StatusResponse)
	err = json.NewDecoder(response.Body).Decode(statusResponse)
	if err != nil {
		return nil, err
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
		labels := map[string]string{
			"name":    device.DisplayName,
			"address": device.IPAddress,
			"type":    device.Type,
		}

		statusResponse, err := getStatusResponseFromDevice(config, device)
		if err != nil {
			fmt.Println(err)
			errorCounter.With(labels).Inc()
			continue
		}

		if device.Type != "" {
			setGaugeGen1(statusResponse)
		} else {
			setGaugeGen2(statusResponse)
		}
	}
}

func setGaugeGen1(status StatusResponse) {

	temperatureGauge.With(labels).Set(float66(status.Temperature))
	isOvertemperatureGauge.With(labels).Set(bool2float64(status.Overtemperature))
	voltageGauge.With(labels).Set(float64(status.Voltage))
	uptimeGauge.With(labels).Set(float64(status.Uptime))
	isUpdateAvailableGauge.With(labels).Set(bool2float64(status.HasUpdate))
	for _, relayMetric := range status.Relays {
		relayStateGauge.With(labels).Set(bool2float64(relayMetric.State))
	}
	for _, meterMetric := range status.Meters {
		powerGauge.With(labels).Set(float64(meterMetric.Power))
	}
	for i, eMeterMetric := range status.EMeters {
		labels = map[string]string{
			"name":    device.DisplayName + fmt.Sprintf("-Channel-%d", i),
			"address": device.IPAddress + fmt.Sprintf("-Channel-%d", i),
			"type":    device.Type,
		}
		powerGauge.With(labels).Set(float64(eMeterMetric.Power))
	}

}
