package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

var (
	dynamicURLPattern = regexp.MustCompile(`^https://www\.miproximopaso\.org/profile/summary.*`)
	mainUrl           = "https://www.miproximopaso.org/find/browse?c=0"
	// mainUrl = "https://www.miproximopaso.org/profile/summary/49-9063.00"
)

type dataMessage struct {
	id       string
	field    string
	value    string
	subField string
}

type sectionCatgegories struct {
	Name  string   `json:"name"`
	Areas []string `json:"areas"`
}

type personality struct {
	Description string   `json:"description"`
	Attributes  []string `json:"attributes"`
}

type career struct {
	Title         string                        `json:"title"`
	Description   string                        `json:"description"`
	Tasks         []string                      `json:"tasks"`
	Knowledge     map[string]sectionCatgegories `json:"knowledge"`
	Capacities    map[string]sectionCatgegories `json:"capacities"`
	Skills        map[string]sectionCatgegories `json:"skills"`
	Technology    map[string]sectionCatgegories `json:"technology"`
	Personality   personality                   `json:"personality"`
	Education     string                        `json:"education"`
	AverageSalary string                        `json:"averageSalary"`
	LowerSalary   string                        `json:"lowerSalary"`
	HighestSalary string                        `json:"highestSalary"`
}

func saveData(e *colly.HTMLElement, field string, value string, dataCollector chan dataMessage, wg *sync.WaitGroup) {
	wg.Add(1)
	currentUrl := e.Request.URL.String()

	splitUrl := strings.Split(currentUrl, "/")

	id := splitUrl[len(splitUrl)-1]

	dm := dataMessage{
		id:    id,
		field: field,
		value: value,
	}

	if value == "" && field != "education" {
		log.Fatalf("error getting %s from %s", field, currentUrl)
	}

	dataCollector <- dm
	wg.Done()
}

func saveSectionCategory(e *colly.HTMLElement, field string, category string, value string, dataCollector chan dataMessage, wg *sync.WaitGroup) {
	wg.Add(1)
	currentUrl := e.Request.URL.String()

	splitUrl := strings.Split(currentUrl, "/")

	id := splitUrl[len(splitUrl)-1]

	dm := dataMessage{
		id:       id,
		field:    field,
		subField: category,
		value:    value,
	}

	if value == "" {
		log.Fatalf("error getting %s's %s from %s", field, category, currentUrl)
	}

	dataCollector <- dm
	wg.Done()
}

func main() {
	c := colly.NewCollector()
	dataCollectorChan := make(chan dataMessage)
	var wg sync.WaitGroup

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

		saveData(e, "title", careerTitle, dataCollectorChan, &wg)
	})

	c.OnHTML("div.fw-bold", func(e *colly.HTMLElement) {
		if !isCareerpage(e) {
			return
		}

		if e.Text == "Lo que hacen:" {
			parentDiv := e.DOM.Parent()

			jobDescription := parentDiv.Contents().Text()

			jobDescription = strings.Replace(jobDescription, "Lo que hacen:", "", 1)

			jobDescription = strings.TrimSpace(jobDescription)

			saveData(e, "description", jobDescription, dataCollectorChan, &wg)
		}
	})

	c.OnHTML("div.fw-bold", func(e *colly.HTMLElement) {
		if !isCareerpage(e) {
			return
		}

		if e.Text == "En el trabajo, usted:" {
			parentDiv := e.DOM.Parent()

			parentDiv.Find("ul li").Each(func(_ int, el *goquery.Selection) {
				task := strings.TrimSpace(el.Text())
				saveData(e, "task", task, dataCollectorChan, &wg)
			})
		}
	})

	c.OnHTML("div.report-section", func(e *colly.HTMLElement) {
		if !isCareerpage(e) {
			return
		}

		sectionTitle := strings.TrimSpace(e.ChildText("h2"))

		bulletedSections := []string{"Conocimiento", "Aptitudes", "Habilidades", "Tecnología"}
		if slices.Contains(bulletedSections, sectionTitle) {
			e.ForEach("h3.seclist", func(_ int, el *colly.HTMLElement) {
				seclistTitle := strings.TrimSpace(el.Text)

				el.DOM.NextFiltered("ul.subsec").Find("li").Each(func(i int, s *goquery.Selection) {
					subsecItem := strings.TrimSpace(s.Text())

					if subsecItem != "" {
						field := ""
						switch sectionTitle {
						case "Conocimiento":
							field = "knowledge"
						case "Aptitudes":
							field = "capacities"
						case "Habilidades":
							field = "skills"
						case "Tecnología":
							field = "technology"
						}

						saveSectionCategory(e, field, seclistTitle, subsecItem, dataCollectorChan, &wg)
					}
				})
			})
		} else if sectionTitle == "Personalidad" {
			sectionDescription := strings.TrimSpace(e.DOM.Find("h2").Next().Text())

			saveData(e, "personality-description", sectionDescription, dataCollectorChan, &wg)

			e.ForEach("div.col-md-6 ul.my-0 li", func(i int, el *colly.HTMLElement) {
				subsecItem := strings.TrimSpace(el.Text)
				saveData(e, "personality-item", subsecItem, dataCollectorChan, &wg)
			})
		} else if sectionTitle == "Educación" {
			educationText := ""

			e.DOM.Find("div.zone-box-wrapper").Children().Each(func(i int, s *goquery.Selection) {
				if !s.HasClass("mb-2") {
					educationText = strings.TrimSpace(s.Text())
				}
			})

			saveData(e, "education", educationText, dataCollectorChan, &wg)
		} else if sectionTitle == "Perspectiva laboral" {
			laboralPerspective := e.ChildText("div.text-center.h4.m-0")
			saveData(e, "laboralPerspective", laboralPerspective, dataCollectorChan, &wg)

			e.DOM.Find("div.salary-wrapper").Each(func(i int, s *goquery.Selection) {
				s.Find("a[role='button']").Each(func(i int, sel *goquery.Selection) {
					salaryText := strings.TrimSpace(sel.Text())
					salaryDescription, exists := sel.Attr("data-bs-title")

					if exists {
						if strings.Contains(salaryDescription, "promedio") {
							saveData(e, "averageSalary", salaryText, dataCollectorChan, &wg)
						} else if strings.Contains(salaryDescription, "menos") {
							saveData(e, "lowerSalary", salaryText, dataCollectorChan, &wg)
						} else if strings.Contains(salaryDescription, "más") {
							saveData(e, "highestSalary", salaryText, dataCollectorChan, &wg)
						}
					}
				})
			})

		}
	})

	// Start data collector

	wg.Add(1)
	go func() {
		dataCollector(dataCollectorChan)
		wg.Done()
	}()

	c.Visit(mainUrl)

	close(dataCollectorChan)

	wg.Wait()
}

func isCareerpage(e *colly.HTMLElement) bool {
	return dynamicURLPattern.MatchString(e.Request.URL.String())
}

func dataCollector(collectorChannel <-chan dataMessage) {
	careers := make(map[string]career, 957)

	for dm := range collectorChannel {
		c, ok := careers[dm.id]

		if !ok {
			c = career{
				Knowledge:  make(map[string]sectionCatgegories),
				Capacities: make(map[string]sectionCatgegories),
				Skills:     make(map[string]sectionCatgegories),
				Technology: make(map[string]sectionCatgegories),
			}
		}

		switch dm.field {
		case "title":
			c.Title = dm.value
		case "description":
			c.Description = dm.value
		case "task":
			c.Tasks = append(c.Tasks, dm.value)
		case "knowledge":
			sectionCategory, ok := c.Knowledge[dm.subField]
			if !ok {
				sectionCategory = sectionCatgegories{
					Name: dm.subField,
				}
			}

			sectionCategory.Areas = append(sectionCategory.Areas, dm.value)

			c.Knowledge[dm.subField] = sectionCategory
		case "capacities":
			sectionCategory, ok := c.Capacities[dm.subField]
			if !ok {
				sectionCategory = sectionCatgegories{
					Name: dm.subField,
				}
			}

			sectionCategory.Areas = append(sectionCategory.Areas, dm.value)

			c.Capacities[dm.subField] = sectionCategory
		case "skills":
			sectionCategory, ok := c.Skills[dm.subField]
			if !ok {
				sectionCategory = sectionCatgegories{
					Name: dm.subField,
				}
			}

			sectionCategory.Areas = append(sectionCategory.Areas, dm.value)

			c.Skills[dm.subField] = sectionCategory
		case "technology":
			sectionCategory, ok := c.Technology[dm.subField]
			if !ok {
				sectionCategory = sectionCatgegories{
					Name: dm.subField,
				}
			}

			sectionCategory.Areas = append(sectionCategory.Areas, dm.value)

			c.Technology[dm.subField] = sectionCategory
		case "personality-description":
			c.Personality.Description = dm.value
		case "personality-item":
			c.Personality.Attributes = append(c.Personality.Attributes, dm.value)
		case "education":
			c.Education = dm.value
		case "averageSalary":
			c.AverageSalary = dm.value
		case "lowerSalary":
			c.LowerSalary = dm.value
		case "highestSalary":
			c.HighestSalary = dm.value
		}

		careers[dm.id] = c
	}

	file, err := os.Create("careers.json")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, c := range careers {
		data, err := json.Marshal(c)
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.WriteString(string(data) + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("Data successfully written to careers.json")
}
