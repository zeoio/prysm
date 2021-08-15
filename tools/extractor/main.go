package main

import (
	"context"
	"encoding/hex"
	"flag"
	"fmt"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stategen"

	types "github.com/prysmaticlabs/eth2-types"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state/interop"
	"github.com/prysmaticlabs/prysm/beacon-chain/db"
	"github.com/prysmaticlabs/prysm/beacon-chain/db/kv"
	"github.com/prysmaticlabs/prysm/shared/featureconfig"
)

var (
	// Required fields
	datadir = flag.String("datadir", "", "Path to data directory.")

	state = flag.Uint("state", 0, "Extract state at this slot.")
)

func main() {
	resetCfg := featureconfig.InitWithReset(&featureconfig.Flags{WriteSSZStateTransitions: true})
	defer resetCfg()
	flag.Parse()
	fmt.Println("Starting process...")
	d, err := db.NewDB(context.Background(), *datadir, &kv.Config{})
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	slot := types.Slot(*state)
	_, roots, err := d.BlockRootsBySlot(ctx, slot)
	if err != nil {
		panic(err)
	}
	if len(roots) != 1 {
		fmt.Printf("Expected 1 block root for slot %d, got %d roots", *state, len(roots))
	}
	var blockRoot [32]byte = roots[0]
	dst := make([]byte, hex.EncodedLen(len(blockRoot)))
	hex.Encode(dst, blockRoot[:])
	fmt.Printf("root for slot %d = %s\n", slot, string(dst))
	stateGen := stategen.New(d)
	s, err := stateGen.StateByRoot(ctx, blockRoot)
	if err != nil {
		panic(err)
	}
	/*
	s, err := d.State(ctx, stateRoot)
	if err != nil {
		panic(err)
	}
	*/

	interop.WriteStateToDisk(s)
	fmt.Println("done")
}
