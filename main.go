package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Number struct {
	Value int `json:"value"`
}

type Sum struct {
	Sum int `json:"sum"`
}

var (
	numbers []int
	mutex   sync.Mutex
)

func main() {
	logFile, err := os.Create("./log.txt")
	if err != nil {
		fmt.Println("Log-file creation error: ", err)
		return
	}
	defer logFile.Close()

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			postHandler(w, r, logFile)

		case http.MethodGet:
			getHandler(w, r, logFile)

		case http.MethodDelete:
			deleteHandler(w)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "Method not allowed")
		}
	})

	err = http.ListenAndServe(":4545", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func postHandler(w http.ResponseWriter, r *http.Request, logFile *os.File) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Wrong Content-Type", http.StatusUnsupportedMediaType)

		date := time.Now().Format(time.DateTime)
		logFile.WriteString(date + " " + r.Method + " Wrong Content-Type\n")

		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)

		date := time.Now().Format(time.DateTime)
		logFile.WriteString(date + " " + r.Method + " Error reading request body\n")

		return
	}
	defer r.Body.Close()

	var num Number
	err = json.Unmarshal(data, &num)
	if err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)

		date := time.Now().Format(time.DateTime)
		logFile.WriteString(date + " " + r.Method + " Invalid JSON format\n")

		return
	}

	mutex.Lock()
	numbers = append(numbers, num.Value)
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
}

func getHandler(w http.ResponseWriter, r *http.Request, logFile *os.File) {
	mutex.Lock()
	defer mutex.Unlock()

	sum := Sum{0}
	for _, value := range numbers {
		sum.Sum += value
	}

	data, err := json.Marshal(sum)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)

		date := time.Now().Format(time.DateTime)
		logFile.WriteString(date + " " + r.Method + " Error encoding response\n")

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func deleteHandler(w http.ResponseWriter) {
	mutex.Lock()
	numbers = []int{}
	mutex.Unlock()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Data has been deleted")
}
