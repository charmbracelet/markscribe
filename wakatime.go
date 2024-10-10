package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func wakatimeData() (WakatimeUserStats, error) {
	req, err := http.NewRequest(http.MethodGet, wakatimeClient.baseurl+"/users/current/stats/last_7_days", nil)
	if err != nil {
		return WakatimeUserStats{}, err // Return empty struct and error
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", wakatimeClient.apikey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return WakatimeUserStats{}, err // Return empty struct and error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return WakatimeUserStats{}, fmt.Errorf("wakatime API returned status code %d: %s", resp.StatusCode, resp.Status)
	}

	// Parse the JSON response and populate WakatimeUserStats struct
	decoder := json.NewDecoder(resp.Body)
	var stats WakatimeDataRes
	err = decoder.Decode(&stats)
	if err != nil {
		return WakatimeUserStats{}, fmt.Errorf("error decoding response: %w", err)
	}

	return stats.Data, nil
}
