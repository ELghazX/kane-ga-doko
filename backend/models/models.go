package models

import "time"

type User struct {
	ID           int    `json:"id"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
	IsAdmin      bool   `json:"is_admin"`
	IsApproved   bool   `json:"is_approved"`
}

type Debt struct {
	ID          int       `json:"id"`
	CreditorID  int       `json:"creditor_id"`
	DebtorID    int       `json:"debtor_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Payment struct {
	ID        int       `json:"id"`
	DebtID    int       `json:"debt_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type DebtSummary struct {
	Debt
	CreditorName string    `json:"creditor_name"`
	DebtorName   string    `json:"debtor_name"`
	TotalPaid    float64   `json:"total_paid"`
	Remaining    float64   `json:"remaining"`
	Payments     []Payment `json:"payments"`
}
