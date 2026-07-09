package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"planner/pck/db"
	"planner/pck/repeat"
)

func writeJson(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func afterNow(now, t time.Time) bool {
	d1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	d2 := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	return d2.After(d1)
}

func checkDate(task *db.Task) error {
	now := time.Now()
	todayStr := now.Format(GlobalDateFormat)

	if task.Date == "" {
		task.Date = todayStr
	}

	t, err := time.Parse(GlobalDateFormat, task.Date)
	if err != nil {
		return errors.New("error date")
	}

	var next string
	if task.Repeat != "" {
		next, err = repeat.NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("error repeat " + err.Error())
		}
	}

	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = todayStr
		} else {
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "error JSON: " + err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "error title"})
		return
	}

	err = checkDate(&task)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": "error database: " + err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]int64{"id": id})
}
