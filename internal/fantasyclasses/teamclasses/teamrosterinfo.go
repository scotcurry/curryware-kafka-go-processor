package teamclasses

type TeamRosterInfo struct {
	TeamKey                 string `json:"TeamKey"`
	TeamId                  int    `json:"TeamId"`
	ManagerId               int    `json:"ManagerId"`
	PlayerKey               string `json:"PlayerKey"`
	PlayerId                int    `json:"PlayerId"`
	TeamFullName            string `json:"TeamFullName"`
	ByeWeek                 int    `json:"ByeWeek"`
	PrimaryPosition         string `json:"PrimaryPosition"`
	HasPlayerNotes          bool   `json:"HasPlayerNotes"`
	LastPlayerNoteTimestamp *int64 `json:"LastPlayerNoteTimestamp"`
}
