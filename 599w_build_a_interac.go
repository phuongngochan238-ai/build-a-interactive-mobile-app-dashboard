package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// Dashboard represents the interactive mobile app dashboard
type Dashboard struct {
	Username string `json:"username"`
	Devices  []string `json:"devices"`
	Tasks    []string `json:"tasks"`
}

var (
	store = sessions.NewCookieStore([]byte("something-very-secret"))
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/login", loginHandler).Methods("POST")
	router.HandleFunc("/dashboard", dashboardHandler).Methods("GET")
	router.HandleFunc("/addDevice", addDeviceHandler).Methods("POST")
	router.HandleFunc("/addTask", addTaskHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "admin" && password == "password" {
		session.Values["username"] = username
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusFound)
		return
	}

	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username, ok := session.Values["username"]
	if !ok {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	dashboard := Dashboard{
		Username: username.(string),
		Devices:  []string{"Device 1", "Device 2"},
		Tasks:    []string{"Task 1", "Task 2"},
	}

	json.NewEncoder(w).Encode(dashboard)
}

func addDeviceHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username, ok := session.Values["username"]
	if !ok {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	device := r.FormValue("device")
	dashboard := Dashboard{
		Username: username.(string),
		Devices:  append([]string{"Device 1", "Device 2"}, device),
		Tasks:    []string{"Task 1", "Task 2"},
	}

	json.NewEncoder(w).Encode(dashboard)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	username, ok := session.Values["username"]
	if !ok {
		http.Error(w, "Not authorized", http.StatusUnauthorized)
		return
	}

	task := r.FormValue("task")
	dashboard := Dashboard{
		Username: username.(string),
		Devices:  []string{"Device 1", "Device 2"},
		Tasks:    append([]string{"Task 1", "Task 2"}, task),
	}

	json.NewEncoder(w).Encode(dashboard)
}