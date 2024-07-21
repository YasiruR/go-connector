# go-connector

## Sample requests

1. Initialize contract request: 
``curl -X POST -d '{"offerId": "offer-id", "providerEndpoint": "http://localhost:9080", "odrlTarget": "test-target", "assigner": "provider1", "assignee": "consumer1", "action": "odrl:use"}' http://localhost:8081/gateway/contract-request``
2. Get negotiation: ``curl -X GET http://localhost:9080/negotiations/{providerPid}``
3. Create policy: ``curl -X POST -d '{"permissions": [{"action": "use", "constraints": [{"leftOperand": "region", "operator": "eq", "rightOperand": "eu"}]}]}' http://localhost:9081/gateway/policy``
4. Create dataset: ``curl -X POST -d '{"title": "sample dataset", "descriptions": ["sample description"], "endpoints": ["http://localhost:9080/datasource"], "policyIds": ["<returned-policy-id>"], "keywords": ["dataspace", "connector"]}' http://localhost:9081/gateway/dataset``
