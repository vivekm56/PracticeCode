package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Signup struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type SessionID struct {
	SID string `json:"id"`
}
type Note struct {
	ID   int32  `json:"id"`
	Note string `json:"note"`
}

// var users map[string]User

// var accountsLogin []Login
var (
	accounts    map[string]Signup
	userSID     SessionID
	notes       []Note
	noteCounter int32
)

func main() {
	accounts = make(map[string]Signup) // Initialize the map
	notes = make([]Note, 0)
	noteCounter = 0
	// Initialize the Gorilla Mux router
	router := mux.NewRouter()

	// Define the API endpoints
	router.HandleFunc("/signup", createSignup).Methods("POST")
	router.HandleFunc("/login", createLogin).Methods("POST")
	router.HandleFunc("/notes", createNotes).Methods("POST")

	// Start the server on localhost:8000
	fmt.Println("Server listening on http://localhost:8000/")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func createSignup(w http.ResponseWriter, r *http.Request) {

	var newSignup Signup
	err := json.NewDecoder(r.Body).Decode(&newSignup)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := accounts[newSignup.Email]; ok {

		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if newSignup.Name == "" || newSignup.Email == "" || newSignup.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		accounts[newSignup.Email] = newSignup
	}
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")
	// Return the created resource in the response
	w.WriteHeader(http.StatusOK)
	// json.NewEncoder(w).Encode(newSignup)
}

func createLogin(w http.ResponseWriter, r *http.Request) {

	var newLogin Login
	err := json.NewDecoder(r.Body).Decode(&newLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if storedUser, ok := accounts[newLogin.Email]; ok {
		if storedUser.Password == newLogin.Password {
			// Convert the input string to bytes
			inputVlues := newLogin.Email + storedUser.Name
			inputBytes := []byte(inputVlues)
			// Calculate the MD5 hash
			hash := md5.Sum(inputBytes)
			// Convert the hash to a hexadecimal string
			userSID.SID = string(hex.EncodeToString(hash[:]))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(userSID.SID)
		} else {
			// http.Error(w, err.Error(), http.StatusUnauthorized)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	} else if newLogin.Email == "" || newLogin.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		// http.Error(w, err.Error(), http.StatusUnauthorized)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
}

func createNotes(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		SessionID string `json:"sid"`
		Note      string `json:"note"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if requestData.SessionID != userSID.SID {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	note := Note{
		ID:   getNextNoteID(),
		Note: requestData.Note,
	}

	notes = append(notes, note)

	responseData := struct {
		ID int32 `json:"id"`
	}{
		ID: note.ID,
	}

	response, err := json.Marshal(responseData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func getNextNoteID() int32 {
	noteCounter++
	return noteCounter
}
