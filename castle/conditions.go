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

func (pc PropertyCondition) ComparisonWeight() int {
	switch pc {
	case Unknown:
		return 0
	case Ruins:
		return 1
	case Damaged:
		return 2
	case Intact:
		return 3
	default:
		return 0
	}
}
