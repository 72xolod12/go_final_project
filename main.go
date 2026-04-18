package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
)

func main() {
	api.InitAuthConfig()
	// 1. Инициализация БД
	dbFile := "scheduler.db"
	if err := db.Init(dbFile); err != nil {

		log.Printf("Ошибка инициализации БД: %v", err)
		return
	}

	defer db.DB.Close()

	// 2. Статика (Фронтенд)
	http.Handle("/", http.FileServer(http.Dir("./web")))

	// 3. API Маршруты
	http.HandleFunc("/api/nextdate", api.NextDateHandler)
	http.HandleFunc("/api/signin", api.SigninHandler)

	// Список задач
	http.HandleFunc("/api/tasks", api.AuthMiddleware(api.TasksHandler))

	// Выполнение задачи (Done)
	http.HandleFunc("/api/task/done", api.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			api.TaskDoneHandler(w, r, api.NextDate)
		} else {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	}))

	// Работа с конкретной задачей (CRUD)
	http.HandleFunc("/api/task", api.AuthMiddleware(api.TaskHandler))

	// 4. Запуск сервера
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Сервер запущен на http://localhost:%s", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Printf("Сервер завершил работу с ошибкой: %v", err)

		return
	}
}
