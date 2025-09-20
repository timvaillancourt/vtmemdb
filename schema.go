package main

import "github.com/hashicorp/go-memdb"

const (
	tabletsTable = "tablets"
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
				tabletsHostnameIndex: {
					Name:    tabletsHostnameIndex,
					Unique:  false,
					Indexer: &memdb.StringFieldIndex{Field: "Hostname"},
				},
			},
		},
	},
}
