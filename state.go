package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
)

const (
	stateGlobalRecovery string = "global_recovery"
)

type stateRecordBase struct {
	Name string `json:"name"`
}

type stateRecordGlobalRecovery struct {
	stateRecordBase
	DisableRecovery bool `json:"disable_recovery"`
}

func newStateRecordGlobalRecovery(disableRecovery bool) *stateRecordGlobalRecovery {
	return &stateRecordGlobalRecovery{
		stateRecordBase{Name: stateGlobalRecovery},
		disableRecovery,
	}
}

func IsGlobalRecoveryDisabled(txn *memdb.Txn) (bool, error) {
	res, err := txn.First(stateTable, stateIndex, stateGlobalRecovery)
	if err != nil {
		return false, err
	}
	gr, ok := res.(*stateRecordGlobalRecovery)
	if !ok {
		return false, fmt.Errorf("data must be *stateRecordGlobalRecovery, got %T", res)
	}
	return gr.DisableRecovery, nil
}

func SetGlobalRecoveryEnabled(txn *memdb.Txn) error {
	return txn.Insert(stateTable, newStateRecordGlobalRecovery(false))
}

func SetGlobalRecoveryDisabled(txn *memdb.Txn) error {
	return txn.Insert(stateTable, newStateRecordGlobalRecovery(true))
}
