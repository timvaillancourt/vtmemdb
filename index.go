package main

import (
	"fmt"

	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
	"vitess.io/vitess/go/vt/topo/topoproto"
)

// TabletAliasIndexer is an indexer for *topodatapb.Tablet using the
// *topodatapb.TabletAlias ("Alias" field) as the index key.
type TabletAliasIndexer struct{}

// FromArgs satisfies the memdb.Indexer interface.
func (tai *TabletAliasIndexer) FromArgs(args ...interface{}) ([]byte, error) {
	if len(args) != 1 {
		return nil, fmt.Errorf("must provide only a single *topodatapb.TabletAlias argument")
	}
	tabletAlias, ok := args[0].(*topodatapb.TabletAlias)
	if !ok {
		return nil, fmt.Errorf("argument must be a *topodatapb.TabletAlias: %+v", args[0])
	}
	val := topoproto.TabletAliasString(tabletAlias)
	val += "\x00" // Add the null character as a terminator
	return []byte(val), nil
}

// FromObject satisfies the memdb.SingleIndexer interface.
func (tai *TabletAliasIndexer) FromObject(obj interface{}) (bool, []byte, error) {
	tablet, ok := obj.(*topodatapb.Tablet)
	if !ok {
		return false, nil, fmt.Errorf("object must be a *topodatapb.Tablet: %+v", obj)
	}
	val := topoproto.TabletAliasString(tablet.Alias)
	val += "\x00" // Add the null character as a terminator
	return true, []byte(val), nil
}
