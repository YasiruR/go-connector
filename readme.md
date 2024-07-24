# go-connector

## Sample requests

1. Create policy: ``curl -X POST -d '{"permissions": [{"action": "use", "constraints": [{"leftOperand": "region", "operator": "eq", "rightOperand": "eu"}]}]}' http://localhost:9081/gateway/create-policy``
2. Create dataset: ``curl -X POST -d '{"title": "sample dataset", "description": ["sample description"], "endpoints": ["http://localhost:9080/datasource"], "offerIds": ["<policy-id>"], "keywords": ["dataspace", "connector"]}' http://localhost:9081/gateway/create-dataset``
3. Get Catalog: ``curl -X POST -d '{"providerEndpoint": "http://localhost:9080"}' http://localhost:8081/gateway/catalog``
4. Get Dataset: ``curl -X POST -d '{"datasetId": "<dataset-id>", "providerEndpoint": "http://localhost:9080"}' http://localhost:8081/gateway/dataset | jq``
5. Request contract: ``curl -X POST -d '{"offerId": "<policy-id>", "providerEndpoint": "http://localhost:9080", "odrlTarget": "test-target", "assigner": "provider1", "assignee": "consumer1", "action": "odrl:use"}' http://localhost:8081/gateway/contract``
6. Get negotiation: ``curl -X GET http://localhost:9080/negotiations/{providerPid} | jq``
7. Agree contract: ``curl -X POST -d '{"offerId": "<policy-id>", "negotiationId": "<providerPid>"}' http://localhost:9081/gateway/agree-contract``
8. Get agreement: ``curl -X GET http://localhost:8081/gateway/agreement/{id}``
9. Verify agreement: ``curl -X POST http://localhost:8081/gateway/verify-agreement/{consumerPid}``
