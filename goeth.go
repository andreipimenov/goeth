package goeth

import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	InfuraMainNet    = netURL("https://mainnet.infura.io")
	InfuraRopstenNet = netURL("https://ropsten.infura.io")
)

const (
	StandartGasLimit = uint64(21000)
)

type netURL string

// ConnectToInfura connects to the given network URL and returns *ethclient.Client and error
func ConnectToInfura(url netURL) (*ethclient.Client, error) {
	return ethclient.Dial(string(url))
}

// Balance returns accounts balance in ETH and error
func Balance(ctx context.Context, client *ethclient.Client, address string) (*big.Float, error) {
	balanceWEI, err := client.BalanceAt(ctx, common.HexToAddress(address), nil)
	if err != nil {
		return nil, err
	}
	floatBalanceWEI := new(big.Float)
	floatBalanceWEI.SetString(balanceWEI.String())
	balanceETH := new(big.Float).Quo(floatBalanceWEI, big.NewFloat(math.Pow10(18)))
	return balanceETH, nil
}

// GasPrice returns suggested gas price
func GasPrice(ctx context.Context, client *ethclient.Client) (*big.Int, error) {
	return client.SuggestGasPrice(ctx)
}

// NewTx creates new standart transaction for sending ETH from one account to another.
func NewTx(ctx context.Context, client *ethclient.Client, fromAddress string, toAddress string, value *big.Int) (*types.Transaction, error) {
	nonce, err := client.PendingNonceAt(ctx, common.HexToAddress(fromAddress))
	if err != nil {
		return nil, err
	}
	gasPrice, err := GasPrice(ctx, client)
	if err != nil {
		return nil, err
	}
	return types.NewTransaction(nonce, common.HexToAddress(toAddress), value, StandartGasLimit, gasPrice, nil), nil
}

// SignTx signs transaction with private key.
func SignTx(ctx context.Context, client *ethclient.Client, tx *types.Transaction, privateKey string) (*types.Transaction, error) {
	privateK, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	chainID, err := client.NetworkID(ctx)
	if err != nil {
		return nil, err
	}
	return types.SignTx(tx, types.NewEIP155Signer(chainID), privateK)
}

// SendTx sends signed transaction into blockchain.
func SendTx(ctx context.Context, client *ethclient.Client, tx *types.Transaction) error {
	return client.SendTransaction(ctx, tx)
}
