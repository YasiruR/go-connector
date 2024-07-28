#!/bin/bash

res=$(curl -X POST -d '{"permissions": [{"action": "use", "constraints": [{"leftOperand": "region", "operator": "eq", "rightOperand": "eu"}]}]}' http://localhost:9081/gateway/create-policy)
policy_id=$(echo "$res" | awk -F[\"\"] '{print $4}')

data='{"title": "sample dataset", "description": ["sample description"], "endpoints": ["http://localhost:9080/datasource"], "offerIds": ["'$policy_id'"], "keywords": ["dataspace", "connector"]}'
res=$(curl -X POST -d "$data" http://localhost:9081/gateway/create-dataset)
dataset_id=$(echo "$res" | awk -F[\"\"] '{print $4}')

data='{"offerId": "'$policy_id'", "providerEndpoint": "http://localhost:9080", "odrlTarget": "test-target", "assigner": "provider1", "assignee": "consumer1", "action": "odrl:use"}'
res=$(curl -X POST -d "$data" http://localhost:8081/gateway/request-contract)
negConsPid=$(echo "$res" | awk -F[\"\"] '{print $4}')

echo "request contract done"

read -p "Contract Negotiation ID (Provider): " negProvPid
data='{"offerId": "'$policy_id'", "contractNegotiationId": "'$negProvPid'"}'
res=$(curl -X POST -d "$data" http://localhost:9081/gateway/agree-contract)
agrId=$(echo "$res" | awk -F[\"\"] '{print $4}')
echo "agree contract done"

curl -X POST http://localhost:8081/gateway/verify-agreement/"$negConsPid"
echo "verify contract done"
curl -X POST http://localhost:9081/gateway/finalize-contract/"$negProvPid"
echo "finalize contract done"

data='{"transferType": "HTTP_PUSH", "agreementId": "'$agrId'", "sinkEndpoint": "http://localhost:8080/datasink", "providerEndpoint": "http://localhost:9080"}'
res=$(curl -X POST -d "$data" http://localhost:8081/gateway/transfer/request)
tpConsPid=$(echo "$res" | awk -F[\"\"] '{print $4}')

read -p "Transfer Process ID (Provider): " tpProvId
data='{"transferProcessId": "'$tpProvId'"}'
curl -X POST -d "$data" http://localhost:9081/gateway/transfer/start
