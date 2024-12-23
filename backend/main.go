package main

import (
    "database/sql"
    "encoding/json"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/gorilla/mux"
    _ "github.com/lib/pq"
)

type Todo struct {
    Id   int    `json:"id"`
    Todo string `json:"todo"`
}

func main() {
    db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    _, err = db.Exec(`CREATE TABLE IF NOT EXISTS todo (
        id   SERIAL PRIMARY KEY,
        todo TEXT
    )`)
    if err != nil {
        log.Fatal(err)
    }

    router := mux.NewRouter()

    // Routes
    router.HandleFunc("/api/go/todo", getTodos(db)).Methods("GET")
    router.HandleFunc("/api/go/todo", createTodo(db)).Methods("POST")
    router.HandleFunc("/api/go/todo/{id}", getTodo(db)).Methods("GET")
    router.HandleFunc("/api/go/todo/{id}", updateTodo(db)).Methods("PUT")
    router.HandleFunc("/api/go/todo/{id}", deleteTodo(db)).Methods("DELETE")

    enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))
    log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}

func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

func createTodo(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        var t Todo
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
            return
        }

        err := db.QueryRow("INSERT INTO todo (todo) VALUES ($1) RETURNING id", t.Todo).Scan(&t.Id)
        if err != nil {
            http.Error(w, `{"error": "Failed to create Todo"}`, http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(t)
    }
}

func getTodos(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        rows, err := db.Query("SELECT id, todo FROM todo")
        if err != nil {
            http.Error(w, `{"error": "Failed to retrieve Todos"}`, http.StatusInternalServerError)
            return
        }
        defer rows.Close()

        var todos []Todo
        for rows.Next() {
            var t Todo
            if err := rows.Scan(&t.Id, &t.Todo); err != nil {
                http.Error(w, `{"error": "Failed to scan Todo"}`, http.StatusInternalServerError)
                return
            }
            todos = append(todos, t)
        }
        if err := rows.Err(); err != nil {
            http.Error(w, `{"error": "Failed to iterate over Todos"}`, http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(todos)
    }
}

func getTodo(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := mux.Vars(r)["id"]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
            return
        }

        var t Todo
        err = db.QueryRow("SELECT id, todo FROM todo WHERE id = $1", id).Scan(&t.Id, &t.Todo)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, `{"error": "Todo not found"}`, http.StatusNotFound)
            } else {
                http.Error(w, `{"error": "Failed to retrieve Todo"}`, http.StatusInternalServerError)
            }
            return
        }

        json.NewEncoder(w).Encode(t)
    }
}

func updateTodo(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := mux.Vars(r)["id"]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
            return
        }

        var t Todo
        if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
            http.Error(w, `{"error": "Invalid request payload"}`, http.StatusBadRequest)
            return
        }

        _, err = db.Exec("UPDATE todo SET todo = $1 WHERE id = $2", t.Todo, id)
        if err != nil {
            http.Error(w, `{"error": "Failed to update Todo"}`, http.StatusInternalServerError)
            return
        }

        // Return the updated record
        var updatedTodo Todo
        err = db.QueryRow("SELECT id, todo FROM todo WHERE id = $1", id).Scan(&updatedTodo.Id, &updatedTodo.Todo)
        if err != nil {
            http.Error(w, `{"error": "Failed to retrieve updated Todo"}`, http.StatusInternalServerError)
            return
        }

        json.NewEncoder(w).Encode(updatedTodo)
    }
}

func deleteTodo(db *sql.DB) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        idStr := mux.Vars(r)["id"]
        id, err := strconv.Atoi(idStr)
        if err != nil {
            http.Error(w, `{"error": "Invalid ID"}`, http.StatusBadRequest)
            return
        }

        var t Todo
        err = db.QueryRow("SELECT id, todo FROM todo WHERE id = $1", id).Scan(&t.Id, &t.Todo)
        if err != nil {
            if err == sql.ErrNoRows {
                http.Error(w, `{"error": "Todo not found"}`, http.StatusNotFound)
            } else {
                http.Error(w, `{"error": "Error retrieving Todo"}`, http.StatusInternalServerError)
            }
            return
        }

        result, err := db.Exec("DELETE FROM todo WHERE id = $1", id)
        if err != nil {
            http.Error(w, `{"error": "Error deleting Todo"}`, http.StatusInternalServerError)
            return
        }

        rowsAffected, err := result.RowsAffected()
        if err != nil {
            http.Error(w, `{"error": "Unable to confirm deletion"}`, http.StatusInternalServerError)
            return
        }
        if rowsAffected == 0 {
            http.Error(w, `{"error": "No Todo was deleted"}`, http.StatusNotFound)
            return
        }

        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "Todo successfully deleted",
            "id":      t.Id,
        })
    }
}
