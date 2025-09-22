package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

type keyspaceRecord struct {
	*topodatapb.Keyspace

	Name string `json:"keyspace"`
}

func SaveKeyspace(txn *memdb.Txn, keyspaceName string, k *topodatapb.Keyspace) error {
	return txn.Insert(keyspacesTable, &keyspaceRecord{
		Keyspace: k,
		Name:     keyspaceName,
	})
}

func ReadKeyspace(txn *memdb.Txn, keyspaceName string) (*topodatapb.Keyspace, error) {
	res, err := txn.First(keyspacesTable, keyspacesNameIndex, keyspaceName)
	if err != nil {
		return nil, err
	}
	k, ok := res.(*keyspaceRecord)
	if !ok {
		return nil, fmt.Errorf("data must be *keyspaceRecord, got %T", res)
	}
	return k.Keyspace, nil
}

func DeleteKeyspace(txn *memdb.Txn, keyspaceName string) error {
	deleted, err := txn.DeleteAll(keyspacesTable, keyspacesNameIndex, keyspaceName)
	if err == nil && deleted == 0 {
		err = fmt.Errorf("zero records deleted")
	}
	return err
}
