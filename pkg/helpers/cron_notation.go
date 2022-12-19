package helpers

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	weeks = []string{"mon", "tue", "wed", "thu", "fri", "sat", "sun"}
)

func isValidMinute(minute string) (res bool) {
	if len(minute) == 0 {
		return
	}
	if minute == "*" {
		res = true
		return
	}
	numberOnly := regexp.MustCompile(`^[\d-]+$`)
	if len(minute) == 1 && !numberOnly.MatchString(minute) {
		return
	}
	if numberOnly.MatchString(minute) {
		val, _ := strconv.Atoi(minute)
		if val < 0 || val > 59 {
			return
		}
	}
	if strings.Contains(minute, "-") {
		splits := strings.Split(minute, "-")
		if len(splits) != 2 {
			return
		}
		minS, maxS := splits[0], splits[1]
		if !isValidMinute(minS) || !isValidMinute(maxS) {
			return
		}
		min, _ := strconv.Atoi(minS)
		max, _ := strconv.Atoi(maxS)
		if max < min {
			return
		}
	}
	if strings.Contains(minute, ",") {
		splits := strings.Split(minute, ",")
		for k, s := range splits {
			if !isValidMinute(s) {
				return
			}
			if k == 0 {
				continue
			}
			prev, _ := strconv.Atoi(splits[k-1])
			curr, _ := strconv.Atoi(s)
			if curr < prev {
				return
			}
		}
	}

	res = true
	return
}

func isValidHour(hour string) (res bool) {
	if len(hour) == 0 {
		return
	}
	if hour == "*" {
		res = true
		return
	}
	validRange := func(val int) bool {
		return !(val < 0 || val > 23)
	}
	numberOnly := regexp.MustCompile(`^[\d-]+$`)
	if len(hour) == 1 && !numberOnly.MatchString(hour) {
		return
	}
	if numberOnly.MatchString(hour) {
		val, _ := strconv.Atoi(hour)
		if !validRange(val) {
			return
		}
	}
	if strings.Contains(hour, "-") {
		splits := strings.Split(hour, "-")
		if len(splits) != 2 {
			return
		}
		minS, maxS := splits[0], splits[1]
		if !isValidHour(minS) || !isValidHour(maxS) {
			return
		}
		min, _ := strconv.Atoi(minS)
		max, _ := strconv.Atoi(maxS)
		if max < min {
			return
		}
	}
	if strings.Contains(hour, ",") {
		splits := strings.Split(hour, ",")
		for k, s := range splits {
			if !isValidHour(s) {
				return
			}
			if k == 0 {
				continue
			}
			prev, _ := strconv.Atoi(splits[k-1])
			curr, _ := strconv.Atoi(s)
			if curr < prev {
				return
			}
		}
	}

	res = true
	return
}
func isValidDate(date string) (res bool) {
	if len(date) == 0 {
		return
	}
	if date == "*" {
		res = true
		return
	}
	numberOnly := regexp.MustCompile(`^\d+$`)
	if len(date) == 1 && !numberOnly.MatchString(date) {
		return
	}
	if numberOnly.MatchString(date) {
		val, _ := strconv.Atoi(date)
		if val < 1 || val > 31 {
			return
		}
	}
	if strings.Contains(date, "-") {
		splits := strings.Split(date, "-")
		if len(splits) != 2 {
			return
		}
		minS, maxS := splits[0], splits[1]
		if !isValidDate(minS) || !isValidDate(maxS) {
			return
		}
		min, _ := strconv.Atoi(minS)
		max, _ := strconv.Atoi(maxS)
		if max < min {
			return
		}
	}
	if strings.Contains(date, ",") {
		splits := strings.Split(date, ",")
		for k, s := range splits {
			if !isValidDate(s) {
				return
			}
			if k == 0 {
				continue
			}
			prev, _ := strconv.Atoi(splits[k-1])
			curr, _ := strconv.Atoi(s)
			if curr < prev {
				return
			}
		}
	}

	res = true
	return
}
func isValidWeekdays(weekday string) (res bool) {
	if len(weekday) == 0 {
		return
	}
	if weekday == "*" {
		res = true
		return
	}
	validRange := func(week int) bool {
		return !(week < 0 || week > 6)
	}
	numberOnly := regexp.MustCompile(`^\d+$`)
	weekdayNamesOnly := regexp.MustCompile(`(?i)^([a-z]{3})$`)
	if len(weekday) == 1 && !numberOnly.MatchString(weekday) {
		return
	}
	if numberOnly.MatchString(weekday) {
		val, _ := getWeekFromShortName(weekday)
		if !validRange(val) {
			return
		}
	}
	if weekdayNamesOnly.MatchString(weekday) {
		// try get the number from the weekday name
		m, err := getWeekFromShortName(weekday)
		if err != nil {
			return
		}
		if !validRange(m) {
			return
		}
	}
	if strings.Contains(weekday, "-") {
		splits := strings.Split(weekday, "-")
		if len(splits) != 2 {
			return
		}
		minS, maxS := splits[0], splits[1]
		if !isValidWeekdays(minS) || !isValidWeekdays(maxS) {
			return
		}
		min, _ := getWeekFromShortName(minS)
		max, _ := getWeekFromShortName(maxS)
		if max < min {
			return
		}
	}
	if strings.Contains(weekday, ",") {
		splits := strings.Split(weekday, ",")
		for k, s := range splits {
			if !isValidWeekdays(s) {
				return
			}
			if k == 0 {
				continue
			}
			prev, _ := getWeekFromShortName(splits[k-1])
			curr, _ := getWeekFromShortName(s)
			if curr < prev {
				return
			}
		}
	}

	res = true
	return
}
func getMonthFromShortName(shortName string) (month int, err error) {
	if regexp.MustCompile(`^\d+$`).MatchString(shortName) {
		return strconv.Atoi(shortName)
	}
	d, err := time.Parse("Jan", cases.Title(language.Und, cases.NoLower).String(shortName))
	if err != nil {
		return
	}
	month = int(d.Month())
	return
}
func getWeekFromShortName(shortName string) (weekday int, err error) {
	if regexp.MustCompile(`^\d+$`).MatchString(shortName) {
		return strconv.Atoi(shortName)
	}
	// invalid week
	weekday = -1
	for k, w := range weeks {
		if w == strings.ToLower(shortName) {
			weekday = k
			break
		}
	}
	return
}
func isValidMonth(month string) (res bool) {
	if len(month) == 0 {
		return
	}
	if month == "*" {
		res = true
		return
	}
	numberOnly := regexp.MustCompile(`^\d+$`)
	monthNamesOnly := regexp.MustCompile(`(?i)^([a-z]{3})$`)
	if len(month) == 1 && !numberOnly.MatchString(month) {
		return
	}
	if numberOnly.MatchString(month) {
		val, _ := strconv.Atoi(month)
		if val < 1 || val > 12 {
			return
		}
	}
	if monthNamesOnly.MatchString(month) {
		// try get the number from the month name
		m, err := getMonthFromShortName(month)
		if err != nil {
			return
		}
		if m < 1 || m > 12 {
			return
		}
	}
	if strings.Contains(month, "-") {
		splits := strings.Split(month, "-")
		if len(splits) != 2 {
			return
		}
		minS, maxS := splits[0], splits[1]
		if !isValidMonth(minS) || !isValidMonth(maxS) {
			return
		}
		min, _ := getMonthFromShortName(minS)
		max, _ := getMonthFromShortName(maxS)
		if max < min {
			return
		}
	}
	if strings.Contains(month, ",") {
		splits := strings.Split(month, ",")
		for k, s := range splits {
			if !isValidMonth(s) {
				return
			}
			if k == 0 {
				continue
			}
			prev, _ := strconv.Atoi(splits[k-1])
			curr, _ := strconv.Atoi(s)
			if curr < prev {
				return
			}
		}
	}

	res = true
	return
}

// IsValidCron will verify if the notation is cron-compatible notation
func IsValidCron(notation string) (res bool) {
	splits := strings.Split(notation, " ")
	if len(splits) != 5 {
		return
	}
	return isValidMinute(splits[0]) &&
		isValidHour(splits[1]) &&
		isValidDate(splits[2]) &&
		isValidMonth(splits[3]) &&
		isValidWeekdays(splits[4])
}
