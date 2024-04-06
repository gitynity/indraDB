package server

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
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
func (s *Storage) CreateOrUpdateDocument(collectionName, documentName string, data interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	collectionPath := filepath.Join(s.basePath, collectionName)
	documentPath := filepath.Join(collectionPath, documentName)

	file, err := os.Create(documentPath)
	if err != nil {
		return fmt.Errorf("failed to create or update document: %v", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(data); err != nil {
		return fmt.Errorf("failed to encode document data: %v", err)
	}

	return nil
}

