package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var MediaFolder string = os.Getenv("MEDIA_FOLDER")

var SupportedVideoFormat = []string{".mp4"}

type FileEntry struct {
	Name     string `json:"name"`
	FilePath string `json:"filepath"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"isdir"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/api/videos/{filepath...}", listFiles)
	http.HandleFunc("/api/videos/stream/{filepath...}", serveVideo)
	http.ListenAndServe(":8090", nil)
}

func checkIsFile(path string) (bool, error) {
	if info, err := os.Stat(path); err != nil {
		return false, err
	} else {
		return !info.IsDir(), nil
	}
}

func listFiles(w http.ResponseWriter, req *http.Request) {
	filePath := req.PathValue("filepath")
	log.Printf("Listing filepath: %v", filePath)
	fullPath := filepath.Join(MediaFolder, filePath)
	files, err := os.ReadDir(fullPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unable to read directory: %q", filePath), http.StatusInternalServerError)
		return
	}

	var fileList []FileEntry
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			http.Error(w, fmt.Sprintf("Unable to read file: %q", file.Name()), http.StatusInternalServerError)
			return
		}

		fileList = append(fileList, FileEntry{
			Name:     fileInfo.Name(),
			Size:     fileInfo.Size(),
			FilePath: filepath.Join(filePath, fileInfo.Name()),
			IsDir:    fileInfo.IsDir(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileList)
}

func serveVideo(w http.ResponseWriter, req *http.Request) {
	filePath := req.PathValue("filepath")
	log.Printf("Serving filepath: %v", filePath)

	fullPath := filepath.Join(MediaFolder, filePath)
	file, err := os.Open(fullPath)
	if err != nil {
		http.NotFound(w, req)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Failed to retrieve file info", http.StatusInternalServerError)
		return
	}

	http.ServeContent(w, req, stat.Name(), stat.ModTime(), file)
}
