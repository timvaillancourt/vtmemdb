package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

type shardRecord struct {
	*topodatapb.Shard

	Name     string `json:"shard"`
	Keyspace string `json:"keyspace"`
}

func SaveShard(txn *memdb.Txn, keyspaceName, shardName string, s *topodatapb.Shard) error {
	return txn.Insert(shardsTable, &shardRecord{
		Name:     shardName,
		Keyspace: keyspaceName,
		Shard:    s,
	})
}

func ReadShard(txn *memdb.Txn, keyspaceName, shardName string) (*topodatapb.Shard, error) {
	res, err := txn.First(shardsTable, shardsKeyspaceShardIndex, keyspaceName, shardName)
	if err != nil {
		return nil, err
	}
	s, ok := res.(*shardRecord)
	if !ok {
		return nil, fmt.Errorf("data must be *shardRecord, got %T", res)
	}
	return s.Shard, nil
}

func DeleteShard(txn *memdb.Txn, keyspaceName, shardName string) error {
	deleted, err := txn.DeleteAll(shardsTable, shardsKeyspaceShardIndex, keyspaceName, shardName)
	if err == nil && deleted == 0 {
		err = fmt.Errorf("zero records deleted")
	}
	return err
}
