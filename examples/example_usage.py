from src.kibana_provider import KibanaKeptnProvider

def main():
    # Initialize the provider
    provider = KibanaKeptnProvider(
        es_host="localhost",
        es_port=9200,
        username="elastic",  # Replace with your Elasticsearch username
        password="changeme"  # Replace with your Elasticsearch password
    )
    
    # Query metrics from Kibana
    metrics_response = provider.query_kibana_metrics(
        index_pattern="metrics-*",
        metric_name="cpu_usage",
        time_range_minutes=30
    )
    
    # Format metrics for Keptn
    keptn_metrics = provider.format_keptn_metrics(metrics_response)
    
    # Send metrics to Keptn
    keptn_endpoint = "http://localhost:8080/api/v1/metrics"  # Replace with your Keptn endpoint
    success = provider.send_to_keptn(keptn_metrics, keptn_endpoint)
    
    if success:
        print("Successfully sent metrics to Keptn")
    else:
        print("Failed to send metrics to Keptn")

if __name__ == "__main__":
    main()