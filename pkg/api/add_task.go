package api

import (
	"encoding/json"
	"net/http"
	"time"

	"go_final_project/pkg/db"
)

func SendError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	resp := map[string]string{"error": message}
	json.NewEncoder(w).Encode(resp)
}

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		SendError(w, "Ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		SendError(w, "Не указан заголовок задачи", http.StatusBadRequest)
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	if task.Date == "" {

		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		SendError(w, "Дата представлена в неверном формате", http.StatusBadRequest)
		return
	}

	if t.Before(now) {
		if task.Repeat == "" {

			task.Date = now.Format(DateFormat)
		} else {
			next, err := NextDate(now, task.Date, task.Repeat)
			if err != nil {
				SendError(w, "Неверный формат правила повторения", http.StatusBadRequest)
				return
			}
			task.Date = next
		}
	}

	id, err := db.AddTask(&task)
	if err != nil {
		SendError(w, "Ошибка при добавлении задачи в БД", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": id})
}
