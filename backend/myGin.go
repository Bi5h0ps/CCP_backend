package backend

import (
	"CCP_backend/backend/stats"
	"github.com/gin-gonic/gin"
	"log"
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
}

func (g *MyGin) Start(addr string) (err error) {
	err = g.Engine.Run(addr)
	if err != nil {
		log.Default().Println(err)
		return err
	}
	return
}
