package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func GetLastCachedFile(cacheDir string) (string, error) {
	files, err := os.ReadDir(cacheDir)
	if err != nil {
		return "", err
	}

	var lastModTime time.Time
	var lastModFile string

	for _, file := range files {
		if !file.IsDir() {
			fileInfo, err := file.Info()
			if err != nil {
				return "", err
			}
			if fileInfo.ModTime().After(lastModTime) {
				lastModTime = fileInfo.ModTime()
				lastModFile = file.Name()
			}
		}
	}

	if lastModFile == "" {
		return "", fmt.Errorf("no files found in cache directory")
	}

	return filepath.Join(cacheDir, lastModFile), nil
}
