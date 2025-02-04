package main

import (
	"time"
)

func main() {
	// Start the SignalFX sender threads
	for index := 0; index < conf.SignalFxSenderCount; index++ {
		go func() {
			logInfo("Starting SignalFX sender thread:", index)
			for {
				SendMetricsChannel(sfxSendChannel)
				logWarn("SignalFX sender thread:", index, " resetting...")
				time.Sleep(1 * time.Second)
			}
		}()
	}
	// Pause to allow SignalFX sender threads to start
	time.Sleep(1 * time.Second)

	// Start the Splunk metric collector threads
	for index, element := range conf.SplunkSources {
		go func() {
			logInfo("Starting Splunk metric feed thread:", index)
			for {
				SplunkMetricQuery(element, index)
				logWarn("Resetting metric feed for Index %d...\n", index)
				time.Sleep(1 * time.Second)
			}
		}()
		// Pause to allow Splunk metric collector threads to start without overlapping
		time.Sleep(250 * time.Millisecond)
	}

	// Start the metrics reporter thread
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			logInfo("Metrics Received: ", metricsReceived, " Metrics Sent: ", metricsSent)
			metricsReceived = 0
			metricsSent = 0
		}
	}()

	// Block forever
	select {}
}

// SplunkMetricQuery queries Splunk Cloud for individual metrics
func SplunkMetricQuery(element SplunkSource, index int) {
	var metric Metric
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		//result, err := QuerySplunkForMetric2(element.SplunkQuery)
		result, err := QuerySplunkForMetric(element.SplunkQuery)
		if err != nil {
			logError("SplunkMetricQuery - Error querying Splunk Cloud for thread:", index, err)
		} else {
			metric.Value = result
			metric.Dimensions = element.Dimensions
			metric.Metric = element.MetricName
			metric.Timestamp = time.Now().UnixNano() / int64(time.Millisecond)
			logDebug("SplunkMetricQuery - Posting Metric to SFX Channel...")
			sfxSendChannel <- []Metric{metric}
		}
	}
}
