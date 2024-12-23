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
	router.HandleFunc("/api/go/todo", getTodos(db)).Methods("GET")
	router.HandleFunc("/api/go/todo", createTodo(db)).Methods("POST")
	router.HandleFunc("/api/go/todo/{id}", getTodo(db)).Methods("GET")

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
//get all Todos
func getTodos(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM todo")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		todo := []Todo{} //array of Todos
		for rows.Next() {
			var t Todo
			if err := rows.Scan(&t.Id, &t.Todo); err != nil {
				log.Fatal(err)
			}
			todo = append(todo, t)
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		json.NewEncoder(w).Encode(todo)
	}
}


func getTodo(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]

		var t Todo
		err := db.QueryRow("SELECT * FROM todo WHERE id = $1", id).Scan(&t.Id,&t.Todo)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(t)
	}
}

