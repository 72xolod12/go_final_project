package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {

	dbFile := "scheduler.db"
	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}
	defer db.DB.Close()

	http.Handle("/", http.FileServer(http.Dir("./web")))

	http.HandleFunc("/api/tasks", api.TasksHandler)

	http.HandleFunc("/api/nextdate", nextDateHandler)

	http.HandleFunc("/api/task", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			id := r.URL.Query().Get("id")
			if id == "" {
				api.SendError(w, "Не указан идентификатор") // используем нашу функцию из пакета api
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

		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Сервер запущен на http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

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
