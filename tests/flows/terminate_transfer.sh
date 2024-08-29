#!/bin/bash

# Flow: Request(C) -> Start(P) -> Suspend(C) -> Start(C) -> Suspend(P) -> Terminate(C)

agrId=$1

data='{"transferType": "HTTP_PUSH", "agreementId": "'$agrId'", "sinkEndpoint": "http://localhost:8080/datasink", "providerEndpoint": "http://localhost:9080"}'
res=$(curl --silent -X POST -d "$data" http://localhost:8081/gateway/transfer/request)
tpConsPid=$(echo "$res" | awk -F[\"\"] '{print $4}')

read -p "Transfer Process ID (Provider): " tpProvId
data='{"provider": true, "transferProcessId": "'$tpProvId'"}'
curl --silent -X POST -d "$data" http://localhost:9081/gateway/transfer/start

data='{"transferProcessId": "'$tpConsPid'", "reasons": ["outdated data", "consumer initiated"]}'
curl --silent -X POST -d "$data" http://localhost:8081/gateway/transfer/suspend

data='{"transferProcessId": "'$tpConsPid'"}'
curl --silent -X POST -d "$data" http://localhost:8081/gateway/transfer/start

data='{"provider": true, "transferProcessId": "'$tpProvId'", "reasons": ["time expired", "provider initiated"]}'
curl --silent -X POST -d "$data" http://localhost:9081/gateway/transfer/suspend

data='{"transferProcessId": "'$tpConsPid'", "reasons": ["provider suspended", "consumer terminated"]}'
curl --silent -X POST -d "$data" http://localhost:8081/gateway/transfer/terminate