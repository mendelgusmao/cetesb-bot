package scraper

import (
	"regexp"

	"github.com/go-rod/rod"
)

const (
	beachQualityURL = "https://qualipraia.cetesb.sp.gov.br/qualidade-da-praia/"

	beachExtractor = `
		() => {
			const elements = [...this.querySelectorAll('h2, table')];

			return elements.reduce(
				(cityBeaches, el, index) => {
					if (el.tagName === "TABLE") {
						const cityName = elements[index - 1].innerText.trim();
						const beaches = [...el.querySelectorAll("tbody > tr")].slice(1).map(
							tr => [...tr.querySelectorAll("td")].map(td => td.innerText.trim())
						)
			
						return {...cityBeaches, [cityName]: beaches};
					}
			
					return cityBeaches;
				},
				{},
			);
		}
	`
	samplingDatesExtractor = `
		() => this.querySelector("h3").innerText.trim();
	`
)

var (
	samplingDatesRE = regexp.MustCompile(`De ([^\s]+) até ([^\s]+)`)
)

type Sampling struct {
	StartDate string
	EndDate   string
}

type City struct {
	Name string
}

type Beach struct {
	City
	Sampling
	Name    string
	Quality string
}

type Scraper struct {
	browser *rod.Browser
}
