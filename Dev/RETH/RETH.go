package RETH

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"time"

	coreEntities "github.com/daoleno/uniswap-sdk-core/entities"
	"github.com/daoleno/uniswapv3-sdk/constants"
	"github.com/daoleno/uniswapv3-sdk/entities"
	"github.com/daoleno/uniswapv3-sdk/examples/helper"
	"github.com/daoleno/uniswapv3-sdk/periphery"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	BASE_TXN_GAS_LIMIT    uint64 = 21000 // gas units for simple txn
	TXN_DATA_NON_ZERO_GAS uint64 = 68
)

func GetMaxGasLimit(data_length uint64) uint64 {
	// TODO: Fix this calculation
	if data_length == 0 {
		return BASE_TXN_GAS_LIMIT
	}

	return BASE_TXN_GAS_LIMIT + TXN_DATA_NON_ZERO_GAS*data_length
}

func GetMaxFeePerGas(client *ethclient.Client, max_gas_tip uint64) *big.Int {
	gas_price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return big.NewInt(int64(2*gas_price.Uint64() + max_gas_tip))
}

// EIP-155
func SendLegacyTxn(url_to_connect string,
	sender_private_key_hex string,
	sender_wei_to_send int64,
	receiver_public_addr_hex string) {
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
		BASE_TXN_GAS_LIMIT,
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

func SendDynamicFeeTxn(client *ethclient.Client,
	sender_wallet *helper.Wallet,
	value *big.Int,
	receiver_address *common.Address,
	data []byte) {
	fmt.Println("Sender addr:", sender_wallet.PublicKey)
	fmt.Println("value:", value)
	fmt.Println("Receiver addr:", receiver_address)
	fmt.Println("data:", data)
	fmt.Println("data length:", len(data))

	suggested_max_gas_tip, err := client.SuggestGasTipCap(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	sender_nonce, err := client.PendingNonceAt(context.Background(), sender_wallet.PublicKey)
	if err != nil {
		log.Fatal(err)
	}

	// EIP-1559 Dynamic Fee Txn
	dynamic_fee_txn := types.NewTx(&types.DynamicFeeTx{
		Nonce:     sender_nonce,
		GasFeeCap: GetMaxFeePerGas(client, suggested_max_gas_tip.Uint64()),
		GasTipCap: suggested_max_gas_tip,
		Gas:       400_000, //GetMaxGasLimit(uint64(len(data))), TODO: Fix this hard-coded gas_limit
		To:        receiver_address,
		Value:     value,
		Data:      data})

	chain_id, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	london_signer := types.NewLondonSigner(chain_id)
	signed_txn, err := types.SignTx(dynamic_fee_txn, london_signer, sender_wallet.PrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.SendTransaction(context.Background(), signed_txn); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully sent Transaction:", signed_txn.Hash().Hex())
}

func SwapTokens(url_to_connect string,
	private_key string,
	source_token *coreEntities.Token,
	target_token *coreEntities.Token,
	amount_to_swap string,
	contract_v3_swaprouter_address string) {
	client, err := ethclient.Dial(url_to_connect)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize a Wallet with it's related Private & Public keys
	wallet := helper.InitWallet(private_key)
	if wallet == nil {
		log.Fatal("InitWallet failed")
	}

	// Initiate a Pool object
	pool, err := helper.ConstructV3Pool(client, source_token, target_token, uint64(constants.FeeMedium))
	if err != nil {
		log.Fatal(err)
	}

	// single trade input
	// single-hop exact input. Route represents a list of pools through which a swap can occur.
	swap_route, err := entities.NewRoute([]*entities.Pool{pool}, source_token, target_token)
	if err != nil {
		log.Fatal(err)
	}

	// Amount to swap
	swap_value := helper.FloatStringToBigInt(amount_to_swap, int(source_token.Decimals()))

	// Simulate a Trade from the SwapRoute
	trade, err := entities.FromRoute(swap_route, coreEntities.FromRawAmount(source_token, swap_value), coreEntities.ExactInput)
	if err != nil {
		log.Fatal(err)
	}

	// Print the Simulation result values, print the floor division for each of the amounts
	fmt.Println("Input:", trade.Swaps[0].InputAmount.Quotient())
	fmt.Println("Output:", trade.Swaps[0].OutputAmount.Wrapped().Quotient())

	// Max Slippage -> 0.5%
	max_slippage_tolerance := coreEntities.NewPercent(big.NewInt(50), big.NewInt(1000))

	// Deadline -> After 15 minutes from Now
	d := time.Now().Add(time.Minute * time.Duration(15)).Unix()
	deadline := big.NewInt(d)

	// Create the Swap Parameters
	params, err := periphery.SwapCallParameters([]*entities.Trade{trade}, &periphery.SwapOptions{
		SlippageTolerance: max_slippage_tolerance,
		Recipient:         wallet.PublicKey,
		Deadline:          deadline,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Params Wei to send = 0x%x\n", params.Value.String())

	swap_router_address := common.HexToAddress(contract_v3_swaprouter_address)

	SendDynamicFeeTxn(client, wallet, swap_value, &swap_router_address, params.Calldata)
}
