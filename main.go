package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type PersonActivity struct {
	Name          string
	DurationTotal time.Duration
	WorkingTime   time.Time
}

var (
	activities = make(map[string]PersonActivity)
	mu         sync.Mutex
)

func fmtDuration(d time.Duration) string {
    d = d.Round(time.Minute)
    h := d / time.Hour
    d -= h * time.Hour
    m := d / time.Minute
    return fmt.Sprintf("%02d:%02d", h, m)
}

func trackEND_Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	params := bytes.Split(body, []byte(","))
	if len(params) < 2 {
		http.Error(w, "Expected format: name,time_end", http.StatusBadRequest)
		return
	}

	name := string(params[0])
	endTime, err := time.Parse(time.RFC3339, string(params[1]))
	if err != nil {
		http.Error(w, "Invalid end time format", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if activity, exists := activities[name]; exists {
		if activity.WorkingTime.IsZero() {
			http.Error(w, "User did not start working", http.StatusBadRequest)
			return
		}
		duration := endTime.Sub(activity.WorkingTime)
		
		activity.DurationTotal += duration
		activity.WorkingTime = time.Time{}
		activities[name] = activity

		
		fmt.Fprintf(w, "Time tracked for %s: %v\n", name, fmtDuration(duration))
	} else {
		http.Error(w, "Activity does not exist", http.StatusBadRequest)
	}
}

func trackSTART_Handler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	params := bytes.Split(body, []byte(","))
	if len(params) < 2 {
		http.Error(w, "Expected format: name,start_time", http.StatusBadRequest)
		return
	}

	name := string(params[0])
	startTime, err := time.Parse(time.RFC3339, string(params[1]))
	if err != nil {
		http.Error(w, "Invalid start time format", http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if activity, exists := activities[name]; exists {
		if activity.WorkingTime.IsZero() {
			activity.WorkingTime = startTime
			activities[name] = activity
			fmt.Fprintf(w, "Work started for %s at %v\n", name, startTime)
		} else {
			http.Error(w, "User is already working", http.StatusBadRequest)
		}
	} else {
		activities[name] = PersonActivity{
			Name:        name,
			WorkingTime: startTime,
		}
		fmt.Fprintf(w, "Work started for %s at %v\n", name, startTime)
	}
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
		fmt.Fprintf(w, "%s has spent %v on the activity\n", activity.Name, activity.DurationTotal)
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
		fmt.Fprintf(w, "Name: %s, Time Spent: %v\n", name, activity.DurationTotal)
	}
}

func main() {
	http.HandleFunc("/track_start", trackSTART_Handler)
	http.HandleFunc("/track_end", trackEND_Handler)
	http.HandleFunc("/get", getTimeHandler)
	http.HandleFunc("/people", getAllPeopleHandler)

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}
