from elasticsearch import Elasticsearch
from datetime import datetime, timedelta
import os
import json
from typing import Dict, List, Optional

class KibanaKeptnProvider:
    def __init__(self, es_host: str, es_port: int = 9200, 
                 username: Optional[str] = None, password: Optional[str] = None):
        self.es_client = Elasticsearch(
            [{'host': es_host, 'port': es_port}],
            http_auth=(username, password) if username and password else None
        )
        
    def query_kibana_metrics(self, index_pattern: str, metric_name: str, 
                           time_range_minutes: int = 60) -> Dict:
        """
        Query metrics from Kibana/Elasticsearch
        
        Args:
            index_pattern: The Elasticsearch index pattern to query
            metric_name: Name of the metric to retrieve
            time_range_minutes: Time range in minutes to look back
            
        Returns:
            Dict containing the query results
        """
        now = datetime.utcnow()
        time_from = now - timedelta(minutes=time_range_minutes)
        
        query = {
            "query": {
                "bool": {
                    "must": [
                        {"match": {"metric_name": metric_name}},
                        {
                            "range": {
                                "@timestamp": {
                                    "gte": time_from.isoformat(),
                                    "lte": now.isoformat()
                                }
                            }
                        }
                    ]
                }
            }
        }
        
        try:
            response = self.es_client.search(
                index=index_pattern,
                body=query
            )
            return response
        except Exception as e:
            print(f"Error querying Elasticsearch: {str(e)}")
            return {}

    def format_keptn_metrics(self, es_response: Dict) -> List[Dict]:
        """
        Format Elasticsearch response into Keptn metrics format
        
        Args:
            es_response: Response from Elasticsearch query
            
        Returns:
            List of formatted metrics for Keptn
        """
        keptn_metrics = []
        
        if 'hits' not in es_response or 'hits' not in es_response['hits']:
            return keptn_metrics
            
        for hit in es_response['hits']['hits']:
            source = hit['_source']
            
            metric = {
                'name': source.get('metric_name'),
                'value': source.get('value'),
                'timestamp': source.get('@timestamp'),
                'labels': {
                    'source': 'kibana',
                    'index': hit['_index']
                }
            }
            keptn_metrics.append(metric)
            
        return keptn_metrics

    def send_to_keptn(self, metrics: List[Dict], keptn_endpoint: str) -> bool:
        """
        Send formatted metrics to Keptn
        
        Args:
            metrics: List of formatted metrics
            keptn_endpoint: Keptn API endpoint
            
        Returns:
            bool indicating success/failure
        """
        try:
            # Here you would implement the actual API call to Keptn
            # This is a placeholder for the actual implementation
            print(f"Sending metrics to Keptn: {json.dumps(metrics, indent=2)}")
            return True
        except Exception as e:
            print(f"Error sending metrics to Keptn: {str(e)}")
            return False