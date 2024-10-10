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
		if float64(i) < percentage/(100/float64(barWidth)) {
			bar += "█"
		} else {
			bar += "░"
		}
	}

	return fmt.Sprintf("%s  %.2f%%", bar, percentage)
}

func wakatimeCategoryBar(count int, category interface{}) string {
	var typedCategory []WakatimeCategoryType

	switch v := category.(type) {
	case []WakatimeCategoryType:
		typedCategory = v
	case []WakatimeMachines:
		// Convert WakatimeMachines to WakatimeCategoryType
		typedCategory = make([]WakatimeCategoryType, len(v))
		for i, machine := range v {
			typedCategory[i] = WakatimeCategoryType{
				Name:         machine.Name,
				TotalSeconds: machine.TotalSeconds,
				Percent:      machine.Percent,
				Digital:      machine.Digital,
				Text:         machine.Text,
				Hours:        machine.Hours,
				Minutes:      machine.Minutes,
				Seconds:      machine.Seconds,
			}
		}
	default:
		panic("unknown category type")
	}

	// sort languages by percentage
	for i := range typedCategory {
		for j := i + 1; j < len(typedCategory); j++ {
			if typedCategory[i].Percent < typedCategory[j].Percent {
				typedCategory[i], typedCategory[j] = typedCategory[j], typedCategory[i]
			}
		}
	}

	// pad the name of the language so that they are all equal in lengh to the longest name plus 2 spaces
	longestName := 0
	longestTime := 0
	for _, c := range typedCategory {
		if len(c.Name) > longestName {
			longestName = len(c.Name)
		}
		time := len(formatTime(c.Hours, c.Minutes, c.Seconds))
		if time > longestTime {
			longestTime = time
		}
	}
	for i, c := range typedCategory {
		typedCategory[i].Name = fmt.Sprintf("%-*s", longestName+2, c.Name)
		typedCategory[i].Digital = fmt.Sprintf("%-*s", longestTime+2, formatTime(c.Hours, c.Minutes, c.Seconds))
	}

	// generate the lines in the format: name bar percent%
	var lines []string
	for _, c := range typedCategory {
		lines = append(lines, fmt.Sprintf("%s %s %s", c.Name, c.Digital, bar(c.Percent, 25)))
	}

	if count < len(lines) {
		return strings.Join(lines[:count], "\n")
	}
	return strings.Join(lines, "\n")
}
