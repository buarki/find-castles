package castle

var (
	castleNameKeywordsToRemoveByLanguage = map[Language][]string{
		English: {
			" castle",
			"castle ",
			" palace",
		},
		Portuguese: {
			"castelo de ",
			"torre",
		},
	}

	keywordsToRemove = map[Country][]string{
		Portugal: castleNameKeywordsToRemoveByLanguage[Portuguese],
		Ireland:  castleNameKeywordsToRemoveByLanguage[English],
		UK:       castleNameKeywordsToRemoveByLanguage[English],
	}
)
