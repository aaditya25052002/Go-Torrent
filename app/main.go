package main

import (
	"embed"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/go-projects/go-torrent/torrent"
)

//go:embed static/index.html
var staticFS embed.FS

func main() {
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/download", handleDownload)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	index, err := staticFS.ReadFile("static/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(index)
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendJsonResponse(w, map[string]string{"error": "Method not allowed"}, http.StatusMethodNotAllowed)
		return
	}

	file, _, err := r.FormFile("torrent")
	if err != nil {
		sendJsonResponse(w, map[string]string{"error": "Bad request"}, http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		sendJsonResponse(w, map[string]string{"error": "Error reading torrent file"}, http.StatusInternalServerError)
		return
	}

	outputBuffer, outputFileName, err := torrent.Run(data)
	if err != nil {
		sendJsonResponse(w, map[string]string{"error": "Error running torrent"}, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", `attachment; filename="`+strings.ReplaceAll(outputFileName, `"`, `%22`)+`"`)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("X-Filename", outputFileName)
	w.Write(outputBuffer)
}

func sendJsonResponse(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
