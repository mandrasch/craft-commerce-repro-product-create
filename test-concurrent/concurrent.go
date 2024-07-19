package main

// TODO: graphql not tested here, only in test-graphl-mutation/

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"sync"
	"time"
	"io"
)

const (
	apiURLGraphQL    = "https://craft-commerce-repro-product-create.ddev.site/graphql"        // Replace with your CraftCMS GraphQL endpoint
	apiURLREST       = "https://craft-commerce-repro-product-create.ddev.site/rest-api/product-categories/create" // Replace with your REST API endpoint
	bearerToken      = "CDwvk4kPKUBTBWsivteRLgd7DtIzcEa6"                                     // Replace with your actual bearer token
	numWorkers       = 3                                                                      // Number of concurrent workers
	entriesPerWorker = 500                                                                    // Number of entries per worker
	authorID         = 1                                                                      // Fixed author ID
)

// GraphQLPayload defines the structure for the GraphQL request payload
type GraphQLPayload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// RESTPayload defines the structure for the REST API request payload
type RESTPayload struct {
	AuthorID int    `json:"authorId"`
	Title    string `json:"title"`
}

// sendGraphQLMutation sends a GraphQL mutation request and checks for errors in the response
func sendGraphQLMutation(title string) (bool, error) {
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
		return false, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Send the request
	req, err := http.NewRequest("POST", apiURLGraphQL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Print the HTTP status code
	fmt.Printf("GraphQL Response Status: %s\n", resp.Status)

	// Read and parse the response body
	var responseBody map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseBody)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %v", err)
	}

	// Check for GraphQL errors
	if errors, ok := responseBody["errors"]; ok {
		return false, fmt.Errorf("graphql errors: %v", errors)
	}

	return resp.StatusCode == http.StatusOK, nil
}

// sendRESTRequest sends a REST API request and checks for errors in the response
func sendRESTRequest(title string) (bool, error) {
	// Create the payload
	payload := RESTPayload{
		AuthorID: authorID,
		Title:    title,
	}

	// Marshal the payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("error marshalling payload: %v", err)
	}

	// Send the request
	req, err := http.NewRequest("POST", apiURLREST, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return false, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	// Read the raw response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response body: %v", err)
	}
	fmt.Printf("Raw Response Body: %s\n", bodyBytes)

	// Decode the response body
	var responseBody map[string]interface{}
	err = json.Unmarshal(bodyBytes, &responseBody)
	if err != nil {
		return false, fmt.Errorf("error unmarshalling response body: %v", err)
	}
	// fmt.Printf("Decoded Response Body: %v\n", responseBody)

	// Check for "success" field
	success, successOK := responseBody["success"].(bool)
	if !successOK || !success {
		// Check for "errors" field only if success is not true
		if errors, errorsOK := responseBody["errors"].([]interface{}); errorsOK {
			if len(errors) > 0 {
				return false, fmt.Errorf("REST API error: errors present: %v", errors)
			}
		} else if responseBody["errors"] != nil {
			return false, fmt.Errorf("REST API error: errors field is not an array")
		}
		return false, fmt.Errorf("REST API error: success is false")
	}

	return true, nil
}



// Worker function to process entries
func worker(workerID int, wg *sync.WaitGroup, entriesPerWorker int, testType string) {
	defer wg.Done()
	var workerWg sync.WaitGroup

	for i := 0; i < entriesPerWorker; i++ {
		title := fmt.Sprintf("Worker %d - Sample Title %d", workerID, time.Now().UnixNano())
		workerWg.Add(1)
		if testType == "graphql" {
			go func() {
				defer workerWg.Done()
				success, err := sendGraphQLMutation(title)
				if err != nil {
					fmt.Printf("GraphQL Mutation Error: %v\n", err)
				} else if !success {
					fmt.Println("GraphQL Mutation was not successful")
				}
			}()
		} else if testType == "rest" {
			go func() {
				defer workerWg.Done()
				success, err := sendRESTRequest(title)
				if err != nil {
					fmt.Printf("REST Request Error: %v\n", err)
				} else if !success {
					fmt.Println("REST Request was not successful")
				}
			}()
		}
	}

	// Wait for all mutations in this worker to complete
	workerWg.Wait()
	fmt.Printf("Worker %d has completed its tasks.\n", workerID)
}

func main() {
	// Parse command-line flag for test type (graphql or rest)
	testType := flag.String("testType", "graphql", "Type of test to perform: graphql or rest")
	flag.Parse()

	// Perform a single initial request to ensure it's successful
	var success bool
	var err error
	title := "Initial Sample Title"

	if *testType == "graphql" {
		success, err = sendGraphQLMutation(title)
	} else if *testType == "rest" {
		success, err = sendRESTRequest(title)
	}

	if !success {
		fmt.Println("Initial request did not return a successful response")
		fmt.Printf("Errors: %v\n", err)
		return
	}

	fmt.Println("Initial request was successful, proceeding with concurrent requests")

	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i+1, &wg, entriesPerWorker, *testType)
	}

	// Wait for all workers to complete
	wg.Wait()

	fmt.Println("All workers have completed their tasks.")
}
