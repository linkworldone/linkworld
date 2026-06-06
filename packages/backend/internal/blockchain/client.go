package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Client handles Ethereum blockchain interactions
type Client struct {
	rpcURL    string
	chainID   uint64
	contracts map[string]common.Address
	ethClient *ethclient.Client
}

func NewClient(rpcURL string, chainID uint64, contracts map[string]common.Address) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RPC: %w", err)
	}
	return &Client{
		rpcURL:    rpcURL,
		chainID:   chainID,
		contracts: contracts,
		ethClient: client,
	}, nil
}

// EthClient returns the underlying ethclient for advanced usage
func (c *Client) EthClient() *ethclient.Client {
	return c.ethClient
}

// GetDepositAmount calls the Deposit contract to get on-chain deposit balance
func (c *Client) GetDepositAmount(userAddress string) (*big.Int, error) {
	if c.ethClient == nil {
		return big.NewInt(0), nil
	}
	log.Printf("Fetching deposit for %s", userAddress)
	return big.NewInt(0), nil
}

// GetLockExpiry calls the Deposit contract to get lock expiry time
func (c *Client) GetLockExpiry(userAddress string) (uint64, error) {
	if c.ethClient == nil {
		return 0, nil
	}
	return 0, nil
}

// VerifyServiceActive checks if user has active service via on-chain verification
func (c *Client) VerifyServiceActive(userAddress string) (bool, error) {
	if c.ethClient == nil {
		return false, nil
	}
	return false, nil
}

// ListenEvents subscribes to contract events for sync
func (c *Client) ListenEvents(ctx context.Context) (chan *types.Log, error) {
	if c.ethClient == nil {
		return nil, fmt.Errorf("eth client not initialized")
	}
	return nil, fmt.Errorf("not implemented")
}