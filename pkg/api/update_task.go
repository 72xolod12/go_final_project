package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request, nextDateFunc func(time.Time, string, string) (string, error)) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		SendError(w, "Ошибка десериализации JSON")
		return
	}

	// Валидация
	if task.ID == "" {
		SendError(w, "Не указан идентификатор задачи")
		return
	}
	if task.Title == "" {
		SendError(w, "Не указан заголовок задачи")
		return
	}

	now := time.Now().Truncate(24 * time.Hour)
	t, err := time.Parse("20060102", task.Date)
	if err != nil {
		SendError(w, "Дата представлена в неверном формате")
		return
	}

	if t.Before(now) {
		if task.Repeat == "" {
			task.Date = now.Format("20060102")
		} else {
			next, err := nextDateFunc(now, task.Date, task.Repeat)
			if err != nil {
				SendError(w, "Неверный формат правила повторения")
				return
			}
			task.Date = next
		}
	}

	if err := db.UpdateTask(task); err != nil {
		SendError(w, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}
