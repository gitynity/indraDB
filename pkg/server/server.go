package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

// Server represents the collection-based database server.
type Server struct {
	addr    string
	storage *Storage
	router  *mux.Router
}

// NewServer creates a new instance of the server.
func NewServer(addr string, basePath string) *Server {
	router := mux.NewRouter()
	server := &Server{
		addr:    addr,
		storage: NewStorage(basePath),
		router:  router,
	}
	server.routes()
	return server
}

// Start starts the server and listens for incoming requests.
func (s *Server) Start() error {
	fmt.Printf("Server is listening on %s\n", s.addr)
	return http.ListenAndServe(s.addr, s.router)
}

func (s *Server) routes() {
	s.router.HandleFunc("/collections", s.listCollections).Methods("GET")
	s.router.HandleFunc("/collections/{collectionName}", s.createCollection).Methods("POST")
	s.router.HandleFunc("/collections/{collectionName}", s.listDocuments).Methods("GET")
	s.router.HandleFunc("/document/{collectionName}/{documentName}", s.getDocument).Methods("GET")
	s.router.HandleFunc("/document/{collectionName}/{documentName}", s.createOrUpdateDocument).Methods("POST")
	s.router.HandleFunc("/document/{collectionName}/{documentName}", s.deleteDocument).Methods("DELETE")
	s.router.HandleFunc("/collections/{collectionName}", s.deleteCollection).Methods("DELETE")
	s.router.HandleFunc("/filterCollections/{collectionName}", s.filterDocuments).Methods("GET")

}

func (s *Server) listCollections(w http.ResponseWriter, r *http.Request) {
	// List all collections
	collections := make([]string, 0)
	files, err := os.ReadDir(s.storage.basePath)
	log.Print(s.storage.basePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			collections = append(collections, file.Name())
		}
	}
	jsonResponse(w, collections)
}

func (s *Server) createCollection(w http.ResponseWriter, r *http.Request) {
	// Create a new collection
	vars := mux.Vars(r)
	collectionName := vars["collectionName"]
	if collectionName == "" {
		http.Error(w, "Invalid collection name", http.StatusBadRequest)
		return
	}

	if err := s.storage.CreateCollection(collectionName); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": fmt.Sprintf("Collection %s created successfully", collectionName)})
}

func (s *Server) listDocuments(w http.ResponseWriter, r *http.Request) {
	// List all documents in a collection
	vars := mux.Vars(r)
	collectionName := vars["collectionName"]

	documents := make([]string, 0)
	files, err := os.ReadDir(filepath.Join(s.storage.basePath, collectionName))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	for _, file := range files {
		documents = append(documents, file.Name())
	}
	jsonResponse(w, documents)
}

func (s *Server) getDocument(w http.ResponseWriter, r *http.Request) {
	// Retrieve a specific document in a collection
	vars := mux.Vars(r)
	collectionName := vars["collectionName"]
	documentName := vars["documentName"]

	documentPath := filepath.Join(s.storage.basePath, collectionName, documentName)
	file, err := os.ReadFile(documentPath)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(file)
}

func (s *Server) createOrUpdateDocument(w http.ResponseWriter, r *http.Request) {
	// Create or update a document in a collection
	vars := mux.Vars(r)
	collectionName := vars["collectionName"]
	documentName := vars["documentName"]

	var documentData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&documentData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := s.storage.CreateOrUpdateDocument(collectionName, documentName, documentData); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponse(w, map[string]string{"message": "Document created/updated successfully"})
}

func (s *Server) deleteDocument(w http.ResponseWriter, r *http.Request) {
	// Delete a document in a collection
	vars := mux.Vars(r)
	collectionName := vars["collectionName"]
	documentName := vars["documentName"]

	documentPath := filepath.Join(s.storage.basePath, collectionName, documentName)
	if err := os.Remove(documentPath); err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, map[string]string{"message": "Document deleted successfully"})
}

func (s *Server) deleteCollection(w http.ResponseWriter, r *http.Request) {
	// Delete a collection
	vars := mux.Vars(r)
	collectionName := vars["collectionName"]

	collectionPath := filepath.Join(s.storage.basePath, collectionName)
	if err := os.RemoveAll(collectionPath); err != nil {
		http.Error(w, "Collection not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, map[string]string{"message": "Document deleted successfully"})
}

func (s *Server) filterDocuments(w http.ResponseWriter, r *http.Request) {
	// Extract collection name from the request URL
	collectionName := mux.Vars(r)["collectionName"]
	if collectionName == "" {
		http.Error(w, "Invalid collection name", http.StatusBadRequest)
		return
	}

	// Parse query parameters to extract filter criteria
	filters := make(map[string]interface{})
	queryParams := r.URL.Query()
	for key, values := range queryParams {
		if len(values) > 0 {
			filters[key] = values[0] // Use the first value if multiple values are present
		}
	}

	// Filter documents in the collection based on the provided criteria
	filteredDocuments, err := s.storage.FilterDocuments(collectionName, filters)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode filtered documents into JSON format and send the response
	jsonResponse(w, filteredDocuments)
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
