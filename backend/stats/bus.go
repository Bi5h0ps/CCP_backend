package stats

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ScheduleData struct {
	Routes           []string
	HongKongToMacau  []string `json:"香港口岸往澳門口岸"`
	MacauToHongKong  []string `json:"澳門口岸往香港口岸"`
	HongKongToZhuHai []string `json:"香港口岸往珠海口岸"`
	ZhuHaiToHongKong []string `json:"珠海口岸往香港口岸"`
}

func Bus() ScheduleData {
	doc, err := goquery.NewDocument("https://wx.hzmbus.com/info/customer/classinfo.html")
	if err != nil {
		fmt.Println("Error:", err)
	}

	lineData := ScheduleData{
		Routes: []string{"Hong Kong-Macau", "Macau-Hong Kong", "Hong Kong-Zhu Hai", "Zhu Hai - Hong Kong"},
	}

	doc.Find("script").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if strings.Contains(s.Text(), "let app = new Vue") {
			lineDataStr := strings.SplitAfterN(s.Text(), "data: {", 2)[1]
			lineDataStr = strings.SplitAfterN(lineDataStr, "lineData:", 2)[1]
			lineDataStr = strings.SplitN(lineDataStr, "},", 2)[0]
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), "\n", "")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), "\t", "")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), ",]", "]")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), ",}", "}")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), "<br/>", " ")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), "至", " to ")
			lineDataStr = strings.ReplaceAll(strings.TrimSpace(lineDataStr), "分钟一班", " Minutes per Denature")

			err := json.Unmarshal([]byte(lineDataStr), &lineData)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return false
		}
		return true
	})
	return lineData
}
