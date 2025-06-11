package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestPostHandler(t *testing.T) {
	logFile, _ := os.Create("testlog.txt")

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	num := Number{42}
	body, _ := json.Marshal(num)
	req := httptest.NewRequest(http.MethodPost, "/data", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	postHandler(w, req, logFile)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
	}

	mutex.Lock()
	if len(numbers) != 1 || numbers[0] != 42 {
		t.Errorf("Expected %v, got %v", num.Value, numbers)
	}
	mutex.Unlock()
}

func TestGetHandler(t *testing.T) {
	logFile, _ := os.Create("testlog.txt")

	defer os.Remove(logFile.Name())
	defer logFile.Close()

	mutex.Lock()
	numbers = []int{1, 2, 3}
	mutex.Unlock()

	req := httptest.NewRequest(http.MethodGet, "/data", nil)
	w := httptest.NewRecorder()
	getHandler(w, req, logFile)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
	}

	var sum Sum
	json.Unmarshal(w.Body.Bytes(), &sum)

	if sum.Sum != 6 {
		t.Errorf("Expected sum 6, got %v", sum.Sum)
	}
}

func TestDeleteHandler(t *testing.T) {
	mutex.Lock()
	numbers = []int{1, 2, 3}
	mutex.Unlock()

	httptest.NewRequest(http.MethodDelete, "/data", nil)
	w := httptest.NewRecorder()
	deleteHandler(w)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %v, got %v", http.StatusOK, w.Code)
	}

	mutex.Lock()
	if len(numbers) != 0 {
		t.Errorf("Expected empty slice, got %v", numbers)
	}
	mutex.Unlock()
}
