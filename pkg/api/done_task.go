package api

import (
	"net/http"
	"time"
	"go_final_project/pkg/db"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request, nextDateFunc func(time.Time, string, string) (string, error)) {
	id := r.URL.Query().Get("id")
	if id == "" {
		SendError(w, "Не указан идентификатор")
		return
	}

	
	task, err := db.GetTask(id)
	if err != nil {
		SendError(w, "Задача не найдена")
		return
	}

	
	if task.Repeat == "" {
		
		err = db.DeleteTask(id)
	} else {
		
		now := time.Now().Truncate(24 * time.Hour)
		next, err := nextDateFunc(now, task.Date, task.Repeat)
		if err != nil {
			SendError(w, "Ошибка вычисления следующей даты: "+err.Error())
			return
		}
		
		task.Date = next
		err = db.UpdateTask(task)
	}

	if err != nil {
		SendError(w, err.Error())
		return
	}

	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Write([]byte("{}"))
}