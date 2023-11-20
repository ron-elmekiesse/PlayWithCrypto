package RETH

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const GAS_LIMIT uint64 = 21000 // gas units

func SendEth(url_to_connect string, sender_private_key_hex string, sender_wei_to_send int64, receiver_public_addr_hex string) {
	fmt.Println("Sending ETH")

	client, err := ethclient.Dial(url_to_connect)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("We have a Connection!")

	sender_private_key, err := crypto.HexToECDSA(sender_private_key_hex)
	if err != nil {
		log.Fatal(err)
	}

	sender_interfaced_public_key := sender_private_key.Public()
	sender_public_key, ok := sender_interfaced_public_key.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed casting public key to ECDSA")
	}

	sender_public_address := crypto.PubkeyToAddress(*sender_public_key)

	sender_nonce, err := client.PendingNonceAt(context.Background(), sender_public_address)
	if err != nil {
		log.Fatal(err)
	}

	suggested_gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	receiver_public_address := common.HexToAddress(receiver_public_addr_hex)

	unsigned_transaction := types.NewTransaction(sender_nonce,
		receiver_public_address,
		big.NewInt(sender_wei_to_send),
		GAS_LIMIT,
		suggested_gas_price,
		nil)

	chain_id, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signed_transaction, err := types.SignTx(unsigned_transaction, types.NewEIP155Signer(chain_id), sender_private_key)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.SendTransaction(context.Background(), signed_transaction); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully sent Transaction:", signed_transaction.Hash().Hex())
}
