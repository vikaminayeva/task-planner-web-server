package api

import (
	"net/http"

	"planner/pck/db"
)

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Missing task identifier"})
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeJson(w, http.StatusNotFound, map[string]string{"error": "Task not found"})
		return
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
