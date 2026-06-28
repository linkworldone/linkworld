// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// ITrafficCardNFTCardInfo is an auto generated low-level Go binding around an user-defined struct.
type ITrafficCardNFTCardInfo struct {
	DataAmount  *big.Int
	CreatedAt   *big.Int
	IsDestroyed bool
}

// TrafficCardNFTMetaData contains all meta data concerning the TrafficCardNFT contract.
var TrafficCardNFTMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_toTokenId\",\"type\":\"uint256\"}],\"name\":\"BatchMetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dataAmount\",\"type\":\"uint256\"}],\"name\":\"CardDestroyed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dataAmount\",\"type\":\"uint256\"}],\"name\":\"CardMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"activationCode\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"smDpAddress\",\"type\":\"string\"}],\"name\":\"ESimRedeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"daysCount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"}],\"name\":\"SimRedeemed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEDUCTION_VALIDITY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getActivationCode\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getCardInfo\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"dataAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isDestroyed\",\"type\":\"bool\"}],\"internalType\":\"structITrafficCardNFT.CardInfo\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserCardCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"dataAmount\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"tokenURI_\",\"type\":\"string\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"to\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"dataAmounts\",\"type\":\"uint256[]\"},{\"internalType\":\"string[]\",\"name\":\"tokenURIs_\",\"type\":\"string[]\"}],\"name\":\"mintBatch\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"tokenIds\",\"type\":\"uint256[]\"}],\"name\":\"redeemForSim\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"daysCount\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_deposit\",\"type\":\"address\"}],\"name\":\"setDepositContract\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_addr\",\"type\":\"string\"}],\"name\":\"setSmDpAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"smDpAddress\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// TrafficCardNFTABI is the input ABI used to generate the binding from.
// Deprecated: Use TrafficCardNFTMetaData.ABI instead.
var TrafficCardNFTABI = TrafficCardNFTMetaData.ABI

// TrafficCardNFT is an auto generated Go binding around an Ethereum contract.
type TrafficCardNFT struct {
	TrafficCardNFTCaller     // Read-only binding to the contract
	TrafficCardNFTTransactor // Write-only binding to the contract
	TrafficCardNFTFilterer   // Log filterer for contract events
}

// TrafficCardNFTCaller is an auto generated read-only Go binding around an Ethereum contract.
type TrafficCardNFTCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrafficCardNFTTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TrafficCardNFTTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrafficCardNFTFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TrafficCardNFTFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TrafficCardNFTSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TrafficCardNFTSession struct {
	Contract     *TrafficCardNFT   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TrafficCardNFTCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TrafficCardNFTCallerSession struct {
	Contract *TrafficCardNFTCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// TrafficCardNFTTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TrafficCardNFTTransactorSession struct {
	Contract     *TrafficCardNFTTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// TrafficCardNFTRaw is an auto generated low-level Go binding around an Ethereum contract.
type TrafficCardNFTRaw struct {
	Contract *TrafficCardNFT // Generic contract binding to access the raw methods on
}

// TrafficCardNFTCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TrafficCardNFTCallerRaw struct {
	Contract *TrafficCardNFTCaller // Generic read-only contract binding to access the raw methods on
}

// TrafficCardNFTTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TrafficCardNFTTransactorRaw struct {
	Contract *TrafficCardNFTTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTrafficCardNFT creates a new instance of TrafficCardNFT, bound to a specific deployed contract.
func NewTrafficCardNFT(address common.Address, backend bind.ContractBackend) (*TrafficCardNFT, error) {
	contract, err := bindTrafficCardNFT(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFT{TrafficCardNFTCaller: TrafficCardNFTCaller{contract: contract}, TrafficCardNFTTransactor: TrafficCardNFTTransactor{contract: contract}, TrafficCardNFTFilterer: TrafficCardNFTFilterer{contract: contract}}, nil
}

// NewTrafficCardNFTCaller creates a new read-only instance of TrafficCardNFT, bound to a specific deployed contract.
func NewTrafficCardNFTCaller(address common.Address, caller bind.ContractCaller) (*TrafficCardNFTCaller, error) {
	contract, err := bindTrafficCardNFT(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTCaller{contract: contract}, nil
}

// NewTrafficCardNFTTransactor creates a new write-only instance of TrafficCardNFT, bound to a specific deployed contract.
func NewTrafficCardNFTTransactor(address common.Address, transactor bind.ContractTransactor) (*TrafficCardNFTTransactor, error) {
	contract, err := bindTrafficCardNFT(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTTransactor{contract: contract}, nil
}

// NewTrafficCardNFTFilterer creates a new log filterer instance of TrafficCardNFT, bound to a specific deployed contract.
func NewTrafficCardNFTFilterer(address common.Address, filterer bind.ContractFilterer) (*TrafficCardNFTFilterer, error) {
	contract, err := bindTrafficCardNFT(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTFilterer{contract: contract}, nil
}

// bindTrafficCardNFT binds a generic wrapper to an already deployed contract.
func bindTrafficCardNFT(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TrafficCardNFTMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TrafficCardNFT *TrafficCardNFTRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TrafficCardNFT.Contract.TrafficCardNFTCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TrafficCardNFT *TrafficCardNFTRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.TrafficCardNFTTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TrafficCardNFT *TrafficCardNFTRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.TrafficCardNFTTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TrafficCardNFT *TrafficCardNFTCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TrafficCardNFT.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TrafficCardNFT *TrafficCardNFTTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TrafficCardNFT *TrafficCardNFTTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.contract.Transact(opts, method, params...)
}

// DEDUCTIONVALIDITY is a free data retrieval call binding the contract method 0xf2551a99.
//
// Solidity: function DEDUCTION_VALIDITY() view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTCaller) DEDUCTIONVALIDITY(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "DEDUCTION_VALIDITY")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DEDUCTIONVALIDITY is a free data retrieval call binding the contract method 0xf2551a99.
//
// Solidity: function DEDUCTION_VALIDITY() view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTSession) DEDUCTIONVALIDITY() (*big.Int, error) {
	return _TrafficCardNFT.Contract.DEDUCTIONVALIDITY(&_TrafficCardNFT.CallOpts)
}

// DEDUCTIONVALIDITY is a free data retrieval call binding the contract method 0xf2551a99.
//
// Solidity: function DEDUCTION_VALIDITY() view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) DEDUCTIONVALIDITY() (*big.Int, error) {
	return _TrafficCardNFT.Contract.DEDUCTIONVALIDITY(&_TrafficCardNFT.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _TrafficCardNFT.Contract.UPGRADEINTERFACEVERSION(&_TrafficCardNFT.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _TrafficCardNFT.Contract.UPGRADEINTERFACEVERSION(&_TrafficCardNFT.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _TrafficCardNFT.Contract.BalanceOf(&_TrafficCardNFT.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _TrafficCardNFT.Contract.BalanceOf(&_TrafficCardNFT.CallOpts, owner)
}

// DepositContract is a free data retrieval call binding the contract method 0xe94ad65b.
//
// Solidity: function depositContract() view returns(address)
func (_TrafficCardNFT *TrafficCardNFTCaller) DepositContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "depositContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// DepositContract is a free data retrieval call binding the contract method 0xe94ad65b.
//
// Solidity: function depositContract() view returns(address)
func (_TrafficCardNFT *TrafficCardNFTSession) DepositContract() (common.Address, error) {
	return _TrafficCardNFT.Contract.DepositContract(&_TrafficCardNFT.CallOpts)
}

// DepositContract is a free data retrieval call binding the contract method 0xe94ad65b.
//
// Solidity: function depositContract() view returns(address)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) DepositContract() (common.Address, error) {
	return _TrafficCardNFT.Contract.DepositContract(&_TrafficCardNFT.CallOpts)
}

// GetActivationCode is a free data retrieval call binding the contract method 0x9e06b06d.
//
// Solidity: function getActivationCode(uint256 tokenId) view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCaller) GetActivationCode(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "getActivationCode", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetActivationCode is a free data retrieval call binding the contract method 0x9e06b06d.
//
// Solidity: function getActivationCode(uint256 tokenId) view returns(string)
func (_TrafficCardNFT *TrafficCardNFTSession) GetActivationCode(tokenId *big.Int) (string, error) {
	return _TrafficCardNFT.Contract.GetActivationCode(&_TrafficCardNFT.CallOpts, tokenId)
}

// GetActivationCode is a free data retrieval call binding the contract method 0x9e06b06d.
//
// Solidity: function getActivationCode(uint256 tokenId) view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) GetActivationCode(tokenId *big.Int) (string, error) {
	return _TrafficCardNFT.Contract.GetActivationCode(&_TrafficCardNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_TrafficCardNFT *TrafficCardNFTCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_TrafficCardNFT *TrafficCardNFTSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _TrafficCardNFT.Contract.GetApproved(&_TrafficCardNFT.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _TrafficCardNFT.Contract.GetApproved(&_TrafficCardNFT.CallOpts, tokenId)
}

// GetCardInfo is a free data retrieval call binding the contract method 0x970129be.
//
// Solidity: function getCardInfo(uint256 tokenId) view returns((uint256,uint256,bool))
func (_TrafficCardNFT *TrafficCardNFTCaller) GetCardInfo(opts *bind.CallOpts, tokenId *big.Int) (ITrafficCardNFTCardInfo, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "getCardInfo", tokenId)

	if err != nil {
		return *new(ITrafficCardNFTCardInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(ITrafficCardNFTCardInfo)).(*ITrafficCardNFTCardInfo)

	return out0, err

}

// GetCardInfo is a free data retrieval call binding the contract method 0x970129be.
//
// Solidity: function getCardInfo(uint256 tokenId) view returns((uint256,uint256,bool))
func (_TrafficCardNFT *TrafficCardNFTSession) GetCardInfo(tokenId *big.Int) (ITrafficCardNFTCardInfo, error) {
	return _TrafficCardNFT.Contract.GetCardInfo(&_TrafficCardNFT.CallOpts, tokenId)
}

// GetCardInfo is a free data retrieval call binding the contract method 0x970129be.
//
// Solidity: function getCardInfo(uint256 tokenId) view returns((uint256,uint256,bool))
func (_TrafficCardNFT *TrafficCardNFTCallerSession) GetCardInfo(tokenId *big.Int) (ITrafficCardNFTCardInfo, error) {
	return _TrafficCardNFT.Contract.GetCardInfo(&_TrafficCardNFT.CallOpts, tokenId)
}

// GetUserCardCount is a free data retrieval call binding the contract method 0xa053f972.
//
// Solidity: function getUserCardCount(address user) view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTCaller) GetUserCardCount(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "getUserCardCount", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetUserCardCount is a free data retrieval call binding the contract method 0xa053f972.
//
// Solidity: function getUserCardCount(address user) view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTSession) GetUserCardCount(user common.Address) (*big.Int, error) {
	return _TrafficCardNFT.Contract.GetUserCardCount(&_TrafficCardNFT.CallOpts, user)
}

// GetUserCardCount is a free data retrieval call binding the contract method 0xa053f972.
//
// Solidity: function getUserCardCount(address user) view returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) GetUserCardCount(user common.Address) (*big.Int, error) {
	return _TrafficCardNFT.Contract.GetUserCardCount(&_TrafficCardNFT.CallOpts, user)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_TrafficCardNFT *TrafficCardNFTCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_TrafficCardNFT *TrafficCardNFTSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _TrafficCardNFT.Contract.IsApprovedForAll(&_TrafficCardNFT.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _TrafficCardNFT.Contract.IsApprovedForAll(&_TrafficCardNFT.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTSession) Name() (string, error) {
	return _TrafficCardNFT.Contract.Name(&_TrafficCardNFT.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) Name() (string, error) {
	return _TrafficCardNFT.Contract.Name(&_TrafficCardNFT.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TrafficCardNFT *TrafficCardNFTCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TrafficCardNFT *TrafficCardNFTSession) Owner() (common.Address, error) {
	return _TrafficCardNFT.Contract.Owner(&_TrafficCardNFT.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) Owner() (common.Address, error) {
	return _TrafficCardNFT.Contract.Owner(&_TrafficCardNFT.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_TrafficCardNFT *TrafficCardNFTCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_TrafficCardNFT *TrafficCardNFTSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _TrafficCardNFT.Contract.OwnerOf(&_TrafficCardNFT.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _TrafficCardNFT.Contract.OwnerOf(&_TrafficCardNFT.CallOpts, tokenId)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TrafficCardNFT *TrafficCardNFTCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TrafficCardNFT *TrafficCardNFTSession) ProxiableUUID() ([32]byte, error) {
	return _TrafficCardNFT.Contract.ProxiableUUID(&_TrafficCardNFT.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) ProxiableUUID() ([32]byte, error) {
	return _TrafficCardNFT.Contract.ProxiableUUID(&_TrafficCardNFT.CallOpts)
}

// SmDpAddress is a free data retrieval call binding the contract method 0x4179cb0a.
//
// Solidity: function smDpAddress() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCaller) SmDpAddress(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "smDpAddress")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// SmDpAddress is a free data retrieval call binding the contract method 0x4179cb0a.
//
// Solidity: function smDpAddress() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTSession) SmDpAddress() (string, error) {
	return _TrafficCardNFT.Contract.SmDpAddress(&_TrafficCardNFT.CallOpts)
}

// SmDpAddress is a free data retrieval call binding the contract method 0x4179cb0a.
//
// Solidity: function smDpAddress() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) SmDpAddress() (string, error) {
	return _TrafficCardNFT.Contract.SmDpAddress(&_TrafficCardNFT.CallOpts)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TrafficCardNFT *TrafficCardNFTCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TrafficCardNFT *TrafficCardNFTSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TrafficCardNFT.Contract.SupportsInterface(&_TrafficCardNFT.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _TrafficCardNFT.Contract.SupportsInterface(&_TrafficCardNFT.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTSession) Symbol() (string, error) {
	return _TrafficCardNFT.Contract.Symbol(&_TrafficCardNFT.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) Symbol() (string, error) {
	return _TrafficCardNFT.Contract.Symbol(&_TrafficCardNFT.CallOpts)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _TrafficCardNFT.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_TrafficCardNFT *TrafficCardNFTSession) TokenURI(tokenId *big.Int) (string, error) {
	return _TrafficCardNFT.Contract.TokenURI(&_TrafficCardNFT.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_TrafficCardNFT *TrafficCardNFTCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _TrafficCardNFT.Contract.TokenURI(&_TrafficCardNFT.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.Approve(&_TrafficCardNFT.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.Approve(&_TrafficCardNFT.TransactOpts, to, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) Burn(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "burn", tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.Burn(&_TrafficCardNFT.TransactOpts, tokenId)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) Burn(tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.Burn(&_TrafficCardNFT.TransactOpts, tokenId)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TrafficCardNFT *TrafficCardNFTSession) Initialize() (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.Initialize(&_TrafficCardNFT.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) Initialize() (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.Initialize(&_TrafficCardNFT.TransactOpts)
}

// Mint is a paid mutator transaction binding the contract method 0xd3fc9864.
//
// Solidity: function mint(address to, uint256 dataAmount, string tokenURI_) returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTTransactor) Mint(opts *bind.TransactOpts, to common.Address, dataAmount *big.Int, tokenURI_ string) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "mint", to, dataAmount, tokenURI_)
}

// Mint is a paid mutator transaction binding the contract method 0xd3fc9864.
//
// Solidity: function mint(address to, uint256 dataAmount, string tokenURI_) returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTSession) Mint(to common.Address, dataAmount *big.Int, tokenURI_ string) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.Mint(&_TrafficCardNFT.TransactOpts, to, dataAmount, tokenURI_)
}

// Mint is a paid mutator transaction binding the contract method 0xd3fc9864.
//
// Solidity: function mint(address to, uint256 dataAmount, string tokenURI_) returns(uint256)
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) Mint(to common.Address, dataAmount *big.Int, tokenURI_ string) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.Mint(&_TrafficCardNFT.TransactOpts, to, dataAmount, tokenURI_)
}

// MintBatch is a paid mutator transaction binding the contract method 0x92315df7.
//
// Solidity: function mintBatch(address[] to, uint256[] dataAmounts, string[] tokenURIs_) returns(uint256[])
func (_TrafficCardNFT *TrafficCardNFTTransactor) MintBatch(opts *bind.TransactOpts, to []common.Address, dataAmounts []*big.Int, tokenURIs_ []string) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "mintBatch", to, dataAmounts, tokenURIs_)
}

// MintBatch is a paid mutator transaction binding the contract method 0x92315df7.
//
// Solidity: function mintBatch(address[] to, uint256[] dataAmounts, string[] tokenURIs_) returns(uint256[])
func (_TrafficCardNFT *TrafficCardNFTSession) MintBatch(to []common.Address, dataAmounts []*big.Int, tokenURIs_ []string) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.MintBatch(&_TrafficCardNFT.TransactOpts, to, dataAmounts, tokenURIs_)
}

// MintBatch is a paid mutator transaction binding the contract method 0x92315df7.
//
// Solidity: function mintBatch(address[] to, uint256[] dataAmounts, string[] tokenURIs_) returns(uint256[])
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) MintBatch(to []common.Address, dataAmounts []*big.Int, tokenURIs_ []string) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.MintBatch(&_TrafficCardNFT.TransactOpts, to, dataAmounts, tokenURIs_)
}

// RedeemForSim is a paid mutator transaction binding the contract method 0xcda54603.
//
// Solidity: function redeemForSim(uint256[] tokenIds) returns(uint256 daysCount)
func (_TrafficCardNFT *TrafficCardNFTTransactor) RedeemForSim(opts *bind.TransactOpts, tokenIds []*big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "redeemForSim", tokenIds)
}

// RedeemForSim is a paid mutator transaction binding the contract method 0xcda54603.
//
// Solidity: function redeemForSim(uint256[] tokenIds) returns(uint256 daysCount)
func (_TrafficCardNFT *TrafficCardNFTSession) RedeemForSim(tokenIds []*big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.RedeemForSim(&_TrafficCardNFT.TransactOpts, tokenIds)
}

// RedeemForSim is a paid mutator transaction binding the contract method 0xcda54603.
//
// Solidity: function redeemForSim(uint256[] tokenIds) returns(uint256 daysCount)
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) RedeemForSim(tokenIds []*big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.RedeemForSim(&_TrafficCardNFT.TransactOpts, tokenIds)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrafficCardNFT *TrafficCardNFTSession) RenounceOwnership() (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.RenounceOwnership(&_TrafficCardNFT.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.RenounceOwnership(&_TrafficCardNFT.TransactOpts)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SafeTransferFrom(&_TrafficCardNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SafeTransferFrom(&_TrafficCardNFT.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SafeTransferFrom0(&_TrafficCardNFT.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SafeTransferFrom0(&_TrafficCardNFT.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SetApprovalForAll(&_TrafficCardNFT.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SetApprovalForAll(&_TrafficCardNFT.TransactOpts, operator, approved)
}

// SetDepositContract is a paid mutator transaction binding the contract method 0x0ec2e821.
//
// Solidity: function setDepositContract(address _deposit) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) SetDepositContract(opts *bind.TransactOpts, _deposit common.Address) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "setDepositContract", _deposit)
}

// SetDepositContract is a paid mutator transaction binding the contract method 0x0ec2e821.
//
// Solidity: function setDepositContract(address _deposit) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) SetDepositContract(_deposit common.Address) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SetDepositContract(&_TrafficCardNFT.TransactOpts, _deposit)
}

// SetDepositContract is a paid mutator transaction binding the contract method 0x0ec2e821.
//
// Solidity: function setDepositContract(address _deposit) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) SetDepositContract(_deposit common.Address) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SetDepositContract(&_TrafficCardNFT.TransactOpts, _deposit)
}

// SetSmDpAddress is a paid mutator transaction binding the contract method 0x29101322.
//
// Solidity: function setSmDpAddress(string _addr) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) SetSmDpAddress(opts *bind.TransactOpts, _addr string) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "setSmDpAddress", _addr)
}

// SetSmDpAddress is a paid mutator transaction binding the contract method 0x29101322.
//
// Solidity: function setSmDpAddress(string _addr) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) SetSmDpAddress(_addr string) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SetSmDpAddress(&_TrafficCardNFT.TransactOpts, _addr)
}

// SetSmDpAddress is a paid mutator transaction binding the contract method 0x29101322.
//
// Solidity: function setSmDpAddress(string _addr) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) SetSmDpAddress(_addr string) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.SetSmDpAddress(&_TrafficCardNFT.TransactOpts, _addr)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.TransferFrom(&_TrafficCardNFT.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.TransferFrom(&_TrafficCardNFT.TransactOpts, from, to, tokenId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TrafficCardNFT *TrafficCardNFTSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.TransferOwnership(&_TrafficCardNFT.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.TransferOwnership(&_TrafficCardNFT.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TrafficCardNFT *TrafficCardNFTTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TrafficCardNFT.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TrafficCardNFT *TrafficCardNFTSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.UpgradeToAndCall(&_TrafficCardNFT.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_TrafficCardNFT *TrafficCardNFTTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _TrafficCardNFT.Contract.UpgradeToAndCall(&_TrafficCardNFT.TransactOpts, newImplementation, data)
}

// TrafficCardNFTApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the TrafficCardNFT contract.
type TrafficCardNFTApprovalIterator struct {
	Event *TrafficCardNFTApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTApproval represents a Approval event raised by the TrafficCardNFT contract.
type TrafficCardNFTApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*TrafficCardNFTApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTApprovalIterator{contract: _TrafficCardNFT.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTApproval)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseApproval(log types.Log) (*TrafficCardNFTApproval, error) {
	event := new(TrafficCardNFTApproval)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the TrafficCardNFT contract.
type TrafficCardNFTApprovalForAllIterator struct {
	Event *TrafficCardNFTApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTApprovalForAll represents a ApprovalForAll event raised by the TrafficCardNFT contract.
type TrafficCardNFTApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*TrafficCardNFTApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTApprovalForAllIterator{contract: _TrafficCardNFT.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTApprovalForAll)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseApprovalForAll(log types.Log) (*TrafficCardNFTApprovalForAll, error) {
	event := new(TrafficCardNFTApprovalForAll)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTBatchMetadataUpdateIterator is returned from FilterBatchMetadataUpdate and is used to iterate over the raw logs and unpacked data for BatchMetadataUpdate events raised by the TrafficCardNFT contract.
type TrafficCardNFTBatchMetadataUpdateIterator struct {
	Event *TrafficCardNFTBatchMetadataUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTBatchMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTBatchMetadataUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTBatchMetadataUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTBatchMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTBatchMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTBatchMetadataUpdate represents a BatchMetadataUpdate event raised by the TrafficCardNFT contract.
type TrafficCardNFTBatchMetadataUpdate struct {
	FromTokenId *big.Int
	ToTokenId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchMetadataUpdate is a free log retrieval operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterBatchMetadataUpdate(opts *bind.FilterOpts) (*TrafficCardNFTBatchMetadataUpdateIterator, error) {

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTBatchMetadataUpdateIterator{contract: _TrafficCardNFT.contract, event: "BatchMetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchBatchMetadataUpdate is a free log subscription operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchBatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTBatchMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTBatchMetadataUpdate)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBatchMetadataUpdate is a log parse operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseBatchMetadataUpdate(log types.Log) (*TrafficCardNFTBatchMetadataUpdate, error) {
	event := new(TrafficCardNFTBatchMetadataUpdate)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTCardDestroyedIterator is returned from FilterCardDestroyed and is used to iterate over the raw logs and unpacked data for CardDestroyed events raised by the TrafficCardNFT contract.
type TrafficCardNFTCardDestroyedIterator struct {
	Event *TrafficCardNFTCardDestroyed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTCardDestroyedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTCardDestroyed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTCardDestroyed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTCardDestroyedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTCardDestroyedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTCardDestroyed represents a CardDestroyed event raised by the TrafficCardNFT contract.
type TrafficCardNFTCardDestroyed struct {
	User       common.Address
	TokenId    *big.Int
	DataAmount *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCardDestroyed is a free log retrieval operation binding the contract event 0xc46168c180dbf8c77c0d472fcbb4bda160e7da889ef45835bd3b4cf13e4757e1.
//
// Solidity: event CardDestroyed(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterCardDestroyed(opts *bind.FilterOpts, user []common.Address) (*TrafficCardNFTCardDestroyedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "CardDestroyed", userRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTCardDestroyedIterator{contract: _TrafficCardNFT.contract, event: "CardDestroyed", logs: logs, sub: sub}, nil
}

// WatchCardDestroyed is a free log subscription operation binding the contract event 0xc46168c180dbf8c77c0d472fcbb4bda160e7da889ef45835bd3b4cf13e4757e1.
//
// Solidity: event CardDestroyed(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchCardDestroyed(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTCardDestroyed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "CardDestroyed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTCardDestroyed)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "CardDestroyed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCardDestroyed is a log parse operation binding the contract event 0xc46168c180dbf8c77c0d472fcbb4bda160e7da889ef45835bd3b4cf13e4757e1.
//
// Solidity: event CardDestroyed(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseCardDestroyed(log types.Log) (*TrafficCardNFTCardDestroyed, error) {
	event := new(TrafficCardNFTCardDestroyed)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "CardDestroyed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTCardMintedIterator is returned from FilterCardMinted and is used to iterate over the raw logs and unpacked data for CardMinted events raised by the TrafficCardNFT contract.
type TrafficCardNFTCardMintedIterator struct {
	Event *TrafficCardNFTCardMinted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTCardMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTCardMinted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTCardMinted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTCardMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTCardMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTCardMinted represents a CardMinted event raised by the TrafficCardNFT contract.
type TrafficCardNFTCardMinted struct {
	User       common.Address
	TokenId    *big.Int
	DataAmount *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterCardMinted is a free log retrieval operation binding the contract event 0xcb1c56d5745b05695241c17b7cfaece9a64f70bc32508c0997990633be8056ff.
//
// Solidity: event CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterCardMinted(opts *bind.FilterOpts, user []common.Address) (*TrafficCardNFTCardMintedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "CardMinted", userRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTCardMintedIterator{contract: _TrafficCardNFT.contract, event: "CardMinted", logs: logs, sub: sub}, nil
}

// WatchCardMinted is a free log subscription operation binding the contract event 0xcb1c56d5745b05695241c17b7cfaece9a64f70bc32508c0997990633be8056ff.
//
// Solidity: event CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchCardMinted(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTCardMinted, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "CardMinted", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTCardMinted)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "CardMinted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseCardMinted is a log parse operation binding the contract event 0xcb1c56d5745b05695241c17b7cfaece9a64f70bc32508c0997990633be8056ff.
//
// Solidity: event CardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseCardMinted(log types.Log) (*TrafficCardNFTCardMinted, error) {
	event := new(TrafficCardNFTCardMinted)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "CardMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTESimRedeemedIterator is returned from FilterESimRedeemed and is used to iterate over the raw logs and unpacked data for ESimRedeemed events raised by the TrafficCardNFT contract.
type TrafficCardNFTESimRedeemedIterator struct {
	Event *TrafficCardNFTESimRedeemed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTESimRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTESimRedeemed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTESimRedeemed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTESimRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTESimRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTESimRedeemed represents a ESimRedeemed event raised by the TrafficCardNFT contract.
type TrafficCardNFTESimRedeemed struct {
	User           common.Address
	TokenId        *big.Int
	ActivationCode string
	SmDpAddress    string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterESimRedeemed is a free log retrieval operation binding the contract event 0x05ee3a1aa9b5412da9d4c0680a87e47a6014024f963aeaa94bc6d831d4d03cf4.
//
// Solidity: event ESimRedeemed(address indexed user, uint256 tokenId, string activationCode, string smDpAddress)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterESimRedeemed(opts *bind.FilterOpts, user []common.Address) (*TrafficCardNFTESimRedeemedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "ESimRedeemed", userRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTESimRedeemedIterator{contract: _TrafficCardNFT.contract, event: "ESimRedeemed", logs: logs, sub: sub}, nil
}

// WatchESimRedeemed is a free log subscription operation binding the contract event 0x05ee3a1aa9b5412da9d4c0680a87e47a6014024f963aeaa94bc6d831d4d03cf4.
//
// Solidity: event ESimRedeemed(address indexed user, uint256 tokenId, string activationCode, string smDpAddress)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchESimRedeemed(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTESimRedeemed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "ESimRedeemed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTESimRedeemed)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "ESimRedeemed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseESimRedeemed is a log parse operation binding the contract event 0x05ee3a1aa9b5412da9d4c0680a87e47a6014024f963aeaa94bc6d831d4d03cf4.
//
// Solidity: event ESimRedeemed(address indexed user, uint256 tokenId, string activationCode, string smDpAddress)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseESimRedeemed(log types.Log) (*TrafficCardNFTESimRedeemed, error) {
	event := new(TrafficCardNFTESimRedeemed)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "ESimRedeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the TrafficCardNFT contract.
type TrafficCardNFTInitializedIterator struct {
	Event *TrafficCardNFTInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTInitialized represents a Initialized event raised by the TrafficCardNFT contract.
type TrafficCardNFTInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterInitialized(opts *bind.FilterOpts) (*TrafficCardNFTInitializedIterator, error) {

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTInitializedIterator{contract: _TrafficCardNFT.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTInitialized) (event.Subscription, error) {

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTInitialized)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseInitialized(log types.Log) (*TrafficCardNFTInitialized, error) {
	event := new(TrafficCardNFTInitialized)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the TrafficCardNFT contract.
type TrafficCardNFTMetadataUpdateIterator struct {
	Event *TrafficCardNFTMetadataUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTMetadataUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTMetadataUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTMetadataUpdate represents a MetadataUpdate event raised by the TrafficCardNFT contract.
type TrafficCardNFTMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*TrafficCardNFTMetadataUpdateIterator, error) {

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTMetadataUpdateIterator{contract: _TrafficCardNFT.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTMetadataUpdate)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseMetadataUpdate(log types.Log) (*TrafficCardNFTMetadataUpdate, error) {
	event := new(TrafficCardNFTMetadataUpdate)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the TrafficCardNFT contract.
type TrafficCardNFTOwnershipTransferredIterator struct {
	Event *TrafficCardNFTOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTOwnershipTransferred represents a OwnershipTransferred event raised by the TrafficCardNFT contract.
type TrafficCardNFTOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*TrafficCardNFTOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTOwnershipTransferredIterator{contract: _TrafficCardNFT.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTOwnershipTransferred)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseOwnershipTransferred(log types.Log) (*TrafficCardNFTOwnershipTransferred, error) {
	event := new(TrafficCardNFTOwnershipTransferred)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTSimRedeemedIterator is returned from FilterSimRedeemed and is used to iterate over the raw logs and unpacked data for SimRedeemed events raised by the TrafficCardNFT contract.
type TrafficCardNFTSimRedeemedIterator struct {
	Event *TrafficCardNFTSimRedeemed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTSimRedeemedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTSimRedeemed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTSimRedeemed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTSimRedeemedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTSimRedeemedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTSimRedeemed represents a SimRedeemed event raised by the TrafficCardNFT contract.
type TrafficCardNFTSimRedeemed struct {
	User      common.Address
	DaysCount *big.Int
	TokenIds  []*big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSimRedeemed is a free log retrieval operation binding the contract event 0x2b22b8bfbca140907fa0a889499309fc6098870b0feb3b4af1993c079007b6b7.
//
// Solidity: event SimRedeemed(address indexed user, uint256 daysCount, uint256[] tokenIds)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterSimRedeemed(opts *bind.FilterOpts, user []common.Address) (*TrafficCardNFTSimRedeemedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "SimRedeemed", userRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTSimRedeemedIterator{contract: _TrafficCardNFT.contract, event: "SimRedeemed", logs: logs, sub: sub}, nil
}

// WatchSimRedeemed is a free log subscription operation binding the contract event 0x2b22b8bfbca140907fa0a889499309fc6098870b0feb3b4af1993c079007b6b7.
//
// Solidity: event SimRedeemed(address indexed user, uint256 daysCount, uint256[] tokenIds)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchSimRedeemed(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTSimRedeemed, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "SimRedeemed", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTSimRedeemed)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "SimRedeemed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSimRedeemed is a log parse operation binding the contract event 0x2b22b8bfbca140907fa0a889499309fc6098870b0feb3b4af1993c079007b6b7.
//
// Solidity: event SimRedeemed(address indexed user, uint256 daysCount, uint256[] tokenIds)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseSimRedeemed(log types.Log) (*TrafficCardNFTSimRedeemed, error) {
	event := new(TrafficCardNFTSimRedeemed)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "SimRedeemed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the TrafficCardNFT contract.
type TrafficCardNFTTransferIterator struct {
	Event *TrafficCardNFTTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTTransfer represents a Transfer event raised by the TrafficCardNFT contract.
type TrafficCardNFTTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*TrafficCardNFTTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTTransferIterator{contract: _TrafficCardNFT.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTTransfer)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseTransfer(log types.Log) (*TrafficCardNFTTransfer, error) {
	event := new(TrafficCardNFTTransfer)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// TrafficCardNFTUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the TrafficCardNFT contract.
type TrafficCardNFTUpgradedIterator struct {
	Event *TrafficCardNFTUpgraded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *TrafficCardNFTUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(TrafficCardNFTUpgraded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(TrafficCardNFTUpgraded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *TrafficCardNFTUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *TrafficCardNFTUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// TrafficCardNFTUpgraded represents a Upgraded event raised by the TrafficCardNFT contract.
type TrafficCardNFTUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TrafficCardNFT *TrafficCardNFTFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*TrafficCardNFTUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &TrafficCardNFTUpgradedIterator{contract: _TrafficCardNFT.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TrafficCardNFT *TrafficCardNFTFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *TrafficCardNFTUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _TrafficCardNFT.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(TrafficCardNFTUpgraded)
				if err := _TrafficCardNFT.contract.UnpackLog(event, "Upgraded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_TrafficCardNFT *TrafficCardNFTFilterer) ParseUpgraded(log types.Log) (*TrafficCardNFTUpgraded, error) {
	event := new(TrafficCardNFTUpgraded)
	if err := _TrafficCardNFT.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
