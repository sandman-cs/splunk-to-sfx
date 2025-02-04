package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var (
	SignalFxURL   = "https://ingest.us0.signalfx.com/v2/datapoint"
	SignalFxToken = "qN4qfGiumrrWEWuKdC8ORQ"
)

// Metric structure to represent the payload format
type Metric struct {
	Metric     string            `json:"metric"`
	Value      float64           `json:"value"`
	Dimensions map[string]string `json:"dimensions"`
	Timestamp  int64             `json:"timestamp"`
}

// Payload structure for the SignalFx API
type Payload struct {
	Datapoints []Metric `json:"gauge"`
}

// Function to send metrics to SignalFx
func SendMetrics(metrics []Metric) {
	// Prepare the JSON payload
	payload := Payload{Datapoints: metrics}

	// Marshal the payload into JSON
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		logError("sendMetrics() - Error marshalling JSON:", err)
		return
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", SignalFxURL, bytes.NewBuffer(payloadJSON))
	if err != nil {
		logError("sendMetrics() - Error creating request:", err)
		return
	}

	// Add headers
	//req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("X-SF-TOKEN", SignalFxToken)
	req.Header.Set("Content-Type", "application/json")

	logDebug("SignalFX Request:", req.Body)
	resp, err := client.Do(req)
	if err != nil {
		logError("SignalFX failed making request:", err)
		logError("SignalFX Response Status: ", resp.Status)
		return
	}
	metricsSent++
	defer resp.Body.Close()

	// Log the response status
	logDebug("SignalFX Response Status: ", resp.Status)
}

// Function to send metrics to SignalFx
func SendMetricsChannel(data chan []Metric) {
	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
	}
	var metrics []Metric

	for {
		metrics = <-data

		// Prepare the JSON payload
		payload := Payload{Datapoints: metrics}

		// Marshal the payload into JSON
		payloadJSON, err := json.Marshal(payload)
		if err != nil {
			logError("sendMetrics() - Error marshalling JSON:", err)
			return
		}

		// Create the HTTP request
		req, err := http.NewRequest("POST", SignalFxURL, bytes.NewBuffer(payloadJSON))
		if err != nil {
			logError("sendMetrics() - Error creating request:", err)
			return
		}

		// Add headers
		//req.Header.Set("Authorization", "Bearer "+accessToken)
		req.Header.Set("X-SF-TOKEN", SignalFxToken)
		req.Header.Set("Content-Type", "application/json")

		logDebug("SignalFX Request:", req.Body)
		resp, err := client.Do(req)
		if err != nil {
			logError("sendMetrics() - SignalFX failed making request:", err)
			logError("sendMetrics() - SignalFX Response Status: ", resp.Status)
			return
		}
		metricsSent++
		resp.Body.Close()

		// Log the response status
		logDebug("sendMetrics() - SignalFX Response Status: ", resp.Status)
	}
}
