package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go_final_project/pkg/db"
)

const DateFormat = "20060102"

func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTask(w, r)
	case http.MethodPost:
		AddTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r)
	case http.MethodDelete:
		deleteTask(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		SendError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		SendError(w, "Задача не найдена", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(task)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		SendError(w, "Не указан идентификатор", http.StatusBadRequest)
		return
	}
	err := db.DeleteTask(id)
	if err != nil {
		SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse(DateFormat, nowStr)
	if err != nil {
		http.Error(w, "Некорректный параметр now", http.StatusBadRequest)
		return
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte(nextDate))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	start, err := time.Parse(DateFormat, dstart)
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
				return res.Format(DateFormat), nil
			}
		}
	case "d":
		if len(parts) < 2 {
			return "", errors.New("не указано количество дней для правила d")
		}
		days, err := strconv.Atoi(parts[1])
		if err != nil || days <= 0 {
			return "", errors.New("некорректный формат дней")
		}
		if days > 400 {
			return "", errors.New("превышен интервал в 400 дней")
		}
		res := start
		for {
			res = res.AddDate(0, 0, days)
			if res.After(now) {
				return res.Format(DateFormat), nil
			}
		}
	default:
		return "", fmt.Errorf("неподдерживаемый формат правила: %s", rule)
	}
}
