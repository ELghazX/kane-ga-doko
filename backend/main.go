package main

import (
	"context"
	"fmt"
	"kane-ga-doko/database"
	"kane-ga-doko/handlers"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		claims := &handlers.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return handlers.JwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func seedAdmin() {
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		_, err := database.DB.Exec("INSERT INTO users (username, password_hash, is_admin) VALUES (?, ?, ?)", "admin", string(hashedPassword), true)
		if err != nil {
			log.Fatal("Failed to seed admin user:", err)
		}
		fmt.Println("Admin user created (username: admin, password: admin123)")
	}
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "database.sqlite"
	}
	database.InitDB(dbPath)
	defer database.DB.Close()

	seedAdmin()

	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("/api/login", handlers.Login)
	mux.HandleFunc("/api/register", handlers.Register)

	// Protected routes
	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			authMiddleware(handlers.CreateUser)(w, r)
		} else if r.Method == http.MethodGet {
			authMiddleware(handlers.ListUsers)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/users/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/approve") && r.Method == http.MethodPut {
			authMiddleware(handlers.ApproveUser)(w, r)
			return
		} else if r.Method == http.MethodPut {
			authMiddleware(handlers.UpdateUser)(w, r)
		} else if r.Method == http.MethodDelete {
			authMiddleware(handlers.DeleteUser)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/contacts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authMiddleware(handlers.GetContacts)(w, r)
		} else if r.Method == http.MethodPost {
			authMiddleware(handlers.AddContact)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/contacts/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			authMiddleware(handlers.RemoveContact)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/debts", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			authMiddleware(handlers.ListDebts)(w, r)
		} else if r.Method == http.MethodPost {
			authMiddleware(handlers.CreateDebt)(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// /api/debts/{id} or /api/debts/{id}/payments
	mux.HandleFunc("/api/debts/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/payments") && r.Method == http.MethodPost {
			authMiddleware(handlers.CreatePayment)(w, r)
			return
		} else if r.Method == http.MethodGet {
			authMiddleware(handlers.GetDebt)(w, r)
			return
		} else if r.Method == http.MethodDelete && !strings.Contains(r.URL.Path, "/payments") {
			authMiddleware(handlers.DeleteDebt)(w, r)
			return
		} else if r.Method == http.MethodPut && strings.Contains(r.URL.Path, "/payments/") {
			authMiddleware(handlers.UpdatePaymentStatus)(w, r)
			return
		}
		http.NotFound(w, r)
	})

	handler := corsMiddleware(mux)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
