package castle

type Country string

func (c Country) String() string {
	return string(c)
}

type PropertyCondition string

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

const (
	Portugal Country = "Portugal"
	UK       Country = "UK"
	Ireland  Country = "Ireland"
	Slovakia Country = "Slovakia"

	Unknown PropertyCondition = "unknown"
	Ruins   PropertyCondition = "ruins"
)

type Model struct {
	Name              string            `json:"name"`
	Link              string            `json:"link"`
	Country           Country           `json:"country"`
	State             string            `json:"state"`
	City              string            `json:"city"`
	District          string            `json:"district"`
	FoundationPeriod  string            `json:"foundationPeriod"`
	PropertyCondition PropertyCondition `json:"propertyCondition"`
	FlagLink          string            `json:"flagLink"`
	Coordinates       Coordinates       `json:"coordinates"`
	RawData           any               `json:"rawData"`
}
