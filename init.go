package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// Configuration File Objects
type configuration struct {
	SplunkURL           string         `json:"splunkURL"`
	SplunkAPIToken      string         `json:"splunkAPIToken"`
	SplunkQueryCount    int            `json:"splunkQueryCount"`
	SignalFxURL         string         `json:"signalFxURL"`
	SignalFxToken       string         `json:"signalFxToken"`
	SignalFxSenderCount int            `json:"signalFxSenderCount"`
	LogLevel            string         `json:"logLevel"`
	SplunkSources       []SplunkSource `json:"splunkSources"`
	MaxIdleConns        int            `json:"maxIdleConns"`
	MaxIdleConnsPerHost int            `json:"maxIdleConnsPerHost"`
	IdleConnTimeout     int            `json:"idleConnTimeout"`
}

// Payload structure for the SignalFx API
type SplunkSource struct {
	SplunkQuery string            `json:"splunkQuery"`
	MetricName  string            `json:"metricName"`
	Dimensions  map[string]string `json:"dimensions"`
}

// Payload structure for the SignalFx API
type SplunkStats struct {
	SplunkQuery string            `json:"splunkQuery"`
	PrefixName  string            `json:"prefixName"`
	Dimensions  map[string]string `json:"dimensions"`
}

var (
	conf            configuration
	sfxSendChannel  = make(chan []Metric, 256)
	metricsReceived = 0
	metricsSent     = 0
	transport       *http.Transport
)

func init() {
	logMessage("Starting Splunk to SignalFx metrics shovel...")

	//Crash catch or panic recovery
	defer func() {
		if err := recover(); err != nil {
			logError("Recovered with exception: ", err)
		}
	}()

	// Load Default Configuration
	conf.LogLevel = "info"
	conf.SignalFxSenderCount = 2
	conf.SplunkQueryCount = 2
	conf.MaxIdleConns = 10
	conf.MaxIdleConnsPerHost = 2
	conf.IdleConnTimeout = 90

	// Get Args from the command line
	nArgs := len(os.Args)
	args := os.Args

	// If config file is not passed as an argument, use default config.json
	configFile := "config.json"
	if nArgs > 2 {
		configFile = args[1]
	}
	// Load Configuration File
	dat, _ := os.ReadFile(configFile)
	err := json.Unmarshal(dat, &conf)

	// Fail and exit if no configuration file available
	logFatal("Error loading config.json: ", err)

	// Log Configuration
	logDebug("Running Config: ", conf.SplunkURL)
	logDebug("Splunk API Token: ", conf.SplunkAPIToken)
	logDebug("SignalFx URL: ", conf.SignalFxURL)
	logDebug("SignalFx Token: ", conf.SignalFxToken)

	// Load global variables from configuration file
	splunkURL = conf.SplunkURL
	splunkAPIToken = conf.SplunkAPIToken
	SignalFxURL = conf.SignalFxURL
	SignalFxToken = conf.SignalFxToken

	// Create Connection Pool for Splunk and SignalFX
	transport = &http.Transport{
		MaxIdleConns:        conf.MaxIdleConns,
		MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(conf.IdleConnTimeout) * time.Second,
	}
}
