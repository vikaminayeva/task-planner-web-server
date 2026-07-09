package api

import (
	"encoding/json"
	"net/http"

	"planner/pck/db"
)

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Missing task identifier"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, http.StatusNotFound, map[string]string{"error": "Task not found"})
		return
	}

	writeJson(w, http.StatusOK, task)
}

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "JSON deserialization error: " + err.Error()})
		return
	}

	if task.ID == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Missing task identifier"})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Task title is required"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err = db.UpdateTask(&task)
	if err != nil {
		writeJson(w, http.StatusNotFound, map[string]string{"error": "Task not found"})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
