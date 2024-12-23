# integrate-kibana-and-keptn
# Kibana-Keptn Integration

A Go-based integration between Kibana/Elasticsearch and Keptn that enables seamless metric data transfer from Kibana to Keptn for evaluation and analysis.

## Features

- Query metrics from Kibana/Elasticsearch using configurable parameters
- Format metrics data into Keptn-compatible format
- Send metrics to Keptn's metrics API
- Support for authentication and secure connections
- Context-aware operations with proper timeout handling
- Comprehensive error handling and logging

## Prerequisites

- Go 1.21 or later
- Running Elasticsearch/Kibana instance
- Running Keptn instance
- Access credentials for both services

## Installation

1. Clone the repository:
```bash
git clone https://github.com/sumukshashidhar/integrate-kibana-and-keptn.git
cd integrate-kibana-and-keptn
```

2. Install dependencies:
```bash
go mod download
```

## Configuration

The provider can be configured using environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| ES_HOST | Elasticsearch host | localhost |
| ES_PORT | Elasticsearch port | 9200 |
| ES_USERNAME | Elasticsearch username | elastic |
| ES_PASSWORD | Elasticsearch password | changeme |
| KEPTN_ENDPOINT | Keptn metrics API endpoint | http://localhost:8080/api/v1/metrics |

## Usage

### As a Standalone Application

1. Build the application:
```bash
go build -o kibana-keptn-provider ./cmd/main.go
```

2. Configure environment variables (optional):
```bash
export ES_HOST=your-elasticsearch-host
export ES_PORT=9200
export ES_USERNAME=your-username
export ES_PASSWORD=your-password
export KEPTN_ENDPOINT=your-keptn-endpoint
```

3. Run the provider:
```bash
./kibana-keptn-provider
```

### As a Library

```go
import (
    "context"
    "log"
    "github.com/sumukshashidhar/integrate-kibana-and-keptn/pkg/provider"
)

func main() {
    // Create provider
    p, err := provider.NewKibanaProvider("localhost", 9200, "elastic", "changeme")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Query metrics
    response, err := p.QueryKibanaMetrics(ctx, "metrics-*", "cpu_usage", 30)
    if err != nil {
        log.Fatal(err)
    }

    // Format metrics for Keptn
    keptnMetrics := p.FormatKeptnMetrics(response)

    // Send metrics to Keptn
    err = p.SendToKeptn(ctx, keptnMetrics, "http://localhost:8080/api/v1/metrics")
    if err != nil {
        log.Fatal(err)
    }
}
```

## Project Structure

```
.
├── cmd/
│   └── main.go                # Main application entry point
├── pkg/
│   ├── models/
│   │   └── metrics.go         # Data models for metrics
│   └── provider/
│       └── kibana_provider.go # Provider implementation
├── go.mod                     # Go module file
├── go.sum                     # Go module checksum
└── README.md                  # This file
```

## API Reference

### KibanaProvider

#### NewKibanaProvider
```go
func NewKibanaProvider(esHost string, esPort int, username, password string) (*KibanaProvider, error)
```
Creates a new instance of the Kibana provider.

#### QueryKibanaMetrics
```go
func (p *KibanaProvider) QueryKibanaMetrics(ctx context.Context, indexPattern, metricName string, timeRangeMinutes int) (*models.ElasticsearchResponse, error)
```
Queries metrics from Kibana/Elasticsearch.

Parameters:
- `ctx`: Context for cancellation and timeouts
- `indexPattern`: Elasticsearch index pattern to query
- `metricName`: Name of the metric to retrieve
- `timeRangeMinutes`: Time range in minutes to look back

#### FormatKeptnMetrics
```go
func (p *KibanaProvider) FormatKeptnMetrics(esResponse *models.ElasticsearchResponse) []models.KeptnMetric
```
Formats Elasticsearch response into Keptn metrics format.

#### SendToKeptn
```go
func (p *KibanaProvider) SendToKeptn(ctx context.Context, metrics []models.KeptnMetric, keptnEndpoint string) error
```
Sends formatted metrics to Keptn.

## Error Handling

The provider implements comprehensive error handling:

- Connection errors are wrapped with context
- Query validation errors are caught early
- Timeout and cancellation are properly handled
- All errors include meaningful messages

Example error handling:
```go
response, err := provider.QueryKibanaMetrics(ctx, "metrics-*", "cpu_usage", 30)
if err != nil {
    switch {
    case errors.Is(err, context.DeadlineExceeded):
        log.Fatal("Query timed out")
    case errors.Is(err, context.Canceled):
        log.Fatal("Query was canceled")
    default:
        log.Fatalf("Query failed: %v", err)
    }
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Keptn team for their excellent metrics API
- Elasticsearch team for their Go client library
