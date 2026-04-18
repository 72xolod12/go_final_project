package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
)

const DefaultTasksLimit = 50

type TasksResp struct {
	
	Tasks []*db.Task `json:"tasks"`
}

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	
	tasks, err := db.Tasks(DefaultTasksLimit)
	if err != nil {
		SendError(w, "Ошибка при получении списка задач", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	
	err = json.NewEncoder(w).Encode(TasksResp{
		Tasks: tasks,
	})
	if err != nil {
		SendError(w, "Ошибка кодирования JSON", http.StatusInternalServerError)
		return
	}
}