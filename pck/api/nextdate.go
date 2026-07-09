package api

import (
	"net/http"
	"strings"
	"time"

	"planner/pck/repeat"
)

func nextDayHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")

	repeatStr := r.FormValue("repeat")

	repeatStr = strings.ReplaceAll(repeatStr, "+", " ")

	var nowTime time.Time
	var err error

	if nowStr == "" {
		nowTime = time.Now()
	} else {
		nowTime, err = time.Parse(GlobalDateFormat, nowStr)
		if err != nil {
			http.Error(w, "Invalid 'now' date format", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := repeat.NextDate(nowTime, dateStr, repeatStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(nextDate))
}
