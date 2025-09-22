package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

type shard struct {
	*topodatapb.Shard
	Name     string `json:"shard"`
	Keyspace string `json:"keyspace"`
}

func SaveShard(txn *memdb.Txn, keyspaceName, shardName string, s *topodatapb.Shard) error {
	return txn.Insert(shardsTable, &shard{
		Name:     shardName,
		Keyspace: keyspaceName,
		Shard:    s,
	})
}

func DeleteShard(txn *memdb.Txn, keyspaceName, shardName string) error {
	deleted, err := txn.DeleteAll(shardsTable, shardsKeyspaceShardIndex, keyspaceName, shardName)
	if err == nil && deleted == 0 {
		err = fmt.Errorf("zero records deleted")
	}
	return err
}
