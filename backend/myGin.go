package backend

import (
	"CCP_backend/backend/datamodel"
	"CCP_backend/backend/stats"
	"CCP_backend/backend/tool"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type MyGin struct {
	Engine *gin.Engine
}

func NewGin() *MyGin {
	return &MyGin{}
}

func (g *MyGin) Init() {
	if g.Engine == nil {
		g.Engine = gin.Default()
	}
	g.registerHandlers()
}

func (g *MyGin) registerHandlers() {
	g.Engine.GET("/data/arrival", func(context *gin.Context) {
		c := stats.NewColly("www.immd.gov.hk")
		data, dates, controlPoints, err := c.CollectInfo()
		if err != nil {
			context.JSON(500, gin.H{
				"data":                nil,
				"dates":               nil,
				"control_point_names": nil,
			})
		} else {
			context.JSON(200, gin.H{
				"data":                data,
				"dates":               dates,
				"control_point_names": controlPoints,
			})
		}
	})

	g.Engine.GET("/data/subway", func(context *gin.Context) {
		// Get parameters from the request
		line := context.Query("line")
		station := context.Query("sta")
		// Send a GET request to the MTR schedule API with the dynamic parameters
		url := fmt.Sprintf("https://rt.data.gov.hk/v1/transport/mtr/getSchedule.php?line=%s&sta=%s", line, station)
		resp, err := http.Get(url)
		if err != nil {
			context.String(http.StatusInternalServerError, "Error: %v", err)
			return
		}
		defer resp.Body.Close()

		// read response body
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		// unmarshal response into struct
		var pojo datamodel.SubwaySchedule
		err = json.Unmarshal(responseBody, &pojo)
		if err != nil {
			fmt.Println("Error unmarshaling response body:", err)
			return
		}

		result := datamodel.SubwayScheduleResponse{}
		if pojo.IsDelay == "N" {
			result.OnTime = true
		}
		for _, lineStation := range pojo.Data {
			up := datamodel.Schedule{}
			for _, trainSchedule := range lineStation.Up {
				up.Start = tool.StationCodeToName(station)
				up.Destination = tool.StationCodeToName(trainSchedule.Dest)
				simplifiedTime := strings.Split(trainSchedule.Time, " ")[1]
				up.DepartureTimes = append(up.DepartureTimes, simplifiedTime)
			}
			down := datamodel.Schedule{}
			for _, trainSchedule := range lineStation.Down {
				down.Start = tool.StationCodeToName(station)
				down.Destination = tool.StationCodeToName(trainSchedule.Dest)
				simplifiedTime := strings.Split(trainSchedule.Time, " ")[1]
				down.DepartureTimes = append(down.DepartureTimes, simplifiedTime)
			}
			result.ScheduleList = []datamodel.Schedule{up, down}
		}
		context.JSON(200, gin.H{
			"data": result,
		})
	})
}

func (g *MyGin) Start(addr string) (err error) {
	err = g.Engine.Run(addr)
	if err != nil {
		log.Default().Println(err)
		return err
	}
	return
}
