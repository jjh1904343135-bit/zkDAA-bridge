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

// BatchDataMigrationMetaData contains all meta data concerning the BatchDataMigration contract.
var BatchDataMigrationMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"BatchRootSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"serialNumber\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"unlocker\",\"type\":\"address\"}],\"name\":\"Unlocked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"}],\"name\":\"VerifierSet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"batchRoots\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"batchTimestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"}],\"name\":\"getBatchRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"}],\"name\":\"getVerifier\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"verifier16\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"verifier64\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"verifier128\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"verifier256\",\"type\":\"address\"}],\"name\":\"setAllVerifiers\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"verifierAddress\",\"type\":\"address\"}],\"name\":\"setVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"submitBatchRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"batchSize\",\"type\":\"uint256\"},{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicInputs\",\"type\":\"uint256[2]\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"usedSerialNumbers\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"verifiers\",\"outputs\":[{\"internalType\":\"contractIVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b50600480546001600160a01b03191633179055610add806100316000396000f3fe608060405234801561001057600080fd5b50600436106100a95760003560e01c806374f375051161007157806374f375051461015d5780638da5cb5b14610190578063ac1eff68146101a3578063b4f7dc49146101cc578063e2350d63146101ec578063fb2cd9371461020c57600080fd5b80631957ba4e146100ae578063301e24a4146100c35780633561bc27146100d6578063437f5ae51461011c57806368e2f9791461012f575b600080fd5b6100c16100bc36600461091f565b61021f565b005b6100c16100d136600461094b565b610378565b6100ff6100e436600461096d565b6000908152602081905260409020546001600160a01b031690565b6040516001600160a01b0390911681526020015b60405180910390f35b6100c161012a366004610986565b610448565b61014f61013d36600461096d565b60016020526000908152604090205481565b604051908152602001610113565b61018061016b36600461096d565b60036020526000908152604090205460ff1681565b6040519015158152602001610113565b6004546100ff906001600160a01b031681565b6100ff6101b136600461096d565b6000602081905290815260409020546001600160a01b031681565b61014f6101da36600461096d565b60026020526000908152604090205481565b61014f6101fa36600461096d565b60009081526001602052604090205490565b6100c161021a3660046109d4565b6105e4565b6004546001600160a01b0316331461026b5760405162461bcd60e51b815260206004820152600a60248201526927b7363c9037bbb732b960b11b60448201526064015b60405180910390fd5b6001600160a01b0381166102c15760405162461bcd60e51b815260206004820152601f60248201527f566572696669657220616464726573732063616e6e6f74206265207a65726f006044820152606401610262565b81601014806102d05750816040145b806102db5750816080145b806102e7575081610100145b6103285760405162461bcd60e51b8152602060048201526012602482015271496e76616c69642062617463682073697a6560701b6044820152606401610262565b6000828152602081815260409182902080546001600160a01b0319166001600160a01b03851690811790915591519182528391600080516020610a88833981519152910160405180910390a25050565b806103bb5760405162461bcd60e51b8152602060048201526013602482015272526f6f742063616e6e6f74206265207a65726f60681b6044820152606401610262565b6000828152602081905260409020546001600160a01b03166103ef5760405162461bcd60e51b815260040161026290610a28565b6000828152600160209081526040808320849055600282529182902042908190559151918252829184917f332a878a15e4c00bf130d2a4dd82a0e2bee22ed9f5a1cdcf9c8d9e035461f22f910160405180910390a35050565b6000838152602081905260409020546001600160a01b03168061047d5760405162461bcd60e51b815260040161026290610a28565b6000848152600160209081526040909120548335918401359082146104da5760405162461bcd60e51b8152602060048201526013602482015272125b9d985b1a590813595c9adb19481c9bdbdd606a1b6044820152606401610262565b60008181526003602052604090205460ff16156105395760405162461bcd60e51b815260206004820152601a60248201527f53657269616c206e756d62657220616c726561647920757365640000000000006044820152606401610262565b604051635fe24f2360e01b81526001600160a01b03841690635fe24f23906105679088908890600401610a6c565b60006040518083038186803b15801561057f57600080fd5b505afa158015610593573d6000803e3d6000fd5b505050600082815260036020526040808220805460ff1916600117905551339250839189917fdf9bbb560a3c97f1125d5e2b956725434c1ba1b7bc6d1cf1119c83cc230d7a509190a4505050505050565b6004546001600160a01b0316331461062b5760405162461bcd60e51b815260206004820152600a60248201526927b7363c9037bbb732b960b11b6044820152606401610262565b6001600160a01b0384166106815760405162461bcd60e51b815260206004820152601960248201527f566572696669657231362063616e6e6f74206265207a65726f000000000000006044820152606401610262565b6001600160a01b0383166106d75760405162461bcd60e51b815260206004820152601960248201527f566572696669657236342063616e6e6f74206265207a65726f000000000000006044820152606401610262565b6001600160a01b03821661072d5760405162461bcd60e51b815260206004820152601a60248201527f56657269666965723132382063616e6e6f74206265207a65726f0000000000006044820152606401610262565b6001600160a01b0381166107835760405162461bcd60e51b815260206004820152601a60248201527f56657269666965723235362063616e6e6f74206265207a65726f0000000000006044820152606401610262565b600060208181527f020abee21eef15c21bc31a406c2b8ac3afc5df94a4b02b38abb286f4334e6c5b80546001600160a01b03199081166001600160a01b038981169182179093557f55a5bd2f1561a275dba97a3e23bf64b586a74cc4c4a2dc9f1c147ffdbbce1122805483168985161790557ff9cab2ea5fc2898ca670240c87501f552f715e74f50142c35920c55eb3c0ea1a805483168885161790556101009094527f2d82c899f9de3d0a1f6f906d748a47eb88f8764f6fcbf2548decd720a33c0f158054909116918516919091179055604051918252601091600080516020610a88833981519152910160405180910390a2604080516001600160a01b0385168152600080516020610a888339815191529060200160405180910390a26040516001600160a01b0383168152608090600080516020610a888339815191529060200160405180910390a26040516001600160a01b038216815261010090600080516020610a888339815191529060200160405180910390a250505050565b80356001600160a01b038116811461091a57600080fd5b919050565b6000806040838503121561093257600080fd5b8235915061094260208401610903565b90509250929050565b6000806040838503121561095e57600080fd5b50508035926020909101359150565b60006020828403121561097f57600080fd5b5035919050565b6000806000610160848603121561099c57600080fd5b833592506101208401858111156109b257600080fd5b60208501925085610160860111156109c957600080fd5b809150509250925092565b600080600080608085870312156109ea57600080fd5b6109f385610903565b9350610a0160208601610903565b9250610a0f60408601610903565b9150610a1d60608601610903565b905092959194509250565b60208082526024908201527f5665726966696572206e6f742073657420666f7220746869732062617463682060408201526373697a6560e01b606082015260800190565b6101408101610100848337604083610100840137939250505056febc291d0e6f60c8ebaeb52dc4380cd7fa1fa6ac795fc0719ea5f57b8acfa28e3ca26469706673582212207f3373a7dd423a820583a3ba8fc0d1c2ecbe525544513215881a518b4c563ded64736f6c634300081c0033",
}

// BatchDataMigrationABI is the input ABI used to generate the binding from.
// Deprecated: Use BatchDataMigrationMetaData.ABI instead.
var BatchDataMigrationABI = BatchDataMigrationMetaData.ABI

// BatchDataMigrationBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use BatchDataMigrationMetaData.Bin instead.
var BatchDataMigrationBin = BatchDataMigrationMetaData.Bin

// DeployBatchDataMigration deploys a new Ethereum contract, binding an instance of BatchDataMigration to it.
func DeployBatchDataMigration(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *BatchDataMigration, error) {
	parsed, err := BatchDataMigrationMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(BatchDataMigrationBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &BatchDataMigration{BatchDataMigrationCaller: BatchDataMigrationCaller{contract: contract}, BatchDataMigrationTransactor: BatchDataMigrationTransactor{contract: contract}, BatchDataMigrationFilterer: BatchDataMigrationFilterer{contract: contract}}, nil
}

// BatchDataMigration is an auto generated Go binding around an Ethereum contract.
type BatchDataMigration struct {
	BatchDataMigrationCaller     // Read-only binding to the contract
	BatchDataMigrationTransactor // Write-only binding to the contract
	BatchDataMigrationFilterer   // Log filterer for contract events
}

// BatchDataMigrationCaller is an auto generated read-only Go binding around an Ethereum contract.
type BatchDataMigrationCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchDataMigrationTransactor is an auto generated write-only Go binding around an Ethereum contract.
type BatchDataMigrationTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchDataMigrationFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type BatchDataMigrationFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// BatchDataMigrationSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type BatchDataMigrationSession struct {
	Contract     *BatchDataMigration // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// BatchDataMigrationCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type BatchDataMigrationCallerSession struct {
	Contract *BatchDataMigrationCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// BatchDataMigrationTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type BatchDataMigrationTransactorSession struct {
	Contract     *BatchDataMigrationTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// BatchDataMigrationRaw is an auto generated low-level Go binding around an Ethereum contract.
type BatchDataMigrationRaw struct {
	Contract *BatchDataMigration // Generic contract binding to access the raw methods on
}

// BatchDataMigrationCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type BatchDataMigrationCallerRaw struct {
	Contract *BatchDataMigrationCaller // Generic read-only contract binding to access the raw methods on
}

// BatchDataMigrationTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type BatchDataMigrationTransactorRaw struct {
	Contract *BatchDataMigrationTransactor // Generic write-only contract binding to access the raw methods on
}

// NewBatchDataMigration creates a new instance of BatchDataMigration, bound to a specific deployed contract.
func NewBatchDataMigration(address common.Address, backend bind.ContractBackend) (*BatchDataMigration, error) {
	contract, err := bindBatchDataMigration(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &BatchDataMigration{BatchDataMigrationCaller: BatchDataMigrationCaller{contract: contract}, BatchDataMigrationTransactor: BatchDataMigrationTransactor{contract: contract}, BatchDataMigrationFilterer: BatchDataMigrationFilterer{contract: contract}}, nil
}

// NewBatchDataMigrationCaller creates a new read-only instance of BatchDataMigration, bound to a specific deployed contract.
func NewBatchDataMigrationCaller(address common.Address, caller bind.ContractCaller) (*BatchDataMigrationCaller, error) {
	contract, err := bindBatchDataMigration(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &BatchDataMigrationCaller{contract: contract}, nil
}

// NewBatchDataMigrationTransactor creates a new write-only instance of BatchDataMigration, bound to a specific deployed contract.
func NewBatchDataMigrationTransactor(address common.Address, transactor bind.ContractTransactor) (*BatchDataMigrationTransactor, error) {
	contract, err := bindBatchDataMigration(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &BatchDataMigrationTransactor{contract: contract}, nil
}

// NewBatchDataMigrationFilterer creates a new log filterer instance of BatchDataMigration, bound to a specific deployed contract.
func NewBatchDataMigrationFilterer(address common.Address, filterer bind.ContractFilterer) (*BatchDataMigrationFilterer, error) {
	contract, err := bindBatchDataMigration(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &BatchDataMigrationFilterer{contract: contract}, nil
}

// bindBatchDataMigration binds a generic wrapper to an already deployed contract.
func bindBatchDataMigration(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := BatchDataMigrationMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchDataMigration *BatchDataMigrationRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchDataMigration.Contract.BatchDataMigrationCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchDataMigration *BatchDataMigrationRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.BatchDataMigrationTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchDataMigration *BatchDataMigrationRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.BatchDataMigrationTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_BatchDataMigration *BatchDataMigrationCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _BatchDataMigration.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_BatchDataMigration *BatchDataMigrationTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_BatchDataMigration *BatchDataMigrationTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.contract.Transact(opts, method, params...)
}

// BatchRoots is a free data retrieval call binding the contract method 0x68e2f979.
//
// Solidity: function batchRoots(uint256 ) view returns(bytes32)
func (_BatchDataMigration *BatchDataMigrationCaller) BatchRoots(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _BatchDataMigration.contract.Call(opts, &out, "batchRoots", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// BatchRoots is a free data retrieval call binding the contract method 0x68e2f979.
//
// Solidity: function batchRoots(uint256 ) view returns(bytes32)
func (_BatchDataMigration *BatchDataMigrationSession) BatchRoots(arg0 *big.Int) ([32]byte, error) {
	return _BatchDataMigration.Contract.BatchRoots(&_BatchDataMigration.CallOpts, arg0)
}

// BatchRoots is a free data retrieval call binding the contract method 0x68e2f979.
//
// Solidity: function batchRoots(uint256 ) view returns(bytes32)
func (_BatchDataMigration *BatchDataMigrationCallerSession) BatchRoots(arg0 *big.Int) ([32]byte, error) {
	return _BatchDataMigration.Contract.BatchRoots(&_BatchDataMigration.CallOpts, arg0)
}

// BatchTimestamps is a free data retrieval call binding the contract method 0xb4f7dc49.
//
// Solidity: function batchTimestamps(uint256 ) view returns(uint256)
func (_BatchDataMigration *BatchDataMigrationCaller) BatchTimestamps(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _BatchDataMigration.contract.Call(opts, &out, "batchTimestamps", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BatchTimestamps is a free data retrieval call binding the contract method 0xb4f7dc49.
//
// Solidity: function batchTimestamps(uint256 ) view returns(uint256)
func (_BatchDataMigration *BatchDataMigrationSession) BatchTimestamps(arg0 *big.Int) (*big.Int, error) {
	return _BatchDataMigration.Contract.BatchTimestamps(&_BatchDataMigration.CallOpts, arg0)
}

// BatchTimestamps is a free data retrieval call binding the contract method 0xb4f7dc49.
//
// Solidity: function batchTimestamps(uint256 ) view returns(uint256)
func (_BatchDataMigration *BatchDataMigrationCallerSession) BatchTimestamps(arg0 *big.Int) (*big.Int, error) {
	return _BatchDataMigration.Contract.BatchTimestamps(&_BatchDataMigration.CallOpts, arg0)
}

// GetBatchRoot is a free data retrieval call binding the contract method 0xe2350d63.
//
// Solidity: function getBatchRoot(uint256 batchSize) view returns(bytes32)
func (_BatchDataMigration *BatchDataMigrationCaller) GetBatchRoot(opts *bind.CallOpts, batchSize *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _BatchDataMigration.contract.Call(opts, &out, "getBatchRoot", batchSize)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetBatchRoot is a free data retrieval call binding the contract method 0xe2350d63.
//
// Solidity: function getBatchRoot(uint256 batchSize) view returns(bytes32)
func (_BatchDataMigration *BatchDataMigrationSession) GetBatchRoot(batchSize *big.Int) ([32]byte, error) {
	return _BatchDataMigration.Contract.GetBatchRoot(&_BatchDataMigration.CallOpts, batchSize)
}

// GetBatchRoot is a free data retrieval call binding the contract method 0xe2350d63.
//
// Solidity: function getBatchRoot(uint256 batchSize) view returns(bytes32)
func (_BatchDataMigration *BatchDataMigrationCallerSession) GetBatchRoot(batchSize *big.Int) ([32]byte, error) {
	return _BatchDataMigration.Contract.GetBatchRoot(&_BatchDataMigration.CallOpts, batchSize)
}

// GetVerifier is a free data retrieval call binding the contract method 0x3561bc27.
//
// Solidity: function getVerifier(uint256 batchSize) view returns(address)
func (_BatchDataMigration *BatchDataMigrationCaller) GetVerifier(opts *bind.CallOpts, batchSize *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BatchDataMigration.contract.Call(opts, &out, "getVerifier", batchSize)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetVerifier is a free data retrieval call binding the contract method 0x3561bc27.
//
// Solidity: function getVerifier(uint256 batchSize) view returns(address)
func (_BatchDataMigration *BatchDataMigrationSession) GetVerifier(batchSize *big.Int) (common.Address, error) {
	return _BatchDataMigration.Contract.GetVerifier(&_BatchDataMigration.CallOpts, batchSize)
}

// GetVerifier is a free data retrieval call binding the contract method 0x3561bc27.
//
// Solidity: function getVerifier(uint256 batchSize) view returns(address)
func (_BatchDataMigration *BatchDataMigrationCallerSession) GetVerifier(batchSize *big.Int) (common.Address, error) {
	return _BatchDataMigration.Contract.GetVerifier(&_BatchDataMigration.CallOpts, batchSize)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BatchDataMigration *BatchDataMigrationCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _BatchDataMigration.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BatchDataMigration *BatchDataMigrationSession) Owner() (common.Address, error) {
	return _BatchDataMigration.Contract.Owner(&_BatchDataMigration.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_BatchDataMigration *BatchDataMigrationCallerSession) Owner() (common.Address, error) {
	return _BatchDataMigration.Contract.Owner(&_BatchDataMigration.CallOpts)
}

// UsedSerialNumbers is a free data retrieval call binding the contract method 0x74f37505.
//
// Solidity: function usedSerialNumbers(uint256 ) view returns(bool)
func (_BatchDataMigration *BatchDataMigrationCaller) UsedSerialNumbers(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var out []interface{}
	err := _BatchDataMigration.contract.Call(opts, &out, "usedSerialNumbers", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UsedSerialNumbers is a free data retrieval call binding the contract method 0x74f37505.
//
// Solidity: function usedSerialNumbers(uint256 ) view returns(bool)
func (_BatchDataMigration *BatchDataMigrationSession) UsedSerialNumbers(arg0 *big.Int) (bool, error) {
	return _BatchDataMigration.Contract.UsedSerialNumbers(&_BatchDataMigration.CallOpts, arg0)
}

// UsedSerialNumbers is a free data retrieval call binding the contract method 0x74f37505.
//
// Solidity: function usedSerialNumbers(uint256 ) view returns(bool)
func (_BatchDataMigration *BatchDataMigrationCallerSession) UsedSerialNumbers(arg0 *big.Int) (bool, error) {
	return _BatchDataMigration.Contract.UsedSerialNumbers(&_BatchDataMigration.CallOpts, arg0)
}

// Verifiers is a free data retrieval call binding the contract method 0xac1eff68.
//
// Solidity: function verifiers(uint256 ) view returns(address)
func (_BatchDataMigration *BatchDataMigrationCaller) Verifiers(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _BatchDataMigration.contract.Call(opts, &out, "verifiers", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Verifiers is a free data retrieval call binding the contract method 0xac1eff68.
//
// Solidity: function verifiers(uint256 ) view returns(address)
func (_BatchDataMigration *BatchDataMigrationSession) Verifiers(arg0 *big.Int) (common.Address, error) {
	return _BatchDataMigration.Contract.Verifiers(&_BatchDataMigration.CallOpts, arg0)
}

// Verifiers is a free data retrieval call binding the contract method 0xac1eff68.
//
// Solidity: function verifiers(uint256 ) view returns(address)
func (_BatchDataMigration *BatchDataMigrationCallerSession) Verifiers(arg0 *big.Int) (common.Address, error) {
	return _BatchDataMigration.Contract.Verifiers(&_BatchDataMigration.CallOpts, arg0)
}

// SetAllVerifiers is a paid mutator transaction binding the contract method 0xfb2cd937.
//
// Solidity: function setAllVerifiers(address verifier16, address verifier64, address verifier128, address verifier256) returns()
func (_BatchDataMigration *BatchDataMigrationTransactor) SetAllVerifiers(opts *bind.TransactOpts, verifier16 common.Address, verifier64 common.Address, verifier128 common.Address, verifier256 common.Address) (*types.Transaction, error) {
	return _BatchDataMigration.contract.Transact(opts, "setAllVerifiers", verifier16, verifier64, verifier128, verifier256)
}

// SetAllVerifiers is a paid mutator transaction binding the contract method 0xfb2cd937.
//
// Solidity: function setAllVerifiers(address verifier16, address verifier64, address verifier128, address verifier256) returns()
func (_BatchDataMigration *BatchDataMigrationSession) SetAllVerifiers(verifier16 common.Address, verifier64 common.Address, verifier128 common.Address, verifier256 common.Address) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.SetAllVerifiers(&_BatchDataMigration.TransactOpts, verifier16, verifier64, verifier128, verifier256)
}

// SetAllVerifiers is a paid mutator transaction binding the contract method 0xfb2cd937.
//
// Solidity: function setAllVerifiers(address verifier16, address verifier64, address verifier128, address verifier256) returns()
func (_BatchDataMigration *BatchDataMigrationTransactorSession) SetAllVerifiers(verifier16 common.Address, verifier64 common.Address, verifier128 common.Address, verifier256 common.Address) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.SetAllVerifiers(&_BatchDataMigration.TransactOpts, verifier16, verifier64, verifier128, verifier256)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x1957ba4e.
//
// Solidity: function setVerifier(uint256 batchSize, address verifierAddress) returns()
func (_BatchDataMigration *BatchDataMigrationTransactor) SetVerifier(opts *bind.TransactOpts, batchSize *big.Int, verifierAddress common.Address) (*types.Transaction, error) {
	return _BatchDataMigration.contract.Transact(opts, "setVerifier", batchSize, verifierAddress)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x1957ba4e.
//
// Solidity: function setVerifier(uint256 batchSize, address verifierAddress) returns()
func (_BatchDataMigration *BatchDataMigrationSession) SetVerifier(batchSize *big.Int, verifierAddress common.Address) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.SetVerifier(&_BatchDataMigration.TransactOpts, batchSize, verifierAddress)
}

// SetVerifier is a paid mutator transaction binding the contract method 0x1957ba4e.
//
// Solidity: function setVerifier(uint256 batchSize, address verifierAddress) returns()
func (_BatchDataMigration *BatchDataMigrationTransactorSession) SetVerifier(batchSize *big.Int, verifierAddress common.Address) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.SetVerifier(&_BatchDataMigration.TransactOpts, batchSize, verifierAddress)
}

// SubmitBatchRoot is a paid mutator transaction binding the contract method 0x301e24a4.
//
// Solidity: function submitBatchRoot(uint256 batchSize, bytes32 root) returns()
func (_BatchDataMigration *BatchDataMigrationTransactor) SubmitBatchRoot(opts *bind.TransactOpts, batchSize *big.Int, root [32]byte) (*types.Transaction, error) {
	return _BatchDataMigration.contract.Transact(opts, "submitBatchRoot", batchSize, root)
}

// SubmitBatchRoot is a paid mutator transaction binding the contract method 0x301e24a4.
//
// Solidity: function submitBatchRoot(uint256 batchSize, bytes32 root) returns()
func (_BatchDataMigration *BatchDataMigrationSession) SubmitBatchRoot(batchSize *big.Int, root [32]byte) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.SubmitBatchRoot(&_BatchDataMigration.TransactOpts, batchSize, root)
}

// SubmitBatchRoot is a paid mutator transaction binding the contract method 0x301e24a4.
//
// Solidity: function submitBatchRoot(uint256 batchSize, bytes32 root) returns()
func (_BatchDataMigration *BatchDataMigrationTransactorSession) SubmitBatchRoot(batchSize *big.Int, root [32]byte) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.SubmitBatchRoot(&_BatchDataMigration.TransactOpts, batchSize, root)
}

// Unlock is a paid mutator transaction binding the contract method 0x437f5ae5.
//
// Solidity: function unlock(uint256 batchSize, uint256[8] proof, uint256[2] publicInputs) returns()
func (_BatchDataMigration *BatchDataMigrationTransactor) Unlock(opts *bind.TransactOpts, batchSize *big.Int, proof [8]*big.Int, publicInputs [2]*big.Int) (*types.Transaction, error) {
	return _BatchDataMigration.contract.Transact(opts, "unlock", batchSize, proof, publicInputs)
}

// Unlock is a paid mutator transaction binding the contract method 0x437f5ae5.
//
// Solidity: function unlock(uint256 batchSize, uint256[8] proof, uint256[2] publicInputs) returns()
func (_BatchDataMigration *BatchDataMigrationSession) Unlock(batchSize *big.Int, proof [8]*big.Int, publicInputs [2]*big.Int) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.Unlock(&_BatchDataMigration.TransactOpts, batchSize, proof, publicInputs)
}

// Unlock is a paid mutator transaction binding the contract method 0x437f5ae5.
//
// Solidity: function unlock(uint256 batchSize, uint256[8] proof, uint256[2] publicInputs) returns()
func (_BatchDataMigration *BatchDataMigrationTransactorSession) Unlock(batchSize *big.Int, proof [8]*big.Int, publicInputs [2]*big.Int) (*types.Transaction, error) {
	return _BatchDataMigration.Contract.Unlock(&_BatchDataMigration.TransactOpts, batchSize, proof, publicInputs)
}

// BatchDataMigrationBatchRootSubmittedIterator is returned from FilterBatchRootSubmitted and is used to iterate over the raw logs and unpacked data for BatchRootSubmitted events raised by the BatchDataMigration contract.
type BatchDataMigrationBatchRootSubmittedIterator struct {
	Event *BatchDataMigrationBatchRootSubmitted // Event containing the contract specifics and raw log

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
func (it *BatchDataMigrationBatchRootSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchDataMigrationBatchRootSubmitted)
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
		it.Event = new(BatchDataMigrationBatchRootSubmitted)
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
func (it *BatchDataMigrationBatchRootSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BatchDataMigrationBatchRootSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BatchDataMigrationBatchRootSubmitted represents a BatchRootSubmitted event raised by the BatchDataMigration contract.
type BatchDataMigrationBatchRootSubmitted struct {
	BatchSize *big.Int
	Root      [32]byte
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterBatchRootSubmitted is a free log retrieval operation binding the contract event 0x332a878a15e4c00bf130d2a4dd82a0e2bee22ed9f5a1cdcf9c8d9e035461f22f.
//
// Solidity: event BatchRootSubmitted(uint256 indexed batchSize, bytes32 indexed root, uint256 timestamp)
func (_BatchDataMigration *BatchDataMigrationFilterer) FilterBatchRootSubmitted(opts *bind.FilterOpts, batchSize []*big.Int, root [][32]byte) (*BatchDataMigrationBatchRootSubmittedIterator, error) {

	var batchSizeRule []interface{}
	for _, batchSizeItem := range batchSize {
		batchSizeRule = append(batchSizeRule, batchSizeItem)
	}
	var rootRule []interface{}
	for _, rootItem := range root {
		rootRule = append(rootRule, rootItem)
	}

	logs, sub, err := _BatchDataMigration.contract.FilterLogs(opts, "BatchRootSubmitted", batchSizeRule, rootRule)
	if err != nil {
		return nil, err
	}
	return &BatchDataMigrationBatchRootSubmittedIterator{contract: _BatchDataMigration.contract, event: "BatchRootSubmitted", logs: logs, sub: sub}, nil
}

// WatchBatchRootSubmitted is a free log subscription operation binding the contract event 0x332a878a15e4c00bf130d2a4dd82a0e2bee22ed9f5a1cdcf9c8d9e035461f22f.
//
// Solidity: event BatchRootSubmitted(uint256 indexed batchSize, bytes32 indexed root, uint256 timestamp)
func (_BatchDataMigration *BatchDataMigrationFilterer) WatchBatchRootSubmitted(opts *bind.WatchOpts, sink chan<- *BatchDataMigrationBatchRootSubmitted, batchSize []*big.Int, root [][32]byte) (event.Subscription, error) {

	var batchSizeRule []interface{}
	for _, batchSizeItem := range batchSize {
		batchSizeRule = append(batchSizeRule, batchSizeItem)
	}
	var rootRule []interface{}
	for _, rootItem := range root {
		rootRule = append(rootRule, rootItem)
	}

	logs, sub, err := _BatchDataMigration.contract.WatchLogs(opts, "BatchRootSubmitted", batchSizeRule, rootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BatchDataMigrationBatchRootSubmitted)
				if err := _BatchDataMigration.contract.UnpackLog(event, "BatchRootSubmitted", log); err != nil {
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

// ParseBatchRootSubmitted is a log parse operation binding the contract event 0x332a878a15e4c00bf130d2a4dd82a0e2bee22ed9f5a1cdcf9c8d9e035461f22f.
//
// Solidity: event BatchRootSubmitted(uint256 indexed batchSize, bytes32 indexed root, uint256 timestamp)
func (_BatchDataMigration *BatchDataMigrationFilterer) ParseBatchRootSubmitted(log types.Log) (*BatchDataMigrationBatchRootSubmitted, error) {
	event := new(BatchDataMigrationBatchRootSubmitted)
	if err := _BatchDataMigration.contract.UnpackLog(event, "BatchRootSubmitted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BatchDataMigrationUnlockedIterator is returned from FilterUnlocked and is used to iterate over the raw logs and unpacked data for Unlocked events raised by the BatchDataMigration contract.
type BatchDataMigrationUnlockedIterator struct {
	Event *BatchDataMigrationUnlocked // Event containing the contract specifics and raw log

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
func (it *BatchDataMigrationUnlockedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchDataMigrationUnlocked)
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
		it.Event = new(BatchDataMigrationUnlocked)
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
func (it *BatchDataMigrationUnlockedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BatchDataMigrationUnlockedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BatchDataMigrationUnlocked represents a Unlocked event raised by the BatchDataMigration contract.
type BatchDataMigrationUnlocked struct {
	BatchSize    *big.Int
	SerialNumber *big.Int
	Unlocker     common.Address
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterUnlocked is a free log retrieval operation binding the contract event 0xdf9bbb560a3c97f1125d5e2b956725434c1ba1b7bc6d1cf1119c83cc230d7a50.
//
// Solidity: event Unlocked(uint256 indexed batchSize, uint256 indexed serialNumber, address indexed unlocker)
func (_BatchDataMigration *BatchDataMigrationFilterer) FilterUnlocked(opts *bind.FilterOpts, batchSize []*big.Int, serialNumber []*big.Int, unlocker []common.Address) (*BatchDataMigrationUnlockedIterator, error) {

	var batchSizeRule []interface{}
	for _, batchSizeItem := range batchSize {
		batchSizeRule = append(batchSizeRule, batchSizeItem)
	}
	var serialNumberRule []interface{}
	for _, serialNumberItem := range serialNumber {
		serialNumberRule = append(serialNumberRule, serialNumberItem)
	}
	var unlockerRule []interface{}
	for _, unlockerItem := range unlocker {
		unlockerRule = append(unlockerRule, unlockerItem)
	}

	logs, sub, err := _BatchDataMigration.contract.FilterLogs(opts, "Unlocked", batchSizeRule, serialNumberRule, unlockerRule)
	if err != nil {
		return nil, err
	}
	return &BatchDataMigrationUnlockedIterator{contract: _BatchDataMigration.contract, event: "Unlocked", logs: logs, sub: sub}, nil
}

// WatchUnlocked is a free log subscription operation binding the contract event 0xdf9bbb560a3c97f1125d5e2b956725434c1ba1b7bc6d1cf1119c83cc230d7a50.
//
// Solidity: event Unlocked(uint256 indexed batchSize, uint256 indexed serialNumber, address indexed unlocker)
func (_BatchDataMigration *BatchDataMigrationFilterer) WatchUnlocked(opts *bind.WatchOpts, sink chan<- *BatchDataMigrationUnlocked, batchSize []*big.Int, serialNumber []*big.Int, unlocker []common.Address) (event.Subscription, error) {

	var batchSizeRule []interface{}
	for _, batchSizeItem := range batchSize {
		batchSizeRule = append(batchSizeRule, batchSizeItem)
	}
	var serialNumberRule []interface{}
	for _, serialNumberItem := range serialNumber {
		serialNumberRule = append(serialNumberRule, serialNumberItem)
	}
	var unlockerRule []interface{}
	for _, unlockerItem := range unlocker {
		unlockerRule = append(unlockerRule, unlockerItem)
	}

	logs, sub, err := _BatchDataMigration.contract.WatchLogs(opts, "Unlocked", batchSizeRule, serialNumberRule, unlockerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BatchDataMigrationUnlocked)
				if err := _BatchDataMigration.contract.UnpackLog(event, "Unlocked", log); err != nil {
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

// ParseUnlocked is a log parse operation binding the contract event 0xdf9bbb560a3c97f1125d5e2b956725434c1ba1b7bc6d1cf1119c83cc230d7a50.
//
// Solidity: event Unlocked(uint256 indexed batchSize, uint256 indexed serialNumber, address indexed unlocker)
func (_BatchDataMigration *BatchDataMigrationFilterer) ParseUnlocked(log types.Log) (*BatchDataMigrationUnlocked, error) {
	event := new(BatchDataMigrationUnlocked)
	if err := _BatchDataMigration.contract.UnpackLog(event, "Unlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// BatchDataMigrationVerifierSetIterator is returned from FilterVerifierSet and is used to iterate over the raw logs and unpacked data for VerifierSet events raised by the BatchDataMigration contract.
type BatchDataMigrationVerifierSetIterator struct {
	Event *BatchDataMigrationVerifierSet // Event containing the contract specifics and raw log

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
func (it *BatchDataMigrationVerifierSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(BatchDataMigrationVerifierSet)
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
		it.Event = new(BatchDataMigrationVerifierSet)
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
func (it *BatchDataMigrationVerifierSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *BatchDataMigrationVerifierSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// BatchDataMigrationVerifierSet represents a VerifierSet event raised by the BatchDataMigration contract.
type BatchDataMigrationVerifierSet struct {
	BatchSize *big.Int
	Verifier  common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterVerifierSet is a free log retrieval operation binding the contract event 0xbc291d0e6f60c8ebaeb52dc4380cd7fa1fa6ac795fc0719ea5f57b8acfa28e3c.
//
// Solidity: event VerifierSet(uint256 indexed batchSize, address verifier)
func (_BatchDataMigration *BatchDataMigrationFilterer) FilterVerifierSet(opts *bind.FilterOpts, batchSize []*big.Int) (*BatchDataMigrationVerifierSetIterator, error) {

	var batchSizeRule []interface{}
	for _, batchSizeItem := range batchSize {
		batchSizeRule = append(batchSizeRule, batchSizeItem)
	}

	logs, sub, err := _BatchDataMigration.contract.FilterLogs(opts, "VerifierSet", batchSizeRule)
	if err != nil {
		return nil, err
	}
	return &BatchDataMigrationVerifierSetIterator{contract: _BatchDataMigration.contract, event: "VerifierSet", logs: logs, sub: sub}, nil
}

// WatchVerifierSet is a free log subscription operation binding the contract event 0xbc291d0e6f60c8ebaeb52dc4380cd7fa1fa6ac795fc0719ea5f57b8acfa28e3c.
//
// Solidity: event VerifierSet(uint256 indexed batchSize, address verifier)
func (_BatchDataMigration *BatchDataMigrationFilterer) WatchVerifierSet(opts *bind.WatchOpts, sink chan<- *BatchDataMigrationVerifierSet, batchSize []*big.Int) (event.Subscription, error) {

	var batchSizeRule []interface{}
	for _, batchSizeItem := range batchSize {
		batchSizeRule = append(batchSizeRule, batchSizeItem)
	}

	logs, sub, err := _BatchDataMigration.contract.WatchLogs(opts, "VerifierSet", batchSizeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(BatchDataMigrationVerifierSet)
				if err := _BatchDataMigration.contract.UnpackLog(event, "VerifierSet", log); err != nil {
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

// ParseVerifierSet is a log parse operation binding the contract event 0xbc291d0e6f60c8ebaeb52dc4380cd7fa1fa6ac795fc0719ea5f57b8acfa28e3c.
//
// Solidity: event VerifierSet(uint256 indexed batchSize, address verifier)
func (_BatchDataMigration *BatchDataMigrationFilterer) ParseVerifierSet(log types.Log) (*BatchDataMigrationVerifierSet, error) {
	event := new(BatchDataMigrationVerifierSet)
	if err := _BatchDataMigration.contract.UnpackLog(event, "VerifierSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
