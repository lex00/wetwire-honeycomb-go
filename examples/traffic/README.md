# Traffic Analysis Examples

This package demonstrates traffic analysis queries for Honeycomb.

## Queries

### RequestsByEndpoint

Shows request volume per endpoint, sorted by frequency.

**Use case:** Identify your most-called endpoints for optimization focus.

### TrafficByService

Request distribution across services with latency context.

**Use case:** Capacity planning and service-level traffic analysis.

### UniqueUsers

Counts distinct users accessing the system.

**Use case:** Track daily active users and usage patterns.

### TrafficByMethod

Request counts grouped by HTTP method (GET, POST, PUT, DELETE).

**Use case:** Understand read vs write traffic patterns.

### ThroughputByRegion

Request and data volume by geographic region.

**Use case:** Load balancing and CDN placement decisions.

## Capacity Planning

Use these queries together to answer:
- Which endpoints get the most traffic?
- How is load distributed across services?
- What are peak traffic patterns?
- Are there geographic hotspots?

## Output

```json
{
  "time_range": 14400,
  "breakdowns": ["http.route"],
  "calculations": [{"op": "COUNT"}],
  "filters": [{"column": "http.route", "op": "exists"}],
  "orders": [{"op": "COUNT", "order": "descending"}],
  "limit": 25
}
```
