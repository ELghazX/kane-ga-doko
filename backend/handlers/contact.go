package handlers

import (
	"encoding/json"
	"kane-ga-doko/database"
	"net/http"
	"strings"
)

type ContactUser struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func AddContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var contactId int
	err := database.DB.QueryRow("SELECT id FROM users WHERE username = ?", req.Username).Scan(&contactId)
	if err != nil {
		http.Error(w, "Username tidak ditemukan", http.StatusNotFound)
		return
	}

	if contactId == claims.UserID {
		http.Error(w, "Tidak bisa menambahkan diri sendiri", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("INSERT INTO contacts (user_id, contact_id) VALUES (?, ?)", claims.UserID, contactId)
	if err != nil {
		http.Error(w, "Pengguna sudah ada di kontak Anda", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kontak berhasil ditambahkan"})
}

func GetContacts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rows, err := database.DB.Query(`
		SELECT u.id, u.username 
		FROM contacts c
		JOIN users u ON c.contact_id = u.id
		WHERE c.user_id = ?
		ORDER BY u.username ASC
	`, claims.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	contacts := []ContactUser{}
	for rows.Next() {
		var c ContactUser
		if err := rows.Scan(&c.ID, &c.Username); err != nil {
			continue
		}
		contacts = append(contacts, c)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(contacts)
}

func RemoveContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	contactId := parts[3]

	_, err := database.DB.Exec("DELETE FROM contacts WHERE user_id = ? AND contact_id = ?", claims.UserID, contactId)
	if err != nil {
		http.Error(w, "Gagal menghapus kontak", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Kontak berhasil dihapus"})
}
