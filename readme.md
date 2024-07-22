# go-connector

## Sample requests

1. Create policy: ``curl -X POST -d '{"permissions": [{"action": "use", "constraints": [{"leftOperand": "leftt", "operator": "eq", "rightOperand": "rightt"}]}]}' http://localhost:9081/gateway/policy``
2. Create dataset: ``curl -X POST -d '{"title": "sample dataset", "description": ["sample description"], "endpoints": ["http://localhost:9080/datasource"], "policyIds": ["<policy-id>"], "keywords": ["dataspace", "connector"]}' http://localhost:9081/gateway/dataset``
3. Contract request: ``curl -X POST -d '{"offerId": "<policy-id>", "providerEndpoint": "http://localhost:9080", "odrlTarget": "test-target", "assigner": "provider1", "assignee": "consumer1", "action": "odrl:use"}' http://localhost:8081/gateway/contract-request``
4. Get negotiation: ``curl -X GET http://localhost:9080/negotiations/{providerPid}``
5. Agree contract: ``curl -X POST -d '{"offerId": "<policy-id>", "negotiationId": "<providerPid>"}' http://localhost:9081/gateway/contract-agreement``
