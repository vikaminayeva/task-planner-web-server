package server

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"planner/pck/api"
)

func Start() {
	webDir := "web"
	path, _ := os.Executable()
	rootPath := filepath.Dir(path)
	webPath := filepath.Join(rootPath, webDir)

	_, err := os.Stat(webPath)
	if err != nil && os.IsNotExist(err) {
		webPath = webDir
	}

	mux := http.NewServeMux()

	api.InitMux(mux)

	fileServer := http.FileServer(http.Dir(webPath))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540"
	}

	log.Printf("Server start :%s...", port)

	err = http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatalf("Error serves: %v", err)
	}
}
