package main

import (
	"testing"

	"github.com/stretchr/testify/require"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
)

func TestTabletAliasIndexer(t *testing.T) {
	tai := &TabletAliasIndexer{}

	// .FromArgs
	_, err := tai.FromArgs("this should fail")
	require.Error(t, err)

	fromArgsData, err := tai.FromArgs(&topodatapb.TabletAlias{
		Cell: "zone1",
		Uid:  123,
	})
	require.NoError(t, err)
	require.Equal(t, "zone1-0000000123\x00", string(fromArgsData))

	// .FromObject
	_, _, err = tai.FromObject(map[string]string{
		"this": "should fail",
	})
	require.Error(t, err)

	ok, fromObjectData, err := tai.FromObject(&topodatapb.Tablet{
		Alias:    &topodatapb.TabletAlias{Cell: "zone1", Uid: 123},
		Hostname: t.Name(),
		Type:     topodatapb.TabletType_REPLICA,
		Keyspace: "ks",
		Shard:    "-",
	})
	require.True(t, ok)
	require.Equal(t, "zone1-0000000123\x00", string(fromObjectData))
	require.NoError(t, err)

	// Check .FromArgs + .FromObject made the same key
	require.Equal(t, string(fromArgsData), string(fromObjectData))
}
