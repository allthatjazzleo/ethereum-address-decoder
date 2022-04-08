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

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

var abiJSON = `[{"inputs":[{"internalType":"string","name":"recipient","type":"string"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"send_to_ibc","outputs":[],"stateMutability":"nonpayable","type":"function"}]`

var blocktimeString = "2022-04-03T21:27:20.094446044Z"

var denom = "transfer/channel-0/basecro"

// Example Legacy tx:
var txJSON = `
[
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
]
`

const (
	// Bech32Prefix defines the Bech32 prefix used for Cronos Accounts
	Bech32Prefix = "crc"

	// Bech32PrefixAccAddr defines the Bech32 prefix of an account's address
	Bech32PrefixAccAddr = Bech32Prefix
	// Bech32PrefixAccPub defines the Bech32 prefix of an account's public key
	Bech32PrefixAccPub = Bech32Prefix + sdk.PrefixPublic
)

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

type DecodedData struct {
	Sender           string      `json:"sender"`
	Recipient        string      `json:"recipient"`
	Amount           interface{} `json:"amount"`
	TimeoutTimestamp int64       `json:"timeout_timestamp"`
	IbcData          string      `json:"ibc_data"`
}

func main() {

	// set var from env if any
	if _abi := os.Getenv("ABI"); _abi != "" {
		abiJSON = _abi
	}
	if _denom := os.Getenv("IBC_DENOM"); _denom != "" {
		denom = _denom
	}
	if _blocktime := os.Getenv("BLOCKTIME"); _blocktime != "" {
		blocktimeString = _blocktime
	}
	if _tx := os.Getenv("TX"); _tx != "" {
		txJSON = _tx
	}

	var txs []Tx
	err := json.Unmarshal([]byte(txJSON), &txs)

	if err != nil {
		fmt.Println("An error occured: %v", err)
		os.Exit(1)
	}

	// get ibc timeout timestamp
	blocktime, err := time.Parse(time.RFC3339Nano, blocktimeString)
	if err != nil {
		fmt.Println("Could not parse time:", err)
	}

	for _, tx := range txs {
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
		var methodName string
		for k, _ := range abi.Methods {
			methodName = k
		}
		input, _ := abi.Methods[methodName].Inputs.Unpack(data[4:])
		if len(input) == 0 {
			continue
		}

		var amount interface{}
		if methodName == "send_cro_to_crypto_org" {
			amount = tx.Data.Value[0 : len(tx.Data.Value)-10]
		} else {
			amount = input[1]
		}

		// get base64 encoded ibc data
		config := sdk.GetConfig()
		config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
		ibcData := fmt.Sprintf(`{"amount":"%v","denom":"%v","receiver":"%v","sender":"%s"}`, amount, denom, input[0], sdk.AccAddress(sender.Bytes()))

		// get base64 encoded ibc data
		ibcDataBase64 := base64.StdEncoding.EncodeToString([]byte(ibcData))

		decodedData := DecodedData{
			Sender:           sdk.AccAddress(sender.Bytes()).String(),
			Recipient:        input[0].(string),
			Amount:           amount,
			TimeoutTimestamp: blocktime.UnixNano() + 86400000000000,
			IbcData:          ibcDataBase64,
		}
		val, err := json.MarshalIndent(decodedData, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", string(val))
	}

}
