package castle

type Country string

const (
	Portugal Country = "Portugal"
	UK       Country = "UK"
	Ireland  Country = "Ireland"
)

type Model struct {
	Name             string  `json:"name"`
	Link             string  `json:"link"`
	Country          Country `json:"country"`
	State            string  `json:"state"`
	City             string  `json:"city"`
	District         string  `json:"district"`
	YearOfFoundation string  `json:"yearOfFoundation"`
	FlagLink         string  `json:"flagLink"`
}
