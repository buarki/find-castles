package htmlfetcher

import "testing"

func TestCharset(t *testing.T) {
	testCases := []struct {
		htmlChunk       []byte
		expectedCharset string
	}{
		{
			htmlChunk: []byte(`
			<!DOCTYPE html>
			<html lang="en-US">
			<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0, viewport-fit=cover" />		<meta name='robots' content='index, follow, max-image-preview:large, max-snippet:-1, max-video-preview:-1' />
			`),
			expectedCharset: "UTF-8",
		},
		{
			htmlChunk: []byte(`
			<!--- HIER CONTENT --><!--Orginaldatei: ../r30/vc_shop/daten/firma73/r30msvcshop_detail_anzeige_hauptdaten.html --><!--[if lt IE 7]>      <html class="no-js lt-ie9 lt-ie8 lt-ie7"> <![endif]--><!--[if IE 7]>         <html class="no-js lt-ie9 lt-ie8"> <![endif]--><!--[if IE 8]>         <html class="no-js lt-ie9"> <![endif]--><!--[if gt IE 8]><!--><html class="no-js"><!--<![endif]--><head>
			<title>-- EBIDAT - Burgendatenbank des Europäischen Burgeninstitutes --</title>
			<meta name="viewport" content="width=device-width, initial-scale=1"></meta>
			<meta name="author" content="M/S VisuCom GmbH"></meta>
			<meta name="keywords" lang="de" content="ebidat,burgenforschung,burgeninstitut,burgeninventar,burgeninventarisierung, europäisches burgeninstitut,bauwerke,bauforschung,burgenvereinigung,philippsburg,braubach,burgen,rhein,donau,imareal,burgenverein"></meta>
			<meta name="description" content="EBIDAT - Burgendatenbank des Europäischen Burgeninstitutes"></meta>
			<meta name="description" content="Startseite"></meta>
			<meta http-equiv="content-type" content="text/html; charset=iso-8859-1"></meta>`),
			expectedCharset: "ISO-8859-1",
		},
		{
			htmlChunk: []byte(`
			<!DOCTYPE html>
			<html lang="pt">
				<head>
					<meta charset="utf-8">
					<meta name="viewport" content="width=device-width, initial-scale=1.0">
					<meta name="description" content="Castelo de Alvor: Castelos de Portugal, suas carater&amp;iacute;sticas e lendas, disperso pelas Dinastias Afonsina, Aviz, Filipina e de Bragança.">
			`),
			expectedCharset: "UTF-8",
		},
	}

	for _, tt := range testCases {
		currentTT := tt
		t.Run(currentTT.expectedCharset, func(t *testing.T) {
			t.Helper()

			received, err := getCharset(currentTT.htmlChunk)

			if err != nil {
				t.Errorf("expected error nil, got %v", err)
			}

			if received != currentTT.expectedCharset {
				t.Errorf("expected charset to be [%s], got [%s]", currentTT.expectedCharset, received)
			}
		})
	}
}
