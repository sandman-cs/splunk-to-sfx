package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

var (
	splunkURL      = "https://http-inputs-<realm>.splunkcloud.com"
	splunkAPIToken = "Splunk <token>"
)

// Struct to hold the response from Splunk search
type SplunkSearchResponse struct {
	Result []map[string]interface{} `json:"results"`
}

// Function to query Splunk Cloud for a metric
func QuerySplunkForMetric(query string) (float64, error) {
	// Query Splunk Cloud for the metric
	result, err := querySplunkCloud(splunkURL, query, splunkAPIToken)
	if err != nil {
		logError("QuerySplunkForMetric - Error querying Splunk Cloud:", err)
		return 0, err
	}

	// Extract the metric value from the search results
	if len(result) == 0 {
		return 0, fmt.Errorf("QuerySplunkForMetric - no results found for query: %s", query)
	}
	str := fmt.Sprintf("%v", result[0]["count"])
	num, err := strconv.ParseFloat(str, 64)
	if err != nil {
		logError("QuerySplunkForMetric - Result received: ", result)
		return 0, err
	}
	logDebug("QuerySplunkForMetric - Metric count received: ", num, " for query: ", query)
	metricsReceived++
	return num, nil
}

// Function to query Splunk Cloud and get the search results using API token
func querySplunkCloud(splunkURL string, query string, apiToken string) ([]map[string]interface{}, error) {
	// Construct the search request
	url := fmt.Sprintf("%s/services/search/jobs", splunkURL)
	data := fmt.Sprintf("search= %s&exec_mode=blocking&output_mode=json", query)

	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(data)))
	if err != nil {
		return nil, err
	}

	// Set HTTP headers including the Authorization token
	req.Header.Set("Authorization", "Bearer "+apiToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		logError("querySplunkCloud() - Error making request to Splunk Cloud:", req.Body, err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read and parse the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logError("querySplunkCloud() - Error reading response from Splunk Cloud:", resp.Body, err)
		return nil, err
	}

	// Check for errors in the Splunk response
	if resp.StatusCode != 201 {
		logError("querySplunkCloud() - Splunk API error:", body, "Status Code:", resp.StatusCode)
		return nil, fmt.Errorf("splunk API error: %s", body)
	}

	// Extract the search job ID from the response
	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		logError("querySplunkCloud() - Error unmarshaling response from Splunk Cloud:", body, err)
		return nil, err
	}
	searchJobID := response["sid"].(string)
	logDebug("Search Job ID: ", searchJobID)

	// Fetch the search results using the job ID
	resultURL := fmt.Sprintf("%s/services/search/jobs/%s/results?count=0", splunkURL, searchJobID)
	// output_mode=raw&count=0' % sid
	// Create a new request using http
	req, _ = http.NewRequest("GET", resultURL, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Authorization", "Bearer "+apiToken)

	resultResp, err := client.Do(req)
	if err != nil {
		logError("querySplunkCloud() - Error making request to Splunk Cloud for results:", resultResp.Body, err)
		return nil, err
	}
	defer resultResp.Body.Close()

	// Read and parse the search results response
	resultBody, err := io.ReadAll(resultResp.Body)
	if err != nil {
		logError("querySplunkCloud - Result Body: ", string(resultBody[:]))
		return nil, err
	}
	// Parse the JSON response to extract results
	var splunkResponse SplunkSearchResponse
	err = json.Unmarshal(resultBody, &splunkResponse)
	if err != nil {
		logError("querySplunkCloud - Splunk Response: ", splunkResponse.Result, err)
		return nil, err
	}
	// Return the results
	return splunkResponse.Result, nil
}
