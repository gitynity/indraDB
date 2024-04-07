package server

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/google/uuid"
)

// Storage represents a simple file-based storage.
type Storage struct {
	basePath string
	mu       sync.Mutex
}

// NewStorage creates a new instance of storage.
func NewStorage(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

// CreateCollection creates a new collection directory.
func (s *Storage) CreateCollection(collectionName string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	collectionPath := filepath.Join(s.basePath, collectionName)
	err := os.Mkdir(collectionPath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create collection: %v", err)
	}
	return nil
}

// CreateOrUpdateDocument creates or updates a document in the specified collection.
func (s *Storage) CreateOrUpdateDocument(collectionName, documentName string, data map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	collectionPath := filepath.Join(s.basePath, collectionName)
	documentPath := filepath.Join(collectionPath, documentName)

	file, err := os.OpenFile(documentPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create document file: %v", err)
	}
	defer file.Close()

	// Decode existing document data, if any
	existingData := make(map[string]interface{})
	if err := json.NewDecoder(file).Decode(&existingData); err != nil && err != io.EOF {
		return fmt.Errorf("failed to decode existing document data: %v", err)
	}
	if len(existingData) == 0 {
		existingData["uuid"], _ = uuid.NewUUID()
	}

	// Update specific key values in the existing document data with the new values provided
	for key, value := range data {
		existingData[key] = value
	}

	// Reset the file offset to the beginning
	if _, err := file.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to seek document file: %v", err)
	}

	// Truncate the file to remove existing content
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("failed to truncate document file: %v", err)
	}

	// Encode and write the updated document data to the file
	if err := json.NewEncoder(file).Encode(existingData); err != nil {
		return fmt.Errorf("failed to encode and write document data: %v", err)
	}

	return nil
}

// FilterDocuments returns documents from the specified collection that match the filter criteria.
func (s *Storage) FilterDocuments(collectionName string, filters map[string]interface{}) ([]map[string]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	collectionPath := filepath.Join(s.basePath, collectionName)
	dirEntries, err := os.ReadDir(collectionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection directory: %v", err)
	}

	var filteredDocuments []map[string]interface{}
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() {
			documentPath := filepath.Join(collectionPath, dirEntry.Name())
			file, err := os.Open(documentPath)
			if err != nil {
				return nil, fmt.Errorf("failed to open document file: %v", err)
			}
			defer file.Close()

			// Decode document data
			var documentData map[string]interface{}
			if err := json.NewDecoder(file).Decode(&documentData); err != nil {
				return nil, fmt.Errorf("failed to decode document data: %v", err)
			}

			// Check if the document matches the filter criteria
			match := true
			for key, value := range filters {
				if documentData[key] != value {
					match = false
					break
				}
			}
			if match {
				filteredDocuments = append(filteredDocuments, documentData)
			}
		}
	}

	return filteredDocuments, nil
}
