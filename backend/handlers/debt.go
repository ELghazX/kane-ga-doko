package handlers

import (
	"encoding/json"
	"kane-ga-doko/database"
	"kane-ga-doko/models"
	"net/http"
	"strconv"
	"strings"
)

func DeleteDebt(w http.ResponseWriter, r *http.Request) {
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
	debtID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid debt ID", http.StatusBadRequest)
		return
	}

	var creditorID int
	var paymentCount int
	err = database.DB.QueryRow("SELECT creditor_id, (SELECT COUNT(*) FROM payments WHERE debt_id = ?) FROM debts WHERE id = ?", debtID, debtID).Scan(&creditorID, &paymentCount)
	if err != nil {
		http.Error(w, "Debt not found", http.StatusNotFound)
		return
	}

	if claims.UserID != creditorID {
		http.Error(w, "Forbidden: only creditor can delete", http.StatusForbidden)
		return
	}

	if paymentCount > 0 {
		http.Error(w, "Cannot delete debt that has payments", http.StatusBadRequest)
		return
	}

	_, err = database.DB.Exec("DELETE FROM debts WHERE id = ?", debtID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Debt deleted successfully"})
}

func GetDebt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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
	debtID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid debt ID", http.StatusBadRequest)
		return
	}

	query := `
		SELECT 
			d.id, d.creditor_id, d.debtor_id, d.amount, d.description, d.created_at,
			c.username as creditor_name, 
			u.username as debtor_name,
			COALESCE(SUM(p.amount), 0) as total_paid
		FROM debts d
		JOIN users c ON d.creditor_id = c.id
		JOIN users u ON d.debtor_id = u.id
		LEFT JOIN payments p ON d.id = p.debt_id AND p.status = 'confirmed'
		WHERE d.id = ? AND (d.creditor_id = ? OR d.debtor_id = ?)
		GROUP BY d.id
	`
	var ds models.DebtSummary
	err = database.DB.QueryRow(query, debtID, claims.UserID, claims.UserID).Scan(
		&ds.ID, &ds.CreditorID, &ds.DebtorID, &ds.Amount, &ds.Description, &ds.CreatedAt,
		&ds.CreditorName, &ds.DebtorName, &ds.TotalPaid,
	)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	ds.Remaining = ds.Amount - ds.TotalPaid

	ds.Payments = []models.Payment{}
	payRows, err := database.DB.Query("SELECT id, amount, status, created_at FROM payments WHERE debt_id = ? ORDER BY created_at DESC", debtID)
	if err == nil {
		defer payRows.Close()
		for payRows.Next() {
			var p models.Payment
			if err := payRows.Scan(&p.ID, &p.Amount, &p.Status, &p.CreatedAt); err == nil {
				ds.Payments = append(ds.Payments, p)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ds)
}

func ListDebts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	query := `
		SELECT 
			d.id, d.creditor_id, d.debtor_id, d.amount, d.description, d.created_at,
			c.username as creditor_name, 
			u.username as debtor_name,
			COALESCE(SUM(p.amount), 0) as total_paid
		FROM debts d
		JOIN users c ON d.creditor_id = c.id
		JOIN users u ON d.debtor_id = u.id
		LEFT JOIN payments p ON d.id = p.debt_id AND p.status = 'confirmed'
		WHERE d.creditor_id = ? OR d.debtor_id = ?
		GROUP BY d.id
	`
	rows, err := database.DB.Query(query, claims.UserID, claims.UserID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var summaries []models.DebtSummary
	for rows.Next() {
		var ds models.DebtSummary
		if err := rows.Scan(
			&ds.ID, &ds.CreditorID, &ds.DebtorID, &ds.Amount, &ds.Description, &ds.CreatedAt,
			&ds.CreditorName, &ds.DebtorName, &ds.TotalPaid,
		); err != nil {
			continue
		}
		ds.Remaining = ds.Amount - ds.TotalPaid
		summaries = append(summaries, ds)
	}

	if summaries == nil {
		summaries = []models.DebtSummary{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

func CreateDebt(w http.ResponseWriter, r *http.Request) {
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
		DebtorID    int     `json:"debtor_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.DebtorID == claims.UserID {
		http.Error(w, "Cannot create debt for yourself", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(req.Description) == "" {
		http.Error(w, "Keterangan wajib diisi", http.StatusBadRequest)
		return
	}

	var isAdmin bool
	err := database.DB.QueryRow("SELECT is_admin FROM users WHERE id = ?", req.DebtorID).Scan(&isAdmin)
	if err != nil || isAdmin {
		http.Error(w, "Cannot create debt for an admin user", http.StatusBadRequest)
		return
	}

	result, err := database.DB.Exec("INSERT INTO debts (creditor_id, debtor_id, amount, description) VALUES (?, ?, ?, ?)",
		claims.UserID, req.DebtorID, req.Amount, req.Description)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Debt created successfully",
		"id":      id,
	})
}

func CreatePayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract ID from URL path (e.g. /api/debts/1/payments)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	debtIDStr := parts[3]
	debtID, err := strconv.Atoi(debtIDStr)
	if err != nil {
		http.Error(w, "Invalid debt ID", http.StatusBadRequest)
		return
	}

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Verify if the user is authorized to pay
	var creditorID, debtorID int
	err = database.DB.QueryRow("SELECT creditor_id, debtor_id FROM debts WHERE id = ?", debtID).Scan(&creditorID, &debtorID)
	if err != nil {
		http.Error(w, "Debt not found", http.StatusNotFound)
		return
	}

	status := "pending"
	if claims.UserID == creditorID {
		status = "confirmed"
	} else if claims.UserID == debtorID {
		status = "pending"
	} else {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	_, err = database.DB.Exec("INSERT INTO payments (debt_id, amount, status) VALUES (?, ?, ?)", debtID, req.Amount, status)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Payment recorded successfully"})
}

func UpdatePaymentStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	claims, ok := r.Context().Value("claims").(*Claims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Extract debtID and paymentID from URL path (e.g. /api/debts/1/payments/2)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 6 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	debtID, _ := strconv.Atoi(parts[3])
	paymentID, _ := strconv.Atoi(parts[5])

	var req struct {
		Status string `json:"status"` // "confirmed" or "rejected"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if req.Status != "confirmed" && req.Status != "rejected" {
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	var creditorID int
	err := database.DB.QueryRow("SELECT creditor_id FROM debts WHERE id = ?", debtID).Scan(&creditorID)
	if err != nil {
		http.Error(w, "Debt not found", http.StatusNotFound)
		return
	}

	if claims.UserID != creditorID {
		http.Error(w, "Forbidden: only creditor can update payment status", http.StatusForbidden)
		return
	}

	if req.Status == "rejected" {
		_, err = database.DB.Exec("DELETE FROM payments WHERE id = ? AND debt_id = ?", paymentID, debtID)
	} else {
		_, err = database.DB.Exec("UPDATE payments SET status = ? WHERE id = ? AND debt_id = ?", req.Status, paymentID, debtID)
	}
	
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Payment status updated successfully"})
}
