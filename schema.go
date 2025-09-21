package main

import "github.com/hashicorp/go-memdb"

const (
	keyspacesTable = "keyspaces"
	shardsTable    = "shards"
	tabletsTable   = "tablets"

	primaryKeyIndex           = "id"
	keyspacesKeyspaceIndex    = primaryKeyIndex
	shardsKeyspaceShardIndex  = primaryKeyIndex
	tabletsAliasIndex         = primaryKeyIndex
	tabletsHostnamePortIndex  = "hostname_port"
	tabletsKeyspaceShardIndex = "keyspace_shard"
)

var dbSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		keyspacesTable: {
			Name: keyspacesTable,
			Indexes: map[string]*memdb.IndexSchema{
				keyspacesKeyspaceIndex: {
					Name:    keyspacesKeyspaceIndex,
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Keyspace"},
				},
			},
		},
		shardsTable: {
			Name: shardsTable,
			Indexes: map[string]*memdb.IndexSchema{
				shardsKeyspaceShardIndex: {
					Name:   shardsKeyspaceShardIndex,
					Unique: true,
					Indexer: &memdb.CompoundIndex{
						Indexes: []memdb.Indexer{
							&memdb.StringFieldIndex{Field: "Keyspace"},
							&memdb.StringFieldIndex{Field: "Shard"},
						},
					},
				},
			},
		},
		tabletsTable: {
			Name: tabletsTable,
			Indexes: map[string]*memdb.IndexSchema{
				tabletsAliasIndex: {
					Name:    tabletsAliasIndex,
					Unique:  true,
					Indexer: &TabletAliasIndexer{},
				},
				tabletsHostnamePortIndex: {
					Name:   tabletsHostnamePortIndex,
					Unique: true,
					Indexer: &memdb.CompoundIndex{
						Indexes: []memdb.Indexer{
							&memdb.StringFieldIndex{Field: "Hostname"},
							&memdb.IntFieldIndex{Field: "MysqlPort"},
						},
					},
				},
				tabletsKeyspaceShardIndex: {
					Name:   tabletsKeyspaceShardIndex,
					Unique: false,
					Indexer: &memdb.CompoundIndex{
						Indexes: []memdb.Indexer{
							&memdb.StringFieldIndex{Field: "Keyspace"},
							&memdb.StringFieldIndex{Field: "Shard"},
						},
					},
				},
			},
		},
	},
}
