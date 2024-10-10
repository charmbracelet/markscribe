package main

// new struct
type Wakatime struct {
	apikey  string
	baseurl string
}

type WakatimeDataRes struct {
	Data WakatimeUserStats `json:"data"`
}

type WakatimeUserStats struct {
	TotalSeconds                                    float64                `json:"total_seconds"`
	TotalSecondsIncludingOtherLanguage              float64                `json:"total_seconds_including_other_language"`
	HumanReadableTotal                              string                 `json:"human_readable_total"`
	HumanReadableTotalIncludingOtherLanguage        string                 `json:"human_readable_total_including_other_language"`
	DailyAverage                                    float64                `json:"daily_average"`
	DailyAverageIncludingOtherLanguage              float64                `json:"daily_average_including_other_language"`
	HumanReadableDailyAverage                       string                 `json:"human_readable_daily_average"`
	HumanReadableDailyAverageIncludingOtherLanguage string                 `json:"human_readable_daily_average_including_other_language"`
	Categories                                      []WakatimeCategoryType `json:"categories"`
	Projects                                        []WakatimeCategoryType `json:"projects"`
	Languages                                       []WakatimeCategoryType `json:"languages"`
	Editors                                         []WakatimeCategoryType `json:"editors"`
	OperatingSystems                                []WakatimeCategoryType `json:"operating_systems"`
	Dependencies                                    []WakatimeCategoryType `json:"dependencies"`
	Machines                                        []WakatimeMachines     `json:"machines"`
	BestDay                                         struct {
		Date         string  `json:"date"`
		Text         string  `json:"text"`
		TotalSeconds float64 `json:"total_seconds"`
	} `json:"best_day"`
	Range                   string `json:"range"`
	HumanReadableRange      string `json:"human_readable_range"`
	Holidays                int    `json:"holidays"`
	DaysIncludingHolidays   int    `json:"days_including_holidays"`
	DaysMinusHolidays       int    `json:"days_minus_holidays"`
	Status                  string `json:"status"`
	PercentCalculated       int    `json:"percent_calculated"`
	IsAlreadyUpdating       bool   `json:"is_already_updating"`
	IsCodingActivityVisible bool   `json:"is_coding_activity_visible"`
	IsOtherUsageVisible     bool   `json:"is_other_usage_visible"`
	IsStuck                 bool   `json:"is_stuck"`
	IsIncludingToday        bool   `json:"is_including_today"`
	IsUpToDate              bool   `json:"is_up_to_date"`
	Start                   string `json:"start"`
	End                     string `json:"end"`
	Timezone                string `json:"timezone"`
	Timeout                 int    `json:"timeout"`
	WritesOnly              bool   `json:"writes_only"`
	UserID                  string `json:"user_id"`
	Username                string `json:"username"`
	CreatedAt               string `json:"created_at"`
	ModifiedAt              string `json:"modified_at"`
}

type WakatimeCategoryType struct {
	Name         string  `json:"name"`
	TotalSeconds float64 `json:"total_seconds"`
	Percent      float64 `json:"percent"`
	Digital      string  `json:"digital"`
	Text         string  `json:"text"`
	Hours        int     `json:"hours"`
	Minutes      int     `json:"minutes"`
	Seconds      int     `json:"seconds"`
}

type WakatimeMachines struct {
	*WakatimeCategoryType
	MachineNameID string `json:"machine_name_id"`
}
