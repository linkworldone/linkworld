package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"linkworld-backend/internal/blockchain/bindings"
)

// ESimInfo holds eSIM activation configuration retrieved from chain
type ESimInfo struct {
	User           common.Address
	TokenId        *big.Int
	ActivationCode string
	SmDpAddress    string
}

// GetESimInfo retrieves activation code and SM-DP address for a tokenId (on-chain lookup)
func (c *Client) GetESimInfo(ctx context.Context, tokenId *big.Int) (*ESimInfo, error) {
	addr, ok := c.contracts["TrafficCardNFT"]
	if !ok {
		return nil, fmt.Errorf("TrafficCardNFT contract address not configured")
	}

	nft, err := bindings.NewTrafficCardNFTCaller(addr, c.backend)
	if err != nil {
		return nil, fmt.Errorf("failed to create TrafficCardNFT caller: %w", err)
	}

	activationCode, err := nft.GetActivationCode(&bind.CallOpts{Context: ctx}, tokenId)
	if err != nil {
		return nil, fmt.Errorf("failed to get activation code: %w", err)
	}

	smDpAddress, err := nft.SmDpAddress(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, fmt.Errorf("failed to get SM-DP address: %w", err)
	}

	return &ESimInfo{
		TokenId:        tokenId,
		ActivationCode: activationCode,
		SmDpAddress:    smDpAddress,
	}, nil
}

// ESimRedeemedHandler callback for ESimRedeemed event processing
type ESimRedeemedHandler func(info ESimInfo)

// WatchESimRedeemed subscribes to ESimRedeemed events from TrafficCardNFT contract.
// Returns a cancel function to stop the subscription.
// Usage: cancel := client.WatchESimRedeemed(func(info ESimInfo) { db.UpdateSimActivation(info) })
func (c *Client) WatchESimRedeemed(handler ESimRedeemedHandler) (func(), error) {
	addr, ok := c.contracts["TrafficCardNFT"]
	if !ok {
		return nil, fmt.Errorf("TrafficCardNFT contract address not configured")
	}
	nft, err := bindings.NewTrafficCardNFTFilterer(addr, c.backend)
	if err != nil {
		return nil, fmt.Errorf("failed to create TrafficCardNFT filterer: %w", err)
	}
	sink := make(chan *bindings.TrafficCardNFTESimRedeemed, 10)
	sub, err := nft.WatchESimRedeemed(nil, sink, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to watch ESimRedeemed: %w", err)
	}
	go func() {
		for {
			ev, ok := <-sink
			if !ok {
				return
			}
			if ev != nil {
				handler(ESimInfo{
					User:           ev.User,
					TokenId:        ev.TokenId,
					ActivationCode: ev.ActivationCode,
					SmDpAddress:    ev.SmDpAddress,
				})
			}
		}
	}()
	log.Printf("ESimRedeemed event listener started for TrafficCardNFT at %s", addr.Hex())
	return sub.Unsubscribe, nil
}

// FormatQRCode generates QR code content string for eSIM profile download
// Format: LPKG://<activationCode>@<smDpAddress>
func (c *Client) FormatQRCode(info ESimInfo) string {
	return fmt.Sprintf("LPKG://%s@%s", info.ActivationCode, info.SmDpAddress)
}
