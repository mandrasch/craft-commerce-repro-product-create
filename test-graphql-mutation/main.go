package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

const (
	apiURL       = "https://craft-commerce-repro-product-create.ddev.site/graphql" // Replace with your CraftCMS GraphQL endpoint
	bearerToken  = "CDwvk4kPKUBTBWsivteRLgd7DtIzcEa6"                              // Replace with your actual bearer token
	numWorkers   = 3                                                              // Number of concurrent workers
	entriesPerWorker = 500                                                         // Number of entries per worker
	authorID      = 1                                                              // Fixed author ID
)

// GraphQLPayload defines the structure for the GraphQL request payload
type GraphQLPayload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// SendMutation sends a GraphQL mutation request and checks for errors in the response
func sendMutation(title string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Define the mutation query
	mutation := `
		mutation InsertCategory($authorId: ID!, $title: String!) {
			save_productCategories_productCategory_Entry(
				authorId: $authorId
				title: $title
			) {
				authorId
				title
			}
		}
	`

	// Create the payload
	payload := GraphQLPayload{
		Query: mutation,
		Variables: map[string]interface{}{
			"authorId": authorID,
			"title":    title,
		},
	}

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshalling payload: %v\n", err)
		return
	}

	// Send the request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Read and parse the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	// Check for GraphQL errors
	if errors, ok := responseBody["errors"]; ok {
		fmt.Printf("GraphQL Errors: %v\n", errors)
	} else {
		fmt.Printf("Response Status: %s\n", resp.Status)
	}
}

// Worker function to process entries
func worker(workerID int, wg *sync.WaitGroup, entriesPerWorker int) {
	defer wg.Done()
	var workerWg sync.WaitGroup

	for i := 0; i < entriesPerWorker; i++ {
		title := fmt.Sprintf("Worker %d - Sample Title %d", workerID, time.Now().UnixNano())
		workerWg.Add(1)
		go sendMutation(title, &workerWg)
	}

	// Wait for all mutations in this worker to complete
	workerWg.Wait()
	fmt.Printf("Worker %d has completed its tasks.\n", workerID)
}

func main() {
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i+1, &wg, entriesPerWorker)
	}

	// Wait for all workers to complete
	wg.Wait()

	fmt.Println("All workers have completed their tasks.")
}