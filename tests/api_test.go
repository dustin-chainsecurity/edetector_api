package tests

import (
	"io"
	"net/http"
	"testing"
)

func TestAPIServer(t *testing.T) {
	// Create an HTTP request to your API endpoint.
	req, err := http.NewRequest("GET", "http://192.168.200.192:5000/setting/whitelist", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer 123")

	// Send the request to the API server.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	// Check the response status code.
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code %d, but got %d", http.StatusOK, resp.StatusCode)
	}

	// Read and process the response body.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	// Define the expected response content.
	expectedResponse := "Expected Response"

	// // Perform assertions on the response body.
	if string(body) != expectedResponse {
		t.Fatalf("Expected response '%s', but got '%s'", expectedResponse, string(body))
	}
}
