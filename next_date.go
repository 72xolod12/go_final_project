package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {

	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	start, err := time.Parse(dateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("некорректная дата начала: %v", err)
	}

	parts := strings.Split(repeat, " ")
	rule := parts[0]

	switch rule {
	case "y":

		res := start
		for {
			res = res.AddDate(1, 0, 0)

			if res.After(now) {
				return res.Format(dateFormat), nil
			}
		}

	case "d":

		if len(parts) < 2 {
			return "", errors.New("не указано количество дней для правила d")
		}

		days, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", errors.New("некорректный формат дней")
		}

		if days > 400 {
			return "", errors.New("превышен интервал в 400 дней")
		}

		res := start
		for {
			res = res.AddDate(0, 0, days)
			if res.After(now) {
				return res.Format(dateFormat), nil
			}
		}

	default:

		return "", fmt.Errorf("неподдерживаемый формат правила: %s", rule)
	}
}
