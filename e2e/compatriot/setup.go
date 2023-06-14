package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	coretypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/magefile/mage/sh"
	"pkg.berachain.dev/polaris/eth/core/types"
)

const POLARIS_RPC = "http://localhost:8545"
const TESTS = "./tests.json"

var client *ethclient.Client

var txHashes []common.Hash

// ConnectToClient connects to an Ethereum client and returns the client instance
func ConnectToClient(url string) (*ethclient.Client, error) {
	client, err := ethclient.Dial(url)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// CreateAccount generates a new Ethereum account and returns the private key and address
func createAccount() (*ecdsa.PrivateKey, common.Address, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return nil, common.Address{}, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, common.Address{}, fmt.Errorf("error casting public key to ECDSA")
	}
	address := crypto.PubkeyToAddress(*publicKeyECDSA)
	return privateKey, address, nil
}

// SignTransaction signs the given transaction with the provided private key and returns the signed transaction object
func signTransaction(tx *types.Transaction, privateKey *ecdsa.PrivateKey) (*types.Transaction, error) {
	chainID, err := getChainID()
	if err != nil {
		return nil, err
	}

	signedTx, err := types.SignTx(tx, coretypes.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return nil, err
	}

	return signedTx, nil
}

// getChainID retrieves the current chain ID from an Ethereum client
func getChainID() (*big.Int, error) {
	client, err := rpc.Dial("http://localhost:8545") // Replace with your Ethereum client URL
	if err != nil {
		return nil, err
	}

	var chainID *big.Int
	err = client.CallContext(context.Background(), &chainID, "eth_chainId")
	if err != nil {
		return nil, err
	}

	return chainID, nil
}

// startPolarisChain starts the Polaris chain
func startPolarisChain() error {
	return sh.RunV("./cosmos/init.sh")
}

func buildTx(address common.Address) *coretypes.Transaction {
	// Build a transaction
	nonce, err := client.PendingNonceAt(context.Background(), address)
	if err != nil {
	}

	toAddress := common.HexToAddress("0x00000000000000000000000000000000DeaDBeef") // Replace with the recipient's Ethereum address
	value := big.NewInt(1000000000000000000)                                       // 1 ETH in wei
	gasLimit := uint64(21000)                                                      // Standard gas limit for a simple transaction
	data := []byte{}                                                               // Optional data for contract interactions

	return coretypes.NewTransaction(nonce, toAddress, value, gasLimit, big.NewInt(0), data), nil

}

// sendTx sends a transaction to the deadbeef address and returns its hash
func sendTx() (common.Hash, error) {
	client, err := ConnectToClient("http://localhost:8545") // Replace with your Ethereum client URL
	if err != nil {
		log.Fatal(err)
	}

	// Create a new account
	privateKey, address, err := createAccount()
	if err != nil {
		log.Fatal(err)
	}

	tx := buildTx(address)

	// Sign the transaction
	signedTx, err := signTransaction(tx, privateKey)
	if err != nil {
		log.Fatal(err)
	}

	// Send the transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SENT TRANSACTIONS TO NETWORK")
	return signedTx.Hash(), nil
}

// setup starts up the chain and spams the transactions
func setup() error {
	// init the chain
	// spam tx
	if err := startPolarisChain(); err != nil {
		return fmt.Errorf("failed to start the chain")
	}
	return nil
}

// submitTransactionsToNetwork submits transactions to the network and returns all the txHashes
func submitTransactionsToNetwork() []common.Hash {
	for i := 0; i < 100; i++ {

		txHash, err := sendTx()
		if err != nil {
			log.Fatal(err)
		}
		txHashes = append(txHashes, txHash)
	}
	return txHashes
}

func queryTheFkingChain() {

	for _, txHash := range txHashes {
		/*
		   sendTx() returns hash of the send transaction

		   GetReceiptsByHash() returns receipts of the transaction

		   get BlockNumber and BlockHash from the Receipt

		   Then call GetBlockByNumber() with the block number

		   Then call GetBlockByHash() with the block hash

		   Then call GetTransactionByHash() with the transaction hash

		   Then call GetReceiptsByHash() with the transaction hash

		   on the first run, these will all work beacuse of the cache

		   then when we stop the node, nuke the cache, and run again, these will all fail because no more cache and historical plugin gone

		*/
		var txReceipt string
		err := client.Client().CallContext(context.Background(), &txReceipt, "eth_getTransactionReceipt", txHash.String())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("txReceippt: ", txReceipt)
	}
}

func main() {
	setup()
	submitTransactionsToNetwork()
	queryTheFkingChain()
}
