#!/bin/bash

res=$(curl --silent -X POST -d '{"permissions": [{"action": "use", "constraints": [{"leftOperand": "region", "operator": "eq", "rightOperand": "eu"}]}]}' http://localhost:9081/gateway/create-policy)
policy_id=$(echo "$res" | awk -F[\"\"] '{print $4}')

data='{"title": "sample dataset", "description": ["sample description"], "endpoints": ["http://localhost:9080/datasource"], "offerIds": ["'$policy_id'"], "keywords": ["dataspace", "connector"]}'
res=$(curl --silent -X POST -d "$data" http://localhost:9081/gateway/create-dataset)
dataset_id=$(echo "$res" | awk -F[\"\"] '{print $4}')

echo "$data" | jq

curl --silent -X POST -d '{"providerEndpoint": "http://localhost:9080"}' http://localhost:8081/gateway/request-catalog

data='{"offerId": "'$policy_id'", "providerEndpoint": "http://localhost:9080", "constraints": {"region": "eu"}}'
res=$(curl --silent -X POST -d "$data" http://localhost:8081/gateway/request-contract)
negConsPid=$(echo "$res" | awk -F[\"\"] '{print $4}')

read -p "Contract Negotiation ID (Provider): " negProvPid
data='{"offerId": "'$policy_id'", "contractNegotiationId": "'$negProvPid'"}'
res=$(curl --silent -X POST -d "$data" http://localhost:9081/gateway/agree-contract)
agrId=$(echo "$res" | awk -F[\"\"] '{print $4}')

curl --silent -X POST http://localhost:8081/gateway/verify-agreement/"$negConsPid"
curl --silent -X POST http://localhost:9081/gateway/finalize-contract/"$negProvPid"

echo $agrId > tmp.txt
