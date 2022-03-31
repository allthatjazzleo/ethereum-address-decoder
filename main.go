package main

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const abiJSON = `[{"inputs":[{"internalType":"string","name":"recipient","type":"string"}],"name":"send_cro_to_crypto_org","outputs":[],"stateMutability":"payable","type":"function"}]`
const blocktimeString = "2022-03-28T15:45:02.835813016Z"
const denom = "transfer/channel-0/basecro"
const methodName = "send_cro_to_crypto_org"

// Example Legacy tx:
const txJSON = `
{
	"@type": "/ethermint.evm.v1.MsgEthereumTx",
	"data": {
		"@type": "/ethermint.evm.v1.LegacyTx",
		"nonce": "332",
		"gas_price": "5000000000000",
		"gas": "33578",
		"to": "0x6b1b50c2223eb31E0d4683b046ea9C6CB0D0ea4F",
		"value": "102030243391546367224",
		"data": "xBzCcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACpjcm8xNzZxOGFtNmM4aHNrNHIyazR2c21xODBoYXBra3U1M3l5NnVoZm0AAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"v": "VQ==",
		"r": "iP76vlA2inv1keUSgS4WA9o3rJQn+H57mnDNJO/FeNY=",
		"s": "T8t898tl+P4aFmsop/KR2JgG8rt8EG7Znhh3IgFN7jo="
	},
	"size": 247,
	"hash": "0xcdd1b6fe327e5c17e11be724768452c41bc9e6906cb105ded82a42cdacbbaaef",
	"from": ""
}
`

type Tx struct {
	Data Data `json:"data"`
}

type Data struct {
	Nonce    string `json:"nonce"`
	GasPrice string `json:"gas_price"`
	Gas      string `json:"gas"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Data     string `json:"data"`
	V        string `json:"v"`
	R        string `json:"r"`
	S        string `json:"s"`
}

func main() {

	var tx Tx
	err := json.Unmarshal([]byte(txJSON), &tx)

	if err != nil {
		fmt.Println("An error occured: %v", err)
		os.Exit(1)
	}

	v, _ := base64.StdEncoding.DecodeString(tx.Data.V)
	v_h := hex.EncodeToString(v)
	r, _ := base64.StdEncoding.DecodeString(tx.Data.R)
	r_h := hex.EncodeToString(r)
	s, _ := base64.StdEncoding.DecodeString(tx.Data.S)
	s_h := hex.EncodeToString(s)

	va, _ := new(big.Int).SetString(v_h, 16)
	ra, _ := new(big.Int).SetString(r_h, 16)
	sa, _ := new(big.Int).SetString(s_h, 16)

	// to decode data, translate base64 to hex
	data_base64, _ := base64.StdEncoding.DecodeString(tx.Data.Data)
	data_h := hex.EncodeToString(data_base64)
	data, _ := hex.DecodeString(data_h)

	address := common.HexToAddress(tx.Data.To)

	vala, _ := new(big.Int).SetString(tx.Data.Value, 10)
	nonce, _ := strconv.Atoi(tx.Data.Nonce)
	gas, _ := strconv.Atoi(tx.Data.Gas)
	gas_price, _ := strconv.Atoi(tx.Data.GasPrice)
	ethTx := types.NewTx(&types.LegacyTx{
		Nonce:    uint64(nonce),
		To:       &address,
		Value:    vala,
		Gas:      uint64(gas),
		GasPrice: big.NewInt(int64(gas_price)),
		Data:     data,
		V:        va,
		R:        ra,
		S:        sa,
	})

	signer := types.NewEIP155Signer(big.NewInt(25))
	sender, _ := signer.Sender(ethTx)

	// load contract ABI
	abi, err := abi.JSON(strings.NewReader(abiJSON))
	if err != nil {
		log.Fatal(err)
	}
	recipient, _ := abi.Methods[methodName].Inputs.Unpack(data[4:])

	// get ibc timeout timestamp
	blocktime, err := time.Parse(time.RFC3339Nano, blocktimeString)
	if err != nil {
		fmt.Println("Could not parse time:", err)
	}

	amount := tx.Data.Value[0 : len(tx.Data.Value)-10]

	// get base64 encoded ibc data
	ibcData := fmt.Sprintf(`{"amount":"%v","denom":"%v","receiver":"%v","sender":"%v"}`, amount, denom, recipient[0], sender)

	// get base64 encoded ibc data
	ibcDataBase64 := base64.StdEncoding.EncodeToString([]byte(ibcData))

	fmt.Printf("sender address: 0x%x\n", sender)
	fmt.Printf("recipient address: %v\n", recipient[0])
	fmt.Printf("amount: %v\n", amount)
	fmt.Printf("timeout timestamp: %v\n", blocktime.UnixNano()+86400000000000)

	fmt.Printf("ibc data: %v\n", ibcDataBase64)
}
