#!/usr/bin/env bash

# change me 
export BLOCK_HEIGHT='2193925'
export TX='
[
  {
    "@type": "/ethermint.evm.v1.MsgEthereumTx",
    "data": {
      "@type": "/ethermint.evm.v1.LegacyTx",
      "nonce": "13",
      "gas_price": "5000000000000",
      "gas": "48880",
      "to": "0xB888d8Dd1733d72681b30c00ee76BDE93ae7aa93",
      "value": "0",
      "data": "pRXLQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAehIAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALWNvc21vczFnbnJ4YzNnZGxqcXJ0MHFhNWd5OGtsd2xncjBybDZ3ZDJka3M2NAAAAAAAAAAAAAAAAAAAAAAAAAA=",
      "v": "Vg==",
      "r": "hID0CIUtn8mE+86LNYJGbmUtXpApV/75UxlV0ossihM=",
      "s": "DVsaa6Xdm2s8gOkGa97EVRIhUmzTPyYlp1oZiqwuMs4="
    },
    "size": 269,
    "hash": "0x7dd322c3da5d87534076cb26444c95219886cb724493a7333724cce5fb2e0ca1",
    "from": ""
  }
]
'


# should be fixed
export BLOCKTIME=$(curl -s "https://rpc.cronos.org/block?height=$BLOCK_HEIGHT" | jq -r .result.block.header.time)
export ABI='[{"inputs":[{"internalType":"string","name":"recipient","type":"string"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"send_to_ibc","outputs":[],"stateMutability":"nonpayable","type":"function"}]'
export IBC_DENOM='transfer/channel-1/uatom'

go run main.go