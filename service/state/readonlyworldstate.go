package state

import (
	"sync"

	"github.com/icon-project/goloop/common/db"
	"github.com/icon-project/goloop/common/errors"
	"github.com/icon-project/goloop/module"
)

type readOnlyWorldState struct {
	WorldSnapshot

	lock           sync.Mutex
	accounts       map[string]AccountState
	validatorState ValidatorState
}

func (ws *readOnlyWorldState) GetAccountState(id []byte) AccountState {
	ws.lock.Lock()
	defer ws.lock.Unlock()

	ids := string(id)
	if as, ok := ws.accounts[ids]; ok {
		return as
	}

	as := newAccountROState(ws.WorldSnapshot.GetAccountSnapshot(id))
	ws.accounts[ids] = as

	return as
}

func (ws *readOnlyWorldState) GetSnapshot() WorldSnapshot {
	return ws.WorldSnapshot
}

func (ws *readOnlyWorldState) GetValidatorState() ValidatorState {
	return ws.validatorState
}

func (ws *readOnlyWorldState) Reset(snapshot WorldSnapshot) error {
	if ws.WorldSnapshot != snapshot {
		return errors.InvalidStateError.New(
			"readOnlyWorldState.Reset() with different snapshot")
	}
	return nil
}

func (ws *readOnlyWorldState) ClearCache() {
	// nothing to do
}

type readonlyValidatorState struct {
	ValidatorSnapshot
}

func (vs *readonlyValidatorState) Set([]module.Validator) error {
	return errors.InvalidStateError.New("ReadOnlyState")
}

func (vs *readonlyValidatorState) Add(v module.Validator) error {
	return errors.InvalidStateError.New("ReadOnlyState")
}

func (vs *readonlyValidatorState) Remove(v module.Validator) bool {
	return false
}

func (vs *readonlyValidatorState) GetSnapshot() ValidatorSnapshot {
	return vs.ValidatorSnapshot
}

func (vs *readonlyValidatorState) Reset(vss ValidatorSnapshot) {
	// do nothing
}

func newReadOnlyValidatorState(vss ValidatorSnapshot) ValidatorState {
	return &readonlyValidatorState{vss}
}

func NewReadOnlyWorldState(database db.Database, stateHash []byte, vs ValidatorSnapshot) WorldState {
	wss := NewWorldSnapshot(database, stateHash, vs)
	return &readOnlyWorldState{
		WorldSnapshot:  wss,
		accounts:       make(map[string]AccountState),
		validatorState: newReadOnlyValidatorState(wss.GetValidatorSnapshot()),
	}
}