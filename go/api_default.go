/*
 * Library
 *
 * Demo Library API
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

import (
	"encoding/json"
	"net/http"
    "strconv"
	"fmt"

	"github.com/gorilla/mux"
)

var nextBookId int32 = 4
var books = map[int32]Book{
	1: {
		ID: 1,
		Title: "Example Book",
		Description: "Example Book Description",
		Author: "John Doe",
	},
	2: {
		ID: 2,
		Title: "Awesome Book",
		Description: "Awesome Book Description",
		Author: "Mark Polo",
	},
	3: {
		ID: 3,
		Title: "Pirate Book",
		Description: "Pirate Book Description",
		Author: "Jack Sparrow",
	},
}

func BrowseBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Collect map values
	values := make([]Book, 0, len(books))
	for _, value := range books {
		values = append(values, value)
	}

    // Marshal the values to JSON and handle error
    if err := json.NewEncoder(w).Encode(values); err != nil {
        http.Error(w, "Failed to encode books to JSON", http.StatusInternalServerError)
        return
    }
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
    var newBook Book

    // Parse the JSON request body
    if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate the book data
    if err := validateBook(newBook); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Add the new book to the map
	newBook.ID = nextBookId
    books[newBook.ID] = newBook
	nextBookId = nextBookId + 1

	// Construct the URL for the newly created book
	baseUrl := r.URL.Scheme + "://" + r.URL.Host + "/v1/books/"
	location := fmt.Sprintf("%s%d", baseUrl, newBook.ID)

    // Respond with the created book and status
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.Header().Set("Location", location)
    w.WriteHeader(http.StatusCreated)
}

func validateBook(book Book) error {
    if book.Title == "" {
        return fmt.Errorf("Title must be provided")
    }
    if book.Description == "" {
        return fmt.Errorf("Description must be provided")
    }
    if book.Author == "" {
        return fmt.Errorf("Author must be provided")
    }
    return nil
}

func getBookByID(w http.ResponseWriter, r *http.Request) (Book, error) {
    // Parse the book ID from the URL
    idStr := mux.Vars(r)["id"]
    id, err := strconv.ParseInt(idStr, 10, 32)
    if err != nil {
        http.Error(w, "Invalid book ID", http.StatusBadRequest)
        return Book{}, err
    }

    // Check if the book exists
    book, exists := books[int32(id)]
    if !exists {
        http.Error(w, "Book not found", http.StatusNotFound)
        return Book{}, fmt.Errorf("book not found")
    }

    return book, nil
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
    book, err := getBookByID(w, r)
    if err != nil {
        return // Error already handled in getBookByID
    }

	delete(books, int32(book.ID))

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func PatchBook(w http.ResponseWriter, r *http.Request) {
    book, err := getBookByID(w, r)
    if err != nil {
        return // Error already handled in getBookByID
    }

    // Extract query parameters
    query := r.URL.Query()

    // Get the "author" parameter (optional)
    author := query.Get("author")
	if author != "" {
		book.Author = author
	}

    // Get the "title" parameter (optional)
    title := query.Get("title")
	if title != "" {
		book.Title = title
	}

    // Get the "description" parameter (optional)
    description := query.Get("description")
	if description != "" {
		book.Description = description
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func ReadBook(w http.ResponseWriter, r *http.Request) {
    book, err := getBookByID(w, r)
    if err != nil {
        return // Error already handled in getBookByID
    }

    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    if err := json.NewEncoder(w).Encode(book); err != nil {
        http.Error(w, "Failed to encode book to JSON", http.StatusInternalServerError)
        return
    }
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

    book, err := getBookByID(w, r)
    if err != nil {
        return // Error already handled in getBookByID
    }

    var newBook Book

    // Parse the JSON request body
    if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }

    // Validate the book data
    if err := validateBook(newBook); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Add the new book to the map
	newBook.ID = book.ID
    books[book.ID] = newBook

    // Respond with the created book and status
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusOK)
}
