package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gitynity/indraDB/pkg/client"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "indra-cli",
		Short: "CLI tool for interacting with the indra database server",
	}

	var baseURL string

	var createCollectionCmd = &cobra.Command{
		Use:   "create-collection <collectionName>",
		Short: "Create a new collection",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c := client.NewClient(baseURL)
			collectionName := args[0]
			if err := c.CreateCollection(collectionName); err != nil {
				log.Fatalf("Failed to create collection: %v", err)
			}
			log.Printf("Collection %s created successfully", collectionName)
		},
	}

	var createDocumentCmd = &cobra.Command{
		Use:   "create-document <collectionName> <documentName> <jsonData>",
		Short: "Create a new document in the specified collection",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			c := client.NewClient(baseURL)
			collectionName := args[0]
			documentName := args[1]
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(args[2]), &jsonData); err != nil {
				log.Fatalf("Failed to parse JSON data: %v", err)
			}
			if err := c.CreateDocument(collectionName, documentName, jsonData); err != nil {
				log.Fatalf("Failed to create document: %v", err)
			}
			log.Printf("Document %s created successfully in collection %s", documentName, collectionName)
		},
	}

	var updateDocumentCmd = &cobra.Command{
		Use:   "update-document <collectionName> <documentName> <jsonData>",
		Short: "Update an existing document in the specified collection",
		Args:  cobra.ExactArgs(3),
		Run: func(cmd *cobra.Command, args []string) {
			c := client.NewClient(baseURL)
			collectionName := args[0]
			documentName := args[1]
			var jsonData map[string]interface{}
			if err := json.Unmarshal([]byte(args[2]), &jsonData); err != nil {
				log.Fatalf("Failed to parse JSON data: %v", err)
			}
			if err := c.UpdateDocument(collectionName, documentName, jsonData); err != nil {
				log.Fatalf("Failed to update document: %v", err)
			}
			log.Printf("Document %s updated successfully in collection %s", documentName, collectionName)
		},
	}

	var deleteDocumentCmd = &cobra.Command{
		Use:   "delete-document <collectionName> <documentName>",
		Short: "Delete a document from the specified collection",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			c := client.NewClient(baseURL)
			collectionName := args[0]
			documentName := args[1]
			if err := c.DeleteDocument(collectionName, documentName); err != nil {
				log.Fatalf("Failed to delete document: %v", err)
			}
			log.Printf("Document %s deleted successfully from collection %s", documentName, collectionName)
		},
	}

	var listCollectionsCmd = &cobra.Command{
		Use:   "list-collections",
		Short: "List all collections",
		Run: func(cmd *cobra.Command, args []string) {
			c := client.NewClient(baseURL)
			collections, err := c.ListCollections()
			if err != nil {
				log.Fatalf("Failed to list collections: %v", err)
			}
			fmt.Println("Collections:")
			for _, collection := range collections {
				fmt.Println(collection)
			}
		},
	}

	var listDocumentsCmd = &cobra.Command{
		Use:   "list-documents <collectionName>",
		Short: "List all documents in the specified collection",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			c := client.NewClient(baseURL)
			collectionName := args[0]
			documents, err := c.ListDocuments(collectionName)
			if err != nil {
				log.Fatalf("Failed to list documents: %v", err)
			}
			fmt.Printf("Documents in collection %s:\n", collectionName)
			for _, document := range documents {
				fmt.Println(document)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&baseURL, "server", "s", "http://localhost:8080", "Server URL")
	rootCmd.AddCommand(
		createCollectionCmd,
		createDocumentCmd,
		updateDocumentCmd,
		deleteDocumentCmd,
		listCollectionsCmd,
		listDocumentsCmd,
	)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
