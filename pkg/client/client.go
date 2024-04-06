package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL string
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
	}
}

func (c *Client) CreateCollection(collectionName string) error {
	url := fmt.Sprintf("%s/collections/%s", c.BaseURL, collectionName)
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to create collection: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create collection: %s", resp.Status)
	}

	return nil
}

func (c *Client) CreateDocument(collectionName, documentName string, data interface{}) error {
	url := fmt.Sprintf("%s/document/%s/%s", c.BaseURL, collectionName, documentName)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal document data: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create document: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to create document: %s", resp.Status)
	}

	return nil
}

func (c *Client) GetDocument(collectionName, documentName string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/document/%s/%s", c.BaseURL, collectionName, documentName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get document: %s", resp.Status)
	}

	var documentData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&documentData); err != nil {
		return nil, fmt.Errorf("failed to decode document data: %v", err)
	}

	return documentData, nil
}

func (c *Client) UpdateDocument(collectionName, documentName string, data interface{}) error {
	url := fmt.Sprintf("%s/document/%s/%s", c.BaseURL, collectionName, documentName)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal document data: %v", err)
	}

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create PUT request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update document: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update document: %s", resp.Status)
	}

	return nil
}

func (c *Client) DeleteDocument(collectionName, documentName string) error {
	url := fmt.Sprintf("%s/document/%s/%s", c.BaseURL, collectionName, documentName)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create DELETE request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete document: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to delete document: %s", resp.Status)
	}

	return nil
}

func (c *Client) ListCollections() ([]string, error) {
	url := fmt.Sprintf("%s/collections", c.BaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list collections: %s", resp.Status)
	}

	var collections []string
	if err := json.NewDecoder(resp.Body).Decode(&collections); err != nil {
		return nil, fmt.Errorf("failed to decode collections: %v", err)
	}

	return collections, nil
}

func (c *Client) ListDocuments(collectionName string) ([]string, error) {
	url := fmt.Sprintf("%s/collections/%s", c.BaseURL, collectionName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to list documents: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list documents: %s", resp.Status)
	}

	var documents []string
	if err := json.NewDecoder(resp.Body).Decode(&documents); err != nil {
		return nil, fmt.Errorf("failed to decode documents: %v", err)
	}

	return documents, nil
}
