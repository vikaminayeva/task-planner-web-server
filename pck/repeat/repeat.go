package repeat

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func afterNow(date, now time.Time) bool {
	d1 := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return d1.After(d2)
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("missing repeat rule")
	}

	repeat = strings.ReplaceAll(repeat, "+", " ")

	t, err := time.Parse("20060102", dstart)
	if err != nil {
		return "", errors.New("invalid start date format")
	}

	parts := strings.Split(repeat, " ")

	if repeat == "y" {
		for {
			t = t.AddDate(1, 0, 0)
			if afterNow(t, now) {
				break
			}
		}
		return t.Format("20060102"), nil
	}

	if parts[0] == "d" {
		if len(parts) < 2 {
			return "", errors.New("missing interval in days")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("invalid days format")
		}

		if days < 1 || days > 400 {
			return "", errors.New("interval must be between 1 and 400 days")
		}

		for {
			t = t.AddDate(0, 0, days)
			if afterNow(t, now) {
				break
			}
		}
		return t.Format("20060102"), nil
	}

	return "", errors.New("unsupported repeat rule format")
}
