package scraper

import (
	"regexp"

	"github.com/go-rod/rod"
)

const (
	qualityMapURL = "https://qualipraia.cetesb.sp.gov.br/qualidade-da-praia/mapa.php"

	citiesExtractor = `
		() => [...this.querySelectorAll("area")].map(area => area.href)
	`
	beachExtractor = `
		() => {
			const tableSelector = 'table[width="550"]:not([bgcolor]) > tbody';
			const table = this.querySelector(tableSelector);

			return [...table.querySelectorAll("tr")].map(
				tr => [...tr.querySelectorAll("td")].map(
					td => /good/.test(td.innerHTML) || td.innerText || false
				)
			).map(
				it => [it.slice(0, 2), it.slice(3, 5)]
			)
			.flat()
			.filter(it => it.length);
		}
	`
	beachExtraInfoExtractor = `
		() => {
			const header = this.querySelector('table[width="550"][bgcolor] > tbody > tr > td').innerText.split("\n");
			const footer = this.querySelector('table[width="560"] > tbody > tr > td[style]').innerText;

			return [...header, footer];
		}
	`
)

var (
	cityRE          = regexp.MustCompile(`Município de (.*)`)
	currentDateRE   = regexp.MustCompile(`Data: (.*)`)
	samplingDatesRE = regexp.MustCompile(`Período de Amostragem: ([^\s]+) - ([^\s]+)`)
)

type Sampling struct {
	CurrentDate string
	StartDate   string
	EndDate     string
}

type City struct {
	Name string
	URL  string
}

type Beach struct {
	City
	Sampling
	Name   string
	Proper bool
}

type Scraper struct {
	browser *rod.Browser
}
