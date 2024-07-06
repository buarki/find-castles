package enricher

type Source string

const (
	EDBIDAT            Source = "EDBIDAT"
	CastelosDePortugal Source = "CastelosDePortugal"
	HeritageIreland    Source = "HeritageIreland"
	MedievalBritain    Source = "MedievalBritain"
)

func (s Source) String() string {
	return string(s)
}
