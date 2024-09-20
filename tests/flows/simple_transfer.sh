#!/bin/bash

agrId=$1

data='{"transferType": "HTTP_PUSH", "agreementId": "'$agrId'", "sinkEndpoint": "http://localhost:8080/datasink", "providerEndpoint": "http://localhost:9080"}'
res=$(curl --silent -X POST -d "$data" http://localhost:8081/gateway/transfer/request)
tpConsPid=$(echo "$res" | awk -F[\"\"] '{print $4}')

read -p "Transfer Process ID (Provider): " tpProvId
data='{"provider": true, "transferProcessId": "'$tpProvId'"}'
curl --silent -X POST -d "$data" http://localhost:9081/gateway/transfer/start

data='{"provider": true, "transferProcessId": "'$tpProvId'"}'
curl --silent -X POST -d "$data" http://localhost:9081/gateway/transfer/complete