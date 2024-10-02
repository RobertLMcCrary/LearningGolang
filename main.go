package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var db *sql.DB

func main() {
	//initialize the database
	var err error
	db, err = sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create books table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
		id TEXT PRIMARY KEY,
		title TEXT,
		author TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}

	//initialize router
	r := mux.NewRouter()

	//routes
	r.HandleFunc("/api/books", getItems).Methods("GET")
	r.HandleFunc("/api/book/{id}", getItem).Methods("GET")
	r.HandleFunc("/api/books", createBook).Methods("POST")
	r.HandleFunc("/api/book/{id}", deleteBook).Methods("DELETE")

	// Start server
	fmt.Println("Server is starting on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}

//r is a pointer that represents incomming http requests (headers, method, etc.)

// get all items
func getItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.Query("SELECT id, title, author FROM books")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []Book
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

// get item by id
func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	var book Book
	err := db.QueryRow("SELECT id, title, author FROM books WHERE id = ?", params["id"]).Scan(&book.ID, &book.Title, &book.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			json.NewEncoder(w).Encode(&Book{})
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// create a book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	_, err := db.Exec("INSERT INTO books (id, title, author) VALUES (?, ?, ?)", book.ID, book.Title, book.Author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(book)
}

// delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	_, err := db.Exec("DELETE FROM books WHERE id = ?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
