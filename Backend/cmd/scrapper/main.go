package main

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var (
	dynamicURLPattern = regexp.MustCompile(`^https://www\.miproximopaso\.org/profile/summary.*`)
	mainUrl           = "https://www.miproximopaso.org/find/browse?c=0"
)

func main() {
	c := colly.NewCollector()

	c.OnHTML(".list-group-item", func(e *colly.HTMLElement) {
		currentUrl := e.Request.URL.String()

		if currentUrl == mainUrl {
			link := e.ChildAttr("a", "href")
			fullLink := e.Request.AbsoluteURL(link)

			e.Request.Visit(fullLink)
		}
	})

	c.OnHTML("h1.h1-pre-printshare", func(e *colly.HTMLElement) {

		if !isCareerpage(e) {
			return
		}

		careerTitle := e.ChildText("span.main")

		fmt.Println("Career Title: ", careerTitle)
	})

	c.OnHTML("div.col-md-6", func(e *colly.HTMLElement) {
		if !isCareerpage(e) {
			return
		}

		if e.ChildText("div.fw-bold") == "Lo que hacen:" {
			jobDescription := e.DOM.Contents().FilterFunction(func(i int, s *goquery.Selection) bool {
				return goquery.NodeName(s) == "#text"
			}).Text()
			fmt.Println("Job Description:", jobDescription)
		}
	})

	c.OnHTML("div.col.col-12.col-md-6.col-print-12", func(e *colly.HTMLElement) {
		if !isCareerpage(e) {
			return
		}

		if e.ChildText("div.fw-bold") == "En el trabajo, usted:" {
			e.ForEach("ul li", func(_ int, el *colly.HTMLElement) {
				task := strings.TrimSpace(el.Text)
				fmt.Println("Task:", task)
			})
		}
	})

	c.OnHTML("div.report-section", func(e *colly.HTMLElement) {
		if !isCareerpage(e) {
			return
		}

		sectionTitle := strings.TrimSpace(e.ChildText("h2"))
		fmt.Println("Section Title: ", sectionTitle)

		bulletedSections := []string{"Conocimiento", "Aptitudes", "Habilidades", "Tecnología"}
		if slices.Contains(bulletedSections, sectionTitle) {
			e.ForEach("h3.seclist", func(_ int, el *colly.HTMLElement) {
				seclistTitle := strings.TrimSpace(el.Text)
				fmt.Println("Seclist Title: ", seclistTitle)

				el.DOM.Next().Find("ul.subsec li").Each(func(i int, s *goquery.Selection) {
					subsecItem := strings.TrimSpace(s.Text())
					fmt.Println("Subsec:", subsecItem)
				})

			})
		} else if sectionTitle == "Personalidad" {
			sectionDescription := strings.TrimSpace(e.DOM.Find("h2").Next().Text())
			fmt.Println("Section Description: ", sectionDescription)

			e.ForEach("div.col-md-6 ul.my-0 li", func(i int, el *colly.HTMLElement) {
				subsecItem := strings.TrimSpace(el.Text)
				fmt.Println("Subsec:", subsecItem)
			})
		} else if sectionTitle == "Educación" {
			educationText := ""

			e.DOM.Find("div.zone-box-wrapper").Children().Each(func(i int, s *goquery.Selection) {
				if !s.HasClass("mb-2") {
					educationText = strings.TrimSpace(s.Text())
				}
			})

			fmt.Println("Education Requirement: ", strings.TrimSpace(educationText))
		} else if sectionTitle == "Perspectiva laboral" {
			laboralPerspective := e.ChildText("div.text-center.h4.m-0.outlook-0")
			fmt.Println("Laboral Perspective: ", laboralPerspective)

			e.DOM.Find("div.salary-wrapper").Each(func(i int, s *goquery.Selection) {
				s.Find("a[role='button']").Each(func(i int, sel *goquery.Selection) {
					salaryText := strings.TrimSpace(sel.Text())
					salaryDescription, exists := sel.Attr("data-bs-title")

					if exists {
						if strings.Contains(salaryDescription, "promedio") {
							fmt.Println("Average Salary:", salaryText)
						} else if strings.Contains(salaryDescription, "menos") {
							fmt.Println("Lower Salary:", salaryText)
						} else if strings.Contains(salaryDescription, "más") {
							fmt.Println("Highest Salary:", salaryText)
						}
					}
				})
			})

		}
	})

	c.Visit(mainUrl)
}

func isCareerpage(e *colly.HTMLElement) bool {
	return dynamicURLPattern.MatchString(e.Request.URL.String())
}
