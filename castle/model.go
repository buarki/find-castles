package castle

import (
	"errors"
	"strings"
)

var (
	ErrCastlesShouldProbablyBeTheSameToReconcile = errors.New("castles to reconcile should be probably the same")
)

type Contact struct {
	Phone string
	Email string
}

type Model struct {
	Name              string            `json:"name"`
	Sources           []string          `json:"sourcs"`
	Country           Country           `json:"country"`
	State             string            `json:"state"`
	City              string            `json:"city"`
	District          string            `json:"district"`
	FoundationPeriod  string            `json:"foundationPeriod"`
	PropertyCondition PropertyCondition `json:"propertyCondition"`
	Coordinates       string            `json:"coordinates"`
	RawData           any               `json:"rawData"`
	MatchingTags      []string          `json:"matchingTags"`
	PictureURL        string            `json:"pictureLink"`
	Contact           *Contact

	CurrentEnrichmentLink string // current link being used on enrichment
}

func (m Model) FilteredName() string {
	cleanedName := strings.ToLower(m.Name)
	for _, keyword := range keywordsToRemove[m.Country] {
		cleanedName = strings.ReplaceAll(cleanedName, keyword, "")
	}
	return cleanedName
}

// Future plan: power it with AI
func (m Model) IsProbably(c Model) bool {
	if c.Country != m.Country {
		return false
	}
	if len(c.Name) == 0 || len(m.Name) == 0 {
		return false
	}
	mFilteredNames := m.FilteredName()
	cFilteredNames := c.FilteredName()
	if !strings.Contains(cFilteredNames, mFilteredNames) && !strings.Contains(mFilteredNames, cFilteredNames) {
		return false
	}

	mState := strings.ToLower(m.State)
	cState := strings.ToLower(c.State)
	if !strings.Contains(cState, mState) && !strings.Contains(mState, cState) {
		return false
	}

	mCity := strings.ToLower(m.City)
	cCity := strings.ToLower(c.City)
	if !strings.Contains(cCity, mCity) && !strings.Contains(mCity, cCity) {
		return false
	}

	mDistrict := strings.ToLower(m.District)
	cDistrict := strings.ToLower(c.District)
	if !strings.Contains(mDistrict, cDistrict) && !strings.Contains(cDistrict, mDistrict) {
		return false
	}

	// TODO handle foundation...

	return true
}

// Idempotent reconciliation of castles
func (m Model) ReconcileWith(c Model) (Model, error) {
	if !m.IsProbably(c) {
		return Model{}, ErrCastlesShouldProbablyBeTheSameToReconcile
	}
	newCastle := m.Copy()

	if newCastle.Name != c.Name {
		if len(newCastle.Name) > len(c.Name) { // always select the smaller name
			newCastle.Name = c.Name
		}
	}

	if newCastle.State == "" {
		newCastle.State = c.State
	} else {
		if len(newCastle.State) < len(c.State) {
			newCastle.State = c.State
		}
	}

	if newCastle.City == "" {
		newCastle.City = c.City
	} else {
		if len(newCastle.City) < len(c.City) {
			newCastle.City = c.City
		}
	}

	if newCastle.District == "" {
		newCastle.District = c.District
	} else {
		if len(newCastle.District) < len(c.District) {
			newCastle.District = c.District
		}
	}

	if newCastle.FoundationPeriod == "" {
		newCastle.FoundationPeriod = c.FoundationPeriod
	} else {
		if len(newCastle.FoundationPeriod) < len(c.FoundationPeriod) {
			newCastle.FoundationPeriod = c.FoundationPeriod
		}
	}

	if newCastle.Coordinates == "" {
		newCastle.Coordinates = c.Coordinates
	} else {
		if len(newCastle.Coordinates) < len(c.Coordinates) {
			newCastle.Coordinates = c.Coordinates
		}
	}

	if newCastle.Contact == nil {
		if c.Contact != nil {
			newCastle.Contact = c.Contact
		}
	} else {
		if c.Contact != nil {
			if len(c.Contact.Phone) > len(newCastle.Contact.Phone) {
				newCastle.Contact.Phone = c.Contact.Phone
			}
			if len(c.Contact.Email) > len(newCastle.Contact.Email) {
				newCastle.Contact.Email = c.Contact.Email
			}
		}
	}

	if newCastle.PropertyCondition == Unknown {
		if m.PropertyCondition != Unknown {
			newCastle.PropertyCondition = m.PropertyCondition
		}
	} else {
		if newCastle.PropertyCondition.ComparisonWeight() < m.PropertyCondition.ComparisonWeight() {
			newCastle.PropertyCondition = m.PropertyCondition
		}
	}

	var newSources []string
	sourceSet := make(map[string]bool)
	for _, source := range newCastle.Sources {
		if !sourceSet[source] {
			newSources = append(newSources, source)
			sourceSet[source] = true
		}
	}
	for _, source := range c.Sources {
		if !sourceSet[source] {
			newSources = append(newSources, source)
			sourceSet[source] = true
		}
	}

	newCastle.Sources = newSources

	return newCastle, nil
}

func (m *Model) CleanFields() {
	// TODO also remove latin symbols
	m.Name = strings.ToLower(m.FilteredName())
	m.State = strings.ToLower(m.State)
	m.City = strings.ToLower(m.City)
	m.District = strings.ToLower(m.District)
}

func (m Model) GetMatchingTags() []string {
	matchingTags := []string{
		m.Country.String(),
		strings.ToLower(m.FilteredName()),
	}

	matchingTags = append(matchingTags, strings.Split(strings.ToLower(m.FilteredName()), " ")...)

	if len(m.State) > 0 {
		matchingTags = append(matchingTags, strings.ToLower(m.State))
		matchingTags = append(matchingTags, strings.Split(m.State, " ")...)
	}

	if len(m.City) > 0 {
		matchingTags = append(matchingTags, strings.ToLower(m.City))
		matchingTags = append(matchingTags, strings.Split(m.City, " ")...)
	}

	if len(m.District) > 0 {
		matchingTags = append(matchingTags, strings.ToLower(m.District))
		matchingTags = append(matchingTags, strings.Split(m.District, " ")...)
	}

	if len(m.FoundationPeriod) > 0 {
		matchingTags = append(matchingTags, strings.ToLower(m.FoundationPeriod))
	}

	return matchingTags
}

func (m Model) Copy() Model {
	var sourcesCopy []string
	if len(m.Sources) > 0 {
		sourcesCopy = make([]string, len(m.Sources))
		copy(sourcesCopy, m.Sources)
	}

	var matchingTagsCopy []string
	if len(m.MatchingTags) > 0 {
		matchingTagsCopy = make([]string, len(m.MatchingTags))
		copy(matchingTagsCopy, m.MatchingTags)
	}

	return Model{
		Country:               m.Country,
		Name:                  m.Name,
		CurrentEnrichmentLink: m.CurrentEnrichmentLink,
		Sources:               sourcesCopy,
		State:                 m.State,
		City:                  m.City,
		District:              m.District,
		FoundationPeriod:      m.FoundationPeriod,
		PropertyCondition:     m.PropertyCondition,
		Coordinates:           m.Coordinates,
		RawData:               m.RawData,
		MatchingTags:          matchingTagsCopy,
		PictureURL:            m.PictureURL,
		Contact:               m.Contact,
	}
}
