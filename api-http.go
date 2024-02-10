package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"net/http"
	"strconv"
)

// getUsers retrieves all users
func getUsers(w http.ResponseWriter, r *http.Request) {
	var users []User
	err := RetrieveAll(&users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// getUser retrieves a user by ID
func getUser(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Path[len("/users/"):])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	var user User
	err = RetrieveByID(&user, uint(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// createUser creates a new user
func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Missing fields in request body", http.StatusBadRequest)
		return
	}

	err = AddUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(gin.H{"message": "User created successfully", "user": user})
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// updateUser updates an existing user
func updateUser(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len("/users/"):]
	var user User
	if err := db.First(&user, id).Error; err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.Save(&user).Error; err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(gin.H{"message": "User updated successfully", "user": user})
	if err != nil {
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// updateUserWithPatch updates an existing user with PATCH method
func updateUserWithPatch(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if GetDbInstance().Save(&user).Error != nil {
		http.Error(w, "User can't be updated", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User updated successfully")
}

// deleteUser deletes a user
func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/users/"):]
	var user User
	db.First(&user, id)
	db.Delete(&user)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User deleted successfully")
}
