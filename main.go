package main

import (
    "context"
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    _ "github.com/jackc/pgx/v4/stdlib"
)

// User represents the user structure
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

func main() {
    connStr := "postgres://postgres:pragal123@localhost:5432/contact_app"
    db, err := sql.Open("pgx", connStr)
    if err != nil {
        log.Fatal("Error connecting to the database:", err)
    }
    defer db.Close()

    ctx := context.Background()
    query := `INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`

    http.HandleFunc("/add-user", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*") // Allow all origins
        w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == http.MethodOptions {
            return // Exit early for preflight requests
        }

        if r.Method != http.MethodPost {
            http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
            return
        }

        var user User
        if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
            http.Error(w, "Invalid request payload", http.StatusBadRequest)
            return
        }

        var id int
        err = db.QueryRowContext(ctx, query, user.Name, user.Email).Scan(&id)
        if err != nil {
            log.Println("Error inserting row:", err)
            http.Error(w, "Error inserting user", http.StatusInternalServerError)
            return
        }

        response := map[string]int{"id": id}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    fmt.Println("Starting server on :8080...")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("Error starting server:", err)
    }
}
