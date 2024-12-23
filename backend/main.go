package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type Todo struct {
    Id    int    `json:"id"`
    Todo  string `json:"todo"`
}

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS todo (id SERIAL PRIMARY KEY, todo TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	//create Router
	router := mux.NewRouter()
	router.HandleFunc("/api/go/todo", createTodo(db)).Methods("POST")

	enhancedRouter := enableCORS(jsonContentTypeMiddleware(router))
	log.Fatal(http.ListenAndServe(":8000", enhancedRouter))
}


func enableCORS(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set CORS headers
        w.Header().Set("Access-Control-Allow-Origin", "*") // Allow any origin
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

        // Check if the request is for CORS preflight
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        // Pass down the request to the next middleware (or final handler)
        next.ServeHTTP(w, r)
    })
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Set JSON Content-Type
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

func createTodo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Todo
		json.NewDecoder(r.Body).Decode(&t)
		err := db.QueryRow("INSERT INTO todo (todo) VALUES ($1) RETURNING id", t.Todo).Scan(&t.Id)
		if err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(t)

	}
}

