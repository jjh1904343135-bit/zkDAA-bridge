// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

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

// DataMigrationMetaData contains all meta data concerning the DataMigration contract.
var DataMigrationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_verifierAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"h\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dataId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeout\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"h\",\"type\":\"bytes32\"}],\"name\":\"Reclaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"h\",\"type\":\"bytes32\"}],\"name\":\"Unlocked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"activeLocks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"dataId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timeout\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_h\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_dataId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_timeoutDuration\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_h\",\"type\":\"bytes32\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicInputs\",\"type\":\"uint256[2]\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifier\",\"outputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DataMigrationABI is the input ABI used to generate the binding from.
// Deprecated: Use DataMigrationMetaData.ABI instead.
var DataMigrationABI = DataMigrationMetaData.ABI

// DataMigration is an auto generated Go binding around an Ethereum contract.
type DataMigration struct {
	DataMigrationCaller     // Read-only binding to the contract
	DataMigrationTransactor // Write-only binding to the contract
	DataMigrationFilterer   // Log filterer for contract events
}

// DataMigrationCaller is an auto generated read-only Go binding around an Ethereum contract.
type DataMigrationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataMigrationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DataMigrationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataMigrationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DataMigrationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataMigrationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DataMigrationSession struct {
	Contract     *DataMigration    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DataMigrationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DataMigrationCallerSession struct {
	Contract *DataMigrationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// DataMigrationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DataMigrationTransactorSession struct {
	Contract     *DataMigrationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// DataMigrationRaw is an auto generated low-level Go binding around an Ethereum contract.
type DataMigrationRaw struct {
	Contract *DataMigration // Generic contract binding to access the raw methods on
}

// DataMigrationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DataMigrationCallerRaw struct {
	Contract *DataMigrationCaller // Generic read-only contract binding to access the raw methods on
}

// DataMigrationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DataMigrationTransactorRaw struct {
	Contract *DataMigrationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDataMigration creates a new instance of DataMigration, bound to a specific deployed contract.
func NewDataMigration(address common.Address, backend bind.ContractBackend) (*DataMigration, error) {
	contract, err := bindDataMigration(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DataMigration{DataMigrationCaller: DataMigrationCaller{contract: contract}, DataMigrationTransactor: DataMigrationTransactor{contract: contract}, DataMigrationFilterer: DataMigrationFilterer{contract: contract}}, nil
}

// NewDataMigrationCaller creates a new read-only instance of DataMigration, bound to a specific deployed contract.
func NewDataMigrationCaller(address common.Address, caller bind.ContractCaller) (*DataMigrationCaller, error) {
	contract, err := bindDataMigration(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataMigrationCaller{contract: contract}, nil
}

// NewDataMigrationTransactor creates a new write-only instance of DataMigration, bound to a specific deployed contract.
func NewDataMigrationTransactor(address common.Address, transactor bind.ContractTransactor) (*DataMigrationTransactor, error) {
	contract, err := bindDataMigration(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataMigrationTransactor{contract: contract}, nil
}

// NewDataMigrationFilterer creates a new log filterer instance of DataMigration, bound to a specific deployed contract.
func NewDataMigrationFilterer(address common.Address, filterer bind.ContractFilterer) (*DataMigrationFilterer, error) {
	contract, err := bindDataMigration(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataMigrationFilterer{contract: contract}, nil
}

// bindDataMigration binds a generic wrapper to an already deployed contract.
func bindDataMigration(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DataMigrationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataMigration *DataMigrationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataMigration.Contract.DataMigrationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataMigration *DataMigrationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataMigration.Contract.DataMigrationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataMigration *DataMigrationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataMigration.Contract.DataMigrationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataMigration *DataMigrationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataMigration.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataMigration *DataMigrationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataMigration.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataMigration *DataMigrationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataMigration.Contract.contract.Transact(opts, method, params...)
}

// ActiveLocks is a free data retrieval call binding the contract method 0xb51a201f.
//
// Solidity: function activeLocks(bytes32 ) view returns(address locker, bytes32 dataId, uint256 timeout)
func (_DataMigration *DataMigrationCaller) ActiveLocks(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Locker  common.Address
	DataId  [32]byte
	Timeout *big.Int
}, error) {
	var out []interface{}
	err := _DataMigration.contract.Call(opts, &out, "activeLocks", arg0)

	outstruct := new(struct {
		Locker  common.Address
		DataId  [32]byte
		Timeout *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Locker = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.DataId = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Timeout = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ActiveLocks is a free data retrieval call binding the contract method 0xb51a201f.
//
// Solidity: function activeLocks(bytes32 ) view returns(address locker, bytes32 dataId, uint256 timeout)
func (_DataMigration *DataMigrationSession) ActiveLocks(arg0 [32]byte) (struct {
	Locker  common.Address
	DataId  [32]byte
	Timeout *big.Int
}, error) {
	return _DataMigration.Contract.ActiveLocks(&_DataMigration.CallOpts, arg0)
}

// ActiveLocks is a free data retrieval call binding the contract method 0xb51a201f.
//
// Solidity: function activeLocks(bytes32 ) view returns(address locker, bytes32 dataId, uint256 timeout)
func (_DataMigration *DataMigrationCallerSession) ActiveLocks(arg0 [32]byte) (struct {
	Locker  common.Address
	DataId  [32]byte
	Timeout *big.Int
}, error) {
	return _DataMigration.Contract.ActiveLocks(&_DataMigration.CallOpts, arg0)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_DataMigration *DataMigrationCaller) Verifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataMigration.contract.Call(opts, &out, "verifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_DataMigration *DataMigrationSession) Verifier() (common.Address, error) {
	return _DataMigration.Contract.Verifier(&_DataMigration.CallOpts)
}

// Verifier is a free data retrieval call binding the contract method 0x2b7ac3f3.
//
// Solidity: function verifier() view returns(address)
func (_DataMigration *DataMigrationCallerSession) Verifier() (common.Address, error) {
	return _DataMigration.Contract.Verifier(&_DataMigration.CallOpts)
}

// Lock is a paid mutator transaction binding the contract method 0x47f2f178.
//
// Solidity: function lock(bytes32 _h, bytes32 _dataId, uint256 _timeoutDuration) returns()
func (_DataMigration *DataMigrationTransactor) Lock(opts *bind.TransactOpts, _h [32]byte, _dataId [32]byte, _timeoutDuration *big.Int) (*types.Transaction, error) {
	return _DataMigration.contract.Transact(opts, "lock", _h, _dataId, _timeoutDuration)
}

// Lock is a paid mutator transaction binding the contract method 0x47f2f178.
//
// Solidity: function lock(bytes32 _h, bytes32 _dataId, uint256 _timeoutDuration) returns()
func (_DataMigration *DataMigrationSession) Lock(_h [32]byte, _dataId [32]byte, _timeoutDuration *big.Int) (*types.Transaction, error) {
	return _DataMigration.Contract.Lock(&_DataMigration.TransactOpts, _h, _dataId, _timeoutDuration)
}

// Lock is a paid mutator transaction binding the contract method 0x47f2f178.
//
// Solidity: function lock(bytes32 _h, bytes32 _dataId, uint256 _timeoutDuration) returns()
func (_DataMigration *DataMigrationTransactorSession) Lock(_h [32]byte, _dataId [32]byte, _timeoutDuration *big.Int) (*types.Transaction, error) {
	return _DataMigration.Contract.Lock(&_DataMigration.TransactOpts, _h, _dataId, _timeoutDuration)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(bytes32 _h) returns()
func (_DataMigration *DataMigrationTransactor) Reclaim(opts *bind.TransactOpts, _h [32]byte) (*types.Transaction, error) {
	return _DataMigration.contract.Transact(opts, "reclaim", _h)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(bytes32 _h) returns()
func (_DataMigration *DataMigrationSession) Reclaim(_h [32]byte) (*types.Transaction, error) {
	return _DataMigration.Contract.Reclaim(&_DataMigration.TransactOpts, _h)
}

// Reclaim is a paid mutator transaction binding the contract method 0x96afb365.
//
// Solidity: function reclaim(bytes32 _h) returns()
func (_DataMigration *DataMigrationTransactorSession) Reclaim(_h [32]byte) (*types.Transaction, error) {
	return _DataMigration.Contract.Reclaim(&_DataMigration.TransactOpts, _h)
}

// Unlock is a paid mutator transaction binding the contract method 0x21915f28.
//
// Solidity: function unlock(uint256[8] proof, uint256[2] publicInputs) returns()
func (_DataMigration *DataMigrationTransactor) Unlock(opts *bind.TransactOpts, proof [8]*big.Int, publicInputs [2]*big.Int) (*types.Transaction, error) {
	return _DataMigration.contract.Transact(opts, "unlock", proof, publicInputs)
}

// Unlock is a paid mutator transaction binding the contract method 0x21915f28.
//
// Solidity: function unlock(uint256[8] proof, uint256[2] publicInputs) returns()
func (_DataMigration *DataMigrationSession) Unlock(proof [8]*big.Int, publicInputs [2]*big.Int) (*types.Transaction, error) {
	return _DataMigration.Contract.Unlock(&_DataMigration.TransactOpts, proof, publicInputs)
}

// Unlock is a paid mutator transaction binding the contract method 0x21915f28.
//
// Solidity: function unlock(uint256[8] proof, uint256[2] publicInputs) returns()
func (_DataMigration *DataMigrationTransactorSession) Unlock(proof [8]*big.Int, publicInputs [2]*big.Int) (*types.Transaction, error) {
	return _DataMigration.Contract.Unlock(&_DataMigration.TransactOpts, proof, publicInputs)
}

// DataMigrationLockedIterator is returned from FilterLocked and is used to iterate over the raw logs and unpacked data for Locked events raised by the DataMigration contract.
type DataMigrationLockedIterator struct {
	Event *DataMigrationLocked // Event containing the contract specifics and raw log

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
func (it *DataMigrationLockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataMigrationLocked)
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
		it.Event = new(DataMigrationLocked)
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
func (it *DataMigrationLockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataMigrationLockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataMigrationLocked represents a Locked event raised by the DataMigration contract.
type DataMigrationLocked struct {
	H       [32]byte
	Locker  common.Address
	DataId  [32]byte
	Timeout *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLocked is a free log retrieval operation binding the contract event 0xafec6b642358ef93f88d9a518cde37080f3ff07afee5022bd614276d0ec7790f.
//
// Solidity: event Locked(bytes32 indexed h, address indexed locker, bytes32 indexed dataId, uint256 timeout)
func (_DataMigration *DataMigrationFilterer) FilterLocked(opts *bind.FilterOpts, h [][32]byte, locker []common.Address, dataId [][32]byte) (*DataMigrationLockedIterator, error) {

	var hRule []interface{}
	for _, hItem := range h {
		hRule = append(hRule, hItem)
	}
	var lockerRule []interface{}
	for _, lockerItem := range locker {
		lockerRule = append(lockerRule, lockerItem)
	}
	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataMigration.contract.FilterLogs(opts, "Locked", hRule, lockerRule, dataIdRule)
	if err != nil {
		return nil, err
	}
	return &DataMigrationLockedIterator{contract: _DataMigration.contract, event: "Locked", logs: logs, sub: sub}, nil
}

// WatchLocked is a free log subscription operation binding the contract event 0xafec6b642358ef93f88d9a518cde37080f3ff07afee5022bd614276d0ec7790f.
//
// Solidity: event Locked(bytes32 indexed h, address indexed locker, bytes32 indexed dataId, uint256 timeout)
func (_DataMigration *DataMigrationFilterer) WatchLocked(opts *bind.WatchOpts, sink chan<- *DataMigrationLocked, h [][32]byte, locker []common.Address, dataId [][32]byte) (event.Subscription, error) {

	var hRule []interface{}
	for _, hItem := range h {
		hRule = append(hRule, hItem)
	}
	var lockerRule []interface{}
	for _, lockerItem := range locker {
		lockerRule = append(lockerRule, lockerItem)
	}
	var dataIdRule []interface{}
	for _, dataIdItem := range dataId {
		dataIdRule = append(dataIdRule, dataIdItem)
	}

	logs, sub, err := _DataMigration.contract.WatchLogs(opts, "Locked", hRule, lockerRule, dataIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataMigrationLocked)
				if err := _DataMigration.contract.UnpackLog(event, "Locked", log); err != nil {
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

// ParseLocked is a log parse operation binding the contract event 0xafec6b642358ef93f88d9a518cde37080f3ff07afee5022bd614276d0ec7790f.
//
// Solidity: event Locked(bytes32 indexed h, address indexed locker, bytes32 indexed dataId, uint256 timeout)
func (_DataMigration *DataMigrationFilterer) ParseLocked(log types.Log) (*DataMigrationLocked, error) {
	event := new(DataMigrationLocked)
	if err := _DataMigration.contract.UnpackLog(event, "Locked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DataMigrationReclaimedIterator is returned from FilterReclaimed and is used to iterate over the raw logs and unpacked data for Reclaimed events raised by the DataMigration contract.
type DataMigrationReclaimedIterator struct {
	Event *DataMigrationReclaimed // Event containing the contract specifics and raw log

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
func (it *DataMigrationReclaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataMigrationReclaimed)
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
		it.Event = new(DataMigrationReclaimed)
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
func (it *DataMigrationReclaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataMigrationReclaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataMigrationReclaimed represents a Reclaimed event raised by the DataMigration contract.
type DataMigrationReclaimed struct {
	H   [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterReclaimed is a free log retrieval operation binding the contract event 0xbe9e485e7f7ace1eaf2897ca5483cdb8bf05d65d8b660c18070acc7596529446.
//
// Solidity: event Reclaimed(bytes32 indexed h)
func (_DataMigration *DataMigrationFilterer) FilterReclaimed(opts *bind.FilterOpts, h [][32]byte) (*DataMigrationReclaimedIterator, error) {

	var hRule []interface{}
	for _, hItem := range h {
		hRule = append(hRule, hItem)
	}

	logs, sub, err := _DataMigration.contract.FilterLogs(opts, "Reclaimed", hRule)
	if err != nil {
		return nil, err
	}
	return &DataMigrationReclaimedIterator{contract: _DataMigration.contract, event: "Reclaimed", logs: logs, sub: sub}, nil
}

// WatchReclaimed is a free log subscription operation binding the contract event 0xbe9e485e7f7ace1eaf2897ca5483cdb8bf05d65d8b660c18070acc7596529446.
//
// Solidity: event Reclaimed(bytes32 indexed h)
func (_DataMigration *DataMigrationFilterer) WatchReclaimed(opts *bind.WatchOpts, sink chan<- *DataMigrationReclaimed, h [][32]byte) (event.Subscription, error) {

	var hRule []interface{}
	for _, hItem := range h {
		hRule = append(hRule, hItem)
	}

	logs, sub, err := _DataMigration.contract.WatchLogs(opts, "Reclaimed", hRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataMigrationReclaimed)
				if err := _DataMigration.contract.UnpackLog(event, "Reclaimed", log); err != nil {
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

// ParseReclaimed is a log parse operation binding the contract event 0xbe9e485e7f7ace1eaf2897ca5483cdb8bf05d65d8b660c18070acc7596529446.
//
// Solidity: event Reclaimed(bytes32 indexed h)
func (_DataMigration *DataMigrationFilterer) ParseReclaimed(log types.Log) (*DataMigrationReclaimed, error) {
	event := new(DataMigrationReclaimed)
	if err := _DataMigration.contract.UnpackLog(event, "Reclaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DataMigrationUnlockedIterator is returned from FilterUnlocked and is used to iterate over the raw logs and unpacked data for Unlocked events raised by the DataMigration contract.
type DataMigrationUnlockedIterator struct {
	Event *DataMigrationUnlocked // Event containing the contract specifics and raw log

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
func (it *DataMigrationUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataMigrationUnlocked)
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
		it.Event = new(DataMigrationUnlocked)
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
func (it *DataMigrationUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataMigrationUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataMigrationUnlocked represents a Unlocked event raised by the DataMigration contract.
type DataMigrationUnlocked struct {
	H   [32]byte
	Raw types.Log // Blockchain specific contextual infos
}

// FilterUnlocked is a free log retrieval operation binding the contract event 0x248ed1ed5e6e28246432a42651d384fadb929f7cefaacbe6a1764513e058ee74.
//
// Solidity: event Unlocked(bytes32 indexed h)
func (_DataMigration *DataMigrationFilterer) FilterUnlocked(opts *bind.FilterOpts, h [][32]byte) (*DataMigrationUnlockedIterator, error) {

	var hRule []interface{}
	for _, hItem := range h {
		hRule = append(hRule, hItem)
	}

	logs, sub, err := _DataMigration.contract.FilterLogs(opts, "Unlocked", hRule)
	if err != nil {
		return nil, err
	}
	return &DataMigrationUnlockedIterator{contract: _DataMigration.contract, event: "Unlocked", logs: logs, sub: sub}, nil
}

// WatchUnlocked is a free log subscription operation binding the contract event 0x248ed1ed5e6e28246432a42651d384fadb929f7cefaacbe6a1764513e058ee74.
//
// Solidity: event Unlocked(bytes32 indexed h)
func (_DataMigration *DataMigrationFilterer) WatchUnlocked(opts *bind.WatchOpts, sink chan<- *DataMigrationUnlocked, h [][32]byte) (event.Subscription, error) {

	var hRule []interface{}
	for _, hItem := range h {
		hRule = append(hRule, hItem)
	}

	logs, sub, err := _DataMigration.contract.WatchLogs(opts, "Unlocked", hRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataMigrationUnlocked)
				if err := _DataMigration.contract.UnpackLog(event, "Unlocked", log); err != nil {
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

// ParseUnlocked is a log parse operation binding the contract event 0x248ed1ed5e6e28246432a42651d384fadb929f7cefaacbe6a1764513e058ee74.
//
// Solidity: event Unlocked(bytes32 indexed h)
func (_DataMigration *DataMigrationFilterer) ParseUnlocked(log types.Log) (*DataMigrationUnlocked, error) {
	event := new(DataMigrationUnlocked)
	if err := _DataMigration.contract.UnpackLog(event, "Unlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
