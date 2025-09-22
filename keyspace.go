package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

type keyspace struct {
	*topodatapb.Keyspace
	Name string `json:"keyspace"`
}

func SaveKeyspace(txn *memdb.Txn, keyspaceName string, k *topodatapb.Keyspace) error {
	return txn.Insert(keyspacesTable, &keyspace{
		Name:     keyspaceName,
		Keyspace: k,
	})
}

func DeleteKeyspace(txn *memdb.Txn, keyspaceName string) error {
	deleted, err := txn.DeleteAll(keyspacesTable, keyspacesNameIndex, keyspaceName)
	if err == nil && deleted == 0 {
		err = fmt.Errorf("zero records deleted")
	}
	return err
}
