package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/sumukshashidhar/integrate-kibana-and-keptn/pkg/provider"
)

func main() {
	// Get configuration from environment variables
	esHost := getEnv("ES_HOST", "localhost")
	esPortStr := getEnv("ES_PORT", "9200")
	esPort, err := strconv.Atoi(esPortStr)
	if err != nil {
		log.Fatalf("Invalid ES_PORT: %v", err)
	}
	esUsername := getEnv("ES_USERNAME", "elastic")
	esPassword := getEnv("ES_PASSWORD", "changeme")
	keptnEndpoint := getEnv("KEPTN_ENDPOINT", "http://localhost:8080/api/v1/metrics")

	// Create provider
	p, err := provider.NewKibanaProvider(esHost, esPort, esUsername, esPassword)
	if err != nil {
		log.Fatalf("Failed to create Kibana provider: %v", err)
	}

	// Query metrics
	ctx := context.Background()
	response, err := p.QueryKibanaMetrics(ctx, "metrics-*", "cpu_usage", 30)
	if err != nil {
		log.Fatalf("Failed to query metrics: %v", err)
	}

	// Format metrics for Keptn
	keptnMetrics := p.FormatKeptnMetrics(response)

	// Send metrics to Keptn
	if err := p.SendToKeptn(ctx, keptnMetrics, keptnEndpoint); err != nil {
		log.Fatalf("Failed to send metrics to Keptn: %v", err)
	}

	log.Println("Successfully sent metrics to Keptn")
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}