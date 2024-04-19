package castle

type Model struct {
	Name             string `json:"name"`
	Link             string `json:"link"`
	Country          string `json:"country"`
	State            string `json:"state"`
	City             string `json:"city"`
	District         string `json:"district"`
	YearOfFoundation string `json:"yearOfFoundation"`
	FlagLink         string `json:"flagLink"`
}
