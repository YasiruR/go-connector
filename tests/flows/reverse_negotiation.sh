#!/bin/bash

res=$(curl --silent -X POST -d '{"permissions": [{"action": "use", "constraints": [{"leftOperand": "region", "operator": "eq", "rightOperand": "eu"}]}]}' http://localhost:9081/gateway/create-policy)
policy_id=$(echo "$res" | awk -F[\"\"] '{print $4}')

data='{"title": "sample dataset", "description": ["sample description"], "endpoints": ["http://localhost:9080/datasource"], "offerIds": ["'$policy_id'"], "keywords": ["dataspace", "connector"]}'
res=$(curl --silent -X POST -d "$data" http://localhost:9081/gateway/create-dataset)
dataset_id=$(echo "$res" | awk -F[\"\"] '{print $4}')

data='{"offerId": "'$policy_id'", "consumerAddr": "http://localhost:8080"}'
res=$(curl --silent -X POST -d "$data" http://localhost:9081/gateway/offer-contract)
negProvPid=$(echo "$res" | awk -F[\"\"] '{print $4}')

read -p "Contract Negotiation ID (Consumer): " negConsPid
data='{"offerId": "'$policy_id'", "consumerPid": "'$negConsPid'", "providerEndpoint": "http://localhost:9080", "odrlTarget": "test-target", "assigner": "provider1", "assignee": "consumer1", "action": "odrl:use"}'
res=$(curl --silent -X POST -d "$data" http://localhost:8081/gateway/request-contract)
#negConsPid=$(echo "$res" | awk -F[\"\"] '{print $4}')

data='{"offerId": "'$policy_id'", "providerPid": "'$negProvPid'"}'
res=$(curl --silent -X POST -d "$data" http://localhost:9081/gateway/offer-contract)

curl --silent -X POST http://localhost:8081/gateway/accept-offer/"$negConsPid"

data='{"offerId": "'$policy_id'", "contractNegotiationId": "'$negProvPid'"}'
res=$(curl --silent -X POST -d "$data" http://localhost:9081/gateway/agree-contract)
agrId=$(echo "$res" | awk -F[\"\"] '{print $4}')

curl --silent -X POST http://localhost:8081/gateway/verify-agreement/"$negConsPid"
curl --silent -X POST http://localhost:9081/gateway/finalize-contract/"$negProvPid"

echo $agrId > tmp.txt
