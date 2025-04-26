package dto

type AgeResp struct {
	Count int    `json:"count"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
}

type GenderResp struct {
	Count       int     `json:"count"`
	Name        string  `json:"name"`
	Gender      string  `json:"gender"`
	Probability float64 `json:"probability"`
}

type NatCountry struct {
	CountryID   string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

type NatResp struct {
	Count   int          `json:"count"`
	Name    string       `json:"name"`
	Country []NatCountry `json:"country"`
}
