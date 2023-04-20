package datamodel

type SubwayScheduleResponse struct {
	ScheduleList []Schedule `json:"scheduleList"`
	OnTime       bool       `json:"onTime"`
}

type Schedule struct {
	Start          string   `json:"start"`
	Destination    string   `json:"destination"`
	DepartureTimes []string `json:"departureTimes"`
}
