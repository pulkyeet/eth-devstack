#!/bin/bash

# The 7 correct signer addresses (from Geth account creation)
SIGNERS=(
  "c6BCfea104dBC7cB26CC990e33e58fEBd62fa533"
  "529166Fe40a2912d369881D8d5c4E63f75e9e1f8"
  "52688b264bd5FFA39A3df49D4BcbE91CC561D9fb"
  "5F67D2ACab6115FD6E337a852fFa7f04e9ADb7ED"
  "44Ec5ABDB6E43C42Da49d960C31af5d550f413c9"
  "0b18dc66d9837E7eED25B6a8CFFBe16a82eA1599"
  "eEf5100728b68504F635150224766861A2888B84"
)

# Build extraData: 32 bytes vanity + addresses + 65 bytes seal
VANITY="0000000000000000000000000000000000000000000000000000000000000000"
SEAL="0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"

EXTRA_DATA="0x${VANITY}"
for addr in "${SIGNERS[@]}"; do
  EXTRA_DATA="${EXTRA_DATA}${addr}"
done
EXTRA_DATA="${EXTRA_DATA}${SEAL}"

# Create genesis.json with correct extraData
cat > genesis.json << GENESIS
{
  "config": {
    "chainId": 1337,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0,
    "berlinBlock": 0,
    "londonBlock": 0,
    "clique": {
      "period": 2,
      "epoch": 30000
    }
  },
  "difficulty": "1",
  "gasLimit": "15000000",
  "extradata": "${EXTRA_DATA}",
  "alloc": {}
}
GENESIS

echo "Genesis created with extraData: ${EXTRA_DATA}"
