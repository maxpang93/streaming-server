package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var MediaFolder string = os.Getenv("MEDIA_FOLDER")

var SupportedVideoFormat = []string{".mp4"}

type FileStat struct {
	Name     string `json:"name"`
	FilePath string `json:"filepath"`
	Size     int64  `json:"size"`
	IsDir    bool   `json:"isdir"`
}

func main() {
	// http.HandleFunc("/videos", listFiles)
	http.HandleFunc("/videos/{filepath...}", serveFolder)
	http.ListenAndServe(":8090", nil)
}

func _fp(path string) string {
	// sanitize file path
	return strings.TrimPrefix(path, MediaFolder)
}

func checkReadPerm(path string) error {
	if _, err := os.Open(path); err != nil {
		return fmt.Errorf("User has no read permission for file: %v", _fp(path))
	}
	return nil
}

func checkIsFile(path string) (bool, error) {
	if info, err := os.Stat(path); err != nil {
		return false, err
	} else {
		return !info.IsDir(), nil
	}
}

func listFiles(prefix string) ([]FileStat, error) {
	files, err := os.ReadDir(prefix)
	if err != nil {
		return nil, err
	}

	var fileList []FileStat
	for _, file := range files {
		fileInfo, err := file.Info()
		if err != nil {
			return nil, err
		}

		fileList = append(fileList, FileStat{
			Name:     fileInfo.Name(),
			Size:     fileInfo.Size(),
			FilePath: _fp(filepath.Join(prefix, fileInfo.Name())),
			IsDir:    fileInfo.IsDir(),
		})
	}
	return fileList, nil
}

func serveFolder(w http.ResponseWriter, req *http.Request) {
	log.Printf("media folder %v", MediaFolder)
	filePath := req.PathValue("filepath")
	log.Printf("filepath requested: %v", filePath)

	path := filepath.Join(MediaFolder, filePath)
	isFile, err := checkIsFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			http.NotFound(w, req)
			return
		}
		log.Printf("checkIsFile(%q): %v", path, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if isFile {
		log.Printf("Serving now: %v", path)
		return
	}

	fileList, err := listFiles(path)
	if err != nil {
		log.Printf("listFiles(%q): %v", path, err)
		http.Error(w, "Unable to read directory", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileList)
}
