package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var Library []Book

func main() {
	//initialize router
	r := mux.NewRouter()

	//adding some sample data
	Library = append(Library, Book{ID: "1", Title: "Harry Potter", Author: "JK Rowling"})

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
	json.NewEncoder(w).Encode(Library)
}

// get item by id
func getItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, book := range Library {
		if book.ID == params["id"] {
			json.NewEncoder(w).Encode(book)
			return
		}
	}

	json.NewEncoder(w).Encode(&Book{})
}

// create a book
func createBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	Library = append(Library, book)
	json.NewEncoder(w).Encode(book)
}

// delete a book
func deleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, book := range Library {
		if book.ID == params["id"] {
			Library = append(Library[:index], Library[:index+1]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
}
