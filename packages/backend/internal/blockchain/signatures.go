package blockchain

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Event signatures (keccak256 hashes)
var (
	// UserRegistered(address,string,uint256) - computed at runtime
	UserRegisteredTopic = crypto.Keccak256Hash([]byte("UserRegistered(address,string,uint256)"))
)

// DepositMadeTopic for Deposit contract
var DepositMadeTopic = crypto.Keccak256Hash([]byte("DepositMade(address,uint256)"))

// BillCreatedTopic for Payment contract
var BillCreatedTopic = crypto.Keccak256Hash([]byte("BillCreated(uint256,address,uint256,uint256,uint256)"))

// BillPaidTopic for Payment contract
var BillPaidTopic = crypto.Keccak256Hash([]byte("BillPaid(uint256,address,uint256,uint256)"))

// TrafficCardMintedTopic for Deposit contract
var TrafficCardMintedTopic = crypto.Keccak256Hash([]byte("TrafficCardMinted(address,uint256,uint256)"))

func init() {
	_ = common.HexToHash("0x") // verify common import works
}

// ComputeUserRegisteredSig returns the event signature for UserRegistered
func ComputeUserRegisteredSig() string {
	return hex.EncodeToString(UserRegisteredTopic.Bytes())
}