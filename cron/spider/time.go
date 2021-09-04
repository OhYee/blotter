package spider

import (
	"regexp"
	"time"

	utime "github.com/OhYee/goutils/time"
)

type timeFinder struct {
	Regexp     *regexp.Regexp
	TimeFormat string
}

var timeFinders = []timeFinder{
	{
		Regexp:     regexp.MustCompile("\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}"),
		TimeFormat: "2006-01-02 15:04:05",
	},
	{
		Regexp:     regexp.MustCompile("\\d{4}-\\d{1,2}-\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}"),
		TimeFormat: "2006-1-2 15:4:5",
	},
	{
		Regexp:     regexp.MustCompile("\\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2}"),
		TimeFormat: "2006/01/02 15:04:05",
	},
	{
		Regexp:     regexp.MustCompile("\\d{4}/\\d{1,2}/\\d{1,2} \\d{1,2}:\\d{1,2}:\\d{1,2}"),
		TimeFormat: "2006/1/2 15:4:5",
	},
	{
		Regexp:     regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}"),
		TimeFormat: "2006-01-02",
	},
	{
		Regexp:     regexp.MustCompile("\\d{4}-\\d{1,2}-\\d{1,2}"),
		TimeFormat: "2006-1-2",
	},
	{
		Regexp:     regexp.MustCompile("\\d{4}/\\d{2}/\\d{2}"),
		TimeFormat: "2006/01/02",
	},
	{
		Regexp:     regexp.MustCompile("\\d{4}/\\d{1,2}/\\d{1,2}"),
		TimeFormat: "2006/1/2",
	},
	{
		Regexp:     regexp.MustCompile("\\d{2}/\\d{1,2}/\\d{1,2}"),
		TimeFormat: "06/1/2",
	},
	{
		Regexp:     regexp.MustCompile("[a-zA-Z]{2,4} \\d{2}, \\d{4}"),
		TimeFormat: "Jan 02, 2006",
	},
	{
		Regexp:     regexp.MustCompile("\\d{2} \\d{2},\\d{4}"),
		TimeFormat: "01 02,2006",
	},
}

func parseTime(s string) *time.Time {
	for _, r := range timeFinders {
		result := r.Regexp.FindAllString(s, -1)
		for _, timeString := range result {
			t, err := time.ParseInLocation(r.TimeFormat, timeString, utime.ChinaTimeZone)
			// output.DebugOutput.Println(timeString, err)
			if err == nil {
				return &t
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
