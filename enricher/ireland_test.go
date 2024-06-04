package enricher

import (
	"testing"

	"github.com/buarki/find-castles/castle"
)

func TestIrishExtraction(t *testing.T) {
	testCases := []struct {
		html   string
		castle castle.Model
	}{
		{
			html: `
			<div id="place--contact">
			<div>
				<h2>Contact</h2>
				<p class="address">Adare Heritage Centre<br />
					Adare<br />
					Co. Limerick<br />
					V94 DWV7</p>
				<p class="phone">061 396 666</p>
				<p class="email"><a href="mailto:reception@adareheritagecentre.ie">reception@adareheritagecentre.ie</a></p>
			</div>
			<div>
			`,
			castle: castle.Model{
				Name:     "Adare",
				District: "Adare Heritage Centre",
				City:     "Adare",
				State:    "Co. Limerick",
			},
		},
		{
			html: `
			<div id="place--contact">
				<div>
					<h2>Contact</h2>
					<p class="address">Trim <br />
						Co Meath<br>C15 HN90</p>
					<p class="phone">046 9438619</p>
					<p class="email"><a href="mailto:trimcastle@opw.ie">trimcastle@opw.ie</a></p>
				</div>
			<div>
			`,
			castle: castle.Model{
				Name:     "Trim",
				District: "Trim",
				City:     "Trim",
				State:    "Co Meath",
			},
		},
		{
			html: `
			<div id="place--contact">
				<div>
					<h2>Contact</h2>
					<p class="address">Ross Castle, <br />
						Ross Road, <br />
						Killarney, <br />
						Co. Kerry<br>V93 V304</p>
					<p class="phone">064 663 5851</p>
					<p class="email"><a href="mailto:rosscastle@opw.ie">rosscastle@opw.ie</a></p>
				</div>
			<div>`,
			castle: castle.Model{
				Name:     "Ross Castle",
				District: "Ross Road",
				City:     "Killarney",
				State:    "Co. Kerr",
			},
		},
	}

	for _, tt := range testCases {
		res, err := extractIrelandCastleInfo(tt.castle, []byte(tt.html))
		if err != nil {
			t.Errorf("expecte err nil, got %v", err)

			if res.City != tt.castle.City {
				t.Errorf("expected city [%s], got [%s]", tt.castle.City, res.City)
			}
			if res.State != tt.castle.State {
				t.Errorf("expected State [%s], got [%s]", tt.castle.State, res.State)
			}
			if res.District != tt.castle.District {
				t.Errorf("expected District [%s], got [%s]", tt.castle.District, res.District)
			}
		}
	}
}
