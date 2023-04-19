package stats

import (
	"CCP_backend/backend/tool"
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
)

type MyColly struct {
	Colly *colly.Collector
}

func NewColly(domains ...string) *MyColly {
	return &MyColly{
		Colly: colly.NewCollector(
			colly.AllowedDomains(domains...),
		),
	}
}

func (c *MyColly) CollectInfo() map[string][]int {
	crawlerResult := make(map[string][]int)
	c.Colly.OnRequest(func(request *colly.Request) {
		//fmt.Printf("Visiting %s\n", request.URL)
	})

	c.Colly.OnError(func(response *colly.Response, err error) {
		fmt.Printf("Error while scraping: %s\n", err.Error())
	})

	c.Colly.OnHTML("tr", func(element *colly.HTMLElement) {
		controlPoint := element.DOM.Find("td[headers='Control_Point']").Text()
		totalArrivalString := element.DOM.Find("td[headers='Total_Arrival']").Text()
		if controlPoint != "" && totalArrivalString != "" {
			totalArrival, err := strconv.Atoi(strings.ReplaceAll(totalArrivalString, ",", ""))
			if err != nil {
				log.Fatal(err)
			}
			if _, ok := crawlerResult[controlPoint]; !ok {
				crawlerResult[controlPoint] = []int{totalArrival}
			} else {
				crawlerResult[controlPoint] = append(crawlerResult[controlPoint], totalArrival)
			}
		}
	})

	timeSlot := tool.LastThirtyDays()
	for _, v := range timeSlot {
		c.Colly.Visit(fmt.Sprintf("https://www.immd.gov.hk/hkt/stat_%s.html", v))
	}
	return crawlerResult
}
