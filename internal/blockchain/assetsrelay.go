// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package blockchain

import (
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
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AssetsRelayABI is the input ABI used to generate the binding from.
const AssetsRelayABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"hash\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"int256\"}],\"name\":\"sign\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"publicKey\",\"type\":\"address\"}],\"name\":\"getPublisherByAddress\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"publicKeys\",\"type\":\"address[]\"}],\"name\":\"disablePublisher\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPublishers\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"aContract\",\"type\":\"address\"}],\"name\":\"setContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getAssetCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"verify\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"int256\"},{\"name\":\"\",\"type\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"string\"}],\"name\":\"getAssetCountForHash\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"publicKeys\",\"type\":\"address[]\"},{\"name\":\"level\",\"type\":\"int256\"}],\"name\":\"setPublisherLevel\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"createdAt\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"string\"},{\"name\":\"assetIndex\",\"type\":\"uint256\"}],\"name\":\"verifyByIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"int256\"},{\"name\":\"\",\"type\":\"int256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"assetsContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hashIndex\",\"type\":\"uint256\"}],\"name\":\"getHashByIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"aContract\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// AssetsRelay is an auto generated Go binding around an Ethereum contract.
type AssetsRelay struct {
	AssetsRelayCaller     // Read-only binding to the contract
	AssetsRelayTransactor // Write-only binding to the contract
	AssetsRelayFilterer   // Log filterer for contract events
}

// AssetsRelayCaller is an auto generated read-only Go binding around an Ethereum contract.
type AssetsRelayCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsRelayTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AssetsRelayTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsRelayFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AssetsRelayFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AssetsRelaySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AssetsRelaySession struct {
	Contract     *AssetsRelay      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AssetsRelayCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AssetsRelayCallerSession struct {
	Contract *AssetsRelayCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// AssetsRelayTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AssetsRelayTransactorSession struct {
	Contract     *AssetsRelayTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AssetsRelayRaw is an auto generated low-level Go binding around an Ethereum contract.
type AssetsRelayRaw struct {
	Contract *AssetsRelay // Generic contract binding to access the raw methods on
}

// AssetsRelayCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AssetsRelayCallerRaw struct {
	Contract *AssetsRelayCaller // Generic read-only contract binding to access the raw methods on
}

// AssetsRelayTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AssetsRelayTransactorRaw struct {
	Contract *AssetsRelayTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAssetsRelay creates a new instance of AssetsRelay, bound to a specific deployed contract.
func NewAssetsRelay(address common.Address, backend bind.ContractBackend) (*AssetsRelay, error) {
	contract, err := bindAssetsRelay(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AssetsRelay{AssetsRelayCaller: AssetsRelayCaller{contract: contract}, AssetsRelayTransactor: AssetsRelayTransactor{contract: contract}, AssetsRelayFilterer: AssetsRelayFilterer{contract: contract}}, nil
}

// NewAssetsRelayCaller creates a new read-only instance of AssetsRelay, bound to a specific deployed contract.
func NewAssetsRelayCaller(address common.Address, caller bind.ContractCaller) (*AssetsRelayCaller, error) {
	contract, err := bindAssetsRelay(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AssetsRelayCaller{contract: contract}, nil
}

// NewAssetsRelayTransactor creates a new write-only instance of AssetsRelay, bound to a specific deployed contract.
func NewAssetsRelayTransactor(address common.Address, transactor bind.ContractTransactor) (*AssetsRelayTransactor, error) {
	contract, err := bindAssetsRelay(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AssetsRelayTransactor{contract: contract}, nil
}

// NewAssetsRelayFilterer creates a new log filterer instance of AssetsRelay, bound to a specific deployed contract.
func NewAssetsRelayFilterer(address common.Address, filterer bind.ContractFilterer) (*AssetsRelayFilterer, error) {
	contract, err := bindAssetsRelay(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AssetsRelayFilterer{contract: contract}, nil
}

// bindAssetsRelay binds a generic wrapper to an already deployed contract.
func bindAssetsRelay(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AssetsRelayABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AssetsRelay *AssetsRelayRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AssetsRelay.Contract.AssetsRelayCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AssetsRelay *AssetsRelayRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AssetsRelay.Contract.AssetsRelayTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AssetsRelay *AssetsRelayRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AssetsRelay.Contract.AssetsRelayTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AssetsRelay *AssetsRelayCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AssetsRelay.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AssetsRelay *AssetsRelayTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AssetsRelay.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AssetsRelay *AssetsRelayTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AssetsRelay.Contract.contract.Transact(opts, method, params...)
}

// AssetsContract is a free data retrieval call binding the contract method 0xddfe5b2d.
//
// Solidity: function assetsContract() constant returns(address)
func (_AssetsRelay *AssetsRelayCaller) AssetsContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AssetsRelay.contract.Call(opts, out, "assetsContract")
	return *ret0, err
}

// AssetsContract is a free data retrieval call binding the contract method 0xddfe5b2d.
//
// Solidity: function assetsContract() constant returns(address)
func (_AssetsRelay *AssetsRelaySession) AssetsContract() (common.Address, error) {
	return _AssetsRelay.Contract.AssetsContract(&_AssetsRelay.CallOpts)
}

// AssetsContract is a free data retrieval call binding the contract method 0xddfe5b2d.
//
// Solidity: function assetsContract() constant returns(address)
func (_AssetsRelay *AssetsRelayCallerSession) AssetsContract() (common.Address, error) {
	return _AssetsRelay.Contract.AssetsContract(&_AssetsRelay.CallOpts)
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() constant returns(uint256)
func (_AssetsRelay *AssetsRelayCaller) CreatedAt(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AssetsRelay.contract.Call(opts, out, "createdAt")
	return *ret0, err
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() constant returns(uint256)
func (_AssetsRelay *AssetsRelaySession) CreatedAt() (*big.Int, error) {
	return _AssetsRelay.Contract.CreatedAt(&_AssetsRelay.CallOpts)
}

// CreatedAt is a free data retrieval call binding the contract method 0xcf09e0d0.
//
// Solidity: function createdAt() constant returns(uint256)
func (_AssetsRelay *AssetsRelayCallerSession) CreatedAt() (*big.Int, error) {
	return _AssetsRelay.Contract.CreatedAt(&_AssetsRelay.CallOpts)
}

// GetAssetCount is a free data retrieval call binding the contract method 0xa0aead4d.
//
// Solidity: function getAssetCount() constant returns(uint256)
func (_AssetsRelay *AssetsRelayCaller) GetAssetCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AssetsRelay.contract.Call(opts, out, "getAssetCount")
	return *ret0, err
}

// GetAssetCount is a free data retrieval call binding the contract method 0xa0aead4d.
//
// Solidity: function getAssetCount() constant returns(uint256)
func (_AssetsRelay *AssetsRelaySession) GetAssetCount() (*big.Int, error) {
	return _AssetsRelay.Contract.GetAssetCount(&_AssetsRelay.CallOpts)
}

// GetAssetCount is a free data retrieval call binding the contract method 0xa0aead4d.
//
// Solidity: function getAssetCount() constant returns(uint256)
func (_AssetsRelay *AssetsRelayCallerSession) GetAssetCount() (*big.Int, error) {
	return _AssetsRelay.Contract.GetAssetCount(&_AssetsRelay.CallOpts)
}

// GetAssetCountForHash is a free data retrieval call binding the contract method 0xbf91e82f.
//
// Solidity: function getAssetCountForHash(string hash) constant returns(uint256)
func (_AssetsRelay *AssetsRelayCaller) GetAssetCountForHash(opts *bind.CallOpts, hash string) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AssetsRelay.contract.Call(opts, out, "getAssetCountForHash", hash)
	return *ret0, err
}

// GetAssetCountForHash is a free data retrieval call binding the contract method 0xbf91e82f.
//
// Solidity: function getAssetCountForHash(string hash) constant returns(uint256)
func (_AssetsRelay *AssetsRelaySession) GetAssetCountForHash(hash string) (*big.Int, error) {
	return _AssetsRelay.Contract.GetAssetCountForHash(&_AssetsRelay.CallOpts, hash)
}

// GetAssetCountForHash is a free data retrieval call binding the contract method 0xbf91e82f.
//
// Solidity: function getAssetCountForHash(string hash) constant returns(uint256)
func (_AssetsRelay *AssetsRelayCallerSession) GetAssetCountForHash(hash string) (*big.Int, error) {
	return _AssetsRelay.Contract.GetAssetCountForHash(&_AssetsRelay.CallOpts, hash)
}

// GetHashByIndex is a free data retrieval call binding the contract method 0xec058186.
//
// Solidity: function getHashByIndex(uint256 hashIndex) constant returns(string)
func (_AssetsRelay *AssetsRelayCaller) GetHashByIndex(opts *bind.CallOpts, hashIndex *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AssetsRelay.contract.Call(opts, out, "getHashByIndex", hashIndex)
	return *ret0, err
}

// GetHashByIndex is a free data retrieval call binding the contract method 0xec058186.
//
// Solidity: function getHashByIndex(uint256 hashIndex) constant returns(string)
func (_AssetsRelay *AssetsRelaySession) GetHashByIndex(hashIndex *big.Int) (string, error) {
	return _AssetsRelay.Contract.GetHashByIndex(&_AssetsRelay.CallOpts, hashIndex)
}

// GetHashByIndex is a free data retrieval call binding the contract method 0xec058186.
//
// Solidity: function getHashByIndex(uint256 hashIndex) constant returns(string)
func (_AssetsRelay *AssetsRelayCallerSession) GetHashByIndex(hashIndex *big.Int) (string, error) {
	return _AssetsRelay.Contract.GetHashByIndex(&_AssetsRelay.CallOpts, hashIndex)
}

// GetPublisherByAddress is a free data retrieval call binding the contract method 0x38052e27.
//
// Solidity: function getPublisherByAddress(address publicKey) constant returns(address, int256, uint256)
func (_AssetsRelay *AssetsRelayCaller) GetPublisherByAddress(opts *bind.CallOpts, publicKey common.Address) (common.Address, *big.Int, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
	}
	err := _AssetsRelay.contract.Call(opts, out, "getPublisherByAddress", publicKey)
	return *ret0, *ret1, *ret2, err
}

// GetPublisherByAddress is a free data retrieval call binding the contract method 0x38052e27.
//
// Solidity: function getPublisherByAddress(address publicKey) constant returns(address, int256, uint256)
func (_AssetsRelay *AssetsRelaySession) GetPublisherByAddress(publicKey common.Address) (common.Address, *big.Int, *big.Int, error) {
	return _AssetsRelay.Contract.GetPublisherByAddress(&_AssetsRelay.CallOpts, publicKey)
}

// GetPublisherByAddress is a free data retrieval call binding the contract method 0x38052e27.
//
// Solidity: function getPublisherByAddress(address publicKey) constant returns(address, int256, uint256)
func (_AssetsRelay *AssetsRelayCallerSession) GetPublisherByAddress(publicKey common.Address) (common.Address, *big.Int, *big.Int, error) {
	return _AssetsRelay.Contract.GetPublisherByAddress(&_AssetsRelay.CallOpts, publicKey)
}

// GetPublishers is a free data retrieval call binding the contract method 0x6c6071aa.
//
// Solidity: function getPublishers() constant returns(address[])
func (_AssetsRelay *AssetsRelayCaller) GetPublishers(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _AssetsRelay.contract.Call(opts, out, "getPublishers")
	return *ret0, err
}

// GetPublishers is a free data retrieval call binding the contract method 0x6c6071aa.
//
// Solidity: function getPublishers() constant returns(address[])
func (_AssetsRelay *AssetsRelaySession) GetPublishers() ([]common.Address, error) {
	return _AssetsRelay.Contract.GetPublishers(&_AssetsRelay.CallOpts)
}

// GetPublishers is a free data retrieval call binding the contract method 0x6c6071aa.
//
// Solidity: function getPublishers() constant returns(address[])
func (_AssetsRelay *AssetsRelayCallerSession) GetPublishers() ([]common.Address, error) {
	return _AssetsRelay.Contract.GetPublishers(&_AssetsRelay.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AssetsRelay *AssetsRelayCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AssetsRelay.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AssetsRelay *AssetsRelaySession) Owner() (common.Address, error) {
	return _AssetsRelay.Contract.Owner(&_AssetsRelay.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_AssetsRelay *AssetsRelayCallerSession) Owner() (common.Address, error) {
	return _AssetsRelay.Contract.Owner(&_AssetsRelay.CallOpts)
}

// Verify is a free data retrieval call binding the contract method 0xbb9c6c3e.
//
// Solidity: function verify(string hash) constant returns(address, int256, int256, uint256)
func (_AssetsRelay *AssetsRelayCaller) Verify(opts *bind.CallOpts, hash string) (common.Address, *big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
		ret3 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _AssetsRelay.contract.Call(opts, out, "verify", hash)
	return *ret0, *ret1, *ret2, *ret3, err
}

// Verify is a free data retrieval call binding the contract method 0xbb9c6c3e.
//
// Solidity: function verify(string hash) constant returns(address, int256, int256, uint256)
func (_AssetsRelay *AssetsRelaySession) Verify(hash string) (common.Address, *big.Int, *big.Int, *big.Int, error) {
	return _AssetsRelay.Contract.Verify(&_AssetsRelay.CallOpts, hash)
}

// Verify is a free data retrieval call binding the contract method 0xbb9c6c3e.
//
// Solidity: function verify(string hash) constant returns(address, int256, int256, uint256)
func (_AssetsRelay *AssetsRelayCallerSession) Verify(hash string) (common.Address, *big.Int, *big.Int, *big.Int, error) {
	return _AssetsRelay.Contract.Verify(&_AssetsRelay.CallOpts, hash)
}

// VerifyByIndex is a free data retrieval call binding the contract method 0xd6ce25a9.
//
// Solidity: function verifyByIndex(string hash, uint256 assetIndex) constant returns(address, int256, int256, uint256)
func (_AssetsRelay *AssetsRelayCaller) VerifyByIndex(opts *bind.CallOpts, hash string, assetIndex *big.Int) (common.Address, *big.Int, *big.Int, *big.Int, error) {
	var (
		ret0 = new(common.Address)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
		ret3 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _AssetsRelay.contract.Call(opts, out, "verifyByIndex", hash, assetIndex)
	return *ret0, *ret1, *ret2, *ret3, err
}

// VerifyByIndex is a free data retrieval call binding the contract method 0xd6ce25a9.
//
// Solidity: function verifyByIndex(string hash, uint256 assetIndex) constant returns(address, int256, int256, uint256)
func (_AssetsRelay *AssetsRelaySession) VerifyByIndex(hash string, assetIndex *big.Int) (common.Address, *big.Int, *big.Int, *big.Int, error) {
	return _AssetsRelay.Contract.VerifyByIndex(&_AssetsRelay.CallOpts, hash, assetIndex)
}

// VerifyByIndex is a free data retrieval call binding the contract method 0xd6ce25a9.
//
// Solidity: function verifyByIndex(string hash, uint256 assetIndex) constant returns(address, int256, int256, uint256)
func (_AssetsRelay *AssetsRelayCallerSession) VerifyByIndex(hash string, assetIndex *big.Int) (common.Address, *big.Int, *big.Int, *big.Int, error) {
	return _AssetsRelay.Contract.VerifyByIndex(&_AssetsRelay.CallOpts, hash, assetIndex)
}

// DisablePublisher is a paid mutator transaction binding the contract method 0x54df1eab.
//
// Solidity: function disablePublisher(address[] publicKeys) returns()
func (_AssetsRelay *AssetsRelayTransactor) DisablePublisher(opts *bind.TransactOpts, publicKeys []common.Address) (*types.Transaction, error) {
	return _AssetsRelay.contract.Transact(opts, "disablePublisher", publicKeys)
}

// DisablePublisher is a paid mutator transaction binding the contract method 0x54df1eab.
//
// Solidity: function disablePublisher(address[] publicKeys) returns()
func (_AssetsRelay *AssetsRelaySession) DisablePublisher(publicKeys []common.Address) (*types.Transaction, error) {
	return _AssetsRelay.Contract.DisablePublisher(&_AssetsRelay.TransactOpts, publicKeys)
}

// DisablePublisher is a paid mutator transaction binding the contract method 0x54df1eab.
//
// Solidity: function disablePublisher(address[] publicKeys) returns()
func (_AssetsRelay *AssetsRelayTransactorSession) DisablePublisher(publicKeys []common.Address) (*types.Transaction, error) {
	return _AssetsRelay.Contract.DisablePublisher(&_AssetsRelay.TransactOpts, publicKeys)
}

// SetContract is a paid mutator transaction binding the contract method 0x75f890ab.
//
// Solidity: function setContract(address aContract) returns()
func (_AssetsRelay *AssetsRelayTransactor) SetContract(opts *bind.TransactOpts, aContract common.Address) (*types.Transaction, error) {
	return _AssetsRelay.contract.Transact(opts, "setContract", aContract)
}

// SetContract is a paid mutator transaction binding the contract method 0x75f890ab.
//
// Solidity: function setContract(address aContract) returns()
func (_AssetsRelay *AssetsRelaySession) SetContract(aContract common.Address) (*types.Transaction, error) {
	return _AssetsRelay.Contract.SetContract(&_AssetsRelay.TransactOpts, aContract)
}

// SetContract is a paid mutator transaction binding the contract method 0x75f890ab.
//
// Solidity: function setContract(address aContract) returns()
func (_AssetsRelay *AssetsRelayTransactorSession) SetContract(aContract common.Address) (*types.Transaction, error) {
	return _AssetsRelay.Contract.SetContract(&_AssetsRelay.TransactOpts, aContract)
}

// SetPublisherLevel is a paid mutator transaction binding the contract method 0xc9e07f09.
//
// Solidity: function setPublisherLevel(address[] publicKeys, int256 level) returns()
func (_AssetsRelay *AssetsRelayTransactor) SetPublisherLevel(opts *bind.TransactOpts, publicKeys []common.Address, level *big.Int) (*types.Transaction, error) {
	return _AssetsRelay.contract.Transact(opts, "setPublisherLevel", publicKeys, level)
}

// SetPublisherLevel is a paid mutator transaction binding the contract method 0xc9e07f09.
//
// Solidity: function setPublisherLevel(address[] publicKeys, int256 level) returns()
func (_AssetsRelay *AssetsRelaySession) SetPublisherLevel(publicKeys []common.Address, level *big.Int) (*types.Transaction, error) {
	return _AssetsRelay.Contract.SetPublisherLevel(&_AssetsRelay.TransactOpts, publicKeys, level)
}

// SetPublisherLevel is a paid mutator transaction binding the contract method 0xc9e07f09.
//
// Solidity: function setPublisherLevel(address[] publicKeys, int256 level) returns()
func (_AssetsRelay *AssetsRelayTransactorSession) SetPublisherLevel(publicKeys []common.Address, level *big.Int) (*types.Transaction, error) {
	return _AssetsRelay.Contract.SetPublisherLevel(&_AssetsRelay.TransactOpts, publicKeys, level)
}

// Sign is a paid mutator transaction binding the contract method 0x28f7f9b8.
//
// Solidity: function sign(string hash, int256 status) returns()
func (_AssetsRelay *AssetsRelayTransactor) Sign(opts *bind.TransactOpts, hash string, status *big.Int) (*types.Transaction, error) {
	return _AssetsRelay.contract.Transact(opts, "sign", hash, status)
}

// Sign is a paid mutator transaction binding the contract method 0x28f7f9b8.
//
// Solidity: function sign(string hash, int256 status) returns()
func (_AssetsRelay *AssetsRelaySession) Sign(hash string, status *big.Int) (*types.Transaction, error) {
	return _AssetsRelay.Contract.Sign(&_AssetsRelay.TransactOpts, hash, status)
}

// Sign is a paid mutator transaction binding the contract method 0x28f7f9b8.
//
// Solidity: function sign(string hash, int256 status) returns()
func (_AssetsRelay *AssetsRelayTransactorSession) Sign(hash string, status *big.Int) (*types.Transaction, error) {
	return _AssetsRelay.Contract.Sign(&_AssetsRelay.TransactOpts, hash, status)
}
