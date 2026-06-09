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

// FeeManagerMetaData contains all meta data concerning the FeeManager contract.
var FeeManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldRate\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newRate\",\"type\":\"uint256\"}],\"name\":\"FeeRateUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"FEE_DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_FEE_RATE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"calculateFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getFeeRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialFeeRate\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newRate\",\"type\":\"uint256\"}],\"name\":\"setFeeRate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// FeeManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeManagerMetaData.ABI instead.
var FeeManagerABI = FeeManagerMetaData.ABI

// FeeManager is an auto generated Go binding around an Ethereum contract.
type FeeManager struct {
	FeeManagerCaller     // Read-only binding to the contract
	FeeManagerTransactor // Write-only binding to the contract
	FeeManagerFilterer   // Log filterer for contract events
}

// FeeManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeManagerSession struct {
	Contract     *FeeManager       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeeManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeManagerCallerSession struct {
	Contract *FeeManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// FeeManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeManagerTransactorSession struct {
	Contract     *FeeManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// FeeManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeManagerRaw struct {
	Contract *FeeManager // Generic contract binding to access the raw methods on
}

// FeeManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeManagerCallerRaw struct {
	Contract *FeeManagerCaller // Generic read-only contract binding to access the raw methods on
}

// FeeManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeManagerTransactorRaw struct {
	Contract *FeeManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeManager creates a new instance of FeeManager, bound to a specific deployed contract.
func NewFeeManager(address common.Address, backend bind.ContractBackend) (*FeeManager, error) {
	contract, err := bindFeeManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeManager{FeeManagerCaller: FeeManagerCaller{contract: contract}, FeeManagerTransactor: FeeManagerTransactor{contract: contract}, FeeManagerFilterer: FeeManagerFilterer{contract: contract}}, nil
}

// NewFeeManagerCaller creates a new read-only instance of FeeManager, bound to a specific deployed contract.
func NewFeeManagerCaller(address common.Address, caller bind.ContractCaller) (*FeeManagerCaller, error) {
	contract, err := bindFeeManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeManagerCaller{contract: contract}, nil
}

// NewFeeManagerTransactor creates a new write-only instance of FeeManager, bound to a specific deployed contract.
func NewFeeManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeManagerTransactor, error) {
	contract, err := bindFeeManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeManagerTransactor{contract: contract}, nil
}

// NewFeeManagerFilterer creates a new log filterer instance of FeeManager, bound to a specific deployed contract.
func NewFeeManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeManagerFilterer, error) {
	contract, err := bindFeeManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeManagerFilterer{contract: contract}, nil
}

// bindFeeManager binds a generic wrapper to an already deployed contract.
func bindFeeManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeManager *FeeManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeManager.Contract.FeeManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeManager *FeeManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeManager.Contract.FeeManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeManager *FeeManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeManager.Contract.FeeManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeManager *FeeManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeManager *FeeManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeManager *FeeManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeManager.Contract.contract.Transact(opts, method, params...)
}

// FEEDENOMINATOR is a free data retrieval call binding the contract method 0xd73792a9.
//
// Solidity: function FEE_DENOMINATOR() view returns(uint256)
func (_FeeManager *FeeManagerCaller) FEEDENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "FEE_DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FEEDENOMINATOR is a free data retrieval call binding the contract method 0xd73792a9.
//
// Solidity: function FEE_DENOMINATOR() view returns(uint256)
func (_FeeManager *FeeManagerSession) FEEDENOMINATOR() (*big.Int, error) {
	return _FeeManager.Contract.FEEDENOMINATOR(&_FeeManager.CallOpts)
}

// FEEDENOMINATOR is a free data retrieval call binding the contract method 0xd73792a9.
//
// Solidity: function FEE_DENOMINATOR() view returns(uint256)
func (_FeeManager *FeeManagerCallerSession) FEEDENOMINATOR() (*big.Int, error) {
	return _FeeManager.Contract.FEEDENOMINATOR(&_FeeManager.CallOpts)
}

// MAXFEERATE is a free data retrieval call binding the contract method 0x92f6576e.
//
// Solidity: function MAX_FEE_RATE() view returns(uint256)
func (_FeeManager *FeeManagerCaller) MAXFEERATE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "MAX_FEE_RATE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXFEERATE is a free data retrieval call binding the contract method 0x92f6576e.
//
// Solidity: function MAX_FEE_RATE() view returns(uint256)
func (_FeeManager *FeeManagerSession) MAXFEERATE() (*big.Int, error) {
	return _FeeManager.Contract.MAXFEERATE(&_FeeManager.CallOpts)
}

// MAXFEERATE is a free data retrieval call binding the contract method 0x92f6576e.
//
// Solidity: function MAX_FEE_RATE() view returns(uint256)
func (_FeeManager *FeeManagerCallerSession) MAXFEERATE() (*big.Int, error) {
	return _FeeManager.Contract.MAXFEERATE(&_FeeManager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FeeManager *FeeManagerCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FeeManager *FeeManagerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _FeeManager.Contract.UPGRADEINTERFACEVERSION(&_FeeManager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_FeeManager *FeeManagerCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _FeeManager.Contract.UPGRADEINTERFACEVERSION(&_FeeManager.CallOpts)
}

// CalculateFee is a free data retrieval call binding the contract method 0x99a5d747.
//
// Solidity: function calculateFee(uint256 amount) view returns(uint256)
func (_FeeManager *FeeManagerCaller) CalculateFee(opts *bind.CallOpts, amount *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "calculateFee", amount)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateFee is a free data retrieval call binding the contract method 0x99a5d747.
//
// Solidity: function calculateFee(uint256 amount) view returns(uint256)
func (_FeeManager *FeeManagerSession) CalculateFee(amount *big.Int) (*big.Int, error) {
	return _FeeManager.Contract.CalculateFee(&_FeeManager.CallOpts, amount)
}

// CalculateFee is a free data retrieval call binding the contract method 0x99a5d747.
//
// Solidity: function calculateFee(uint256 amount) view returns(uint256)
func (_FeeManager *FeeManagerCallerSession) CalculateFee(amount *big.Int) (*big.Int, error) {
	return _FeeManager.Contract.CalculateFee(&_FeeManager.CallOpts, amount)
}

// GetFeeRate is a free data retrieval call binding the contract method 0x84e5eed0.
//
// Solidity: function getFeeRate() view returns(uint256)
func (_FeeManager *FeeManagerCaller) GetFeeRate(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "getFeeRate")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetFeeRate is a free data retrieval call binding the contract method 0x84e5eed0.
//
// Solidity: function getFeeRate() view returns(uint256)
func (_FeeManager *FeeManagerSession) GetFeeRate() (*big.Int, error) {
	return _FeeManager.Contract.GetFeeRate(&_FeeManager.CallOpts)
}

// GetFeeRate is a free data retrieval call binding the contract method 0x84e5eed0.
//
// Solidity: function getFeeRate() view returns(uint256)
func (_FeeManager *FeeManagerCallerSession) GetFeeRate() (*big.Int, error) {
	return _FeeManager.Contract.GetFeeRate(&_FeeManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeManager *FeeManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeManager *FeeManagerSession) Owner() (common.Address, error) {
	return _FeeManager.Contract.Owner(&_FeeManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeManager *FeeManagerCallerSession) Owner() (common.Address, error) {
	return _FeeManager.Contract.Owner(&_FeeManager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FeeManager *FeeManagerCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _FeeManager.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FeeManager *FeeManagerSession) ProxiableUUID() ([32]byte, error) {
	return _FeeManager.Contract.ProxiableUUID(&_FeeManager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_FeeManager *FeeManagerCallerSession) ProxiableUUID() ([32]byte, error) {
	return _FeeManager.Contract.ProxiableUUID(&_FeeManager.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 initialFeeRate) returns()
func (_FeeManager *FeeManagerTransactor) Initialize(opts *bind.TransactOpts, initialFeeRate *big.Int) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "initialize", initialFeeRate)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 initialFeeRate) returns()
func (_FeeManager *FeeManagerSession) Initialize(initialFeeRate *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.Initialize(&_FeeManager.TransactOpts, initialFeeRate)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 initialFeeRate) returns()
func (_FeeManager *FeeManagerTransactorSession) Initialize(initialFeeRate *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.Initialize(&_FeeManager.TransactOpts, initialFeeRate)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeManager *FeeManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeManager *FeeManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _FeeManager.Contract.RenounceOwnership(&_FeeManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeManager *FeeManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FeeManager.Contract.RenounceOwnership(&_FeeManager.TransactOpts)
}

// SetFeeRate is a paid mutator transaction binding the contract method 0x45596e2e.
//
// Solidity: function setFeeRate(uint256 newRate) returns()
func (_FeeManager *FeeManagerTransactor) SetFeeRate(opts *bind.TransactOpts, newRate *big.Int) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "setFeeRate", newRate)
}

// SetFeeRate is a paid mutator transaction binding the contract method 0x45596e2e.
//
// Solidity: function setFeeRate(uint256 newRate) returns()
func (_FeeManager *FeeManagerSession) SetFeeRate(newRate *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.SetFeeRate(&_FeeManager.TransactOpts, newRate)
}

// SetFeeRate is a paid mutator transaction binding the contract method 0x45596e2e.
//
// Solidity: function setFeeRate(uint256 newRate) returns()
func (_FeeManager *FeeManagerTransactorSession) SetFeeRate(newRate *big.Int) (*types.Transaction, error) {
	return _FeeManager.Contract.SetFeeRate(&_FeeManager.TransactOpts, newRate)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeManager *FeeManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeManager *FeeManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FeeManager.Contract.TransferOwnership(&_FeeManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeManager *FeeManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FeeManager.Contract.TransferOwnership(&_FeeManager.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FeeManager *FeeManagerTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FeeManager.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FeeManager *FeeManagerSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FeeManager.Contract.UpgradeToAndCall(&_FeeManager.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_FeeManager *FeeManagerTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _FeeManager.Contract.UpgradeToAndCall(&_FeeManager.TransactOpts, newImplementation, data)
}

// FeeManagerFeeRateUpdatedIterator is returned from FilterFeeRateUpdated and is used to iterate over the raw logs and unpacked data for FeeRateUpdated events raised by the FeeManager contract.
type FeeManagerFeeRateUpdatedIterator struct {
	Event *FeeManagerFeeRateUpdated // Event containing the contract specifics and raw log

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
func (it *FeeManagerFeeRateUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerFeeRateUpdated)
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
		it.Event = new(FeeManagerFeeRateUpdated)
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
func (it *FeeManagerFeeRateUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeManagerFeeRateUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeManagerFeeRateUpdated represents a FeeRateUpdated event raised by the FeeManager contract.
type FeeManagerFeeRateUpdated struct {
	OldRate *big.Int
	NewRate *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterFeeRateUpdated is a free log retrieval operation binding the contract event 0x14914da2bf76024616fbe1859783fcd4dbddcb179b1f3a854949fbf920dcb957.
//
// Solidity: event FeeRateUpdated(uint256 oldRate, uint256 newRate)
func (_FeeManager *FeeManagerFilterer) FilterFeeRateUpdated(opts *bind.FilterOpts) (*FeeManagerFeeRateUpdatedIterator, error) {

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "FeeRateUpdated")
	if err != nil {
		return nil, err
	}
	return &FeeManagerFeeRateUpdatedIterator{contract: _FeeManager.contract, event: "FeeRateUpdated", logs: logs, sub: sub}, nil
}

// WatchFeeRateUpdated is a free log subscription operation binding the contract event 0x14914da2bf76024616fbe1859783fcd4dbddcb179b1f3a854949fbf920dcb957.
//
// Solidity: event FeeRateUpdated(uint256 oldRate, uint256 newRate)
func (_FeeManager *FeeManagerFilterer) WatchFeeRateUpdated(opts *bind.WatchOpts, sink chan<- *FeeManagerFeeRateUpdated) (event.Subscription, error) {

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "FeeRateUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeManagerFeeRateUpdated)
				if err := _FeeManager.contract.UnpackLog(event, "FeeRateUpdated", log); err != nil {
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

// ParseFeeRateUpdated is a log parse operation binding the contract event 0x14914da2bf76024616fbe1859783fcd4dbddcb179b1f3a854949fbf920dcb957.
//
// Solidity: event FeeRateUpdated(uint256 oldRate, uint256 newRate)
func (_FeeManager *FeeManagerFilterer) ParseFeeRateUpdated(log types.Log) (*FeeManagerFeeRateUpdated, error) {
	event := new(FeeManagerFeeRateUpdated)
	if err := _FeeManager.contract.UnpackLog(event, "FeeRateUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the FeeManager contract.
type FeeManagerInitializedIterator struct {
	Event *FeeManagerInitialized // Event containing the contract specifics and raw log

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
func (it *FeeManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerInitialized)
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
		it.Event = new(FeeManagerInitialized)
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
func (it *FeeManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeManagerInitialized represents a Initialized event raised by the FeeManager contract.
type FeeManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FeeManager *FeeManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*FeeManagerInitializedIterator, error) {

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FeeManagerInitializedIterator{contract: _FeeManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_FeeManager *FeeManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FeeManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeManagerInitialized)
				if err := _FeeManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_FeeManager *FeeManagerFilterer) ParseInitialized(log types.Log) (*FeeManagerInitialized, error) {
	event := new(FeeManagerInitialized)
	if err := _FeeManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FeeManager contract.
type FeeManagerOwnershipTransferredIterator struct {
	Event *FeeManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *FeeManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerOwnershipTransferred)
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
		it.Event = new(FeeManagerOwnershipTransferred)
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
func (it *FeeManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeManagerOwnershipTransferred represents a OwnershipTransferred event raised by the FeeManager contract.
type FeeManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FeeManager *FeeManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FeeManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FeeManagerOwnershipTransferredIterator{contract: _FeeManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FeeManager *FeeManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeManagerOwnershipTransferred)
				if err := _FeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_FeeManager *FeeManagerFilterer) ParseOwnershipTransferred(log types.Log) (*FeeManagerOwnershipTransferred, error) {
	event := new(FeeManagerOwnershipTransferred)
	if err := _FeeManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeManagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the FeeManager contract.
type FeeManagerUpgradedIterator struct {
	Event *FeeManagerUpgraded // Event containing the contract specifics and raw log

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
func (it *FeeManagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeManagerUpgraded)
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
		it.Event = new(FeeManagerUpgraded)
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
func (it *FeeManagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeManagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeManagerUpgraded represents a Upgraded event raised by the FeeManager contract.
type FeeManagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FeeManager *FeeManagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*FeeManagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FeeManager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &FeeManagerUpgradedIterator{contract: _FeeManager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_FeeManager *FeeManagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *FeeManagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _FeeManager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeManagerUpgraded)
				if err := _FeeManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_FeeManager *FeeManagerFilterer) ParseUpgraded(log types.Log) (*FeeManagerUpgraded, error) {
	event := new(FeeManagerUpgraded)
	if err := _FeeManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
