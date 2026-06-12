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

// IDepositTranche is an auto generated low-level Go binding around an user-defined struct.
type IDepositTranche struct {
	Amount    *big.Int
	UnlockAt  *big.Int
	Withdrawn bool
}

// DepositMetaData contains all meta data concerning the Deposit contract.
var DepositMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"length\",\"type\":\"uint256\"}],\"name\":\"StringsInsufficientHexLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"DepositMade\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"trancheIndex\",\"type\":\"uint256\"}],\"name\":\"DepositWithdrawn\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"dataAmount\",\"type\":\"uint256\"}],\"name\":\"TrafficCardMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"LOCK_PERIOD\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"generateTokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getDepositAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getLockExpiry\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getTrancheCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getTranches\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"unlockAt\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"withdrawn\",\"type\":\"bool\"}],\"internalType\":\"structIDeposit.Tranche[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_userRegistry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_usdt\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_trafficCardNFT\",\"type\":\"address\"}],\"name\":\"setTrafficCardNFT\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trafficCardNFT\",\"outputs\":[{\"internalType\":\"contractITrafficCardNFT\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"usdt\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"userRegistry\",\"outputs\":[{\"internalType\":\"contractIUserRegistry\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"trancheIndex\",\"type\":\"uint256\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// DepositABI is the input ABI used to generate the binding from.
// Deprecated: Use DepositMetaData.ABI instead.
var DepositABI = DepositMetaData.ABI

// Deposit is an auto generated Go binding around an Ethereum contract.
type Deposit struct {
	DepositCaller     // Read-only binding to the contract
	DepositTransactor // Write-only binding to the contract
	DepositFilterer   // Log filterer for contract events
}

// DepositCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositSession struct {
	Contract     *Deposit          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositCallerSession struct {
	Contract *DepositCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// DepositTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositTransactorSession struct {
	Contract     *DepositTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// DepositRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositRaw struct {
	Contract *Deposit // Generic contract binding to access the raw methods on
}

// DepositCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositCallerRaw struct {
	Contract *DepositCaller // Generic read-only contract binding to access the raw methods on
}

// DepositTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositTransactorRaw struct {
	Contract *DepositTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDeposit creates a new instance of Deposit, bound to a specific deployed contract.
func NewDeposit(address common.Address, backend bind.ContractBackend) (*Deposit, error) {
	contract, err := bindDeposit(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Deposit{DepositCaller: DepositCaller{contract: contract}, DepositTransactor: DepositTransactor{contract: contract}, DepositFilterer: DepositFilterer{contract: contract}}, nil
}

// NewDepositCaller creates a new read-only instance of Deposit, bound to a specific deployed contract.
func NewDepositCaller(address common.Address, caller bind.ContractCaller) (*DepositCaller, error) {
	contract, err := bindDeposit(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositCaller{contract: contract}, nil
}

// NewDepositTransactor creates a new write-only instance of Deposit, bound to a specific deployed contract.
func NewDepositTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositTransactor, error) {
	contract, err := bindDeposit(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositTransactor{contract: contract}, nil
}

// NewDepositFilterer creates a new log filterer instance of Deposit, bound to a specific deployed contract.
func NewDepositFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositFilterer, error) {
	contract, err := bindDeposit(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositFilterer{contract: contract}, nil
}

// bindDeposit binds a generic wrapper to an already deployed contract.
func bindDeposit(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DepositMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.DepositCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.DepositTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Deposit *DepositCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Deposit.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Deposit *DepositTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Deposit *DepositTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Deposit.Contract.contract.Transact(opts, method, params...)
}

// LOCKPERIOD is a free data retrieval call binding the contract method 0x1820cabb.
//
// Solidity: function LOCK_PERIOD() view returns(uint256)
func (_Deposit *DepositCaller) LOCKPERIOD(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "LOCK_PERIOD")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LOCKPERIOD is a free data retrieval call binding the contract method 0x1820cabb.
//
// Solidity: function LOCK_PERIOD() view returns(uint256)
func (_Deposit *DepositSession) LOCKPERIOD() (*big.Int, error) {
	return _Deposit.Contract.LOCKPERIOD(&_Deposit.CallOpts)
}

// LOCKPERIOD is a free data retrieval call binding the contract method 0x1820cabb.
//
// Solidity: function LOCK_PERIOD() view returns(uint256)
func (_Deposit *DepositCallerSession) LOCKPERIOD() (*big.Int, error) {
	return _Deposit.Contract.LOCKPERIOD(&_Deposit.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Deposit *DepositCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Deposit *DepositSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Deposit.Contract.UPGRADEINTERFACEVERSION(&_Deposit.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Deposit *DepositCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Deposit.Contract.UPGRADEINTERFACEVERSION(&_Deposit.CallOpts)
}

// GenerateTokenURI is a free data retrieval call binding the contract method 0x13613d4d.
//
// Solidity: function generateTokenURI(address user, uint256 index) view returns(string)
func (_Deposit *DepositCaller) GenerateTokenURI(opts *bind.CallOpts, user common.Address, index *big.Int) (string, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "generateTokenURI", user, index)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GenerateTokenURI is a free data retrieval call binding the contract method 0x13613d4d.
//
// Solidity: function generateTokenURI(address user, uint256 index) view returns(string)
func (_Deposit *DepositSession) GenerateTokenURI(user common.Address, index *big.Int) (string, error) {
	return _Deposit.Contract.GenerateTokenURI(&_Deposit.CallOpts, user, index)
}

// GenerateTokenURI is a free data retrieval call binding the contract method 0x13613d4d.
//
// Solidity: function generateTokenURI(address user, uint256 index) view returns(string)
func (_Deposit *DepositCallerSession) GenerateTokenURI(user common.Address, index *big.Int) (string, error) {
	return _Deposit.Contract.GenerateTokenURI(&_Deposit.CallOpts, user, index)
}

// GetDepositAmount is a free data retrieval call binding the contract method 0xb8ba16fd.
//
// Solidity: function getDepositAmount(address user) view returns(uint256)
func (_Deposit *DepositCaller) GetDepositAmount(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "getDepositAmount", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDepositAmount is a free data retrieval call binding the contract method 0xb8ba16fd.
//
// Solidity: function getDepositAmount(address user) view returns(uint256)
func (_Deposit *DepositSession) GetDepositAmount(user common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetDepositAmount(&_Deposit.CallOpts, user)
}

// GetDepositAmount is a free data retrieval call binding the contract method 0xb8ba16fd.
//
// Solidity: function getDepositAmount(address user) view returns(uint256)
func (_Deposit *DepositCallerSession) GetDepositAmount(user common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetDepositAmount(&_Deposit.CallOpts, user)
}

// GetLockExpiry is a free data retrieval call binding the contract method 0xbdab6bf9.
//
// Solidity: function getLockExpiry(address user) view returns(uint256)
func (_Deposit *DepositCaller) GetLockExpiry(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "getLockExpiry", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLockExpiry is a free data retrieval call binding the contract method 0xbdab6bf9.
//
// Solidity: function getLockExpiry(address user) view returns(uint256)
func (_Deposit *DepositSession) GetLockExpiry(user common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetLockExpiry(&_Deposit.CallOpts, user)
}

// GetLockExpiry is a free data retrieval call binding the contract method 0xbdab6bf9.
//
// Solidity: function getLockExpiry(address user) view returns(uint256)
func (_Deposit *DepositCallerSession) GetLockExpiry(user common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetLockExpiry(&_Deposit.CallOpts, user)
}

// GetTrancheCount is a free data retrieval call binding the contract method 0x49cfece1.
//
// Solidity: function getTrancheCount(address user) view returns(uint256)
func (_Deposit *DepositCaller) GetTrancheCount(opts *bind.CallOpts, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "getTrancheCount", user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTrancheCount is a free data retrieval call binding the contract method 0x49cfece1.
//
// Solidity: function getTrancheCount(address user) view returns(uint256)
func (_Deposit *DepositSession) GetTrancheCount(user common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetTrancheCount(&_Deposit.CallOpts, user)
}

// GetTrancheCount is a free data retrieval call binding the contract method 0x49cfece1.
//
// Solidity: function getTrancheCount(address user) view returns(uint256)
func (_Deposit *DepositCallerSession) GetTrancheCount(user common.Address) (*big.Int, error) {
	return _Deposit.Contract.GetTrancheCount(&_Deposit.CallOpts, user)
}

// GetTranches is a free data retrieval call binding the contract method 0x39732bcc.
//
// Solidity: function getTranches(address user) view returns((uint256,uint256,bool)[])
func (_Deposit *DepositCaller) GetTranches(opts *bind.CallOpts, user common.Address) ([]IDepositTranche, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "getTranches", user)

	if err != nil {
		return *new([]IDepositTranche), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDepositTranche)).(*[]IDepositTranche)

	return out0, err

}

// GetTranches is a free data retrieval call binding the contract method 0x39732bcc.
//
// Solidity: function getTranches(address user) view returns((uint256,uint256,bool)[])
func (_Deposit *DepositSession) GetTranches(user common.Address) ([]IDepositTranche, error) {
	return _Deposit.Contract.GetTranches(&_Deposit.CallOpts, user)
}

// GetTranches is a free data retrieval call binding the contract method 0x39732bcc.
//
// Solidity: function getTranches(address user) view returns((uint256,uint256,bool)[])
func (_Deposit *DepositCallerSession) GetTranches(user common.Address) ([]IDepositTranche, error) {
	return _Deposit.Contract.GetTranches(&_Deposit.CallOpts, user)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Deposit *DepositCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Deposit *DepositSession) Oracle() (common.Address, error) {
	return _Deposit.Contract.Oracle(&_Deposit.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Deposit *DepositCallerSession) Oracle() (common.Address, error) {
	return _Deposit.Contract.Oracle(&_Deposit.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositSession) Owner() (common.Address, error) {
	return _Deposit.Contract.Owner(&_Deposit.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Deposit *DepositCallerSession) Owner() (common.Address, error) {
	return _Deposit.Contract.Owner(&_Deposit.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Deposit *DepositCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Deposit *DepositSession) ProxiableUUID() ([32]byte, error) {
	return _Deposit.Contract.ProxiableUUID(&_Deposit.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Deposit *DepositCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Deposit.Contract.ProxiableUUID(&_Deposit.CallOpts)
}

// TrafficCardNFT is a free data retrieval call binding the contract method 0xe6ae2af1.
//
// Solidity: function trafficCardNFT() view returns(address)
func (_Deposit *DepositCaller) TrafficCardNFT(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "trafficCardNFT")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TrafficCardNFT is a free data retrieval call binding the contract method 0xe6ae2af1.
//
// Solidity: function trafficCardNFT() view returns(address)
func (_Deposit *DepositSession) TrafficCardNFT() (common.Address, error) {
	return _Deposit.Contract.TrafficCardNFT(&_Deposit.CallOpts)
}

// TrafficCardNFT is a free data retrieval call binding the contract method 0xe6ae2af1.
//
// Solidity: function trafficCardNFT() view returns(address)
func (_Deposit *DepositCallerSession) TrafficCardNFT() (common.Address, error) {
	return _Deposit.Contract.TrafficCardNFT(&_Deposit.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_Deposit *DepositCaller) Usdt(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "usdt")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_Deposit *DepositSession) Usdt() (common.Address, error) {
	return _Deposit.Contract.Usdt(&_Deposit.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_Deposit *DepositCallerSession) Usdt() (common.Address, error) {
	return _Deposit.Contract.Usdt(&_Deposit.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Deposit *DepositCaller) UserRegistry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Deposit.contract.Call(opts, &out, "userRegistry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Deposit *DepositSession) UserRegistry() (common.Address, error) {
	return _Deposit.Contract.UserRegistry(&_Deposit.CallOpts)
}

// UserRegistry is a free data retrieval call binding the contract method 0x5c7460d6.
//
// Solidity: function userRegistry() view returns(address)
func (_Deposit *DepositCallerSession) UserRegistry() (common.Address, error) {
	return _Deposit.Contract.UserRegistry(&_Deposit.CallOpts)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Deposit *DepositTransactor) Deposit(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "deposit", amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Deposit *DepositSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.Deposit(&_Deposit.TransactOpts, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0xb6b55f25.
//
// Solidity: function deposit(uint256 amount) returns()
func (_Deposit *DepositTransactorSession) Deposit(amount *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.Deposit(&_Deposit.TransactOpts, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _userRegistry, address _usdt) returns()
func (_Deposit *DepositTransactor) Initialize(opts *bind.TransactOpts, _userRegistry common.Address, _usdt common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "initialize", _userRegistry, _usdt)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _userRegistry, address _usdt) returns()
func (_Deposit *DepositSession) Initialize(_userRegistry common.Address, _usdt common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.Initialize(&_Deposit.TransactOpts, _userRegistry, _usdt)
}

// Initialize is a paid mutator transaction binding the contract method 0x485cc955.
//
// Solidity: function initialize(address _userRegistry, address _usdt) returns()
func (_Deposit *DepositTransactorSession) Initialize(_userRegistry common.Address, _usdt common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.Initialize(&_Deposit.TransactOpts, _userRegistry, _usdt)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deposit *DepositTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deposit *DepositSession) RenounceOwnership() (*types.Transaction, error) {
	return _Deposit.Contract.RenounceOwnership(&_Deposit.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Deposit *DepositTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Deposit.Contract.RenounceOwnership(&_Deposit.TransactOpts)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_Deposit *DepositTransactor) SetOracle(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "setOracle", _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_Deposit *DepositSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.SetOracle(&_Deposit.TransactOpts, _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_Deposit *DepositTransactorSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.SetOracle(&_Deposit.TransactOpts, _oracle)
}

// SetTrafficCardNFT is a paid mutator transaction binding the contract method 0xa41e3436.
//
// Solidity: function setTrafficCardNFT(address _trafficCardNFT) returns()
func (_Deposit *DepositTransactor) SetTrafficCardNFT(opts *bind.TransactOpts, _trafficCardNFT common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "setTrafficCardNFT", _trafficCardNFT)
}

// SetTrafficCardNFT is a paid mutator transaction binding the contract method 0xa41e3436.
//
// Solidity: function setTrafficCardNFT(address _trafficCardNFT) returns()
func (_Deposit *DepositSession) SetTrafficCardNFT(_trafficCardNFT common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.SetTrafficCardNFT(&_Deposit.TransactOpts, _trafficCardNFT)
}

// SetTrafficCardNFT is a paid mutator transaction binding the contract method 0xa41e3436.
//
// Solidity: function setTrafficCardNFT(address _trafficCardNFT) returns()
func (_Deposit *DepositTransactorSession) SetTrafficCardNFT(_trafficCardNFT common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.SetTrafficCardNFT(&_Deposit.TransactOpts, _trafficCardNFT)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deposit *DepositTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deposit *DepositSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.TransferOwnership(&_Deposit.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Deposit *DepositTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Deposit.Contract.TransferOwnership(&_Deposit.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Deposit *DepositTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Deposit *DepositSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Deposit.Contract.UpgradeToAndCall(&_Deposit.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Deposit *DepositTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Deposit.Contract.UpgradeToAndCall(&_Deposit.TransactOpts, newImplementation, data)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 trancheIndex) returns()
func (_Deposit *DepositTransactor) Withdraw(opts *bind.TransactOpts, trancheIndex *big.Int) (*types.Transaction, error) {
	return _Deposit.contract.Transact(opts, "withdraw", trancheIndex)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 trancheIndex) returns()
func (_Deposit *DepositSession) Withdraw(trancheIndex *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.Withdraw(&_Deposit.TransactOpts, trancheIndex)
}

// Withdraw is a paid mutator transaction binding the contract method 0x2e1a7d4d.
//
// Solidity: function withdraw(uint256 trancheIndex) returns()
func (_Deposit *DepositTransactorSession) Withdraw(trancheIndex *big.Int) (*types.Transaction, error) {
	return _Deposit.Contract.Withdraw(&_Deposit.TransactOpts, trancheIndex)
}

// DepositDepositMadeIterator is returned from FilterDepositMade and is used to iterate over the raw logs and unpacked data for DepositMade events raised by the Deposit contract.
type DepositDepositMadeIterator struct {
	Event *DepositDepositMade // Event containing the contract specifics and raw log

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
func (it *DepositDepositMadeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositDepositMade)
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
		it.Event = new(DepositDepositMade)
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
func (it *DepositDepositMadeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositDepositMadeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositDepositMade represents a DepositMade event raised by the Deposit contract.
type DepositDepositMade struct {
	User   common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDepositMade is a free log retrieval operation binding the contract event 0xd15c9547ea5c06670c0010ce19bc32d54682a4b3801ece7f3ab0c3f17106b4bb.
//
// Solidity: event DepositMade(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) FilterDepositMade(opts *bind.FilterOpts, user []common.Address) (*DepositDepositMadeIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "DepositMade", userRule)
	if err != nil {
		return nil, err
	}
	return &DepositDepositMadeIterator{contract: _Deposit.contract, event: "DepositMade", logs: logs, sub: sub}, nil
}

// WatchDepositMade is a free log subscription operation binding the contract event 0xd15c9547ea5c06670c0010ce19bc32d54682a4b3801ece7f3ab0c3f17106b4bb.
//
// Solidity: event DepositMade(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) WatchDepositMade(opts *bind.WatchOpts, sink chan<- *DepositDepositMade, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "DepositMade", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositDepositMade)
				if err := _Deposit.contract.UnpackLog(event, "DepositMade", log); err != nil {
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

// ParseDepositMade is a log parse operation binding the contract event 0xd15c9547ea5c06670c0010ce19bc32d54682a4b3801ece7f3ab0c3f17106b4bb.
//
// Solidity: event DepositMade(address indexed user, uint256 amount)
func (_Deposit *DepositFilterer) ParseDepositMade(log types.Log) (*DepositDepositMade, error) {
	event := new(DepositDepositMade)
	if err := _Deposit.contract.UnpackLog(event, "DepositMade", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositDepositWithdrawnIterator is returned from FilterDepositWithdrawn and is used to iterate over the raw logs and unpacked data for DepositWithdrawn events raised by the Deposit contract.
type DepositDepositWithdrawnIterator struct {
	Event *DepositDepositWithdrawn // Event containing the contract specifics and raw log

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
func (it *DepositDepositWithdrawnIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositDepositWithdrawn)
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
		it.Event = new(DepositDepositWithdrawn)
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
func (it *DepositDepositWithdrawnIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositDepositWithdrawnIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositDepositWithdrawn represents a DepositWithdrawn event raised by the Deposit contract.
type DepositDepositWithdrawn struct {
	User         common.Address
	Amount       *big.Int
	TrancheIndex *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDepositWithdrawn is a free log retrieval operation binding the contract event 0x7719804546c0185709e60c90d164447ff251a5ba29af0216faa921350f6bebf7.
//
// Solidity: event DepositWithdrawn(address indexed user, uint256 amount, uint256 trancheIndex)
func (_Deposit *DepositFilterer) FilterDepositWithdrawn(opts *bind.FilterOpts, user []common.Address) (*DepositDepositWithdrawnIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "DepositWithdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return &DepositDepositWithdrawnIterator{contract: _Deposit.contract, event: "DepositWithdrawn", logs: logs, sub: sub}, nil
}

// WatchDepositWithdrawn is a free log subscription operation binding the contract event 0x7719804546c0185709e60c90d164447ff251a5ba29af0216faa921350f6bebf7.
//
// Solidity: event DepositWithdrawn(address indexed user, uint256 amount, uint256 trancheIndex)
func (_Deposit *DepositFilterer) WatchDepositWithdrawn(opts *bind.WatchOpts, sink chan<- *DepositDepositWithdrawn, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "DepositWithdrawn", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositDepositWithdrawn)
				if err := _Deposit.contract.UnpackLog(event, "DepositWithdrawn", log); err != nil {
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

// ParseDepositWithdrawn is a log parse operation binding the contract event 0x7719804546c0185709e60c90d164447ff251a5ba29af0216faa921350f6bebf7.
//
// Solidity: event DepositWithdrawn(address indexed user, uint256 amount, uint256 trancheIndex)
func (_Deposit *DepositFilterer) ParseDepositWithdrawn(log types.Log) (*DepositDepositWithdrawn, error) {
	event := new(DepositDepositWithdrawn)
	if err := _Deposit.contract.UnpackLog(event, "DepositWithdrawn", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Deposit contract.
type DepositInitializedIterator struct {
	Event *DepositInitialized // Event containing the contract specifics and raw log

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
func (it *DepositInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositInitialized)
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
		it.Event = new(DepositInitialized)
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
func (it *DepositInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositInitialized represents a Initialized event raised by the Deposit contract.
type DepositInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Deposit *DepositFilterer) FilterInitialized(opts *bind.FilterOpts) (*DepositInitializedIterator, error) {

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DepositInitializedIterator{contract: _Deposit.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Deposit *DepositFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DepositInitialized) (event.Subscription, error) {

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositInitialized)
				if err := _Deposit.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Deposit *DepositFilterer) ParseInitialized(log types.Log) (*DepositInitialized, error) {
	event := new(DepositInitialized)
	if err := _Deposit.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Deposit contract.
type DepositOwnershipTransferredIterator struct {
	Event *DepositOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DepositOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositOwnershipTransferred)
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
		it.Event = new(DepositOwnershipTransferred)
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
func (it *DepositOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositOwnershipTransferred represents a OwnershipTransferred event raised by the Deposit contract.
type DepositOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Deposit *DepositFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DepositOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DepositOwnershipTransferredIterator{contract: _Deposit.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Deposit *DepositFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DepositOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositOwnershipTransferred)
				if err := _Deposit.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Deposit *DepositFilterer) ParseOwnershipTransferred(log types.Log) (*DepositOwnershipTransferred, error) {
	event := new(DepositOwnershipTransferred)
	if err := _Deposit.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositTrafficCardMintedIterator is returned from FilterTrafficCardMinted and is used to iterate over the raw logs and unpacked data for TrafficCardMinted events raised by the Deposit contract.
type DepositTrafficCardMintedIterator struct {
	Event *DepositTrafficCardMinted // Event containing the contract specifics and raw log

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
func (it *DepositTrafficCardMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositTrafficCardMinted)
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
		it.Event = new(DepositTrafficCardMinted)
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
func (it *DepositTrafficCardMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositTrafficCardMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositTrafficCardMinted represents a TrafficCardMinted event raised by the Deposit contract.
type DepositTrafficCardMinted struct {
	User       common.Address
	TokenId    *big.Int
	DataAmount *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTrafficCardMinted is a free log retrieval operation binding the contract event 0x3072026fc6418657755f94a3ef0972a28bb4266b72de75d96c3e999d4ff7067d.
//
// Solidity: event TrafficCardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_Deposit *DepositFilterer) FilterTrafficCardMinted(opts *bind.FilterOpts, user []common.Address) (*DepositTrafficCardMintedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "TrafficCardMinted", userRule)
	if err != nil {
		return nil, err
	}
	return &DepositTrafficCardMintedIterator{contract: _Deposit.contract, event: "TrafficCardMinted", logs: logs, sub: sub}, nil
}

// WatchTrafficCardMinted is a free log subscription operation binding the contract event 0x3072026fc6418657755f94a3ef0972a28bb4266b72de75d96c3e999d4ff7067d.
//
// Solidity: event TrafficCardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_Deposit *DepositFilterer) WatchTrafficCardMinted(opts *bind.WatchOpts, sink chan<- *DepositTrafficCardMinted, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "TrafficCardMinted", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositTrafficCardMinted)
				if err := _Deposit.contract.UnpackLog(event, "TrafficCardMinted", log); err != nil {
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

// ParseTrafficCardMinted is a log parse operation binding the contract event 0x3072026fc6418657755f94a3ef0972a28bb4266b72de75d96c3e999d4ff7067d.
//
// Solidity: event TrafficCardMinted(address indexed user, uint256 tokenId, uint256 dataAmount)
func (_Deposit *DepositFilterer) ParseTrafficCardMinted(log types.Log) (*DepositTrafficCardMinted, error) {
	event := new(DepositTrafficCardMinted)
	if err := _Deposit.contract.UnpackLog(event, "TrafficCardMinted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DepositUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Deposit contract.
type DepositUpgradedIterator struct {
	Event *DepositUpgraded // Event containing the contract specifics and raw log

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
func (it *DepositUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositUpgraded)
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
		it.Event = new(DepositUpgraded)
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
func (it *DepositUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositUpgraded represents a Upgraded event raised by the Deposit contract.
type DepositUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Deposit *DepositFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DepositUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Deposit.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DepositUpgradedIterator{contract: _Deposit.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Deposit *DepositFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DepositUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Deposit.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositUpgraded)
				if err := _Deposit.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Deposit *DepositFilterer) ParseUpgraded(log types.Log) (*DepositUpgraded, error) {
	event := new(DepositUpgraded)
	if err := _Deposit.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
