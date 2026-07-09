package api

import (
	"net/http"
)

const GlobalDateFormat = "20060102"

func InitMux(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", nextDayHandler)
	mux.HandleFunc("/api/task", Auth(taskHandler))
	mux.HandleFunc("/api/tasks", Auth(tasksHandler))
	mux.HandleFunc("/api/task/done", Auth(taskDoneHandler))
	mux.HandleFunc("/api/signin", signinHandler)
}
