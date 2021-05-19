package cowin

type CalendarByDistrictResponse struct {
	Centers []Center `json:"centers"`
}

type Center struct {
	CenterId     int       `json:"center_id"`
	Address      string    `json:"address"`
	StateName    string    `json:"state_name"`
	DistrictName string    `json:"district_name"`
	BlockName    string    `json:block_name`
	Pincode      int       `json:"pincode"`
	Sessions     []Session `json:"sessions"`
}

type Session struct {
	SessionId               string   `json:"session_id"`
	Date                    string   `json:"date"`
	AvailableCapacity       int      `json:"available_capacity"`
	AgeLimit                int      `json:"min_age_limit"`
	Slots                   []string `json:"slots"`
	Vaccine                 string   `json:"vaccine"`
	AvailableCapacity4Dose1 int      `json:"available_capacity_dose1"`
	AvailableCapacity4Dose2 int      `json:"available_capacity_dose2"`
}
