# go-connector

_go-connector_ is a Golang framework of a data space connector which adheres to IDS Protocols.

## Sample requests

Sample requests provided in this document assume the following endpoints to be up and running.

- Consumer DSP API: 8080
- Consumer gateway API: 8081
- Provider DSP API: 9080
- Provider gateway API: 9081

### Catalog Protocol

1. Create policy (Provider): ``curl -X POST -d '{"permissions": [{"action": "use", "constraints": [{"leftOperand": "region", "operator": "eq", "rightOperand": "eu"}]}]}' http://localhost:9081/gateway/create-policy``
2. Create dataset (Provider): ``curl -X POST -d '{"title": "sample dataset", "description": ["sample description"], "endpoints": ["http://localhost:9080/datasource"], "policyIds": ["<policy-id>"], "keywords": ["dataspace", "connector"]}' http://localhost:9081/gateway/create-dataset``
3. Request catalog (Consumer): ``curl -X POST -d '{"providerEndpoint": "http://localhost:9080"}' http://localhost:8081/gateway/request-catalog | jq``
4. Request dataset (Consumer): ``curl -X POST -d '{"datasetId": "<dataset-id>", "providerEndpoint": "http://localhost:9080"}' http://localhost:8081/gateway/request-dataset | jq``
5. Get stored catalogs: ````

### Contract Negotiation

1. Request contract (Consumer): ``curl -X POST -d '{"offerId": "<offer-id>", "providerEndpoint": "http://localhost:9080"}' http://localhost:8081/gateway/request-contract``
2. Offer contract (Provider): ``curl -X POST -d '{"offerId": "<offer-id>", "consumerAddr": "http://localhost:8080"}' http://localhost:9081/gateway/offer-contract``
3. Accept contract (Consumer): ``curl -X POST http://localhost:8081/gateway/accept-offer/<consumerPid>``
4. Get negotiation (Provider): ``curl -X GET http://localhost:9080/negotiations/{providerPid} | jq`` 
5. Agree contract (Provider): ``curl -X POST -d '{"offerId": "<offer-id>", "contractNegotiationId": "<providerPid>"}' http://localhost:9081/gateway/agree-contract``
6. Get agreement (Consumer): ``curl -X GET http://localhost:8081/gateway/agreement/{id}``
7. Verify agreement (Consumer): ``curl -X POST http://localhost:8081/gateway/verify-agreement/{consumerPid}``
8. Finalize contract (Provider): ``curl -X POST http://localhost:9081/gateway/finalize-contract/{providerPid}`` 
9. Request contract as a response to an offer (Consumer): ``curl -X POST -d '{"consumerPid": "<consumerPid>", "offerId": "<offer-id>", "providerEndpoint": "http://localhost:9080", "odrlTarget": "test-target", "assigner": "provider1", "assignee": "consumer1", "action": "odrl:use"}' http://localhost:8081/gateway/request-contract``

### Transfer Process

1. Request transfer (Consumer): ``curl -X POST -d '{"transferType": "HTTP_PUSH", "<agreement-id>": "<agreement-id>", "sinkEndpoint": "http://localhost:8080/datasink", "providerEndpoint": "http://localhost:9080"}' http://localhost:8081/gateway/transfer/request``
2. Start transfer (Provider): ``curl -X POST -d '{"transferProcessId": "<providerPid>"}' http://localhost:9081/gateway/transfer/start``
3. Suspend transfer (Consumer/Provider): ``curl -X POST -d '{"provider": false, "<transfer-process-id>": "<consumerPid>", "code": "2400", "Reasons": ["invalid data", "incompatible syntax"]}' http://localhost:8081/gateway/transfer/suspend``
4. Complete transfer (Consumer/Provider): ``curl -X POST -d '{"provider": true, "<transfer-process-id>": "<providerPid>"}' http://localhost:8081/gateway/transfer/complete`` 
5. Terminate transfer (Consumer/Provider): ``curl -X POST -d '{"transferProcessId": "<transfer-process-id>", "code": "2333", "reasons": ["outdated data"]}' http://localhost:8081/gateway/transfer/terminate``
