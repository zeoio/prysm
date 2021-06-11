package slasherkv

import (
	"context"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/prysmaticlabs/prysm/shared/testutil/require"
	"github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(ioutil.Discard)
	m.Run()
}

func TestInvestigate(t *testing.T) {
	ctx := context.Background()
	slasherDB, err := NewKVStore(ctx, "/home/raul_prysmaticlabs_com/themount", &Config{})
	require.NoError(t, err)
	err = slasherDB.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(attestedEpochsByValidator)
		fmt.Println("Attested epochs by val", bkt.Stats().KeyN)
	})
	require.NoError(t, err)
	err = slasherDB.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(attestationDataRootsBucket)
		fmt.Println("Data roots", bkt.Stats().KeyN)
	})
	require.NoError(t, err)
	err = slasherDB.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(proposalRecordsBucket)
		fmt.Println("Proposal records", bkt.Stats().KeyN)
	})
	require.NoError(t, err)
	err = slasherDB.db.View(func(tx *bolt.Tx) error {
		bkt := tx.Bucket(slasherChunksBucket)
		fmt.Println("Slasher chunks", bkt.Stats().KeyN)
	})
	require.NoError(t, err)
}
