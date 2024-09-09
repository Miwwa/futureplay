package futureplay_assignment

type PlayerId string

type PlayerData struct {
	Id          PlayerId `json:"id"`
	Level       int      `json:"level"`
	CountryCode string   `json:"country_code"`
}
