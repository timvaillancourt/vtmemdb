package main

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/hashicorp/go-memdb"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
	"vitess.io/vitess/go/vt/vtctl/reparentutil/policy"
)

func main() {
	// Create a new data base
	db, err := memdb.NewMemDB(dbSchema)
	if err != nil {
		panic(err)
	}

	// Populate table
	keyspace := "ks"
	shard := "-"
	wTxn := db.Txn(true)

	if err := SaveKeyspace(wTxn, keyspace, &topodatapb.Keyspace{
		DurabilityPolicy: policy.DurabilityCrossCell,
	}); err != nil {
		panic(err)
	}

	if err := SaveShard(wTxn, keyspace, shard, &topodatapb.Shard{
		PrimaryAlias: &topodatapb.TabletAlias{
			Cell: "zone1",
			Uid:  0,
		},
	}); err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		tabletType := topodatapb.TabletType_REPLICA
		if i == 0 {
			tabletType = topodatapb.TabletType_PRIMARY
		}
		t := &topodatapb.Tablet{
			Alias: &topodatapb.TabletAlias{
				Cell: "zone1",
				Uid:  uint32(i),
			},
			Keyspace:  keyspace,
			Shard:     shard,
			Hostname:  fmt.Sprintf("host%d.zone1.local", i),
			MysqlPort: 3306,
			Type:      tabletType,
		}
		if err := SaveTablet(wTxn, t); err != nil {
			panic(err)
		}
	}
	wTxn.Commit()

	// Read txn
	txn := db.Txn(false)
	defer txn.Abort()

	k, err := ReadKeyspace(txn, keyspace)
	if err != nil {
		panic(err)
	}
	fmt.Printf("res keyspace: %+v\n", k)

	s, err := ReadShard(txn, keyspace, shard)
	if err != nil {
		panic(err)
	}
	fmt.Printf("res shard: %+v\n", s)

	tablet, err := ReadTablet(txn, &topodatapb.TabletAlias{Cell: "zone1", Uid: 9})
	if err != nil {
		panic(err)
	}
	fmt.Printf("res tablet: %+v\n", tablet)

	tablets, err := ReadTabletsByHostname(txn, "host3.zone1.local")
	if err != nil {
		panic(err)
	}
	fmt.Printf("res tablets: %+v\n", tablets)

	tablet, err = ReadTabletByHostnameAndMysqlPort(txn, "host5.zone1.local", int32(3306))
	if err != nil {
		panic(err)
	}
	fmt.Printf("res tablet: %+v\n", tablet)

	wTxn = db.Txn(true)
	if err := SetGlobalRecoveryDisabled(wTxn); err != nil {
		panic(err)
	}
	wTxn.Commit()

	txn = db.Txn(false)

	recoveryDisabled, err := IsGlobalRecoveryDisabled(txn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("res recoveryDisabled: %+v\n", recoveryDisabled)

	ds, err := GetDatabaseState(txn)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(ds))
}

func GetDatabaseState(txn *memdb.Txn) ([]byte, error) {
	output := make(map[string]any, 0)
	for tableName := range maps.Keys(dbSchema.Tables) {
		items := make([]any, 0)
		it, err := txn.Get(tableName, primaryKeyIndex)
		if err != nil {
			return nil, err
		}
		for obj := it.Next(); obj != nil; obj = it.Next() {
			items = append(items, obj)
		}
		output[tableName] = items
	}
	return json.MarshalIndent(output, "", "    ")
}
