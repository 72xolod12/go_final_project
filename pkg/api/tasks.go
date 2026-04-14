package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {

	tasks, err := db.Tasks(50)
	if err != nil {
		SendError(w, "Ошибка при получении списка задач")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err = json.NewEncoder(w).Encode(TasksResp{
		Tasks: tasks,
	})
	if err != nil {
		SendError(w, "Ошибка кодирования JSON")
		return
	}
}
