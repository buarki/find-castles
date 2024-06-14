# ebidat

## Link
Link: https://www.ebidat.de/

## Countries
- Germany: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=1;
- Denmark: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=2;
- Finland: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=3;
- Latvia: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=4;
- Netherlands: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=5;
- Austria: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=9;
- Slovakia: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=6;
- Czech Republic: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=7;
- Hungary: https://www.ebidat.de/cgi-bin/ebidat.pl?a=a&te53=8;

## Data Access Information Status

|Aspect|Status|Note|
|--|--|--|
|Easily accessible by machines|Partially|The data is embedded within HTML and can be extracted using web scraping techniques, but it's not as straightforward or reliable as accessing structured JSON, XML, or CSV. Any changes in the HTML structure could break the data collection process.|
|Follow Good practices of web development|Partially|The HTML uses some structured elements (like section, div, a tags with appropriate classes). However, it relies heavily on inline styles and JavaScript for some functionality, which is not ideal for maintainability and separation of concerns. It also does not seem to use any microdata, RDFa, or JSON-LD, which are standards for embedding semantic metadata within HTML documents. Without these semantic annotations, the data is less discoverable and harder to link to other datasets.|
|Option to select data format|No|The data appears to be available only through HTML, which requires scraping for extraction. There is no option to access the data in other formats like JSON, XML, or CSV.|
|Consistent format|Yes|The HTML structure for the castle entries seems consistent (e.g., each entry has a div with class burgenanreisser containing the name and location details), which allows for reliable scraping as long as the structure remains unchanged.|


## Extracting data

The main challenge here is iterating through the pages, bacause it has a particular way to do the pagination. A chunk of the document that is in charge of is bellow one:

```html
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
<input name="var_datei_selektionen" type="hidden" value="20240614/212718770666c0c0ba6e21d.dat">
<input name="var_uebergabe1" type="hidden" value="">
<input name="var_anzahl_angezeigte_saetze" type="hidden" value="10">
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
```

As we could see:
- the current page is always between <b> tags;
- by knowning the current page we can find the next by doing current page + 1;
- to find out the link to next page link we need to find the form with name "formseite + nextPage_number";
- once the form is found we need to find the input with name "var_datei_selektionen" and extract the value;
- and to properly add the page on query param we use param "var_anzahl_angezeigte_saetze" ad the value is the (page - 1) * 10;

```txt
https://www.ebidat.de/cgi-bin/r30msvcshop_anzeige.pl?var_hauptpfad=..%2Fr30%2Fvc_shop%2F&var_datei_selektionen=20240613%2F212718770666ad15f4ef739.dat&var_anzahl_angezeigte_saetze=10

host: https://www.ebidat.de
resource: cgi-bin/r30msvcshop_anzeige.pl
var_hauptpfad: ..%2Fr30%2Fvc_shop%2F (../r30/vc_shop/)
var_datei_selektionen: 20240613%2F212718770666ad15f4ef739.dat
var_anzahl_angezeigte_saetze: 10
```
