package main

import (
	"fmt"

	"PlayWithCrypto.com/RETH"
)

const (
	URL_TO_CONNECT          string = "https://sepolia.infura.io/<ENTER_API_KEY>"
	SENDER_PRIVATE_KEY      string = "<ENTER_SENDER_PRIVATE_KEY>"
	WEI_TO_SEND             int64  = 0
	RECEIVER_PUBLIC_ADDRESS string = "<ENTER_RECEIVER_PUBLIC_ADDRESS>"
)

func main() {
	fmt.Println(URL_TO_CONNECT)
	fmt.Println(SENDER_PRIVATE_KEY, "(", WEI_TO_SEND, ")")
	fmt.Println("\t|")
	fmt.Println("\tV")
	fmt.Println(RECEIVER_PUBLIC_ADDRESS)

	RETH.SendEth(URL_TO_CONNECT, SENDER_PRIVATE_KEY, WEI_TO_SEND, RECEIVER_PUBLIC_ADDRESS)
}
