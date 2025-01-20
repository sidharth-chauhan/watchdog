package utils

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGetLastCachedFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "cache")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	file3 := filepath.Join(tmpDir, "file3.txt")

	createFileWithModTime(t, file1, time.Now().Add(-2*time.Hour))
	createFileWithModTime(t, file2, time.Now().Add(-1*time.Hour))
	createFileWithModTime(t, file3, time.Now())

	lastFile, err := GetLastCachedFile(tmpDir)
	if err != nil {
		t.Fatalf("GetLastCachedFile failed: %v", err)
	}

	expectedFile := file3
	if lastFile != expectedFile {
		t.Errorf("Expected last file to be %s, got %s", expectedFile, lastFile)
	}
}

func createFileWithModTime(t *testing.T, path string, modTime time.Time) {
	t.Helper()

	file, err := os.Create(path)
	if err != nil {
		t.Fatalf("Failed to create file %s: %v", path, err)
	}
	defer file.Close()

	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("Failed to set modification time for file %s: %v", path, err)
	}
}

func TestDownloadGTFSBundle(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "gtfs-bundle-*.zip")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tmpFilePath := tmpFile.Name()
	defer os.Remove(tmpFilePath)

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("mock GTFS data"))
	}))
	defer mockServer.Close()

	err = DownloadGTFSBundle(mockServer.URL, tmpFilePath)
	if err != nil {
		t.Fatalf("DownloadGTFSBundle failed: %v", err)
	}

	fileContent, err := os.ReadFile(tmpFilePath)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	expectedContent := "mock GTFS data"
	if string(fileContent) != expectedContent {
		t.Errorf("Expected file content to be %s, got %s", expectedContent, string(fileContent))
	}
}
