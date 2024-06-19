package castle

type PropertyCondition string

const (
	Unknown PropertyCondition = "unknown"
	Ruins   PropertyCondition = "ruins"
	Damaged PropertyCondition = "damaged"
	Intact  PropertyCondition = "intact"
)

func (pc PropertyCondition) String() string {
	return string(pc)
}
