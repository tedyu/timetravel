package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAPIScenarios(t *testing.T) {
	baseURL := "http://127.0.0.1:8000/api"

	// Simulate a POST request to /v2/records/1 with '{"hello":"x"}'
	postData1 := map[string]string{"hello": "x"}
	resp1, err := sendPostRequest(baseURL+"/v2/records/1", postData1)
	if err != nil {
		t.Fatalf("POST /v2/records/1 failed: %v", err)
	}
	defer resp1.Body.Close()
	checkResponseCode(t, resp1, http.StatusOK)

	// Simulate a POST request to /v2/records/1 with '{"hello":"z"}'
	postData2 := map[string]string{"hello": "z"}
	resp2, err := sendPostRequest(baseURL+"/v2/records/1", postData2)
	if err != nil {
		t.Fatalf("POST /v2/records/1 failed: %v", err)
	}
	defer resp2.Body.Close()
	checkResponseCode(t, resp2, http.StatusOK)

	// Simulate a GET request to /v2/records/1/version/5 (should return error)
	resp3, err := sendGetRequest(baseURL + "/v2/records/1/version/5")
	if err != nil {
		t.Fatalf("GET /v2/records/1/version/5 failed: %v", err)
	}
	defer resp3.Body.Close()
	checkResponseCode(t, resp3, http.StatusBadRequest)
	expectedError := `{"error":"record of id 1 does not exist"}`
	checkResponseBody(t, resp3, expectedError)

	// Simulate a GET request to /v2/records/1/version/1 (should return first version)
	resp4, err := sendGetRequest(baseURL + "/v2/records/1/version/1")
	if err != nil {
		t.Fatalf("GET /v2/records/1/version/1 failed: %v", err)
	}
	defer resp4.Body.Close()
	checkResponseCode(t, resp4, http.StatusOK)
	expectedBody1 := `{"id":1,"data":{"hello":"x"},"Ver":1}`
	checkResponseBody(t, resp4, expectedBody1)

	// Simulate a GET request to /v2/records/1/version/2 (should return second version)
	resp5, err := sendGetRequest(baseURL + "/v2/records/1/version/2")
	if err != nil {
		t.Fatalf("GET /v2/records/1/version/2 failed: %v", err)
	}
	defer resp5.Body.Close()
	checkResponseCode(t, resp5, http.StatusOK)
	expectedBody2 := `{"id":1,"data":{"hello":"z"},"Ver":2}`
	checkResponseBody(t, resp5, expectedBody2)

	// Simulate a POST request to /v1/records/2 with '{"hell":"z"}'
	postData3 := map[string]string{"hell": "z"}
	resp6, err := sendPostRequest(baseURL+"/v1/records/2", postData3)
	if err != nil {
		t.Fatalf("POST /v1/records/2 failed: %v", err)
	}
	defer resp6.Body.Close()
	checkResponseCode(t, resp6, http.StatusOK)

	// Simulate a POST request to /v1/records/2 with '{"hell":"y"}'
	postData4 := map[string]string{"hell": "y"}
	resp7, err := sendPostRequest(baseURL+"/v1/records/2", postData4)
	if err != nil {
		t.Fatalf("POST /v1/records/2 failed: %v", err)
	}
	defer resp7.Body.Close()
	checkResponseCode(t, resp7, http.StatusOK)

	// Simulate a GET request to /v1/records/2 (should return updated record)
	resp8, err := sendGetRequest(baseURL + "/v1/records/2")
	if err != nil {
		t.Fatalf("GET /v1/records/2 failed: %v", err)
	}
	defer resp8.Body.Close()
	checkResponseCode(t, resp8, http.StatusOK)
	expectedBody3 := `{"id":2,"data":{"hell":"y"},"Ver":1}`
	checkResponseBody(t, resp8, expectedBody3)

	// Simulate a GET request to /v2/records/1/latest (should return latest version)
	resp9, err := sendGetRequest(baseURL + "/v2/records/1/latest")
	if err != nil {
		t.Fatalf("GET /v2/records/1/latest failed: %v", err)
	}
	defer resp9.Body.Close()
	checkResponseCode(t, resp9, http.StatusOK)
	expectedLatest := `"{version:2}"`
	checkResponseBody(t, resp9, expectedLatest)
}

// Helper function to send POST request
func sendPostRequest(url string, data map[string]string) (*http.Response, error) {
	jsonData, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	return client.Do(req)
}

// Helper function to send GET request
func sendGetRequest(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return client.Do(req)
}

// Helper function to check the response code
func checkResponseCode(t *testing.T, resp *http.Response, expected int) {
	if resp.StatusCode != expected {
		t.Errorf("Expected response code %d but got %d", expected, resp.StatusCode)
	}
}

// Helper function to check the response body
func checkResponseBody(t *testing.T, resp *http.Response, expected string) {
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	body := strings.TrimSpace(string(bodyBytes))
	if body != expected {
		t.Errorf("Expected response body %s but got %s", expected, body)
	}
}

