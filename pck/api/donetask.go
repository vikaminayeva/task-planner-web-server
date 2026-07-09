package api

import (
	"net/http"
	"time"

	"planner/pck/db"
	"planner/pck/repeat"
)

func taskDoneHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}

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

	if task.Repeat == "" {
		err = db.DeleteTask(id)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	} else {
		nextDate, err := repeat.NextDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "Invalid repeat format: " + err.Error()})
			return
		}

		err = db.UpdateDate(nextDate, id)
		if err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	writeJson(w, http.StatusOK, map[string]string{})
}
