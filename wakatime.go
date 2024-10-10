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

func wakatimeCategoryBar(count int, category WakatimeCategoryType) string {
	// sort languages by percentage
	for i := range category {
		for j := i + 1; j < len(category); j++ {
			if category[i].Percent < category[j].Percent {
				category[i], category[j] = category[j], category[i]
			}
		}
	}

	// pad the name of the language so that they are all equal in lengh to the longest name plus 2 spaces
	longestLanguage := 0
	longestTime := 0
	for _, c := range category {
		if len(c.Name) > longestLanguage {
			longestLanguage = len(c.Name)
		}
		time := len(formatTime(c.Hours, c.Minutes, c.Seconds))
		if time > longestTime {
			longestTime = time
		}
	}
	for i, c := range category {
		category[i].Name = fmt.Sprintf("%-*s", longestLanguage+2, c.Name)
		category[i].Digital = fmt.Sprintf("%-*s", longestTime+2, formatTime(c.Hours, c.Minutes, c.Seconds))
	}

	// generate the lines in the format: name bar percent%
	var lines []string
	for _, c := range category {
		lines = append(lines, fmt.Sprintf("%s %s %s", c.Name, c.Digital, bar(c.Percent, 25)))
	}

	if count < len(lines) {
		return strings.Join(lines[:count], "\n")
	}
	return strings.Join(lines, "\n")
}
