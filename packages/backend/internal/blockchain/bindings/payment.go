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

// IPaymentBill is an auto generated low-level Go binding around an user-defined struct.
type IPaymentBill struct {
	Id          *big.Int
	User        common.Address
	OperatorId  *big.Int
	Amount      *big.Int
	PlatformFee *big.Int
	CreatedAt   *big.Int
	IsPaid      bool
}

// PaymentMetaData contains all meta data concerning the Payment contract.
var PaymentMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"SafeERC20FailedOperation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"billId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"platformFee\",\"type\":\"uint256\"}],\"name\":\"BillCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"billId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"operatorAmount\",\"type\":\"uint256\"}],\"name\":\"BillPaid\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"billId\",\"type\":\"uint256\"}],\"name\":\"TrafficCardApplied\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"billId\",\"type\":\"uint256\"}],\"name\":\"applyTrafficCardToBill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"createBill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"feeManager\",\"outputs\":[{\"internalType\":\"contractIFeeManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUnpaidBills\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"platformFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPaid\",\"type\":\"bool\"}],\"internalType\":\"structIPayment.Bill[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserBills\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"platformFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isPaid\",\"type\":\"bool\"}],\"internalType\":\"structIPayment.Bill[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_platformWallet\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_usdt\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_serviceManager\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"billId\",\"type\":\"uint256\"}],\"name\":\"payBill\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"platformWallet\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"serviceManager\",\"outputs\":[{\"internalType\":\"contractIServiceManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_feeManager\",\"type\":\"address\"}],\"name\":\"setFeeManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_platformWallet\",\"type\":\"address\"}],\"name\":\"setPlatformWallet\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_serviceManager\",\"type\":\"address\"}],\"name\":\"setServiceManager\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"usdt\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PaymentABI is the input ABI used to generate the binding from.
// Deprecated: Use PaymentMetaData.ABI instead.
var PaymentABI = PaymentMetaData.ABI

// Payment is an auto generated Go binding around an Ethereum contract.
type Payment struct {
	PaymentCaller     // Read-only binding to the contract
	PaymentTransactor // Write-only binding to the contract
	PaymentFilterer   // Log filterer for contract events
}

// PaymentCaller is an auto generated read-only Go binding around an Ethereum contract.
type PaymentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PaymentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PaymentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PaymentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PaymentSession struct {
	Contract     *Payment          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PaymentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PaymentCallerSession struct {
	Contract *PaymentCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// PaymentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PaymentTransactorSession struct {
	Contract     *PaymentTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// PaymentRaw is an auto generated low-level Go binding around an Ethereum contract.
type PaymentRaw struct {
	Contract *Payment // Generic contract binding to access the raw methods on
}

// PaymentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PaymentCallerRaw struct {
	Contract *PaymentCaller // Generic read-only contract binding to access the raw methods on
}

// PaymentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PaymentTransactorRaw struct {
	Contract *PaymentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPayment creates a new instance of Payment, bound to a specific deployed contract.
func NewPayment(address common.Address, backend bind.ContractBackend) (*Payment, error) {
	contract, err := bindPayment(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Payment{PaymentCaller: PaymentCaller{contract: contract}, PaymentTransactor: PaymentTransactor{contract: contract}, PaymentFilterer: PaymentFilterer{contract: contract}}, nil
}

// NewPaymentCaller creates a new read-only instance of Payment, bound to a specific deployed contract.
func NewPaymentCaller(address common.Address, caller bind.ContractCaller) (*PaymentCaller, error) {
	contract, err := bindPayment(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentCaller{contract: contract}, nil
}

// NewPaymentTransactor creates a new write-only instance of Payment, bound to a specific deployed contract.
func NewPaymentTransactor(address common.Address, transactor bind.ContractTransactor) (*PaymentTransactor, error) {
	contract, err := bindPayment(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PaymentTransactor{contract: contract}, nil
}

// NewPaymentFilterer creates a new log filterer instance of Payment, bound to a specific deployed contract.
func NewPaymentFilterer(address common.Address, filterer bind.ContractFilterer) (*PaymentFilterer, error) {
	contract, err := bindPayment(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PaymentFilterer{contract: contract}, nil
}

// bindPayment binds a generic wrapper to an already deployed contract.
func bindPayment(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PaymentMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Payment *PaymentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Payment.Contract.PaymentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Payment *PaymentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Payment.Contract.PaymentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Payment *PaymentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Payment.Contract.PaymentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Payment *PaymentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Payment.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Payment *PaymentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Payment.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Payment *PaymentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Payment.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Payment *PaymentCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Payment *PaymentSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Payment.Contract.UPGRADEINTERFACEVERSION(&_Payment.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_Payment *PaymentCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _Payment.Contract.UPGRADEINTERFACEVERSION(&_Payment.CallOpts)
}

// FeeManager is a free data retrieval call binding the contract method 0xd0fb0203.
//
// Solidity: function feeManager() view returns(address)
func (_Payment *PaymentCaller) FeeManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "feeManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeManager is a free data retrieval call binding the contract method 0xd0fb0203.
//
// Solidity: function feeManager() view returns(address)
func (_Payment *PaymentSession) FeeManager() (common.Address, error) {
	return _Payment.Contract.FeeManager(&_Payment.CallOpts)
}

// FeeManager is a free data retrieval call binding the contract method 0xd0fb0203.
//
// Solidity: function feeManager() view returns(address)
func (_Payment *PaymentCallerSession) FeeManager() (common.Address, error) {
	return _Payment.Contract.FeeManager(&_Payment.CallOpts)
}

// GetUnpaidBills is a free data retrieval call binding the contract method 0xf0b8e016.
//
// Solidity: function getUnpaidBills(address user) view returns((uint256,address,uint256,uint256,uint256,uint256,bool)[])
func (_Payment *PaymentCaller) GetUnpaidBills(opts *bind.CallOpts, user common.Address) ([]IPaymentBill, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "getUnpaidBills", user)

	if err != nil {
		return *new([]IPaymentBill), err
	}

	out0 := *abi.ConvertType(out[0], new([]IPaymentBill)).(*[]IPaymentBill)

	return out0, err

}

// GetUnpaidBills is a free data retrieval call binding the contract method 0xf0b8e016.
//
// Solidity: function getUnpaidBills(address user) view returns((uint256,address,uint256,uint256,uint256,uint256,bool)[])
func (_Payment *PaymentSession) GetUnpaidBills(user common.Address) ([]IPaymentBill, error) {
	return _Payment.Contract.GetUnpaidBills(&_Payment.CallOpts, user)
}

// GetUnpaidBills is a free data retrieval call binding the contract method 0xf0b8e016.
//
// Solidity: function getUnpaidBills(address user) view returns((uint256,address,uint256,uint256,uint256,uint256,bool)[])
func (_Payment *PaymentCallerSession) GetUnpaidBills(user common.Address) ([]IPaymentBill, error) {
	return _Payment.Contract.GetUnpaidBills(&_Payment.CallOpts, user)
}

// GetUserBills is a free data retrieval call binding the contract method 0xe65ef7f5.
//
// Solidity: function getUserBills(address user) view returns((uint256,address,uint256,uint256,uint256,uint256,bool)[])
func (_Payment *PaymentCaller) GetUserBills(opts *bind.CallOpts, user common.Address) ([]IPaymentBill, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "getUserBills", user)

	if err != nil {
		return *new([]IPaymentBill), err
	}

	out0 := *abi.ConvertType(out[0], new([]IPaymentBill)).(*[]IPaymentBill)

	return out0, err

}

// GetUserBills is a free data retrieval call binding the contract method 0xe65ef7f5.
//
// Solidity: function getUserBills(address user) view returns((uint256,address,uint256,uint256,uint256,uint256,bool)[])
func (_Payment *PaymentSession) GetUserBills(user common.Address) ([]IPaymentBill, error) {
	return _Payment.Contract.GetUserBills(&_Payment.CallOpts, user)
}

// GetUserBills is a free data retrieval call binding the contract method 0xe65ef7f5.
//
// Solidity: function getUserBills(address user) view returns((uint256,address,uint256,uint256,uint256,uint256,bool)[])
func (_Payment *PaymentCallerSession) GetUserBills(user common.Address) ([]IPaymentBill, error) {
	return _Payment.Contract.GetUserBills(&_Payment.CallOpts, user)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Payment *PaymentCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Payment *PaymentSession) Oracle() (common.Address, error) {
	return _Payment.Contract.Oracle(&_Payment.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_Payment *PaymentCallerSession) Oracle() (common.Address, error) {
	return _Payment.Contract.Oracle(&_Payment.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Payment *PaymentCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Payment *PaymentSession) Owner() (common.Address, error) {
	return _Payment.Contract.Owner(&_Payment.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Payment *PaymentCallerSession) Owner() (common.Address, error) {
	return _Payment.Contract.Owner(&_Payment.CallOpts)
}

// PlatformWallet is a free data retrieval call binding the contract method 0xfa2af9da.
//
// Solidity: function platformWallet() view returns(address)
func (_Payment *PaymentCaller) PlatformWallet(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "platformWallet")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PlatformWallet is a free data retrieval call binding the contract method 0xfa2af9da.
//
// Solidity: function platformWallet() view returns(address)
func (_Payment *PaymentSession) PlatformWallet() (common.Address, error) {
	return _Payment.Contract.PlatformWallet(&_Payment.CallOpts)
}

// PlatformWallet is a free data retrieval call binding the contract method 0xfa2af9da.
//
// Solidity: function platformWallet() view returns(address)
func (_Payment *PaymentCallerSession) PlatformWallet() (common.Address, error) {
	return _Payment.Contract.PlatformWallet(&_Payment.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Payment *PaymentCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Payment *PaymentSession) ProxiableUUID() ([32]byte, error) {
	return _Payment.Contract.ProxiableUUID(&_Payment.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_Payment *PaymentCallerSession) ProxiableUUID() ([32]byte, error) {
	return _Payment.Contract.ProxiableUUID(&_Payment.CallOpts)
}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_Payment *PaymentCaller) ServiceManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "serviceManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_Payment *PaymentSession) ServiceManager() (common.Address, error) {
	return _Payment.Contract.ServiceManager(&_Payment.CallOpts)
}

// ServiceManager is a free data retrieval call binding the contract method 0x3998fdd3.
//
// Solidity: function serviceManager() view returns(address)
func (_Payment *PaymentCallerSession) ServiceManager() (common.Address, error) {
	return _Payment.Contract.ServiceManager(&_Payment.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_Payment *PaymentCaller) Usdt(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Payment.contract.Call(opts, &out, "usdt")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_Payment *PaymentSession) Usdt() (common.Address, error) {
	return _Payment.Contract.Usdt(&_Payment.CallOpts)
}

// Usdt is a free data retrieval call binding the contract method 0x2f48ab7d.
//
// Solidity: function usdt() view returns(address)
func (_Payment *PaymentCallerSession) Usdt() (common.Address, error) {
	return _Payment.Contract.Usdt(&_Payment.CallOpts)
}

// ApplyTrafficCardToBill is a paid mutator transaction binding the contract method 0x0d907c34.
//
// Solidity: function applyTrafficCardToBill(uint256 billId) returns()
func (_Payment *PaymentTransactor) ApplyTrafficCardToBill(opts *bind.TransactOpts, billId *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "applyTrafficCardToBill", billId)
}

// ApplyTrafficCardToBill is a paid mutator transaction binding the contract method 0x0d907c34.
//
// Solidity: function applyTrafficCardToBill(uint256 billId) returns()
func (_Payment *PaymentSession) ApplyTrafficCardToBill(billId *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.ApplyTrafficCardToBill(&_Payment.TransactOpts, billId)
}

// ApplyTrafficCardToBill is a paid mutator transaction binding the contract method 0x0d907c34.
//
// Solidity: function applyTrafficCardToBill(uint256 billId) returns()
func (_Payment *PaymentTransactorSession) ApplyTrafficCardToBill(billId *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.ApplyTrafficCardToBill(&_Payment.TransactOpts, billId)
}

// CreateBill is a paid mutator transaction binding the contract method 0xceb323e8.
//
// Solidity: function createBill(address user, uint256 operatorId, uint256 amount) returns()
func (_Payment *PaymentTransactor) CreateBill(opts *bind.TransactOpts, user common.Address, operatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "createBill", user, operatorId, amount)
}

// CreateBill is a paid mutator transaction binding the contract method 0xceb323e8.
//
// Solidity: function createBill(address user, uint256 operatorId, uint256 amount) returns()
func (_Payment *PaymentSession) CreateBill(user common.Address, operatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateBill(&_Payment.TransactOpts, user, operatorId, amount)
}

// CreateBill is a paid mutator transaction binding the contract method 0xceb323e8.
//
// Solidity: function createBill(address user, uint256 operatorId, uint256 amount) returns()
func (_Payment *PaymentTransactorSession) CreateBill(user common.Address, operatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.CreateBill(&_Payment.TransactOpts, user, operatorId, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address _feeManager, address _platformWallet, address _usdt, address _serviceManager) returns()
func (_Payment *PaymentTransactor) Initialize(opts *bind.TransactOpts, _feeManager common.Address, _platformWallet common.Address, _usdt common.Address, _serviceManager common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "initialize", _feeManager, _platformWallet, _usdt, _serviceManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address _feeManager, address _platformWallet, address _usdt, address _serviceManager) returns()
func (_Payment *PaymentSession) Initialize(_feeManager common.Address, _platformWallet common.Address, _usdt common.Address, _serviceManager common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Initialize(&_Payment.TransactOpts, _feeManager, _platformWallet, _usdt, _serviceManager)
}

// Initialize is a paid mutator transaction binding the contract method 0xf8c8765e.
//
// Solidity: function initialize(address _feeManager, address _platformWallet, address _usdt, address _serviceManager) returns()
func (_Payment *PaymentTransactorSession) Initialize(_feeManager common.Address, _platformWallet common.Address, _usdt common.Address, _serviceManager common.Address) (*types.Transaction, error) {
	return _Payment.Contract.Initialize(&_Payment.TransactOpts, _feeManager, _platformWallet, _usdt, _serviceManager)
}

// PayBill is a paid mutator transaction binding the contract method 0xf0975190.
//
// Solidity: function payBill(uint256 billId) returns()
func (_Payment *PaymentTransactor) PayBill(opts *bind.TransactOpts, billId *big.Int) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "payBill", billId)
}

// PayBill is a paid mutator transaction binding the contract method 0xf0975190.
//
// Solidity: function payBill(uint256 billId) returns()
func (_Payment *PaymentSession) PayBill(billId *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.PayBill(&_Payment.TransactOpts, billId)
}

// PayBill is a paid mutator transaction binding the contract method 0xf0975190.
//
// Solidity: function payBill(uint256 billId) returns()
func (_Payment *PaymentTransactorSession) PayBill(billId *big.Int) (*types.Transaction, error) {
	return _Payment.Contract.PayBill(&_Payment.TransactOpts, billId)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Payment *PaymentTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Payment *PaymentSession) RenounceOwnership() (*types.Transaction, error) {
	return _Payment.Contract.RenounceOwnership(&_Payment.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Payment *PaymentTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Payment.Contract.RenounceOwnership(&_Payment.TransactOpts)
}

// SetFeeManager is a paid mutator transaction binding the contract method 0x472d35b9.
//
// Solidity: function setFeeManager(address _feeManager) returns()
func (_Payment *PaymentTransactor) SetFeeManager(opts *bind.TransactOpts, _feeManager common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "setFeeManager", _feeManager)
}

// SetFeeManager is a paid mutator transaction binding the contract method 0x472d35b9.
//
// Solidity: function setFeeManager(address _feeManager) returns()
func (_Payment *PaymentSession) SetFeeManager(_feeManager common.Address) (*types.Transaction, error) {
	return _Payment.Contract.SetFeeManager(&_Payment.TransactOpts, _feeManager)
}

// SetFeeManager is a paid mutator transaction binding the contract method 0x472d35b9.
//
// Solidity: function setFeeManager(address _feeManager) returns()
func (_Payment *PaymentTransactorSession) SetFeeManager(_feeManager common.Address) (*types.Transaction, error) {
	return _Payment.Contract.SetFeeManager(&_Payment.TransactOpts, _feeManager)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_Payment *PaymentTransactor) SetOracle(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "setOracle", _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_Payment *PaymentSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _Payment.Contract.SetOracle(&_Payment.TransactOpts, _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_Payment *PaymentTransactorSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _Payment.Contract.SetOracle(&_Payment.TransactOpts, _oracle)
}

// SetPlatformWallet is a paid mutator transaction binding the contract method 0x8831e9cf.
//
// Solidity: function setPlatformWallet(address _platformWallet) returns()
func (_Payment *PaymentTransactor) SetPlatformWallet(opts *bind.TransactOpts, _platformWallet common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "setPlatformWallet", _platformWallet)
}

// SetPlatformWallet is a paid mutator transaction binding the contract method 0x8831e9cf.
//
// Solidity: function setPlatformWallet(address _platformWallet) returns()
func (_Payment *PaymentSession) SetPlatformWallet(_platformWallet common.Address) (*types.Transaction, error) {
	return _Payment.Contract.SetPlatformWallet(&_Payment.TransactOpts, _platformWallet)
}

// SetPlatformWallet is a paid mutator transaction binding the contract method 0x8831e9cf.
//
// Solidity: function setPlatformWallet(address _platformWallet) returns()
func (_Payment *PaymentTransactorSession) SetPlatformWallet(_platformWallet common.Address) (*types.Transaction, error) {
	return _Payment.Contract.SetPlatformWallet(&_Payment.TransactOpts, _platformWallet)
}

// SetServiceManager is a paid mutator transaction binding the contract method 0x9b41bf23.
//
// Solidity: function setServiceManager(address _serviceManager) returns()
func (_Payment *PaymentTransactor) SetServiceManager(opts *bind.TransactOpts, _serviceManager common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "setServiceManager", _serviceManager)
}

// SetServiceManager is a paid mutator transaction binding the contract method 0x9b41bf23.
//
// Solidity: function setServiceManager(address _serviceManager) returns()
func (_Payment *PaymentSession) SetServiceManager(_serviceManager common.Address) (*types.Transaction, error) {
	return _Payment.Contract.SetServiceManager(&_Payment.TransactOpts, _serviceManager)
}

// SetServiceManager is a paid mutator transaction binding the contract method 0x9b41bf23.
//
// Solidity: function setServiceManager(address _serviceManager) returns()
func (_Payment *PaymentTransactorSession) SetServiceManager(_serviceManager common.Address) (*types.Transaction, error) {
	return _Payment.Contract.SetServiceManager(&_Payment.TransactOpts, _serviceManager)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Payment *PaymentTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Payment *PaymentSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Payment.Contract.TransferOwnership(&_Payment.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Payment *PaymentTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Payment.Contract.TransferOwnership(&_Payment.TransactOpts, newOwner)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Payment *PaymentTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Payment.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Payment *PaymentSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Payment.Contract.UpgradeToAndCall(&_Payment.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_Payment *PaymentTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _Payment.Contract.UpgradeToAndCall(&_Payment.TransactOpts, newImplementation, data)
}

// PaymentBillCreatedIterator is returned from FilterBillCreated and is used to iterate over the raw logs and unpacked data for BillCreated events raised by the Payment contract.
type PaymentBillCreatedIterator struct {
	Event *PaymentBillCreated // Event containing the contract specifics and raw log

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
func (it *PaymentBillCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentBillCreated)
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
		it.Event = new(PaymentBillCreated)
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
func (it *PaymentBillCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentBillCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentBillCreated represents a BillCreated event raised by the Payment contract.
type PaymentBillCreated struct {
	BillId      *big.Int
	User        common.Address
	TotalAmount *big.Int
	PlatformFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBillCreated is a free log retrieval operation binding the contract event 0xcdfdeecd9f301cb609cbfd87c3a7f1e4d3da395ff0ba5084da583d3b6deced21.
//
// Solidity: event BillCreated(uint256 indexed billId, address indexed user, uint256 totalAmount, uint256 platformFee)
func (_Payment *PaymentFilterer) FilterBillCreated(opts *bind.FilterOpts, billId []*big.Int, user []common.Address) (*PaymentBillCreatedIterator, error) {

	var billIdRule []interface{}
	for _, billIdItem := range billId {
		billIdRule = append(billIdRule, billIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "BillCreated", billIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PaymentBillCreatedIterator{contract: _Payment.contract, event: "BillCreated", logs: logs, sub: sub}, nil
}

// WatchBillCreated is a free log subscription operation binding the contract event 0xcdfdeecd9f301cb609cbfd87c3a7f1e4d3da395ff0ba5084da583d3b6deced21.
//
// Solidity: event BillCreated(uint256 indexed billId, address indexed user, uint256 totalAmount, uint256 platformFee)
func (_Payment *PaymentFilterer) WatchBillCreated(opts *bind.WatchOpts, sink chan<- *PaymentBillCreated, billId []*big.Int, user []common.Address) (event.Subscription, error) {

	var billIdRule []interface{}
	for _, billIdItem := range billId {
		billIdRule = append(billIdRule, billIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "BillCreated", billIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentBillCreated)
				if err := _Payment.contract.UnpackLog(event, "BillCreated", log); err != nil {
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

// ParseBillCreated is a log parse operation binding the contract event 0xcdfdeecd9f301cb609cbfd87c3a7f1e4d3da395ff0ba5084da583d3b6deced21.
//
// Solidity: event BillCreated(uint256 indexed billId, address indexed user, uint256 totalAmount, uint256 platformFee)
func (_Payment *PaymentFilterer) ParseBillCreated(log types.Log) (*PaymentBillCreated, error) {
	event := new(PaymentBillCreated)
	if err := _Payment.contract.UnpackLog(event, "BillCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentBillPaidIterator is returned from FilterBillPaid and is used to iterate over the raw logs and unpacked data for BillPaid events raised by the Payment contract.
type PaymentBillPaidIterator struct {
	Event *PaymentBillPaid // Event containing the contract specifics and raw log

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
func (it *PaymentBillPaidIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentBillPaid)
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
		it.Event = new(PaymentBillPaid)
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
func (it *PaymentBillPaidIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentBillPaidIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentBillPaid represents a BillPaid event raised by the Payment contract.
type PaymentBillPaid struct {
	BillId         *big.Int
	User           common.Address
	TotalAmount    *big.Int
	OperatorAmount *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterBillPaid is a free log retrieval operation binding the contract event 0x53646f88205e1dd1de6fdaa8898d840fbe17129cb98729533d3c83a3a0e6045a.
//
// Solidity: event BillPaid(uint256 indexed billId, address indexed user, uint256 totalAmount, uint256 operatorAmount)
func (_Payment *PaymentFilterer) FilterBillPaid(opts *bind.FilterOpts, billId []*big.Int, user []common.Address) (*PaymentBillPaidIterator, error) {

	var billIdRule []interface{}
	for _, billIdItem := range billId {
		billIdRule = append(billIdRule, billIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "BillPaid", billIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return &PaymentBillPaidIterator{contract: _Payment.contract, event: "BillPaid", logs: logs, sub: sub}, nil
}

// WatchBillPaid is a free log subscription operation binding the contract event 0x53646f88205e1dd1de6fdaa8898d840fbe17129cb98729533d3c83a3a0e6045a.
//
// Solidity: event BillPaid(uint256 indexed billId, address indexed user, uint256 totalAmount, uint256 operatorAmount)
func (_Payment *PaymentFilterer) WatchBillPaid(opts *bind.WatchOpts, sink chan<- *PaymentBillPaid, billId []*big.Int, user []common.Address) (event.Subscription, error) {

	var billIdRule []interface{}
	for _, billIdItem := range billId {
		billIdRule = append(billIdRule, billIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "BillPaid", billIdRule, userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentBillPaid)
				if err := _Payment.contract.UnpackLog(event, "BillPaid", log); err != nil {
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

// ParseBillPaid is a log parse operation binding the contract event 0x53646f88205e1dd1de6fdaa8898d840fbe17129cb98729533d3c83a3a0e6045a.
//
// Solidity: event BillPaid(uint256 indexed billId, address indexed user, uint256 totalAmount, uint256 operatorAmount)
func (_Payment *PaymentFilterer) ParseBillPaid(log types.Log) (*PaymentBillPaid, error) {
	event := new(PaymentBillPaid)
	if err := _Payment.contract.UnpackLog(event, "BillPaid", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Payment contract.
type PaymentInitializedIterator struct {
	Event *PaymentInitialized // Event containing the contract specifics and raw log

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
func (it *PaymentInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentInitialized)
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
		it.Event = new(PaymentInitialized)
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
func (it *PaymentInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentInitialized represents a Initialized event raised by the Payment contract.
type PaymentInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Payment *PaymentFilterer) FilterInitialized(opts *bind.FilterOpts) (*PaymentInitializedIterator, error) {

	logs, sub, err := _Payment.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PaymentInitializedIterator{contract: _Payment.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_Payment *PaymentFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PaymentInitialized) (event.Subscription, error) {

	logs, sub, err := _Payment.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentInitialized)
				if err := _Payment.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Payment *PaymentFilterer) ParseInitialized(log types.Log) (*PaymentInitialized, error) {
	event := new(PaymentInitialized)
	if err := _Payment.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Payment contract.
type PaymentOwnershipTransferredIterator struct {
	Event *PaymentOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PaymentOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentOwnershipTransferred)
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
		it.Event = new(PaymentOwnershipTransferred)
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
func (it *PaymentOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentOwnershipTransferred represents a OwnershipTransferred event raised by the Payment contract.
type PaymentOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Payment *PaymentFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PaymentOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PaymentOwnershipTransferredIterator{contract: _Payment.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Payment *PaymentFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PaymentOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentOwnershipTransferred)
				if err := _Payment.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Payment *PaymentFilterer) ParseOwnershipTransferred(log types.Log) (*PaymentOwnershipTransferred, error) {
	event := new(PaymentOwnershipTransferred)
	if err := _Payment.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentTrafficCardAppliedIterator is returned from FilterTrafficCardApplied and is used to iterate over the raw logs and unpacked data for TrafficCardApplied events raised by the Payment contract.
type PaymentTrafficCardAppliedIterator struct {
	Event *PaymentTrafficCardApplied // Event containing the contract specifics and raw log

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
func (it *PaymentTrafficCardAppliedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentTrafficCardApplied)
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
		it.Event = new(PaymentTrafficCardApplied)
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
func (it *PaymentTrafficCardAppliedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentTrafficCardAppliedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentTrafficCardApplied represents a TrafficCardApplied event raised by the Payment contract.
type PaymentTrafficCardApplied struct {
	BillId *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTrafficCardApplied is a free log retrieval operation binding the contract event 0x6ee1062e7611525ff44a2bea1e6ccffff047903a965c69aa96334ad680cd8701.
//
// Solidity: event TrafficCardApplied(uint256 indexed billId)
func (_Payment *PaymentFilterer) FilterTrafficCardApplied(opts *bind.FilterOpts, billId []*big.Int) (*PaymentTrafficCardAppliedIterator, error) {

	var billIdRule []interface{}
	for _, billIdItem := range billId {
		billIdRule = append(billIdRule, billIdItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "TrafficCardApplied", billIdRule)
	if err != nil {
		return nil, err
	}
	return &PaymentTrafficCardAppliedIterator{contract: _Payment.contract, event: "TrafficCardApplied", logs: logs, sub: sub}, nil
}

// WatchTrafficCardApplied is a free log subscription operation binding the contract event 0x6ee1062e7611525ff44a2bea1e6ccffff047903a965c69aa96334ad680cd8701.
//
// Solidity: event TrafficCardApplied(uint256 indexed billId)
func (_Payment *PaymentFilterer) WatchTrafficCardApplied(opts *bind.WatchOpts, sink chan<- *PaymentTrafficCardApplied, billId []*big.Int) (event.Subscription, error) {

	var billIdRule []interface{}
	for _, billIdItem := range billId {
		billIdRule = append(billIdRule, billIdItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "TrafficCardApplied", billIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentTrafficCardApplied)
				if err := _Payment.contract.UnpackLog(event, "TrafficCardApplied", log); err != nil {
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

// ParseTrafficCardApplied is a log parse operation binding the contract event 0x6ee1062e7611525ff44a2bea1e6ccffff047903a965c69aa96334ad680cd8701.
//
// Solidity: event TrafficCardApplied(uint256 indexed billId)
func (_Payment *PaymentFilterer) ParseTrafficCardApplied(log types.Log) (*PaymentTrafficCardApplied, error) {
	event := new(PaymentTrafficCardApplied)
	if err := _Payment.contract.UnpackLog(event, "TrafficCardApplied", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PaymentUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the Payment contract.
type PaymentUpgradedIterator struct {
	Event *PaymentUpgraded // Event containing the contract specifics and raw log

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
func (it *PaymentUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PaymentUpgraded)
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
		it.Event = new(PaymentUpgraded)
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
func (it *PaymentUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PaymentUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PaymentUpgraded represents a Upgraded event raised by the Payment contract.
type PaymentUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Payment *PaymentFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*PaymentUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Payment.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &PaymentUpgradedIterator{contract: _Payment.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_Payment *PaymentFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *PaymentUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _Payment.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PaymentUpgraded)
				if err := _Payment.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_Payment *PaymentFilterer) ParseUpgraded(log types.Log) (*PaymentUpgraded, error) {
	event := new(PaymentUpgraded)
	if err := _Payment.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
