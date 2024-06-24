package castle

type Country string

const (
	Portugal Country = "pt"
	UK       Country = "uk"
	Ireland  Country = "ie"
	Slovakia Country = "sk"
)

func (c Country) String() string {
	return string(c)
}
