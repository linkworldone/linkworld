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

// IServiceManagerOperator is an auto generated low-level Go binding around an user-defined struct.
type IServiceManagerOperator struct {
	Id              *big.Int
	Name            string
	Region          string
	CountryCode     string
	RequiredDeposit *big.Int
	IsActive        bool
	PaymentAddress  common.Address
}

// ServiceManagerMetaData contains all meta data concerning the ServiceManager contract.
var ServiceManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"target\",\"type\":\"address\"}],\"name\":\"AddressEmptyCode\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"ERC1967InvalidImplementation\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ERC1967NonPayable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FailedCall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UUPSUnauthorizedCallContext\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"UUPSUnsupportedProxiableUUID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"region\",\"type\":\"string\"}],\"name\":\"OperatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"}],\"name\":\"OperatorDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"paymentAddress\",\"type\":\"address\"}],\"name\":\"OperatorPaymentAddressSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"}],\"name\":\"OperatorUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"UPGRADE_INTERFACE_VERSION\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"region\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"countryCode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"requiredDeposit\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"paymentAddress\",\"type\":\"address\"}],\"name\":\"addOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"}],\"name\":\"deactivateOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getActiveOperators\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"region\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"countryCode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"requiredDeposit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"paymentAddress\",\"type\":\"address\"}],\"internalType\":\"structIServiceManager.Operator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"}],\"name\":\"getOperator\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"region\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"countryCode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"requiredDeposit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"paymentAddress\",\"type\":\"address\"}],\"internalType\":\"structIServiceManager.Operator\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"countryCode\",\"type\":\"string\"}],\"name\":\"getOperatorsByCountry\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"region\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"countryCode\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"requiredDeposit\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"paymentAddress\",\"type\":\"address\"}],\"internalType\":\"structIServiceManager.Operator[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"proxiableUUID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"paymentAddress\",\"type\":\"address\"}],\"name\":\"setOperatorPaymentAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"operatorId\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"name\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"region\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"requiredDeposit\",\"type\":\"uint256\"}],\"name\":\"updateOperator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newImplementation\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"upgradeToAndCall\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
}

// ServiceManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use ServiceManagerMetaData.ABI instead.
var ServiceManagerABI = ServiceManagerMetaData.ABI

// ServiceManager is an auto generated Go binding around an Ethereum contract.
type ServiceManager struct {
	ServiceManagerCaller     // Read-only binding to the contract
	ServiceManagerTransactor // Write-only binding to the contract
	ServiceManagerFilterer   // Log filterer for contract events
}

// ServiceManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type ServiceManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ServiceManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ServiceManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ServiceManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ServiceManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ServiceManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ServiceManagerSession struct {
	Contract     *ServiceManager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ServiceManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ServiceManagerCallerSession struct {
	Contract *ServiceManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// ServiceManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ServiceManagerTransactorSession struct {
	Contract     *ServiceManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// ServiceManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type ServiceManagerRaw struct {
	Contract *ServiceManager // Generic contract binding to access the raw methods on
}

// ServiceManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ServiceManagerCallerRaw struct {
	Contract *ServiceManagerCaller // Generic read-only contract binding to access the raw methods on
}

// ServiceManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ServiceManagerTransactorRaw struct {
	Contract *ServiceManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewServiceManager creates a new instance of ServiceManager, bound to a specific deployed contract.
func NewServiceManager(address common.Address, backend bind.ContractBackend) (*ServiceManager, error) {
	contract, err := bindServiceManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ServiceManager{ServiceManagerCaller: ServiceManagerCaller{contract: contract}, ServiceManagerTransactor: ServiceManagerTransactor{contract: contract}, ServiceManagerFilterer: ServiceManagerFilterer{contract: contract}}, nil
}

// NewServiceManagerCaller creates a new read-only instance of ServiceManager, bound to a specific deployed contract.
func NewServiceManagerCaller(address common.Address, caller bind.ContractCaller) (*ServiceManagerCaller, error) {
	contract, err := bindServiceManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerCaller{contract: contract}, nil
}

// NewServiceManagerTransactor creates a new write-only instance of ServiceManager, bound to a specific deployed contract.
func NewServiceManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*ServiceManagerTransactor, error) {
	contract, err := bindServiceManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerTransactor{contract: contract}, nil
}

// NewServiceManagerFilterer creates a new log filterer instance of ServiceManager, bound to a specific deployed contract.
func NewServiceManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*ServiceManagerFilterer, error) {
	contract, err := bindServiceManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerFilterer{contract: contract}, nil
}

// bindServiceManager binds a generic wrapper to an already deployed contract.
func bindServiceManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ServiceManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ServiceManager *ServiceManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ServiceManager.Contract.ServiceManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ServiceManager *ServiceManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ServiceManager.Contract.ServiceManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ServiceManager *ServiceManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ServiceManager.Contract.ServiceManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ServiceManager *ServiceManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ServiceManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ServiceManager *ServiceManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ServiceManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ServiceManager *ServiceManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ServiceManager.Contract.contract.Transact(opts, method, params...)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ServiceManager *ServiceManagerCaller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _ServiceManager.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ServiceManager *ServiceManagerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ServiceManager.Contract.UPGRADEINTERFACEVERSION(&_ServiceManager.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_ServiceManager *ServiceManagerCallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _ServiceManager.Contract.UPGRADEINTERFACEVERSION(&_ServiceManager.CallOpts)
}

// GetActiveOperators is a free data retrieval call binding the contract method 0x64bdc67e.
//
// Solidity: function getActiveOperators() view returns((uint256,string,string,string,uint256,bool,address)[])
func (_ServiceManager *ServiceManagerCaller) GetActiveOperators(opts *bind.CallOpts) ([]IServiceManagerOperator, error) {
	var out []interface{}
	err := _ServiceManager.contract.Call(opts, &out, "getActiveOperators")

	if err != nil {
		return *new([]IServiceManagerOperator), err
	}

	out0 := *abi.ConvertType(out[0], new([]IServiceManagerOperator)).(*[]IServiceManagerOperator)

	return out0, err

}

// GetActiveOperators is a free data retrieval call binding the contract method 0x64bdc67e.
//
// Solidity: function getActiveOperators() view returns((uint256,string,string,string,uint256,bool,address)[])
func (_ServiceManager *ServiceManagerSession) GetActiveOperators() ([]IServiceManagerOperator, error) {
	return _ServiceManager.Contract.GetActiveOperators(&_ServiceManager.CallOpts)
}

// GetActiveOperators is a free data retrieval call binding the contract method 0x64bdc67e.
//
// Solidity: function getActiveOperators() view returns((uint256,string,string,string,uint256,bool,address)[])
func (_ServiceManager *ServiceManagerCallerSession) GetActiveOperators() ([]IServiceManagerOperator, error) {
	return _ServiceManager.Contract.GetActiveOperators(&_ServiceManager.CallOpts)
}

// GetOperator is a free data retrieval call binding the contract method 0x05f63c8a.
//
// Solidity: function getOperator(uint256 operatorId) view returns((uint256,string,string,string,uint256,bool,address))
func (_ServiceManager *ServiceManagerCaller) GetOperator(opts *bind.CallOpts, operatorId *big.Int) (IServiceManagerOperator, error) {
	var out []interface{}
	err := _ServiceManager.contract.Call(opts, &out, "getOperator", operatorId)

	if err != nil {
		return *new(IServiceManagerOperator), err
	}

	out0 := *abi.ConvertType(out[0], new(IServiceManagerOperator)).(*IServiceManagerOperator)

	return out0, err

}

// GetOperator is a free data retrieval call binding the contract method 0x05f63c8a.
//
// Solidity: function getOperator(uint256 operatorId) view returns((uint256,string,string,string,uint256,bool,address))
func (_ServiceManager *ServiceManagerSession) GetOperator(operatorId *big.Int) (IServiceManagerOperator, error) {
	return _ServiceManager.Contract.GetOperator(&_ServiceManager.CallOpts, operatorId)
}

// GetOperator is a free data retrieval call binding the contract method 0x05f63c8a.
//
// Solidity: function getOperator(uint256 operatorId) view returns((uint256,string,string,string,uint256,bool,address))
func (_ServiceManager *ServiceManagerCallerSession) GetOperator(operatorId *big.Int) (IServiceManagerOperator, error) {
	return _ServiceManager.Contract.GetOperator(&_ServiceManager.CallOpts, operatorId)
}

// GetOperatorsByCountry is a free data retrieval call binding the contract method 0xb7e9889c.
//
// Solidity: function getOperatorsByCountry(string countryCode) view returns((uint256,string,string,string,uint256,bool,address)[])
func (_ServiceManager *ServiceManagerCaller) GetOperatorsByCountry(opts *bind.CallOpts, countryCode string) ([]IServiceManagerOperator, error) {
	var out []interface{}
	err := _ServiceManager.contract.Call(opts, &out, "getOperatorsByCountry", countryCode)

	if err != nil {
		return *new([]IServiceManagerOperator), err
	}

	out0 := *abi.ConvertType(out[0], new([]IServiceManagerOperator)).(*[]IServiceManagerOperator)

	return out0, err

}

// GetOperatorsByCountry is a free data retrieval call binding the contract method 0xb7e9889c.
//
// Solidity: function getOperatorsByCountry(string countryCode) view returns((uint256,string,string,string,uint256,bool,address)[])
func (_ServiceManager *ServiceManagerSession) GetOperatorsByCountry(countryCode string) ([]IServiceManagerOperator, error) {
	return _ServiceManager.Contract.GetOperatorsByCountry(&_ServiceManager.CallOpts, countryCode)
}

// GetOperatorsByCountry is a free data retrieval call binding the contract method 0xb7e9889c.
//
// Solidity: function getOperatorsByCountry(string countryCode) view returns((uint256,string,string,string,uint256,bool,address)[])
func (_ServiceManager *ServiceManagerCallerSession) GetOperatorsByCountry(countryCode string) ([]IServiceManagerOperator, error) {
	return _ServiceManager.Contract.GetOperatorsByCountry(&_ServiceManager.CallOpts, countryCode)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ServiceManager *ServiceManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _ServiceManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ServiceManager *ServiceManagerSession) Owner() (common.Address, error) {
	return _ServiceManager.Contract.Owner(&_ServiceManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_ServiceManager *ServiceManagerCallerSession) Owner() (common.Address, error) {
	return _ServiceManager.Contract.Owner(&_ServiceManager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ServiceManager *ServiceManagerCaller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _ServiceManager.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ServiceManager *ServiceManagerSession) ProxiableUUID() ([32]byte, error) {
	return _ServiceManager.Contract.ProxiableUUID(&_ServiceManager.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_ServiceManager *ServiceManagerCallerSession) ProxiableUUID() ([32]byte, error) {
	return _ServiceManager.Contract.ProxiableUUID(&_ServiceManager.CallOpts)
}

// AddOperator is a paid mutator transaction binding the contract method 0x1505b1f0.
//
// Solidity: function addOperator(string name, string region, string countryCode, uint256 requiredDeposit, address paymentAddress) returns()
func (_ServiceManager *ServiceManagerTransactor) AddOperator(opts *bind.TransactOpts, name string, region string, countryCode string, requiredDeposit *big.Int, paymentAddress common.Address) (*types.Transaction, error) {
	return _ServiceManager.contract.Transact(opts, "addOperator", name, region, countryCode, requiredDeposit, paymentAddress)
}

// AddOperator is a paid mutator transaction binding the contract method 0x1505b1f0.
//
// Solidity: function addOperator(string name, string region, string countryCode, uint256 requiredDeposit, address paymentAddress) returns()
func (_ServiceManager *ServiceManagerSession) AddOperator(name string, region string, countryCode string, requiredDeposit *big.Int, paymentAddress common.Address) (*types.Transaction, error) {
	return _ServiceManager.Contract.AddOperator(&_ServiceManager.TransactOpts, name, region, countryCode, requiredDeposit, paymentAddress)
}

// AddOperator is a paid mutator transaction binding the contract method 0x1505b1f0.
//
// Solidity: function addOperator(string name, string region, string countryCode, uint256 requiredDeposit, address paymentAddress) returns()
func (_ServiceManager *ServiceManagerTransactorSession) AddOperator(name string, region string, countryCode string, requiredDeposit *big.Int, paymentAddress common.Address) (*types.Transaction, error) {
	return _ServiceManager.Contract.AddOperator(&_ServiceManager.TransactOpts, name, region, countryCode, requiredDeposit, paymentAddress)
}

// DeactivateOperator is a paid mutator transaction binding the contract method 0x299cab4e.
//
// Solidity: function deactivateOperator(uint256 operatorId) returns()
func (_ServiceManager *ServiceManagerTransactor) DeactivateOperator(opts *bind.TransactOpts, operatorId *big.Int) (*types.Transaction, error) {
	return _ServiceManager.contract.Transact(opts, "deactivateOperator", operatorId)
}

// DeactivateOperator is a paid mutator transaction binding the contract method 0x299cab4e.
//
// Solidity: function deactivateOperator(uint256 operatorId) returns()
func (_ServiceManager *ServiceManagerSession) DeactivateOperator(operatorId *big.Int) (*types.Transaction, error) {
	return _ServiceManager.Contract.DeactivateOperator(&_ServiceManager.TransactOpts, operatorId)
}

// DeactivateOperator is a paid mutator transaction binding the contract method 0x299cab4e.
//
// Solidity: function deactivateOperator(uint256 operatorId) returns()
func (_ServiceManager *ServiceManagerTransactorSession) DeactivateOperator(operatorId *big.Int) (*types.Transaction, error) {
	return _ServiceManager.Contract.DeactivateOperator(&_ServiceManager.TransactOpts, operatorId)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ServiceManager *ServiceManagerTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ServiceManager.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ServiceManager *ServiceManagerSession) Initialize() (*types.Transaction, error) {
	return _ServiceManager.Contract.Initialize(&_ServiceManager.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_ServiceManager *ServiceManagerTransactorSession) Initialize() (*types.Transaction, error) {
	return _ServiceManager.Contract.Initialize(&_ServiceManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ServiceManager *ServiceManagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ServiceManager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ServiceManager *ServiceManagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _ServiceManager.Contract.RenounceOwnership(&_ServiceManager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_ServiceManager *ServiceManagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _ServiceManager.Contract.RenounceOwnership(&_ServiceManager.TransactOpts)
}

// SetOperatorPaymentAddress is a paid mutator transaction binding the contract method 0xadb76801.
//
// Solidity: function setOperatorPaymentAddress(uint256 operatorId, address paymentAddress) returns()
func (_ServiceManager *ServiceManagerTransactor) SetOperatorPaymentAddress(opts *bind.TransactOpts, operatorId *big.Int, paymentAddress common.Address) (*types.Transaction, error) {
	return _ServiceManager.contract.Transact(opts, "setOperatorPaymentAddress", operatorId, paymentAddress)
}

// SetOperatorPaymentAddress is a paid mutator transaction binding the contract method 0xadb76801.
//
// Solidity: function setOperatorPaymentAddress(uint256 operatorId, address paymentAddress) returns()
func (_ServiceManager *ServiceManagerSession) SetOperatorPaymentAddress(operatorId *big.Int, paymentAddress common.Address) (*types.Transaction, error) {
	return _ServiceManager.Contract.SetOperatorPaymentAddress(&_ServiceManager.TransactOpts, operatorId, paymentAddress)
}

// SetOperatorPaymentAddress is a paid mutator transaction binding the contract method 0xadb76801.
//
// Solidity: function setOperatorPaymentAddress(uint256 operatorId, address paymentAddress) returns()
func (_ServiceManager *ServiceManagerTransactorSession) SetOperatorPaymentAddress(operatorId *big.Int, paymentAddress common.Address) (*types.Transaction, error) {
	return _ServiceManager.Contract.SetOperatorPaymentAddress(&_ServiceManager.TransactOpts, operatorId, paymentAddress)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ServiceManager *ServiceManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _ServiceManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ServiceManager *ServiceManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ServiceManager.Contract.TransferOwnership(&_ServiceManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_ServiceManager *ServiceManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _ServiceManager.Contract.TransferOwnership(&_ServiceManager.TransactOpts, newOwner)
}

// UpdateOperator is a paid mutator transaction binding the contract method 0x82ba9917.
//
// Solidity: function updateOperator(uint256 operatorId, string name, string region, uint256 requiredDeposit) returns()
func (_ServiceManager *ServiceManagerTransactor) UpdateOperator(opts *bind.TransactOpts, operatorId *big.Int, name string, region string, requiredDeposit *big.Int) (*types.Transaction, error) {
	return _ServiceManager.contract.Transact(opts, "updateOperator", operatorId, name, region, requiredDeposit)
}

// UpdateOperator is a paid mutator transaction binding the contract method 0x82ba9917.
//
// Solidity: function updateOperator(uint256 operatorId, string name, string region, uint256 requiredDeposit) returns()
func (_ServiceManager *ServiceManagerSession) UpdateOperator(operatorId *big.Int, name string, region string, requiredDeposit *big.Int) (*types.Transaction, error) {
	return _ServiceManager.Contract.UpdateOperator(&_ServiceManager.TransactOpts, operatorId, name, region, requiredDeposit)
}

// UpdateOperator is a paid mutator transaction binding the contract method 0x82ba9917.
//
// Solidity: function updateOperator(uint256 operatorId, string name, string region, uint256 requiredDeposit) returns()
func (_ServiceManager *ServiceManagerTransactorSession) UpdateOperator(operatorId *big.Int, name string, region string, requiredDeposit *big.Int) (*types.Transaction, error) {
	return _ServiceManager.Contract.UpdateOperator(&_ServiceManager.TransactOpts, operatorId, name, region, requiredDeposit)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ServiceManager *ServiceManagerTransactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ServiceManager.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ServiceManager *ServiceManagerSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ServiceManager.Contract.UpgradeToAndCall(&_ServiceManager.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_ServiceManager *ServiceManagerTransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _ServiceManager.Contract.UpgradeToAndCall(&_ServiceManager.TransactOpts, newImplementation, data)
}

// ServiceManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the ServiceManager contract.
type ServiceManagerInitializedIterator struct {
	Event *ServiceManagerInitialized // Event containing the contract specifics and raw log

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
func (it *ServiceManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceManagerInitialized)
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
		it.Event = new(ServiceManagerInitialized)
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
func (it *ServiceManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceManagerInitialized represents a Initialized event raised by the ServiceManager contract.
type ServiceManagerInitialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ServiceManager *ServiceManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*ServiceManagerInitializedIterator, error) {

	logs, sub, err := _ServiceManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ServiceManagerInitializedIterator{contract: _ServiceManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_ServiceManager *ServiceManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ServiceManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _ServiceManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceManagerInitialized)
				if err := _ServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_ServiceManager *ServiceManagerFilterer) ParseInitialized(log types.Log) (*ServiceManagerInitialized, error) {
	event := new(ServiceManagerInitialized)
	if err := _ServiceManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ServiceManagerOperatorAddedIterator is returned from FilterOperatorAdded and is used to iterate over the raw logs and unpacked data for OperatorAdded events raised by the ServiceManager contract.
type ServiceManagerOperatorAddedIterator struct {
	Event *ServiceManagerOperatorAdded // Event containing the contract specifics and raw log

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
func (it *ServiceManagerOperatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceManagerOperatorAdded)
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
		it.Event = new(ServiceManagerOperatorAdded)
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
func (it *ServiceManagerOperatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceManagerOperatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceManagerOperatorAdded represents a OperatorAdded event raised by the ServiceManager contract.
type ServiceManagerOperatorAdded struct {
	OperatorId *big.Int
	Name       string
	Region     string
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOperatorAdded is a free log retrieval operation binding the contract event 0x5b79ce62dd937310f2ad88f3c0d314b45dc6dacff18e0a8025d9f8cb37fc6900.
//
// Solidity: event OperatorAdded(uint256 indexed operatorId, string name, string region)
func (_ServiceManager *ServiceManagerFilterer) FilterOperatorAdded(opts *bind.FilterOpts, operatorId []*big.Int) (*ServiceManagerOperatorAddedIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ServiceManager.contract.FilterLogs(opts, "OperatorAdded", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerOperatorAddedIterator{contract: _ServiceManager.contract, event: "OperatorAdded", logs: logs, sub: sub}, nil
}

// WatchOperatorAdded is a free log subscription operation binding the contract event 0x5b79ce62dd937310f2ad88f3c0d314b45dc6dacff18e0a8025d9f8cb37fc6900.
//
// Solidity: event OperatorAdded(uint256 indexed operatorId, string name, string region)
func (_ServiceManager *ServiceManagerFilterer) WatchOperatorAdded(opts *bind.WatchOpts, sink chan<- *ServiceManagerOperatorAdded, operatorId []*big.Int) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ServiceManager.contract.WatchLogs(opts, "OperatorAdded", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceManagerOperatorAdded)
				if err := _ServiceManager.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
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

// ParseOperatorAdded is a log parse operation binding the contract event 0x5b79ce62dd937310f2ad88f3c0d314b45dc6dacff18e0a8025d9f8cb37fc6900.
//
// Solidity: event OperatorAdded(uint256 indexed operatorId, string name, string region)
func (_ServiceManager *ServiceManagerFilterer) ParseOperatorAdded(log types.Log) (*ServiceManagerOperatorAdded, error) {
	event := new(ServiceManagerOperatorAdded)
	if err := _ServiceManager.contract.UnpackLog(event, "OperatorAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ServiceManagerOperatorDeactivatedIterator is returned from FilterOperatorDeactivated and is used to iterate over the raw logs and unpacked data for OperatorDeactivated events raised by the ServiceManager contract.
type ServiceManagerOperatorDeactivatedIterator struct {
	Event *ServiceManagerOperatorDeactivated // Event containing the contract specifics and raw log

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
func (it *ServiceManagerOperatorDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceManagerOperatorDeactivated)
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
		it.Event = new(ServiceManagerOperatorDeactivated)
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
func (it *ServiceManagerOperatorDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceManagerOperatorDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceManagerOperatorDeactivated represents a OperatorDeactivated event raised by the ServiceManager contract.
type ServiceManagerOperatorDeactivated struct {
	OperatorId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOperatorDeactivated is a free log retrieval operation binding the contract event 0x6ff0b7b14b65a91172fa7e0a7b49d909add202dadcef1b08c80f3c136a914fa8.
//
// Solidity: event OperatorDeactivated(uint256 indexed operatorId)
func (_ServiceManager *ServiceManagerFilterer) FilterOperatorDeactivated(opts *bind.FilterOpts, operatorId []*big.Int) (*ServiceManagerOperatorDeactivatedIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ServiceManager.contract.FilterLogs(opts, "OperatorDeactivated", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerOperatorDeactivatedIterator{contract: _ServiceManager.contract, event: "OperatorDeactivated", logs: logs, sub: sub}, nil
}

// WatchOperatorDeactivated is a free log subscription operation binding the contract event 0x6ff0b7b14b65a91172fa7e0a7b49d909add202dadcef1b08c80f3c136a914fa8.
//
// Solidity: event OperatorDeactivated(uint256 indexed operatorId)
func (_ServiceManager *ServiceManagerFilterer) WatchOperatorDeactivated(opts *bind.WatchOpts, sink chan<- *ServiceManagerOperatorDeactivated, operatorId []*big.Int) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ServiceManager.contract.WatchLogs(opts, "OperatorDeactivated", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceManagerOperatorDeactivated)
				if err := _ServiceManager.contract.UnpackLog(event, "OperatorDeactivated", log); err != nil {
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

// ParseOperatorDeactivated is a log parse operation binding the contract event 0x6ff0b7b14b65a91172fa7e0a7b49d909add202dadcef1b08c80f3c136a914fa8.
//
// Solidity: event OperatorDeactivated(uint256 indexed operatorId)
func (_ServiceManager *ServiceManagerFilterer) ParseOperatorDeactivated(log types.Log) (*ServiceManagerOperatorDeactivated, error) {
	event := new(ServiceManagerOperatorDeactivated)
	if err := _ServiceManager.contract.UnpackLog(event, "OperatorDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ServiceManagerOperatorPaymentAddressSetIterator is returned from FilterOperatorPaymentAddressSet and is used to iterate over the raw logs and unpacked data for OperatorPaymentAddressSet events raised by the ServiceManager contract.
type ServiceManagerOperatorPaymentAddressSetIterator struct {
	Event *ServiceManagerOperatorPaymentAddressSet // Event containing the contract specifics and raw log

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
func (it *ServiceManagerOperatorPaymentAddressSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceManagerOperatorPaymentAddressSet)
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
		it.Event = new(ServiceManagerOperatorPaymentAddressSet)
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
func (it *ServiceManagerOperatorPaymentAddressSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceManagerOperatorPaymentAddressSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceManagerOperatorPaymentAddressSet represents a OperatorPaymentAddressSet event raised by the ServiceManager contract.
type ServiceManagerOperatorPaymentAddressSet struct {
	OperatorId     *big.Int
	PaymentAddress common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterOperatorPaymentAddressSet is a free log retrieval operation binding the contract event 0xe6cb5a1aae8c82a1d14a9911ce3192769335cf5eead87c345f3d2bd7305acc3a.
//
// Solidity: event OperatorPaymentAddressSet(uint256 indexed operatorId, address paymentAddress)
func (_ServiceManager *ServiceManagerFilterer) FilterOperatorPaymentAddressSet(opts *bind.FilterOpts, operatorId []*big.Int) (*ServiceManagerOperatorPaymentAddressSetIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ServiceManager.contract.FilterLogs(opts, "OperatorPaymentAddressSet", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerOperatorPaymentAddressSetIterator{contract: _ServiceManager.contract, event: "OperatorPaymentAddressSet", logs: logs, sub: sub}, nil
}

// WatchOperatorPaymentAddressSet is a free log subscription operation binding the contract event 0xe6cb5a1aae8c82a1d14a9911ce3192769335cf5eead87c345f3d2bd7305acc3a.
//
// Solidity: event OperatorPaymentAddressSet(uint256 indexed operatorId, address paymentAddress)
func (_ServiceManager *ServiceManagerFilterer) WatchOperatorPaymentAddressSet(opts *bind.WatchOpts, sink chan<- *ServiceManagerOperatorPaymentAddressSet, operatorId []*big.Int) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ServiceManager.contract.WatchLogs(opts, "OperatorPaymentAddressSet", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceManagerOperatorPaymentAddressSet)
				if err := _ServiceManager.contract.UnpackLog(event, "OperatorPaymentAddressSet", log); err != nil {
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

// ParseOperatorPaymentAddressSet is a log parse operation binding the contract event 0xe6cb5a1aae8c82a1d14a9911ce3192769335cf5eead87c345f3d2bd7305acc3a.
//
// Solidity: event OperatorPaymentAddressSet(uint256 indexed operatorId, address paymentAddress)
func (_ServiceManager *ServiceManagerFilterer) ParseOperatorPaymentAddressSet(log types.Log) (*ServiceManagerOperatorPaymentAddressSet, error) {
	event := new(ServiceManagerOperatorPaymentAddressSet)
	if err := _ServiceManager.contract.UnpackLog(event, "OperatorPaymentAddressSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ServiceManagerOperatorUpdatedIterator is returned from FilterOperatorUpdated and is used to iterate over the raw logs and unpacked data for OperatorUpdated events raised by the ServiceManager contract.
type ServiceManagerOperatorUpdatedIterator struct {
	Event *ServiceManagerOperatorUpdated // Event containing the contract specifics and raw log

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
func (it *ServiceManagerOperatorUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceManagerOperatorUpdated)
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
		it.Event = new(ServiceManagerOperatorUpdated)
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
func (it *ServiceManagerOperatorUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceManagerOperatorUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceManagerOperatorUpdated represents a OperatorUpdated event raised by the ServiceManager contract.
type ServiceManagerOperatorUpdated struct {
	OperatorId *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOperatorUpdated is a free log retrieval operation binding the contract event 0xb0207706fc9ff2640909bbc2cd4ba7252c482b276bd367624c5771895de35e5f.
//
// Solidity: event OperatorUpdated(uint256 indexed operatorId)
func (_ServiceManager *ServiceManagerFilterer) FilterOperatorUpdated(opts *bind.FilterOpts, operatorId []*big.Int) (*ServiceManagerOperatorUpdatedIterator, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ServiceManager.contract.FilterLogs(opts, "OperatorUpdated", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerOperatorUpdatedIterator{contract: _ServiceManager.contract, event: "OperatorUpdated", logs: logs, sub: sub}, nil
}

// WatchOperatorUpdated is a free log subscription operation binding the contract event 0xb0207706fc9ff2640909bbc2cd4ba7252c482b276bd367624c5771895de35e5f.
//
// Solidity: event OperatorUpdated(uint256 indexed operatorId)
func (_ServiceManager *ServiceManagerFilterer) WatchOperatorUpdated(opts *bind.WatchOpts, sink chan<- *ServiceManagerOperatorUpdated, operatorId []*big.Int) (event.Subscription, error) {

	var operatorIdRule []interface{}
	for _, operatorIdItem := range operatorId {
		operatorIdRule = append(operatorIdRule, operatorIdItem)
	}

	logs, sub, err := _ServiceManager.contract.WatchLogs(opts, "OperatorUpdated", operatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceManagerOperatorUpdated)
				if err := _ServiceManager.contract.UnpackLog(event, "OperatorUpdated", log); err != nil {
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

// ParseOperatorUpdated is a log parse operation binding the contract event 0xb0207706fc9ff2640909bbc2cd4ba7252c482b276bd367624c5771895de35e5f.
//
// Solidity: event OperatorUpdated(uint256 indexed operatorId)
func (_ServiceManager *ServiceManagerFilterer) ParseOperatorUpdated(log types.Log) (*ServiceManagerOperatorUpdated, error) {
	event := new(ServiceManagerOperatorUpdated)
	if err := _ServiceManager.contract.UnpackLog(event, "OperatorUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ServiceManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the ServiceManager contract.
type ServiceManagerOwnershipTransferredIterator struct {
	Event *ServiceManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *ServiceManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceManagerOwnershipTransferred)
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
		it.Event = new(ServiceManagerOwnershipTransferred)
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
func (it *ServiceManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceManagerOwnershipTransferred represents a OwnershipTransferred event raised by the ServiceManager contract.
type ServiceManagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ServiceManager *ServiceManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*ServiceManagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ServiceManager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerOwnershipTransferredIterator{contract: _ServiceManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_ServiceManager *ServiceManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *ServiceManagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _ServiceManager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceManagerOwnershipTransferred)
				if err := _ServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_ServiceManager *ServiceManagerFilterer) ParseOwnershipTransferred(log types.Log) (*ServiceManagerOwnershipTransferred, error) {
	event := new(ServiceManagerOwnershipTransferred)
	if err := _ServiceManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ServiceManagerUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the ServiceManager contract.
type ServiceManagerUpgradedIterator struct {
	Event *ServiceManagerUpgraded // Event containing the contract specifics and raw log

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
func (it *ServiceManagerUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ServiceManagerUpgraded)
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
		it.Event = new(ServiceManagerUpgraded)
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
func (it *ServiceManagerUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ServiceManagerUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ServiceManagerUpgraded represents a Upgraded event raised by the ServiceManager contract.
type ServiceManagerUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ServiceManager *ServiceManagerFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*ServiceManagerUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ServiceManager.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &ServiceManagerUpgradedIterator{contract: _ServiceManager.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_ServiceManager *ServiceManagerFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *ServiceManagerUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _ServiceManager.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ServiceManagerUpgraded)
				if err := _ServiceManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_ServiceManager *ServiceManagerFilterer) ParseUpgraded(log types.Log) (*ServiceManagerUpgraded, error) {
	event := new(ServiceManagerUpgraded)
	if err := _ServiceManager.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
