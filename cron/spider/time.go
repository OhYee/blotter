package spider

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	utime "github.com/OhYee/goutils/time"
)

type timeFinder struct {
	Regexp         *regexp.Regexp
	TimeFormat     string
	TimeFormatFunc func(string) *time.Time
}

var timeFinders = []timeFinder{
	{
		Regexp: regexp.MustCompile("^\\d{2,4}-\\d{1,2}-\\d{1,2}T\\d{1,2}:\\d{1,2}:\\d{1,2}(.\\d{3})*(Z([+\\-]\\d+)*)*$"),
		TimeFormatFunc: func(s string) *time.Time {
			datePart := s
			timeZone := 0
			withTimeZone := false

			// 分离时区
			lst := strings.Split(s, "Z")
			if len(lst) > 1 {
				datePart = lst[0]
				timeZonePart := lst[1]
				withTimeZone = true
				if len(timeZonePart) > 1 {
					timeZone = toInt(timeZonePart[1:])
					if timeZonePart[0] == '-' {
						timeZone = -timeZone
					}
				}
			}

			// 忽略毫秒
			lst = strings.Split(datePart, ".")
			if len(lst) > 1 {
				datePart = lst[0]
			}

			t := splitDateTime(datePart, "-", ":", "T")
			if t != nil && withTimeZone {
				tt := t.Add(time.Duration(8-timeZone) * time.Hour)
				return &tt
			}
			return t
		},
	},
	{
		Regexp:         regexp.MustCompile("^\\d{2,4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}$"),
		TimeFormatFunc: func(s string) *time.Time { return splitDateTime(s, "-", ":", " ") },
	},
	{
		Regexp:         regexp.MustCompile("^\\d{2,4}/\\d{1,2}/\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}$"),
		TimeFormatFunc: func(s string) *time.Time { return splitDateTime(s, "/", ":", " ") },
	},
	{
		Regexp:         regexp.MustCompile("^\\d{2,4}-\\d{1,2}-\\d{1,2}$"),
		TimeFormatFunc: func(s string) *time.Time { return splitDate(s, "-") },
	},
	{
		Regexp:         regexp.MustCompile("^\\d{2,4}/\\d{1,2}/\\d{1,2}$"),
		TimeFormatFunc: func(s string) *time.Time { return splitDate(s, "/") },
	},
	{
		Regexp:     regexp.MustCompile("^[a-zA-Z]{2,4} \\d{2}, \\d{4}$"),
		TimeFormat: "Jan 02, 2006",
	},
	{
		Regexp:     regexp.MustCompile("^\\d{2} \\d{2},\\d{4}$"),
		TimeFormat: "01 02,2006",
	},
}

func parseTime(s string) *time.Time {
	for _, r := range timeFinders {
		result := r.Regexp.FindAllString(s, -1)
		for _, timeString := range result {
			if r.TimeFormatFunc == nil {
				if t, err := time.ParseInLocation(r.TimeFormat, timeString, utime.ChinaTimeZone); err == nil {
					return &t
				}
			} else {
				if t := r.TimeFormatFunc(timeString); t != nil {
					return t
				}
			}
		}
	}
	return nil
}

func toUnix(t *time.Time) int64 {
	if t == nil {
		return 0
	}
	return t.Unix()
}

func toInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

func splitInts(s string, char string) (int, int, int) {
	slice := strings.Split(s, char)
	if len(slice) != 3 {
		return 0, 0, 0
	}
	return toInt(slice[0]), toInt(slice[1]), toInt(slice[2])
}

func splitDate(s string, char string) *time.Time {
	slice := strings.Split(s, char)
	if len(slice) != 3 {
		return nil
	}

	year, month, day := splitInts(s, char)
	if year == 0 || month < 1 || month > 12 || day < 1 || day > 31 {
		return nil
	}
	if year < 100 {
		year += 2000
	}

	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, utime.ChinaTimeZone)
	return &t
}

func splitDateTime(s string, char string, char2 string, char3 string) *time.Time {
	slice := strings.Split(s, char)
	if len(slice) != 3 {
		return nil
	}

	dateSlice := strings.Split(s, char3)
	if len(dateSlice) != 2 {
		return nil
	}

	year, month, day := splitInts(dateSlice[0], char)
	hour, minute, second := splitInts(dateSlice[1], char2)
	if year == 0 || month < 1 || month > 12 || day < 1 || day > 31 ||
		hour < 0 || hour > 23 || minute < 0 || minute > 59 || second < 0 || second > 59 {
		return nil
	}
	if year < 100 {
		year += 2000
	}

	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, utime.ChinaTimeZone)
	return &t
}
