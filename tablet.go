package main

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-memdb"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

type tablet struct {
	*topodatapb.Tablet
	UpdatedAt time.Time `json:"updated_at"`
}

func SaveTablet(txn *memdb.Txn, t *topodatapb.Tablet) error {
	return txn.Insert(tabletsTable, &tablet{
		Tablet:    t,
		UpdatedAt: time.Now().UTC(),
	})
}

func ReadTablet(txn *memdb.Txn, tabletAlias *topodatapb.TabletAlias) (*topodatapb.Tablet, error) {
	res, err := txn.First(tabletsTable, tabletsAliasIndex, tabletAlias)
	if err != nil {
		return nil, err
	}
	t, ok := res.(*tablet)
	if !ok {
		return nil, fmt.Errorf("data must be *tablet, got %T", res)
	}
	return t.Tablet, nil
}

func ReadTabletsByHostname(txn *memdb.Txn, hostname string) ([]*topodatapb.Tablet, error) {
	// do prefix scan on tabletsHostnameMysqlPortIndex using hostname
	it, err := txn.Get(tabletsTable, tabletsHostnameMysqlPortIndex+"_prefix", hostname)
	if err != nil {
		return nil, err
	}
	tablets := make([]*topodatapb.Tablet, 0)
	for obj := it.Next(); obj != nil; obj = it.Next() {
		t, ok := obj.(*tablet)
		if !ok {
			return nil, fmt.Errorf("data must be *tablet, got %T", obj)
		}
		tablets = append(tablets, t.Tablet)
	}
	return tablets, nil
}

func ReadTabletByHostnameAndMysqlPort(txn *memdb.Txn, hostname string, mysqlPort int32) (*topodatapb.Tablet, error) {
	res, err := txn.First(tabletsTable, tabletsHostnameMysqlPortIndex, hostname, mysqlPort)
	if err != nil {
		return nil, err
	}
	t, ok := res.(*tablet)
	if !ok {
		return nil, fmt.Errorf("data must be *tablet, got %T", res)
	}
	return t.Tablet, nil
}
