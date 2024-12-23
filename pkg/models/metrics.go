package models

import "time"

type KibanaMetric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels"`
}

type KeptnMetric struct {
	Name      string            `json:"name"`
	Value     float64           `json:"value"`
	Timestamp time.Time         `json:"timestamp"`
	Labels    map[string]string `json:"labels"`
}

type ElasticsearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source struct {
				MetricName string    `json:"metric_name"`
				Value      float64   `json:"value"`
				Timestamp  time.Time `json:"@timestamp"`
			} `json:"_source"`
			Index string `json:"_index"`
		} `json:"hits"`
	} `json:"hits"`
}