package stats

import (
	"CCP_backend/backend/datamodel"
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

func (c *MyColly) CollectInfo() (arrivalCount []*datamodel.ControlPointInfo, dates []string, cpNames []string, err error) {
	crawlerResultMap := make(map[string][]int)
	cpNameSet := make(map[string]struct{})
	c.Colly.OnRequest(func(request *colly.Request) {
		//fmt.Printf("Visiting %s\n", request.URL)
	})

	c.Colly.OnError(func(response *colly.Response, scrappingErr error) {
		fmt.Printf("Error while scrapping: %s\n", err.Error())
		err = scrappingErr
	})

	c.Colly.OnHTML("tr", func(element *colly.HTMLElement) {
		controlPoint := element.DOM.Find("td[headers='Control_Point']").Text()
		cpNameSet[controlPoint] = struct{}{}
		totalArrivalString := element.DOM.Find("td[headers='Total_Arrival']").Text()
		if controlPoint != "" && totalArrivalString != "" {
			totalArrival, err := strconv.Atoi(strings.ReplaceAll(totalArrivalString, ",", ""))
			if err != nil {
				log.Fatal(err)
			}
			if _, ok := crawlerResultMap[controlPoint]; !ok {
				crawlerResultMap[controlPoint] = []int{totalArrival}
			} else {
				crawlerResultMap[controlPoint] = append(crawlerResultMap[controlPoint], totalArrival)
			}
		}
	})

	timeSlot := tool.LastThirtyDays()
	for _, v := range timeSlot {
		c.Colly.Visit(fmt.Sprintf("https://www.immd.gov.hk/eng/stat_%s.html", v))
	}
	crawlerResult := make([]*datamodel.ControlPointInfo, 0)
	for k, v := range crawlerResultMap {
		crawlerResult = append(crawlerResult,
			&datamodel.ControlPointInfo{
				ControlPointName: k,
				ArrivalCount:     v,
			},
		)
	}

	for k, _ := range cpNameSet {
		if k != "" {
			cpNames = append(cpNames, k)
		}
	}
	cpNames = tool.SortSlice(cpNames, "Total")

	return crawlerResult, timeSlot, cpNames, err
}
