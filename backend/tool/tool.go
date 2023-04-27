package tool

import (
	"sort"
	"time"
)

func LastThirtyDays() (result []string) {
	// Get the current time
	now := time.Now()
	reversed := make([]string, 0)

	// Loop through the last 30 days
	for i := 1; i < 31; i++ {
		// Calculate the date by subtracting the appropriate number of days
		date := now.AddDate(0, 0, -i)

		// Format the date as "YYYYMMDD"
		formattedDate := date.Format("20060102")

		// Print the formatted date to the console
		reversed = append(reversed, formattedDate)
	}

	//do not include today
	for i := len(reversed) - 1; i >= 0; i-- {
		result = append(result, reversed[i])
	}

	return
}

func SortSlice(slice []string, firstValue string) []string {
	// sort slice using custom sort function
	sort.Slice(slice, func(i, j int) bool {
		if slice[i] == firstValue {
			return true
		} else if slice[j] == firstValue {
			return false
		} else {
			return slice[i] < slice[j]
		}
	})
	return slice
}

func GetHongKongMidnight() time.Time {
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		panic(err)
	}
	now := time.Now().In(loc)
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	return midnight
}

func StationCodeToName(code string) string {
	switch code {
	case "AIR":
		return "Airport"
	case "HOK":
		return "Hong Kong"
	case "AWE":
		return "AsiaWorld Expo"
	case "SHS":
		return "Sheung Shui"
	case "LMC":
		return "Lok Ma Chau"
	case "LOW":
		return "Lo Wu"
	case "ADM":
		return "Admiralty"
	case "KOW":
		return "Kowloon"
	case "TUC":
		return "Tung Chung"
	case "TUM":
		return "Tuen Mun"
	case "WKS":
		return "Wu Kai Sha"
	case "TIS":
		return "Tin Shui Wa"
	default:
		return "Unknown"
	}
}
