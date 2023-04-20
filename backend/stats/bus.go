package stats

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"encoding/json"
)

type ScheduleData struct {
    HongKongToMacau     []string `json:"香港口岸往澳門口岸"`
    MacauToHongKong     []string `json:"澳門口岸往香港口岸"`
    HongKongToZhuHai    []string `json:"香港口岸往珠海口岸"`
    ZhuHaiToHongKong    []string `json:"珠海口岸往香港口岸"`
}

func bus() {
	doc, err := goquery.NewDocument("https://wx.hzmbus.com/info/customer/classinfo.html")
	if err != nil {
		fmt.Println("Error:", err)
	}

	var lineData ScheduleData

	doc.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "let app = new Vue") {
			lineDataStr := strings.SplitAfterN(s.Text(), "data: {", 2)[1]
			lineDataStr = strings.SplitAfterN(lineDataStr, "lineData:", 2)[1]
			lineDataStr = strings.SplitN(lineDataStr, "},", 2)[0]
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), "\n", "")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), "\t", "")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), ",]", "]")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), ",}", "}")
			err := json.Unmarshal([]byte(lineDataStr), &lineData)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return false
		}
		return true
	})
	
	fmt.Printf("%+v", lineData)
	//看呙哥怎么用
	
}

