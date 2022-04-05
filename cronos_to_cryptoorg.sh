#!/usr/bin/env bash

# change me 
export BLOCK_HEIGHT='2193925'
export TX='
            {
              "@type": "/ethermint.evm.v1.MsgEthereumTx",
              "data": {
                "@type": "/ethermint.evm.v1.LegacyTx",
                "nonce": "0",
                "gas_price": "5000000000000",
                "gas": "33578",
                "to": "0x6b1b50c2223eb31E0d4683b046ea9C6CB0D0ea4F",
                "value": "185592978360000000000",
                "data": "xBzCcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACpjcm8xZzB5amw3dWVydTdsYTR1M2U1Nnl3ZjdnMDh0ZXZmM21ndDBjbW0AAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
                "v": "Vg==",
                "r": "fayuevvXuPmKjcsj6i6dvIUeiKoHX+8TAsfW2wi+uxo=",
                "s": "AxZPhgjFNcIHI07EpZy2M7yqQIfy8f1WYeNlK6NFP4M="
              },
              "size": 245,
              "hash": "0x76da9033b17ae43e7d4776960fda06f1c7635511902d7747ad849365eec228f3",
              "from": ""
            }
'


# should be fixed
export BLOCKTIME=$(curl -s "https://rpc.cronos.org/block?height=$BLOCK_HEIGHT" | jq -r .result.block.header.time)
export ABI='[{"inputs":[{"internalType":"string","name":"recipient","type":"string"}],"name":"send_cro_to_crypto_org","outputs":[],"stateMutability":"payable","type":"function"}]'
export DENOM='transfer/channel-0/basecro'

go run main.go