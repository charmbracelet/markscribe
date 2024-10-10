package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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

func formatTime(hours int, minutes int, seconds int) string {
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

func bar(percentage float64, barWidth int) string {
	bar := ""
	for i := 0; i < barWidth; i++ {
		if float64(i) < percentage/100/float64(barWidth) {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return fmt.Sprintf("%s  %.2f%%", bar, percentage)
}

func wakatimeLanguagesBar(count int) string {
	data, err := wakatimeData()
	if err != nil {
		return fmt.Sprintf("Error: %s", err)
	}

	// sort languages by percentage
	languages := data.Languages
	for i := range languages {
		for j := i + 1; j < len(languages); j++ {
			if languages[i].Percent < languages[j].Percent {
				languages[i], languages[j] = languages[j], languages[i]
			}
		}
	}

	// pad the name of the language so that they are all equal in lengh to the longest name plus 2 spaces
	longestLanguage := 0
	longestTime := 0
	for _, l := range languages {
		if len(l.Name) > longestLanguage {
			longestLanguage = len(l.Name)
		}
		time := len(formatTime(l.Hours, l.Minutes, l.Seconds))
		if time > longestTime {
			longestTime = time
		}
	}
	for i, l := range languages {
		languages[i].Name = fmt.Sprintf("%-*s", longestLanguage+2, l.Name)
		languages[i].Digital = fmt.Sprintf("%-*s", longestTime+2, formatTime(l.Hours, l.Minutes, l.Seconds))
	}

	// generate the lines in the format: name bar percent%
	var lines []string
	for _, l := range languages {
		lines = append(lines, fmt.Sprintf("%s %s %s", l.Name, l.Digital, bar(l.Percent, 25)))
	}

	return strings.Join(lines[:count], "\n")
}
