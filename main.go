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

const abiJSON = `[{"inputs":[{"internalType":"string","name":"recipient","type":"string"}],"name":"send_cro_to_crypto_org","outputs":[],"stateMutability":"payable","type":"function"}]`
const blocktimeString = "2022-04-02T03:21:33.94019933Z"
const denom = "transfer/channel-0/basecro"
const methodName = "send_cro_to_crypto_org"

// Example Legacy tx:
const txJSON = `
{
	"@type": "/ethermint.evm.v1.MsgEthereumTx",
	"data": {
		"@type": "/ethermint.evm.v1.LegacyTx",
		"nonce": "91",
		"gas_price": "5000000000000",
		"gas": "33578",
		"to": "0x6b1b50c2223eb31E0d4683b046ea9C6CB0D0ea4F",
		"value": "20000000000000000000",
		"data": "xBzCcAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACpjcm8xNjZxc2cwenFlOGM2bmE3ZnZjem1ydmprdzRoemNoemV2M3R4ZzgAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"v": "VQ==",
		"r": "14gWo3532OdP+rZmqQ+iAEFTPftdaikzu6AINO2iAyI=",
		"s": "Q9Qjhgpx/PC55nsnp2fPGoXbaqNA1g1rqKSWY6XvlNk="
	},
	"size": 245,
	"hash": "0xf6c22dbd922b437c8e30263ed1ff6253909931f88f87f70d88b3e6ae91930308",
	"from": ""
}
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
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(Bech32PrefixAccAddr, Bech32PrefixAccPub)
	ibcData := fmt.Sprintf(`{"amount":"%v","denom":"%v","receiver":"%v","sender":"%s"}`, amount, denom, recipient[0], sdk.AccAddress(sender.Bytes()))

	// get base64 encoded ibc data
	ibcDataBase64 := base64.StdEncoding.EncodeToString([]byte(ibcData))

	fmt.Printf("sender address: %s\n", sdk.AccAddress(sender.Bytes()))
	fmt.Printf("recipient address: %v\n", recipient[0])
	fmt.Printf("amount: %v\n", amount)
	fmt.Printf("timeout timestamp: %v\n", blocktime.UnixNano()+86400000000000)

	fmt.Printf("ibc data: %v\n", ibcDataBase64)
}
