#!/bin/bash

readFile() {
  while IFS= read -r line; do
    agrId=$line
  done < tmp.txt
  echo $agrId
}

printf "# Test 1: Consumer initiates the negotiation flow + simple transfer\n"

bash flows/simple_negotiation.sh > tmp.txt
agrId=$(readFile)
bash flows/simple_transfer.sh $agrId

printf "# Test 2: Provider initiates the negotiation flow + simple transfer\n"

bash flows/reverse_negotiation.sh > tmp.txt
agrId=$(readFile)
bash flows/simple_transfer.sh $agrId

printf "# Test 3: Consumer initiates the negotiation flow + suspend and terminate transfer\n"

bash flows/simple_negotiation.sh > tmp.txt
agrId=$(readFile)
bash flows/terminate_transfer.sh $agrId

rm tmp.txt
