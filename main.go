package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"bytes"
	"io/ioutil"
)

type PersonActivity struct {
	Name     string
	Duration time.Duration
}

var (
	activities = make(map[string]PersonActivity)
	mu         sync.Mutex
)

func trackCommandHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	params := bytes.Split(body, []byte(","))
	if len(params) < 3 {
		http.Error(w, "Expected format: name,start,end", http.StatusBadRequest)
		return
	}

	name := string(params[0])
	startTime, err := time.Parse(time.RFC3339, string(params[1]))
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}
	endTime, err := time.Parse(time.RFC3339, string(params[2]))
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	duration := endTime.Sub(startTime)

	mu.Lock()
	defer mu.Unlock()

	if activity, exists := activities[name]; exists {
		activity.Duration += duration
		activities[name] = activity
	} else {
		activities[name] = PersonActivity{
			Name:     name,
			Duration: duration,
		}
	}

	fmt.Fprintf(w, "Time tracked for %s: %v\n", name, duration)
}

func getTimeHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")

	if name == "" {
		http.Error(w, "Missing 'name' parameter", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if activity, exists := activities[name]; exists {
		fmt.Fprintf(w, "%s has spent %v on the activity\n", activity.Name, activity.Duration)
	} else {
		http.Error(w, "Person not found", http.StatusNotFound)
	}
}

func getAllPeopleHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	if len(activities) == 0 {
		fmt.Fprintf(w, "No records found")
		return
	}

	for name, activity := range activities {
		fmt.Fprintf(w, "Name: %s, Time Spent: %v\n", name, activity.Duration)
	}
}

func main() {
	http.HandleFunc("/track", trackCommandHandler)
	http.HandleFunc("/get", getTimeHandler)
	http.HandleFunc("/people", getAllPeopleHandler)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
