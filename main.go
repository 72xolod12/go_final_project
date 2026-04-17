package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

const dateFormat = "20060102"

func main() {
	// 1. Инициализация БД
	dbFile := "scheduler.db"
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	// 2. Статика (Фронтенд)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// 3. Публичные API (доступны без пароля)
	http.HandleFunc("/api/nextdate", nextDateHandler)
	http.HandleFunc("/api/signin", api.SigninHandler)

	// 4. Защищенные API (через AuthMiddleware)

	// Список задач
	http.HandleFunc("/api/tasks", api.AuthMiddleware(api.TasksHandler))

	// Выполнение задачи (Done)
	http.HandleFunc("/api/task/done", api.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.TaskDoneHandler(w, r, NextDate)
		} else {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}))

	// Основные действия с задачей (GET, POST, PUT, DELETE)
	// Оборачиваем весь большой обработчик в AuthMiddleware
	http.HandleFunc("/api/task", api.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id == "" {
				api.SendError(w, "Не указан идентификатор")
				return
			}
			task, err := db.GetTask(id)
			if err != nil {
				api.SendError(w, "Задача не найдена")
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			json.NewEncoder(w).Encode(task)

		case http.MethodPost:
			api.AddTaskHandler(w, r, NextDate)

		case http.MethodPut:
			api.UpdateTaskHandler(w, r, NextDate)

		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			if id == "" {
				api.SendError(w, "Не указан идентификатор")
				return
			}
			err := db.DeleteTask(id)
			if err != nil {
				api.SendError(w, err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write([]byte("{}"))

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}))

	// 5. Запуск сервера
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Сервер запущен на http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// Вспомогательный обработчик
func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.URL.Query().Get("now")
	dateStr := r.URL.Query().Get("date")
	repeat := r.URL.Query().Get("repeat")

	now, err := time.Parse("20060102", nowStr)
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
