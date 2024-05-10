package timeutil

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var (
	valueRegex = regexp.MustCompile(`^\d+`)
	unitRegex  = regexp.MustCompile(`\w$`)
)

// ParseDuration attempts to parse a duration string. If successful, it returns a time.Duration, otherwise returning
// an error in its second value.
func ParseDuration(str string) (time.Duration, error) {
	matches := valueRegex.FindAllString(str, 1)
	if len(matches) <= 0 {
		return 0, fmt.Errorf("parse duration: invalid value")
	}
	value, err := strconv.Atoi(matches[0])
	if err != nil {
		return 0, err
	}
	if value <= 0 {
		return 0, fmt.Errorf("parse duration: value must be greater than zero")
	}
	unit := unitRegex.FindAllString(str, 1)[0]
	if duration, ok := map[string]time.Duration{
		"s": time.Second,
		"m": time.Minute,
		"h": time.Hour,
		"d": time.Hour * 24,
	}[unit]; ok {
		return time.Duration(value) * duration, nil
	}
	return 0, fmt.Errorf("parse duration: invalid unit")
}

// FormatDuration formats a duration to a string.
func FormatDuration(d time.Duration) string {
	if d == 0 {
		return "0s"
	}
	var str string
	if d.Hours() >= 24 {
		str += fmt.Sprintf("%dd", int(d.Hours()/24))
		d -= time.Hour * 24 * time.Duration(int(d.Hours()/24))
	}
	if d.Hours() >= 1 {
		str += fmt.Sprintf("%dh", int(d.Hours()))
		d -= time.Hour * time.Duration(int(d.Hours()))
	}
	if d.Minutes() >= 1 {
		str += fmt.Sprintf("%dm", int(d.Minutes()))
		d -= time.Minute * time.Duration(int(d.Minutes()))
	}
	if d.Seconds() >= 1 {
		str += fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return str
}
