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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_unlockVerifier\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_auditVerifier\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"h\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"dataId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeout\",\"type\":\"uint256\"}],\"name\":\"Locked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"h\",\"type\":\"bytes32\"}],\"name\":\"Reclaimed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"h\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"mechanic\",\"type\":\"string\"}],\"name\":\"Unlocked\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"activeLocks\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"locker\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"dataId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"timeout\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256[3]\",\"name\":\"publicInputs\",\"type\":\"uint256[3]\"}],\"name\":\"auditUnlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"auditVerifier\",\"outputs\":[{\"internalType\":\"contractIAuditVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_h\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_dataId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_timeoutDuration\",\"type\":\"uint256\"}],\"name\":\"lock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_h\",\"type\":\"bytes32\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256[8]\",\"name\":\"proof\",\"type\":\"uint256[8]\"},{\"internalType\":\"uint256[2]\",\"name\":\"publicInputs\",\"type\":\"uint256[2]\"}],\"name\":\"unlock\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"unlockVerifier\",\"outputs\":[{\"internalType\":\"contractIUnlockVerifier\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b50604051611681380380611681833981810160405281019061003291906101e2565b600073ffffffffffffffffffffffffffffffffffffffff168273ffffffffffffffffffffffffffffffffffffffff16036100a1576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100989061027f565b60405180910390fd5b600073ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1603610110576040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610107906102eb565b60405180910390fd5b8173ffffffffffffffffffffffffffffffffffffffff1660808173ffffffffffffffffffffffffffffffffffffffff16815250508073ffffffffffffffffffffffffffffffffffffffff1660a08173ffffffffffffffffffffffffffffffffffffffff1681525050505061030b565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b60006101af82610184565b9050919050565b6101bf816101a4565b81146101ca57600080fd5b50565b6000815190506101dc816101b6565b92915050565b600080604083850312156101f9576101f861017f565b5b6000610207858286016101cd565b9250506020610218858286016101cd565b9150509250929050565b600082825260208201905092915050565b7f556e6c6f636b2056657269666965722063616e6e6f74206265207a65726f0000600082015250565b6000610269601e83610222565b915061027482610233565b602082019050919050565b600060208201905081810360008301526102988161025c565b9050919050565b7f41756469742056657269666965722063616e6e6f74206265207a65726f000000600082015250565b60006102d5601d83610222565b91506102e08261029f565b602082019050919050565b60006020820190508181036000830152610304816102c8565b9050919050565b60805160a05161134361033e6000396000818161038801526104700152600081816102260152610a4001526113436000f3fe608060405234801561001057600080fd5b506004361061007d5760003560e01c806347f2f1781161005b57806347f2f178146100d857806396afb365146100f4578063b51a201f14610110578063c0f3029b146101425761007d565b806321915f2814610082578063265297ac1461009e57806338466c66146100bc575b600080fd5b61009c60048036038101906100979190610ab0565b610160565b005b6100a6610386565b6040516100b39190610b71565b60405180910390f35b6100d660048036038101906100d19190610bae565b6103aa565b005b6100f260048036038101906100ed9190610c5c565b6105d0565b005b61010e60048036038101906101099190610caf565b6107c0565b005b61012a60048036038101906101259190610caf565b6109f4565b60405161013993929190610d1b565b60405180910390f35b61014a610a3e565b6040516101579190610d73565b60405180910390f35b60008160006002811061017657610175610d8e565b5b602002013560001b9050600073ffffffffffffffffffffffffffffffffffffffff1660008083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1603610224576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161021b90610e1a565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635fe24f2384846040518363ffffffff1660e01b815260040161027f929190610e64565b602060405180830381865afa15801561029c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906102c09190610ec7565b6102ff576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016102f690610f40565b60405180910390fd5b600080828152602001908152602001600020600080820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600182016000905560028201600090555050807f702d94ddc1e4f1e432b55d42e9531acdfbcfc079c640901654a219e48863127c60405161037990610fac565b60405180910390a2505050565b7f000000000000000000000000000000000000000000000000000000000000000081565b6000816002600381106103c0576103bf610d8e565b5b602002013560001b9050600073ffffffffffffffffffffffffffffffffffffffff1660008083815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff160361046e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161046590610e1a565b60405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166365c0325984846040518363ffffffff1660e01b81526004016104c9929190610fdc565b602060405180830381865afa1580156104e6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061050a9190610ec7565b610549576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161054090611053565b60405180910390fd5b600080828152602001908152602001600020600080820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600182016000905560028201600090555050807f702d94ddc1e4f1e432b55d42e9531acdfbcfc079c640901654a219e48863127c6040516105c3906110bf565b60405180910390a2505050565b600073ffffffffffffffffffffffffffffffffffffffff1660008085815260200190815260200160002060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1614610674576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161066b9061112b565b60405180910390fd5b6000801b83036106b9576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016106b090611197565b60405180910390fd5b600081426106c791906111e6565b905060405180606001604052803373ffffffffffffffffffffffffffffffffffffffff1681526020018481526020018281525060008086815260200190815260200160002060008201518160000160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506020820151816001015560408201518160020155905050823373ffffffffffffffffffffffffffffffffffffffff16857fafec6b642358ef93f88d9a518cde37080f3ff07afee5022bd614276d0ec7790f846040516107b2919061121a565b60405180910390a450505050565b60008060008381526020019081526020016000206040518060600160405290816000820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001600182015481526020016002820154815250509050600073ffffffffffffffffffffffffffffffffffffffff16816000015173ffffffffffffffffffffffffffffffffffffffff16036108be576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016108b590610e1a565b60405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff16816000015173ffffffffffffffffffffffffffffffffffffffff1614610930576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161092790611281565b60405180910390fd5b8060400151421015610977576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161096e906112ed565b60405180910390fd5b600080838152602001908152602001600020600080820160006101000a81549073ffffffffffffffffffffffffffffffffffffffff0219169055600182016000905560028201600090555050817fbe9e485e7f7ace1eaf2897ca5483cdb8bf05d65d8b660c18070acc759652944660405160405180910390a25050565b60006020528060005260406000206000915090508060000160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16908060010154908060020154905083565b7f000000000000000000000000000000000000000000000000000000000000000081565b600080fd5b600080fd5b600081905082602060080282011115610a8857610a87610a67565b5b92915050565b600081905082602060020282011115610aaa57610aa9610a67565b5b92915050565b6000806101408385031215610ac857610ac7610a62565b5b6000610ad685828601610a6c565b925050610100610ae885828601610a8e565b9150509250929050565b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000819050919050565b6000610b37610b32610b2d84610af2565b610b12565b610af2565b9050919050565b6000610b4982610b1c565b9050919050565b6000610b5b82610b3e565b9050919050565b610b6b81610b50565b82525050565b6000602082019050610b866000830184610b62565b92915050565b600081905082602060030282011115610ba857610ba7610a67565b5b92915050565b6000806101608385031215610bc657610bc5610a62565b5b6000610bd485828601610a6c565b925050610100610be685828601610b8c565b9150509250929050565b6000819050919050565b610c0381610bf0565b8114610c0e57600080fd5b50565b600081359050610c2081610bfa565b92915050565b6000819050919050565b610c3981610c26565b8114610c4457600080fd5b50565b600081359050610c5681610c30565b92915050565b600080600060608486031215610c7557610c74610a62565b5b6000610c8386828701610c11565b9350506020610c9486828701610c11565b9250506040610ca586828701610c47565b9150509250925092565b600060208284031215610cc557610cc4610a62565b5b6000610cd384828501610c11565b91505092915050565b6000610ce782610af2565b9050919050565b610cf781610cdc565b82525050565b610d0681610bf0565b82525050565b610d1581610c26565b82525050565b6000606082019050610d306000830186610cee565b610d3d6020830185610cfd565b610d4a6040830184610d0c565b949350505050565b6000610d5d82610b3e565b9050919050565b610d6d81610d52565b82525050565b6000602082019050610d886000830184610d64565b92915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082825260208201905092915050565b7f4c6f636b20646f6573206e6f7420657869737400000000000000000000000000600082015250565b6000610e04601383610dbd565b9150610e0f82610dce565b602082019050919050565b60006020820190508181036000830152610e3381610df7565b9050919050565b82818337505050565b610e506101008383610e3a565b5050565b610e6060408383610e3a565b5050565b600061014082019050610e7a6000830185610e43565b610e88610100830184610e54565b9392505050565b60008115159050919050565b610ea481610e8f565b8114610eaf57600080fd5b50565b600081519050610ec181610e9b565b92915050565b600060208284031215610edd57610edc610a62565b5b6000610eeb84828501610eb2565b91505092915050565b7f496e76616c696420556e6c6f636b2050726f6f66000000000000000000000000600082015250565b6000610f2a601483610dbd565b9150610f3582610ef4565b602082019050919050565b60006020820190508181036000830152610f5981610f1d565b9050919050565b7f445350415f556e6c6f636b000000000000000000000000000000000000000000600082015250565b6000610f96600b83610dbd565b9150610fa182610f60565b602082019050919050565b60006020820190508181036000830152610fc581610f89565b9050919050565b610fd860608383610e3a565b5050565b600061016082019050610ff26000830185610e43565b611000610100830184610fcc565b9392505050565b7f496e76616c69642041756469742050726f6f6600000000000000000000000000600082015250565b600061103d601383610dbd565b915061104882611007565b602082019050919050565b6000602082019050818103600083015261106c81611030565b9050919050565b7f445350425f4175646974556e6c6f636b00000000000000000000000000000000600082015250565b60006110a9601083610dbd565b91506110b482611073565b602082019050919050565b600060208201905081810360008301526110d88161109c565b9050919050565b7f4c6f636b20616c72656164792065786973747300000000000000000000000000600082015250565b6000611115601383610dbd565b9150611120826110df565b602082019050919050565b6000602082019050818103600083015261114481611108565b9050919050565b7f486173682063616e6e6f74206265207a65726f00000000000000000000000000600082015250565b6000611181601383610dbd565b915061118c8261114b565b602082019050919050565b600060208201905081810360008301526111b081611174565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60006111f182610c26565b91506111fc83610c26565b9250828201905080821115611214576112136111b7565b5b92915050565b600060208201905061122f6000830184610d0c565b92915050565b7f4f6e6c79206c6f636b65722063616e207265636c61696d000000000000000000600082015250565b600061126b601783610dbd565b915061127682611235565b602082019050919050565b6000602082019050818103600083015261129a8161125e565b9050919050565b7f54696d656f7574206e6f74207265616368656400000000000000000000000000600082015250565b60006112d7601383610dbd565b91506112e2826112a1565b602082019050919050565b60006020820190508181036000830152611306816112ca565b905091905056fea2646970667358221220999bb1832b926ace8681d148e82531d933b530df282e0a710c093a8f60a735bb64736f6c634300081c0033",
}

// DataMigrationABI is the input ABI used to generate the binding from.
// Deprecated: Use DataMigrationMetaData.ABI instead.
var DataMigrationABI = DataMigrationMetaData.ABI

// DataMigrationBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DataMigrationMetaData.Bin instead.
var DataMigrationBin = DataMigrationMetaData.Bin

// DeployDataMigration deploys a new Ethereum contract, binding an instance of DataMigration to it.
func DeployDataMigration(auth *bind.TransactOpts, backend bind.ContractBackend, _unlockVerifier common.Address, _auditVerifier common.Address) (common.Address, *types.Transaction, *DataMigration, error) {
	parsed, err := DataMigrationMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DataMigrationBin), backend, _unlockVerifier, _auditVerifier)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DataMigration{DataMigrationCaller: DataMigrationCaller{contract: contract}, DataMigrationTransactor: DataMigrationTransactor{contract: contract}, DataMigrationFilterer: DataMigrationFilterer{contract: contract}}, nil
}

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

// AuditVerifier is a free data retrieval call binding the contract method 0x265297ac.
//
// Solidity: function auditVerifier() view returns(address)
func (_DataMigration *DataMigrationCaller) AuditVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataMigration.contract.Call(opts, &out, "auditVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AuditVerifier is a free data retrieval call binding the contract method 0x265297ac.
//
// Solidity: function auditVerifier() view returns(address)
func (_DataMigration *DataMigrationSession) AuditVerifier() (common.Address, error) {
	return _DataMigration.Contract.AuditVerifier(&_DataMigration.CallOpts)
}

// AuditVerifier is a free data retrieval call binding the contract method 0x265297ac.
//
// Solidity: function auditVerifier() view returns(address)
func (_DataMigration *DataMigrationCallerSession) AuditVerifier() (common.Address, error) {
	return _DataMigration.Contract.AuditVerifier(&_DataMigration.CallOpts)
}

// UnlockVerifier is a free data retrieval call binding the contract method 0xc0f3029b.
//
// Solidity: function unlockVerifier() view returns(address)
func (_DataMigration *DataMigrationCaller) UnlockVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataMigration.contract.Call(opts, &out, "unlockVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UnlockVerifier is a free data retrieval call binding the contract method 0xc0f3029b.
//
// Solidity: function unlockVerifier() view returns(address)
func (_DataMigration *DataMigrationSession) UnlockVerifier() (common.Address, error) {
	return _DataMigration.Contract.UnlockVerifier(&_DataMigration.CallOpts)
}

// UnlockVerifier is a free data retrieval call binding the contract method 0xc0f3029b.
//
// Solidity: function unlockVerifier() view returns(address)
func (_DataMigration *DataMigrationCallerSession) UnlockVerifier() (common.Address, error) {
	return _DataMigration.Contract.UnlockVerifier(&_DataMigration.CallOpts)
}

// AuditUnlock is a paid mutator transaction binding the contract method 0x38466c66.
//
// Solidity: function auditUnlock(uint256[8] proof, uint256[3] publicInputs) returns()
func (_DataMigration *DataMigrationTransactor) AuditUnlock(opts *bind.TransactOpts, proof [8]*big.Int, publicInputs [3]*big.Int) (*types.Transaction, error) {
	return _DataMigration.contract.Transact(opts, "auditUnlock", proof, publicInputs)
}

// AuditUnlock is a paid mutator transaction binding the contract method 0x38466c66.
//
// Solidity: function auditUnlock(uint256[8] proof, uint256[3] publicInputs) returns()
func (_DataMigration *DataMigrationSession) AuditUnlock(proof [8]*big.Int, publicInputs [3]*big.Int) (*types.Transaction, error) {
	return _DataMigration.Contract.AuditUnlock(&_DataMigration.TransactOpts, proof, publicInputs)
}

// AuditUnlock is a paid mutator transaction binding the contract method 0x38466c66.
//
// Solidity: function auditUnlock(uint256[8] proof, uint256[3] publicInputs) returns()
func (_DataMigration *DataMigrationTransactorSession) AuditUnlock(proof [8]*big.Int, publicInputs [3]*big.Int) (*types.Transaction, error) {
	return _DataMigration.Contract.AuditUnlock(&_DataMigration.TransactOpts, proof, publicInputs)
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
	H        [32]byte
	Mechanic string
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterUnlocked is a free log retrieval operation binding the contract event 0x702d94ddc1e4f1e432b55d42e9531acdfbcfc079c640901654a219e48863127c.
//
// Solidity: event Unlocked(bytes32 indexed h, string mechanic)
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

// WatchUnlocked is a free log subscription operation binding the contract event 0x702d94ddc1e4f1e432b55d42e9531acdfbcfc079c640901654a219e48863127c.
//
// Solidity: event Unlocked(bytes32 indexed h, string mechanic)
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

// ParseUnlocked is a log parse operation binding the contract event 0x702d94ddc1e4f1e432b55d42e9531acdfbcfc079c640901654a219e48863127c.
//
// Solidity: event Unlocked(bytes32 indexed h, string mechanic)
func (_DataMigration *DataMigrationFilterer) ParseUnlocked(log types.Log) (*DataMigrationUnlocked, error) {
	event := new(DataMigrationUnlocked)
	if err := _DataMigration.contract.UnpackLog(event, "Unlocked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
