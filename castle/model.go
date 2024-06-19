package castle

import (
	"errors"
	"strings"
)

type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

var (
	ErrCastlesShouldProbablyBeTheSameToReconcile = errors.New("castles to reconcile should be probably the same")
)

type Model struct {
	Name              string            `json:"name"`
	Link              string            `json:"link"`
	Sources           []string          `json:"sourcs"`
	Country           Country           `json:"country"`
	State             string            `json:"state"`
	City              string            `json:"city"`
	District          string            `json:"district"`
	FoundationPeriod  string            `json:"foundationPeriod"`
	PropertyCondition PropertyCondition `json:"propertyCondition"`
	Coordinates       Coordinates       `json:"coordinates"`
	RawData           any               `json:"rawData"`
	MatchingTags      []string          `json:"matchingTags"`
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
func (m Model) ReconcileWith(c Model) (Model, error) {
	if !m.IsProbably(c) {
		return Model{}, ErrCastlesShouldProbablyBeTheSameToReconcile
	}

	newCastle := m.Copy()

	if newCastle.Name != c.Name {
		// always select the smaller name
		if len(newCastle.Name) > len(c.Name) {
			newCastle.Name = c.Name
		}
	}

	if newCastle.State == "" {
		newCastle.State = c.State
	}

	if len(newCastle.State) > len(c.State) {
		newCastle.State = c.State
	}

	if newCastle.City == "" {
		newCastle.City = c.City
	}

	if newCastle.District == "" {
		newCastle.District = c.District
	}

	if newCastle.FoundationPeriod == "" {
		newCastle.FoundationPeriod = c.FoundationPeriod
	}

	// TODO handle property condition, lat and long...

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
	if len(m.Sources) >= 0 {
		copy(sourcesCopy, m.Sources)
	}

	var matchingTagsCopy []string
	if len(m.MatchingTags) > 0 {
		copy(matchingTagsCopy, m.MatchingTags)
	}

	coordinatesCopy := Coordinates{
		Latitude:  m.Coordinates.Latitude,
		Longitude: m.Coordinates.Longitude,
	}

	return Model{
		Country:           m.Country,
		Name:              m.Name,
		Link:              m.Link,
		Sources:           sourcesCopy,
		State:             m.State,
		City:              m.City,
		District:          m.District,
		FoundationPeriod:  m.FoundationPeriod,
		PropertyCondition: m.PropertyCondition,
		Coordinates:       coordinatesCopy,
		RawData:           m.RawData,
		MatchingTags:      matchingTagsCopy,
	}
}
