package castle

type Country string

const (
	Portugal Country = "pt"
	UK       Country = "uk"
	Ireland  Country = "ie"
	Slovakia Country = "sk"
	Denmark  Country = "dk"
)

func (c Country) String() string {
	return string(c)
}
