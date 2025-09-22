package main

import "github.com/hashicorp/go-memdb"

const (
	// keyspaceTable stores keyspace records
	keyspacesTable = "keyspaces"

	// shardsTable stores shard records
	shardsTable = "shards"

	// stateTable stores VTOrc state
	stateTable = "state"

	// tabletsTable stores tablet records
	tabletsTable = "tablets"
)

const (
	// primary key indexes
	primaryKeyIndex          = "id"
	keyspacesNameIndex       = primaryKeyIndex
	shardsKeyspaceShardIndex = primaryKeyIndex
	stateIndex               = primaryKeyIndex
	tabletsAliasIndex        = primaryKeyIndex

	// secondary indexes
	tabletsHostnameMysqlPortIndex = "hostname_mysqlPort"
	tabletsKeyspaceShardIndex     = "keyspace_shard"
)

var dbSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
		keyspacesTable: {
			Name: keyspacesTable,
			Indexes: map[string]*memdb.IndexSchema{
				keyspacesNameIndex: {
					Name:    keyspacesNameIndex,
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Name"},
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
							&memdb.StringFieldIndex{Field: "Name"},
						},
					},
				},
			},
		},
		stateTable: {
			Name: stateTable,
			Indexes: map[string]*memdb.IndexSchema{
				stateIndex: {
					Name:    stateIndex,
					Unique:  true,
					Indexer: &memdb.StringFieldIndex{Field: "Name"},
				},
			},
		},
		tabletsTable: {
			Name: tabletsTable,
			Indexes: map[string]*memdb.IndexSchema{
				tabletsAliasIndex: {
					Name:    tabletsAliasIndex,
					Unique:  true,
					Indexer: NewTabletAliasIndexer(),
				},
				tabletsHostnameMysqlPortIndex: {
					Name:   tabletsHostnameMysqlPortIndex,
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
