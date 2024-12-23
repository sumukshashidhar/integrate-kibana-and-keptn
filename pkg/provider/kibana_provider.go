package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/sumukshashidhar/integrate-kibana-and-keptn/pkg/models"
)

type KibanaProvider struct {
	esClient *elasticsearch.Client
}

func NewKibanaProvider(esHost string, esPort int, username, password string) (*KibanaProvider, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{fmt.Sprintf("http://%s:%d", esHost, esPort)},
		Username:  username,
		Password:  password,
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("error creating elasticsearch client: %w", err)
	}

	return &KibanaProvider{
		esClient: client,
	}, nil
}

func (p *KibanaProvider) QueryKibanaMetrics(ctx context.Context, indexPattern, metricName string, timeRangeMinutes int) (*models.ElasticsearchResponse, error) {
	now := time.Now().UTC()
	timeFrom := now.Add(-time.Duration(timeRangeMinutes) * time.Minute)

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"metric_name": metricName,
						},
					},
					{
						"range": map[string]interface{}{
							"@timestamp": map[string]interface{}{
								"gte": timeFrom.Format(time.RFC3339),
								"lte": now.Format(time.RFC3339),
							},
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, fmt.Errorf("error encoding query: %w", err)
	}

	res, err := p.esClient.Search(
		p.esClient.Search.WithContext(ctx),
		p.esClient.Search.WithIndex(indexPattern),
		p.esClient.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, fmt.Errorf("error executing search: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	var response models.ElasticsearchResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	return &response, nil
}

func (p *KibanaProvider) FormatKeptnMetrics(esResponse *models.ElasticsearchResponse) []models.KeptnMetric {
	var keptnMetrics []models.KeptnMetric

	for _, hit := range esResponse.Hits.Hits {
		metric := models.KeptnMetric{
			Name:      hit.Source.MetricName,
			Value:     hit.Source.Value,
			Timestamp: hit.Source.Timestamp,
			Labels: map[string]string{
				"source": "kibana",
				"index":  hit.Index,
			},
		}
		keptnMetrics = append(keptnMetrics, metric)
	}

	return keptnMetrics
}

func (p *KibanaProvider) SendToKeptn(ctx context.Context, metrics []models.KeptnMetric, keptnEndpoint string) error {
	payload, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("error marshaling metrics: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", keptnEndpoint, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending metrics to Keptn: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from Keptn: %d", resp.StatusCode)
	}

	return nil
}