package tool

import "time"

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
