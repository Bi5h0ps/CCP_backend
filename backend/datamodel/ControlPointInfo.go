package datamodel

type ControlPointInfo struct {
	ControlPointName string `json:"control_point_name"`
	ArrivalCount     []int  `json:"arrival_count"`
}
