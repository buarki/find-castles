package castle

type Country string

const (
	Portugal Country = "pt"
	UK       Country = "uk"
	Ireland  Country = "ir"
	Slovakia Country = "sv"
)

func (c Country) String() string {
	return string(c)
}
