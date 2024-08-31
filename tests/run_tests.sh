#!/bin/bash

readFile() {
  while IFS= read -r line; do
    agrId=$line
  done < tmp.txt
  echo $agrId
}

printf "# Test 1: Consumer initiates the negotiation flow + simple transfer\n"
printf "  - Catalog flow: Create offer (P) -> Create dataset (P)\n"
printf "  - Negotiation flow: Request contract (C) -> Agree contract (P) -> Verify agreement (C) -> Finalize contract (P)\n"
printf "  - Transfer flow: Request transfer (C) -> Start transfer (P) -> Complete transfer (C)\n\n"

bash flows/simple_negotiation.sh > tmp.txt
agrId=$(readFile)
bash flows/simple_transfer.sh $agrId

printf "\n# Test 2: Provider initiates the negotiation flow + simple transfer\n"
printf "  - Catalog flow: Create offer (P) -> Create dataset (P)\n"
printf "  - Negotiation flow: Request contract (C) -> Offer contract (P) -> Accept offer (C) -> Agree contract (P) -> Verify agreement (C) -> Finalize contract (P)\n"
printf "  - Transfer flow: Request transfer (C) -> Start transfer (P) -> Complete transfer (C)\n\n"

bash flows/reverse_negotiation.sh > tmp.txt
agrId=$(readFile)
bash flows/simple_transfer.sh $agrId

printf "\n# Test 3: Consumer initiates the negotiation flow + suspend and terminate transfer\n"
printf "  - Catalog flow: Create offer (P) -> Create dataset (P)\n"
printf "  - Negotiation flow: Request contract (C) -> Agree contract (P) -> Verify agreement (C) -> Finalize contract (P)\n"
printf "  - Transfer flow: Request transfer (C) -> Start transfer (P) -> Suspend transfer (C) -> Start transfer (C) -> Suspend transfer (P) -> Terminate transfer (C)\n\n"

bash flows/simple_negotiation.sh > tmp.txt
agrId=$(readFile)
bash flows/terminate_transfer.sh $agrId

rm tmp.txt
