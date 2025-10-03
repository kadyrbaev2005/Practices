package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// GET and POST handler for /user
func UserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		handleGetUser(w, r)
	case http.MethodPost:
		handleCreateUser(w, r)
	default:
		http.Error(w, `{"error": "method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
		return
	}

	resp := map[string]int{"user_id": id}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Name == "" {
		http.Error(w, `{"error":"invalid name"}`, http.StatusBadRequest)
		return
	}

	resp := map[string]string{"created": body.Name}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}