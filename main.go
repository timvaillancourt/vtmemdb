package main

import (
	"fmt"

	"github.com/hashicorp/go-memdb"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

func main() {
	// Create a new data base
	db, err := memdb.NewMemDB(dbSchema)
	if err != nil {
		panic(err)
	}

	idxr := &TabletAliasIndexer{}
	idxBytes, err := idxr.FromArgs(&topodatapb.TabletAlias{Cell: "zone1", Uid: 127})
	if err != nil {
		panic(err)
	}
	fmt.Printf("idx key: %s\n", string(idxBytes))

	txn := db.Txn(true)
	for i := 0; i < 10000; i++ {
		t := &topodatapb.Tablet{
			Alias: &topodatapb.TabletAlias{
				Cell: "zone1",
				Uid:  uint32(i),
			},
			Hostname:  fmt.Sprintf("host%d.zone1.local", i),
			MysqlPort: 3306,
			Type:      topodatapb.TabletType_REPLICA,
		}
		if err := txn.Insert(tabletsTable, t); err != nil {
			panic(err)
		}
	}
	txn.Commit()

	txn = db.Txn(false)
	defer txn.Abort()

	res, err := txn.First(tabletsTable, tabletsAliasIndex, &topodatapb.TabletAlias{Cell: "zone1", Uid: 9999})
	if err != nil {
		panic(err)
	} else if res != nil {
		tablet, ok := res.(*topodatapb.Tablet)
		if !ok {
			panic(fmt.Errorf("data must be *topodatapb.Tablet, got %T", res))
		}
		fmt.Printf("res tablet: %+v\n", tablet)
	}

	res, err = txn.First(tabletsTable, tabletsHostnamePortIndex, "host123.zone1.local", int32(3306))
	if err != nil {
		panic(err)
	} else if res != nil {
		tablet, ok := res.(*topodatapb.Tablet)
		if !ok {
			panic(fmt.Errorf("data must be *topodatapb.Tablet, got %T", res))
		}
		fmt.Printf("res tablet: %+v\n", tablet)
	}
}
