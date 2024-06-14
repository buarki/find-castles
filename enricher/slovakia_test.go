package enricher

import (
	"testing"

	"github.com/buarki/find-castles/castle"
)

func TestExtractAccessCode(t *testing.T) {
	e := slovakEnricher{}

	testCases := []struct {
		html       []byte
		found      bool
		formName   string
		accessCode string
	}{
		{
			html: []byte(`
			<section class="ergebnis">
					<ul>
						<li>
							<label>Ergebnis: 46</label> [ <b>1</b>
							<a  href="javascript:document.formseite2.submit()">2</a>
							<a  href="javascript:document.formseite3.submit()">3</a>
							<a  href="javascript:document.formseite4.submit()">4</a>
							<a  href="javascript:document.formseite5.submit()">5</a>
							]
							[<a  href="javascript:document.formseitevor.submit()"> &raquo; </a>]

														<FORM METHOD="GET" ACTION="/cgi-bin/r30msvcshop_anzeige.pl" name="formseite2">
							<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
							<input name="var_datei_selektionen" type="hidden" value="20240613/212718770666ad15f4ef739.dat">
			`),
			found:      true,
			accessCode: `20240613/212718770666ad15f4ef739.dat`,
			formName:   "formseite2",
		},
		{
			html: []byte(`
			<section class="ergebnis">
					<ul>
						<li>
							<label>Ergebnis: 46</label>[ <a  href="javascript:document.formseitezurueck.submit()"> &laquo; </a> ]
							[ <a  href="javascript:document.formseite1.submit()">1</a>
							<a  href="javascript:document.formseite2.submit()">2</a>
							<a  href="javascript:document.formseite3.submit()">3</a>
							<a  href="javascript:document.formseite4.submit()">4</a>
							<b>5</b>
							]

							<FORM METHOD="GET" ACTION="/cgi-bin/r30msvcshop_anzeige.pl" name="formseite1">
							<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
							<input name="var_datei_selektionen" type="hidden" value="20240613/212718770666b68a4ae7faa.dat">
			`),
			found:      false,
			accessCode: ``,
			formName:   "formseite6",
		},
	}
	for _, tt := range testCases {
		found, nonce := e.getNonce(tt.html, tt.formName)
		if tt.found != found {
			t.Errorf("expected to have found [%v], got [%v]", tt.found, found)
		}
		if nonce != tt.accessCode {
			t.Errorf("expected access code to be [%s], got [%s]", tt.accessCode, nonce)
		}
	}
}

func TestGetCurrentPage(t *testing.T) {
	e := slovakEnricher{}
	testCases := []struct {
		html []byte
		page int
		err  bool
	}{
		{
			html: []byte(`
			<section class="ergebnis">
					<ul>
						<li>
							<label>Ergebnis: 46</label> [ <b>1</b>
							<a  href="javascript:document.formseite2.submit()">2</a>
							<a  href="javascript:document.formseite3.submit()">3</a>
							<a  href="javascript:document.formseite4.submit()">4</a>
							<a  href="javascript:document.formseite5.submit()">5</a>
							]
							[<a  href="javascript:document.formseitevor.submit()"> &raquo; </a>]

														<FORM METHOD="GET" ACTION="/cgi-bin/r30msvcshop_anzeige.pl" name="formseite2">
							<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
							<input name="var_datei_selektionen" type="hidden" value="20240613/212718770666ad15f4ef739.dat">
			`),
			err:  false,
			page: 1,
		},
		{
			html: []byte(`
			<section class="ergebnis">
					<ul>
						<li>
							<label>Ergebnis: 46</label>[ <a  href="javascript:document.formseitezurueck.submit()"> &laquo; </a> ]
							[ <a  href="javascript:document.formseite1.submit()">1</a>
							<a  href="javascript:document.formseite2.submit()">2</a>
							<a  href="javascript:document.formseite3.submit()">3</a>
							<a  href="javascript:document.formseite4.submit()">4</a>
							<b>5</b>
							]

							<FORM METHOD="GET" ACTION="/cgi-bin/r30msvcshop_anzeige.pl" name="formseite1">
							<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
							<input name="var_datei_selektionen" type="hidden" value="20240613/212718770666b68a4ae7faa.dat">
			`),
			err:  false,
			page: 5,
		},
		{
			html: []byte(`
			<section class="ergebnis">
					<ul>
						<li>
							<label>Ergebnis: 8154</label>[ <a  href="javascript:document.formseitezurueck.submit()"> &laquo; </a> ]
							[ <a  href="javascript:document.formseite11.submit()">11</a>
							<a  href="javascript:document.formseite12.submit()">12</a>
							<a  href="javascript:document.formseite13.submit()">13</a>
							<a  href="javascript:document.formseite14.submit()">14</a>
							<a  href="javascript:document.formseite15.submit()">15</a>
							<a  href="javascript:document.formseite16.submit()">16</a>
							<a  href="javascript:document.formseite17.submit()">17</a>
							<a  href="javascript:document.formseite18.submit()">18</a>
							<a  href="javascript:document.formseite19.submit()">19</a>
							<a  href="javascript:document.formseite20.submit()">20</a>
							<b>21</b>
							<a  href="javascript:document.formseite22.submit()">22</a>
							<a  href="javascript:document.formseite23.submit()">23</a>
							<a  href="javascript:document.formseite24.submit()">24</a>
							<a  href="javascript:document.formseite25.submit()">25</a>
							<a  href="javascript:document.formseite26.submit()">26</a>
							<a  href="javascript:document.formseite27.submit()">27</a>
							<a  href="javascript:document.formseite28.submit()">28</a>
							<a  href="javascript:document.formseite29.submit()">29</a>
							<a  href="javascript:document.formseite30.submit()">30</a>
							<a  href="javascript:document.formseite31.submit()">31</a>
							]
							[<a  href="javascript:document.formseitevor.submit()"> &raquo; </a>]

														<FORM METHOD="GET" ACTION="/cgi-bin/r30msvcshop_anzeige.pl" name="formseite11">
							<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
							<input name="var_datei_selektionen" type="hidden" value="20240613/212718770666b693223699a.dat">
			`),
			err:  false,
			page: 21,
		},
	}

	for _, tt := range testCases {
		page, err := e.getCurrentPage(tt.html)
		if !tt.err && err != nil {
			t.Errorf("expected err nil, got %v", err)
		}
		if page != tt.page {
			t.Errorf("expected page [%d], got [%d]", tt.page, page)
		}
	}
}

func TestParseNumber(t *testing.T) {
	testCases := []struct {
		currentPage  int
		expectedCode string
	}{
		{currentPage: 1, expectedCode: "00"},
		{currentPage: 2, expectedCode: "10"},
		{currentPage: 3, expectedCode: "20"},
		{currentPage: 19, expectedCode: "180"},
		{currentPage: 22, expectedCode: "210"},
	}

	e := &slovakEnricher{}

	for _, tt := range testCases {
		receivedCode := e.parsePageNumber(tt.currentPage)
		if receivedCode != tt.expectedCode {
			t.Errorf("expected [%s], got [%s]", tt.expectedCode, receivedCode)
		}
	}
}

func TestCheckForNextPage(t *testing.T) {
	testCases := []struct {
		html         []byte
		found        bool
		nextPageLink string
	}{
		{
			html: []byte(`
			<section class="ergebnis">
					<ul>
						<li>
							<label>Ergebnis: 46</label>[ <a  href="javascript:document.formseitezurueck.submit()"> &laquo; </a> ]
							[ <a  href="javascript:document.formseite1.submit()">1</a>
							<a  href="javascript:document.formseite2.submit()">2</a>
							<a  href="javascript:document.formseite3.submit()">3</a>
							<a  href="javascript:document.formseite4.submit()">4</a>
							<b>5</b>
							]

							<FORM METHOD="GET" ACTION="/cgi-bin/r30msvcshop_anzeige.pl" name="formseite1">
							<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
							<input name="var_datei_selektionen" type="hidden" value="20240613/212718770666b68a4ae7faa.dat">
			`),
			found:        false,
			nextPageLink: "",
		},
		{
			html: []byte(`
			<section class="ergebnis">
					<ul>
						<li>
							<label>Ergebnis: 46</label> [ <b>1</b>
							<a  href="javascript:document.formseite2.submit()">2</a>
							<a  href="javascript:document.formseite3.submit()">3</a>
							<a  href="javascript:document.formseite4.submit()">4</a>
							<a  href="javascript:document.formseite5.submit()">5</a>
							]
							[<a  href="javascript:document.formseitevor.submit()"> &raquo; </a>]

														<FORM METHOD="GET" ACTION="/cgi-bin/r30msvcshop_anzeige.pl" name="formseite2">
							<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
							<input name="var_datei_selektionen" type="hidden" value="20240613/212718770666ad15f4ef739.dat">
			`),
			found:        true,
			nextPageLink: "https://www.ebidat.de/cgi-bin/r30msvcshop_anzeige.pl?var_hauptpfad=../r30/vc_shop/&var_datei_selektionen=20240613/212718770666ad15f4ef739.dat&var_anzahl_angezeigte_saetze=10",
		},
		{
			html: []byte(`
			<section class="ergebnis">
			<ul>
				<li>
								<label>Ergebnis: 8154</label>[ <a  href="javascript:document.formseitezurueck.submit()"> &laquo; </a> ]
			[ <a  href="javascript:document.formseite1.submit()">1</a>
			<a  href="javascript:document.formseite2.submit()">2</a>
			<a  href="javascript:document.formseite3.submit()">3</a>
			<a  href="javascript:document.formseite4.submit()">4</a>
			<a  href="javascript:document.formseite5.submit()">5</a>
			<a  href="javascript:document.formseite6.submit()">6</a>
			<a  href="javascript:document.formseite7.submit()">7</a>
			<a  href="javascript:document.formseite8.submit()">8</a>
			<a  href="javascript:document.formseite9.submit()">9</a>
			<a  href="javascript:document.formseite10.submit()">10</a>
			<b>11</b>
			<a  href="javascript:document.formseite12.submit()">12</a>
			<a  href="javascript:document.formseite13.submit()">13</a>
			<a  href="javascript:document.formseite14.submit()">14</a>
			<a  href="javascript:document.formseite15.submit()">15</a>
			<a  href="javascript:document.formseite16.submit()">16</a>
			<a  href="javascript:document.formseite17.submit()">17</a>
			<a  href="javascript:document.formseite18.submit()">18</a>
			<a  href="javascript:document.formseite19.submit()">19</a>
			<a  href="javascript:document.formseite20.submit()">20</a>
			<a  href="javascript:document.formseite21.submit()">21</a>
			]
			[<a  href="javascript:document.formseitevor.submit()"> &raquo; </a>]

								<FORM METHOD="GET" ACTION="/cgi-bin/r30msvcshop_anzeige.pl" name="formseite12">
			<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
			<input name="var_datei_selektionen" type="hidden" value="20240614/212718770666b730228b76b.dat">
			<input name="var_uebergabe1" type="hidden" value="">
			<input name="var_anzahl_angezeigte_saetze" type="hidden" value="0">
			<input name="var_letzter_suchbegriff" type="hidden" value="">
			<input name="var_suchbegriff" type="hidden" value="">
			<input name="var_html_folgemaske" type="hidden" value="r30msvcshop_anzeige.html">
			<input name="var_navigation_iframe" type="hidden" value="">
			<input name="var_transaktionsnr" type="hidden" value="">
			<input name="var_sprache" type="hidden" value="">
			<input name="var_sql_befehl_orderby" type="hidden" value="">
			<input name="var_maximale_anzeige_zeilen_aus_maske" type="hidden" value="">
			<input name="var_geschuetzte_artikel_anzeigen" type="hidden" value="">
			<input name="var_kz_listen" type="hidden" value="">
			</form>
			`),
			found:        true,
			nextPageLink: "https://www.ebidat.de/cgi-bin/r30msvcshop_anzeige.pl?var_hauptpfad=../r30/vc_shop/&var_datei_selektionen=20240614/212718770666b730228b76b.dat&var_anzahl_angezeigte_saetze=110",
		},
	}
	e := &slovakEnricher{}

	for _, tt := range testCases {
		found, nextPageURL := e.checkForNextPage(tt.html)
		if found != tt.found {
			t.Errorf("expected to find [%v], got [%v]", tt.found, found)
		}
		if nextPageURL != tt.nextPageLink {
			t.Errorf("expected to receive [%s], got [%s]", tt.nextPageLink, nextPageURL)
		}
	}
}

func TestCollectCastleNameAndLinks(t *testing.T) {
	testCases := []struct {
		html  []byte
		c     castle.Model
		error bool
	}{
		{
			html: []byte(`
			<div class="main">
			<div style="padding:2.53%">
				<button class="submit-btn" style="float:left;" onclick="window.location.href='../cgi-bin/ebidat.pl'">neue Suche</button>
			</div>
			<div class="mainContent">
				<h2>&Uuml;bersicht</h2>
				<section>
					<form id="shopanzeige_suche" method="POST" action="/cgi-bin/r30msvcshop_anzeige.pl" name="weiter">
						<input name="var_hauptpfad" type="hidden" value="../r30/vc_shop/">
						<input name="var_fa1_select" type="hidden" value="var_fa1_select||73|">
						<input name="var_anzahl_angezeigte_saetze" type="hidden" value="0">
						<input name="var_html_folgemaske" type="hidden" value="r30msvcshop_anzeige.html">
						<input name="var_letzter_suchbegriff" type="hidden" value="">
						<input name="var_datei_selektionen" type="hidden" value="20240614/212718770666b77537b5b21.dat">
						<fieldset>
							<div class="formset type-text">
								<label>Name</label>
								<input name="var_suchbegriff" type="text" value="">
							</div>
							<input class="submit-btn" type="submit" value="suchen">
						</fieldset>
					</form>
				</section>
				<section class="google_earth">
					<script language="JavaScript" type="text/javascript">
						//>> LS 15.01.2008
						//var var_actual_page;
						function get_aktuelle_seite(var_actual_page) {
							// Link soll nur in best.iten angezeigt werden !
							// die Variablen werden in der Kundenspez. Aussprungs-Routine gesetzt (Z:\www\anwendungen\r30\vc_shop\daten\firma73\cgi-bin\r30msvcshop_anzeige_sub30_posten.pm)
							if (("search_tmp_212718770666b77537b5b21" != "") && ("J" == "J") || (var_actual_page > 1)) {
								document.getElementById("gediv").innerHTML = '<a href="/cgi-bin/r30msvcxxx_ebidat_kml_download.pl?obj=search_tmp_212718770666b77537b5b21"><b>Google Earth (aktuelle Suche)</b></a>' +
									'<br>[m&ouml;glicherweise gibt es nicht zu allen Objekten die Google-Earth Daten!]';
							}
							else {
								document.getElementById("gediv").innerHTML = '';
							}
						}
						// Platzhalter fï¿½gle-Earth Link (kann erst im Fussbereich ausgewertet werden!, wegen mï¿½cher Folgeseiten der Suche)
						document.write('<div id="gediv"></div>');
						//>> LS 15.01.2008
					</script>
				</section>
				<section class="burgenanzeige">
					<div class="burgenanreisser">
						<img src="../r30/vc_shop/bilder/firma73/navigation/slowakei.gif">&nbsp;&nbsp;<a
							href="/cgi-bin/ebidat.pl?id=2015"><b>Biely Kamen</b></a>
						<br>Neštich
						<br>Pezinok
						<br>Bratislava
						<script language="JavaScript" type="text/javascript">
							if ('Bratislava' == 'Nordrhein-Westfalen') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 300px;">Erfassung gef&ouml;rdert durch die NRW-Stiftung</span>');
							}
							else if ('Bratislava' == 'Niedersachsen') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die VGH-Stiftung und <br>die Landschaften</span>');
							}
							else if ('215' == '137') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die <br> Sparkassenstiftung Dillenburg</span>');
							}
							else if ('215' == '76') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die KAL</span>');
							}
							else if ('215' == '301') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Stiftung der Sparkasse Oberhessen</span>');
							}
							else if ('215' == '37') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Sparkasse Pforzheim-Calw</span>');
							}
							else if ('215' == '65') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Sparkasse Pforzheim-Calw</span>');
							}
							else if ('215' == '210') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Sparkasse Pforzheim-Calw</span>');
							}
							else if ('215' == '244') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Stiftung der Kreissparkasse Rottweil</span>');
							}
							else if (('215' == '102') || ('215' == '103')) {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Kulturstiftung der Kreissparkasse Heilbronn</span>');
							}
							else if ('215' == '335') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Stiftung der Volksbank Hohenzollern-Balingen eG</span>');
							}
							else if ('215' == '250') {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Stiftung der Kreissparkasse Rottweil f&uuml;r Kunst-,Kultur- und Denkmalpflege</span>');
							}
							else if (('215' == '216') || ('215' == '38')) {
								document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die Sparkasse Pforzheim Calw - Anstalt des &Ouml;ffentlichen Rechts</span>');
							}
						</script>
					</div>
					<div class="burgengrafik">
						<a href="/cgi-bin/ebidat.pl?a=d&m=b&te1=2015"
							onClick="void window.open('/cgi-bin/ebidat.pl?a=d&m=b&te1=2015', '', 'width=400,height=473,resizeable=no,dependent=yes');return false;"
							target="_blank"><img src="../r30/vc_shop/bilder/firma73/65003a_400_473_klein.jpg" border=0
								alt="Biely Kamen, Mauerreste in nordwestlichen Teil der Burg."
								title="Biely Kamen, Mauerreste in nordwestlichen Teil der Burg." width="84" height="100"></a>
					</div>
				</section>
			</main>
			`),
			error: false,
			c: castle.Model{
				Name:     "Biely Kamen",
				District: "Neštich",
				City:     "Pezinok",
				State:    "Bratislava",
				Country:  castle.Slovakia,
				Link:     "",
			},
		},
	}
	e := &slovakEnricher{}

	for _, tt := range testCases {
		foundCastle, err := e.collectCastleNameAndLinks(tt.html)
		if !tt.error && err != nil {
			t.Errorf("expected to have err nil, got [%v]", err)
		}
		if foundCastle[0].Country != tt.c.Country {
			t.Errorf("expected to have country [%s], got [%s]", tt.c.Country, foundCastle[0].Country)
		}
		if foundCastle[0].Name != tt.c.Name {
			t.Errorf("expected to have Name [%s], got [%s]", tt.c.Name, foundCastle[0].Name)
		}
		if foundCastle[0].District != tt.c.District {
			t.Errorf("expected to have District [%s], got [%s]", tt.c.District, foundCastle[0].District)
		}
		if foundCastle[0].City != tt.c.City {
			t.Errorf("expected to have City [%s], got [%s]", tt.c.City, foundCastle[0].City)
		}
		if foundCastle[0].State != tt.c.State {
			t.Errorf("expected to have State [%s], got [%s]", tt.c.State, foundCastle[0].State)
		}
	}
}

func TestExtractDistrictCityAndState(t *testing.T) {
	e := &slovakEnricher{}

	testCases := []struct {
		html string
		c    castle.Model
	}{
		{
			html: `
			Biely Kamen
			Neštich
			Pezinok
			Bratislava

				if ('Bratislava' == 'Nordrhein-Westfalen') {
					document.write('<br><span style="font-style:italic; font-weight: bold; width: 300px;">Erfassung gef&ouml;rdert durch die NRW-Stiftung</span>');
				}
				else if ('Bratislava' == 'Niedersachsen') {
					document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die VGH-Stiftung und <br>die Landschaften</span>');
				}
			`,
			c: castle.Model{
				State:    "Bratislava",
				City:     "Pezinok",
				District: "Neštich",
			},
		},
		{
			html: `
			Budatínsky zámok
			Žilina
			Žilina

			if ('Bratislava' == 'Nordrhein-Westfalen') {
					document.write('<br><span style="font-style:italic; font-weight: bold; width: 300px;">Erfassung gef&ouml;rdert durch die NRW-Stiftung</span>');
				}
				else if ('Bratislava' == 'Niedersachsen') {
					document.write('<br><span style="font-style:italic; font-weight: bold; width: 400px;">Erfassung gef&ouml;rdert durch die VGH-Stiftung und <br>die Landschaften</span>');
				}
			`,
			c: castle.Model{
				State:    "Žilina",
				City:     "Žilina",
				District: "Budatínsky zámok",
			},
		},
	}

	/*
	 */
	for _, tt := range testCases {
		extracted := e.extractDistrictCityAndState(tt.html)
		if extracted.district != tt.c.District {
			t.Errorf("expected district to be [%s], got [%s]", tt.c.District, extracted.district)
		}
		if extracted.city != tt.c.City {
			t.Errorf("expected district to be [%s], got [%s]", tt.c.City, extracted.city)
		}
		if extracted.state != tt.c.State {
			t.Errorf("expected district to be [%s], got [%s]", tt.c.State, extracted.state)
		}
	}
}
