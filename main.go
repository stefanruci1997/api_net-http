package main

import (
	json "encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	_ "time"
)

func init() {
	// Initialize the database connection
	db = InitializePostgresDD()

}

func main() {
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			fmt.Println("Failed to get the database connection pool:", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			fmt.Println("Failed to close the database connection:", err)
		}
	}()

	router := mux.NewRouter()

	// Define routes using appropriate HTTP methods
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			// Handle GET request for root path
			response, _ := json.Marshal("hello API")
			_, err := w.Write(response)
			if err != nil {
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getUsers(w, r)
		} else if r.Method == http.MethodPost {
			createUser(w, r)

		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	router.HandleFunc("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method)

		if r.Method == http.MethodGet {
			getUser(w, r)
		} else if r.Method == http.MethodPut {
			updateUser(w, r)
		} else if r.Method == http.MethodDelete {
			deleteUser(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the server
	log.Println("Server starting on port 8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Fatal(err)
	}

}
