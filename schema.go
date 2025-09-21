package main

import "github.com/hashicorp/go-memdb"

const (
	tabletsTable = "tablets"

	tabletsAliasIndex        = "id"
	tabletsHostnamePortIndex = "hostname_port"
)

var dbSchema = &memdb.DBSchema{
	Tables: map[string]*memdb.TableSchema{
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
			},
		},
	},
}
