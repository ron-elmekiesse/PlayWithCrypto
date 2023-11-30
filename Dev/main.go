package main

import (
	"fmt"

	"PlayWithCrypto.com/RETH"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/ethereum/go-ethereum/common"
)

const (
	URL_TO_CONNECT          string = "https://goerli.infura.io/<ENTER_API_KEY>"
	SENDER_PRIVATE_KEY      string = "<ENTER_SENDER_PRIVATE_KEY>"
	WEI_TO_SEND             int64  = 0
	RECEIVER_PUBLIC_ADDRESS string = "<ENTER_RECEIVER_PUBLIC_ADDRESS>"

	GOERLI_CHAIN_ID                       uint   = 5
	GOERLI_WETH_CONTRACT_ADDRESS          string = "0xB4FBF271143F4FBf7B91A5ded31805e42b2208d6"
	GOERLI_USDC_CONTRACT_ADDRESS          string = "0x07865c6E87B9F70255377e024ace6630C1Eaa37F"
	GOERLI_CONTRACT_V3_SWAPROUTER_ADDRESS string = "0xE592427A0AEce92De3Edee1F18E0157C05861564"
)

var (
	GOERLI_WETH_TOKEN = coreEntities.NewToken(GOERLI_CHAIN_ID, common.HexToAddress(GOERLI_WETH_CONTRACT_ADDRESS), 18, "WETH", "Wrapped Ether")
	GOERLI_USDC_TOKEN = coreEntities.NewToken(GOERLI_CHAIN_ID, common.HexToAddress(GOERLI_USDC_CONTRACT_ADDRESS), 6, "USDC", "USD//C")
)

func send_eth() {
	fmt.Println(URL_TO_CONNECT)
	fmt.Println(SENDER_PRIVATE_KEY, "(", WEI_TO_SEND, ")")
	fmt.Println("\t|")
	fmt.Println("\tV")
	fmt.Println(RECEIVER_PUBLIC_ADDRESS)

	RETH.SendLegacyTxn(URL_TO_CONNECT, SENDER_PRIVATE_KEY, WEI_TO_SEND, RECEIVER_PUBLIC_ADDRESS)
}

func swap_weth_to_usdc() {
	RETH.SwapTokens(URL_TO_CONNECT, SENDER_PRIVATE_KEY, GOERLI_WETH_TOKEN, GOERLI_USDC_TOKEN, "0.00001", GOERLI_CONTRACT_V3_SWAPROUTER_ADDRESS)
}

func main() {
	swap_weth_to_usdc()
}
