#!/usr/bin/env bash

# change me 
export BLOCK_HEIGHT='2233392'
export TX='
[
  {
    "@type": "/ethermint.evm.v1.MsgEthereumTx",
    "data": {
      "@type": "/ethermint.evm.v1.LegacyTx",
      "nonce": "0",
      "gas_price": "5000000000000",
      "gas": "33578",
      "to": "0x6b1b50c2223eb31E0d4683b046ea9C6CB0D0ea4F",
      "value": "2100000000000000000000",
      "data": "xBzCcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACpjcm8xNHFxOHB2dGUzdm5yYXI2MjJ5ZjZkcTU0NzRtem4zcXlyMzIyNmYAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
      "v": "Vg==",
      "r": "2ZTq4BbxuHdgij/msRHAchwLVnCKrgg12JsCRQ625nU=",
      "s": "fPB/dwctr9bb1QjS2257IE5Tt4zf25Q/0VcftzXK2a8="
    },
    "size": 245,
    "hash": "0x928362f02ec187ff20efdd9c4b42bdddcfa9b356d6013fe69a74411f7d6c5975",
    "from": ""
  },
  {
    "@type": "/ethermint.evm.v1.MsgEthereumTx",
    "data": {
      "@type": "/ethermint.evm.v1.LegacyTx",
      "nonce": "559",
      "gas_price": "5000000000000",
      "gas": "215585",
      "to": "0x145677FC4d9b8F19B5D56d1820c48e0443049a30",
      "value": "0",
      "data": "OO0XOQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAABHVQ8QZgJ/v0AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAnVACuUBY8sYAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAoAAAAAAAAAAAAAAAACyiqI0ZTqT/OC+n5DP2wkrsZKUEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAGJN5bsAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAJd0nJth+HiogN/jEtJZSuB67XZWAAAAAAAAAAAAAAAAISMx4UNajfIwcV20wCsqOgq/jGE=",
      "v": "VQ==",
      "r": "CPYln4g86jWP0ZVCp+dlnXUvC5jGh/4YR9nL8ZmoFek=",
      "s": "NsU3+bWUAcnqVyJYR2CaaEDgxjb6Me/KJXcs61kNswo="
    },
    "size": 369,
    "hash": "0xa7bd6c4753558f2be23ba79dc7f7852ad86c4a8363127b20da79cfe9d73887b9",
    "from": ""
  }
]
'


# should be fixed
export BLOCKTIME=$(curl -s "https://rpc.cronos.org/block?height=$BLOCK_HEIGHT" | jq -r .result.block.header.time)
export ABI='[{"inputs":[{"internalType":"string","name":"recipient","type":"string"}],"name":"send_cro_to_crypto_org","outputs":[],"stateMutability":"payable","type":"function"}]'
export IBC_DENOM='transfer/channel-0/basecro'

go run main.go