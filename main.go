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
	info, err := os.Stat(path)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return false, fmt.Errorf("File does not exists: %v. Error: %v", path, err)
	}
	return !info.IsDir(), nil
}

func listFiles(prefix string) []FileStat {
	files, err := os.ReadDir(prefix)
	if err != nil {
		log.Fatal(err)
	}

	var fileList []FileStat
	for _, file := range files {
		filePath := filepath.Join(prefix, file.Name())
		fmt.Println(filePath)
		if err := checkReadPerm(filePath); err != nil {
			log.Printf("Skipping.. Reason: %v", err)
			continue
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil {
			log.Fatal(err)
		}

		fileList = append(fileList, FileStat{
			Name:     fileInfo.Name(),
			Size:     fileInfo.Size(),
			FilePath: _fp(filePath),
			IsDir:    fileInfo.IsDir(),
		})
	}
	return fileList
}

func serveFolder(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("media folder %v\n", MediaFolder)
	filePath := req.PathValue("filepath")
	fmt.Printf("filepath requested: %v\n", filePath)

	path := filepath.Join(MediaFolder, filePath)
	isFile, err := checkIsFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if isFile {
		fmt.Printf("Serving now: %v", path)
		return
	}

	fileList := listFiles(path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(fileList)
}
