#!/usr/bin/env bash

export ABI='[{"inputs":[{"internalType":"string","name":"recipient","type":"string"}],"name":"send_cro_to_crypto_org","outputs":[],"stateMutability":"payable","type":"function"}]'

export DENOM='transfer/channel-0/basecro'

export BLOCKTIME='2022-04-03T20:15:22.338515757Z'

export TX='
            {
              "@type": "/ethermint.evm.v1.MsgEthereumTx",
              "data": {
                "@type": "/ethermint.evm.v1.LegacyTx",
                "nonce": "0",
                "gas_price": "5000000000000",
                "gas": "27982",
                "to": "0x6b1b50c2223eb31E0d4683b046ea9C6CB0D0ea4F",
                "value": "660090000000000000",
                "data": "xBzCcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACpjcm8xeGg4emZjMDQyZXMyY2VsZjI3YzhmcDYzYzJnenBjcHB5YXFmN2oAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
                "v": "VQ==",
                "r": "be4QYrYWLnWz5YAcevirXJPZ4pzEqk7dXaArXEkvk9k=",
                "s": "H880+5qcXP3ImO1xu0BJEVgmUz8zuEZVCNiSf4hJpdw="
              },
              "size": 244,
              "hash": "0xb08c8cbc9151164aa42551929629f860bd6a19dd832295d338f481ef062cf6c7",
              "from": ""
            }
'

go run main.go