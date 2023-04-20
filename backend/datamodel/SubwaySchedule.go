package datamodel

type SubwaySchedule struct {
	Status   int                    `json:"status"`
	Message  string                 `json:"message"`
	SysTime  string                 `json:"sys_time"`
	CurrTime string                 `json:"curr_time"`
	Data     map[string]LineStation `json:"data"`
	IsDelay  string                 `json:"isdelay"`
}

type LineStation struct {
	CurrTime string          `json:"curr_time"`
	SysTime  string          `json:"sys_time"`
	Up       []TrainSchedule `json:"UP"`
	Down     []TrainSchedule `json:"DOWN"`
}

type TrainSchedule struct {
	Ttnt   string `json:"ttnt"`
	Valid  string `json:"valid"`
	Plat   string `json:"plat"`
	Time   string `json:"time"`
	Source string `json:"source"`
	Dest   string `json:"dest"`
	Seq    string `json:"seq"`
}
