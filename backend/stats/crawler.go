package stats

import (
	"CCP_backend/backend/datamodel"
	"CCP_backend/backend/tool"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gocolly/colly"
	"log"
	"strconv"
	"strings"
	"time"
)

type MyColly struct {
	Colly       *colly.Collector
	RedisClient *redis.Client
}

func NewColly(redisClient *redis.Client, domains ...string) *MyColly {
	return &MyColly{
		Colly: colly.NewCollector(
			colly.AllowedDomains(domains...)),
		RedisClient: redisClient,
	}
}

func (c *MyColly) GetInfo() (arrivalCount []*datamodel.ControlPointInfo, dates []string, cpNames []string, err error) {
	val, err := c.RedisClient.HGetAll("controlPoints").Result()
	if err != nil {
		return
	}
	if val != nil && len(val) != 0 {
		//if cached, return from cache
		arrivalCount, dates, cpNames, err = c.retrieveInfo(val)
		if err != nil {
			return
		}
	}
	//if cache expired or does not have value set, crawl data from internet
	return c.crawlInfo()
}

func (c *MyColly) retrieveInfo(rawData map[string]string) (arrivalCount []*datamodel.ControlPointInfo, dates []string, cpNames []string, err error) {
	var controlPoints map[string][]int
	var timeSlots []string
	err = json.Unmarshal([]byte(rawData["controlPoints"]), &controlPoints)
	if err != nil {
		return
	}
	err = json.Unmarshal([]byte(rawData["dates"]), &timeSlots)
	if err != nil {
		return
	}
	for k, v := range controlPoints {
		cpNames = append(cpNames, k)
		arrivalCount = append(arrivalCount,
			&datamodel.ControlPointInfo{
				ControlPointName: k,
				ArrivalCount:     v,
			})

	}
	for _, v := range timeSlots {
		dates = append(dates, v)
	}
	return
}

func (c *MyColly) crawlInfo() (arrivalCount []*datamodel.ControlPointInfo, dates []string, cpNames []string, err error) {
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

	//update the cache
	byteControlPoints, err := json.Marshal(crawlerResult)
	if err != nil {
		return
	}
	byteTimeSlot, err := json.Marshal(timeSlot)
	if err != nil {
		return
	}
	c.RedisClient.HSet("controlPointsData", "controlPoints", byteControlPoints)
	c.RedisClient.HSet("controlPointsData", "dates", byteTimeSlot)
	c.RedisClient.Expire("controlPointsData", 2*time.Minute)
	return crawlerResult, timeSlot, cpNames, err
}
