package module

import (
	"math/big"
)

// TransitionCallback provides transition change notifications. All functions
// are called back with the same Transition instance for the convenience.
type TransitionCallback interface {
	// Called if validation is done.
	OnValidate(Transition, error)

	// Called if execution is done.
	OnExecute(Transition, error)
}

// Block information used by service manager.
type BlockInfo interface {
	Height() int64
	Timestamp() int64
}

type Transaction interface {
	Group() TransactionGroup
	ID() []byte
	Bytes() []byte
	Hash() []byte
	Verify() error
	Version() int
	ToJSON(version int) (interface{}, error)

	// Version() int
	// From() Address
	// To() Address
	// Value() *big.Int
	// StepLimit() *big.Int
	// Timestamp() int64
	// NID() int
	// Nonce() int64
	// Signature() []byte
}

type TransactionIterator interface {
	Has() bool
	Next() error
	Get() (Transaction, int, error)
}

type TransactionList interface {
	Get(int) (Transaction, error)
	Iterator() TransactionIterator
	Hash() []byte
	Equal(TransactionList) bool
	Flush() error
}

type Status int

const (
	StatusSuccess      = 0
	StatusNotPayable   = 0x7d64
	StatusOutOfBalance = 0x7f58
	StatusSystemError  = 0x7000
	StatusTimeout      = 0x7001
)

type Receipt interface {
	Bytes() []byte
	To() Address
	CumulativeStepUsed() *big.Int
	StepPrice() *big.Int
	StepUsed() *big.Int
	Status() Status
	SCOREAddress() Address
	Check(r Receipt) error
	ToJSON(int) (interface{}, error)
}

type ReceiptIterator interface {
	Has() bool
	Next() error
	Get() (Receipt, error)
}

type ReceiptList interface {
	Get(int) (Receipt, error)
	Iterator() ReceiptIterator
	Hash() []byte
	Flush() error
}

type Transition interface {
	PatchTransactions() TransactionList
	NormalTransactions() TransactionList

	// Execute executes this transition.
	// The result is asynchronously notified by cb. canceler can be used
	// to cancel it after calling Execute. After canceler returns true,
	// all succeeding cb functions may not be called back.
	// REMARK: It is assumed to be called once. Any additional call returns
	// error.
	Execute(cb TransitionCallback) (canceler func() bool, err error)

	// Result returns service manager defined result bytes.
	// For example, it can be "[world_state_hash][patch_tx_hash][normal_tx_hash]".
	Result() []byte

	// NextValidators returns the addresses of validators as a result of
	// transaction processing.
	// It may return nil before cb.OnExecute is called back by Execute.
	NextValidators() ValidatorList

	// LogBloom returns log bloom filter for this transition.
	// It may return nil before cb.OnExecute is called back by Execute.
	LogBloom() []byte
}

// Options for finalize
const (
	FinalizeNormalTransaction = 1 << iota
	FinalizePatchTransaction
	FinalizeResult

	// TODO It's only necessary if storing receipt index is determined by
	// block manager. The current service manager determines by itself according
	// to version, so it doesn't use it.
	FinalizeWriteReceiptIndex
)

// ServiceManager provides Service APIs.
// For a block proposal, it is usually called as follows:
// 		1. GetPatches
//		2. if any changes of patches exist from GetPatches
//			2.1 PatchTransaction
//			2.2 Transition.Execute
// 		3. ProposeTransition
//		4. Transition.Execute
// For a block validation,
//		1. if any changes of patches are detected from a new block
//			1.1 PatchTransition
//			1.2 Transition.Execute
//		2. create Transaction instances by TransactionFromBytes
//		3. CreateTransition with TransactionList
//		4. Transition.Execute
type ServiceManager interface {
	// ProposeTransition proposes a Transition following the parent Transition.
	// Returned Transition always passes validation.
	ProposeTransition(parent Transition) (Transition, error)
	// ProposeGenesisTransition proposes a Transition for Genesis
	// with transactions of Genesis.
	ProposeGenesisTransition(parent Transition) (Transition, error)
	// CreateInitialTransition creates an initial Transition
	// Height is the height of the block which includes transactions.
	// e.g. If the transactions are included in block n and results are in
	// block n+1, Height is n. It can be -1 if it is the initial state for
	// a genesis block.
	CreateInitialTransition(result []byte, nextValidators ValidatorList, height int64) (Transition, error)
	// CreateTransition creates a Transition following parent Transition.
	CreateTransition(parent Transition, txs TransactionList) (Transition, error)
	// GetPatches returns all patch transactions based on the parent transition.
	GetPatches(parent Transition) TransactionList
	// PatchTransition creates a Transition by overwriting patches on the transition.
	PatchTransition(transition Transition, patches TransactionList) Transition

	// Finalize finalizes data related to the transition. It usually stores
	// data to a persistent storage. opt indicates which data are finalized.
	// It should be called for every transition.
	Finalize(transition Transition, opt int)

	// TransactionFromBytes returns a Transaction instance from bytes.
	TransactionFromBytes(b []byte, blockVersion int) (Transaction, error)

	// TransactionListFromHash returns a TransactionList instance from
	// the hash of transactions or nil when no transactions exist.
	// It assumes it's called only by new version block, so it doesn't receive
	// version value.
	TransactionListFromHash(hash []byte) TransactionList

	// TransactionListFromSlice returns list of transactions.
	TransactionListFromSlice(txs []Transaction, version int) TransactionList

	// ReceiptListFromResult returns list of receipts from result.
	ReceiptListFromResult(result []byte, g TransactionGroup) ReceiptList

	// SendTransaction adds transaction to a transaction pool.
	SendTransaction(tx interface{}) ([]byte, error)
	// SendTransaction(tx Transaction) ([]byte, error)

	// ValidatorListFromHash returns ValidatorList from hash.
	ValidatorListFromHash(hash []byte) ValidatorList

	// GetBalance get balance of the account
	GetBalance(result []byte, addr Address) *big.Int
}
