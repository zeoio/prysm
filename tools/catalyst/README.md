## Prysm Catalyst Setup


**Inside of the Prysm repository**, setup Catalyst configuration for a devnet
```text
export PRYSM_DIR=$(pwd)
"$PRYSM_DIR/scripts/catalyst.sh" init
```

Clone bazel-go-ethereum to your folder of choice with `git clone -b merge git@github.com:prysmaticlabs/bazel-go-ethereum`, then, setup go-ethereum for Catalyst mode.
```text
bazel run //cmd/geth -- \
 --datadir=/tmp/catalyst \
 init \
 "$PRYSM_DIR/tools/catalyst/eth1_config.json"
```

Run go-ethereum in catalyst mode
```text
bazel run //cmd/geth -- \
 --catalyst \
 --rpc \
 --rpcapi net,eth,consensus \
 --nodiscover \
 --miner.etherbase 0x1000000000000000000000000000000000000000 \
 --datadir=/tmp/catalyst
```

Run the Prysm beacon for Catalyst mode
```text
"$PRYSM_DIR/scripts/catalyst.sh" beacon-chain
```

Run the Prysm validator for Catalyst mode
```text
"$PRYSM_DIR/scripts/catalyst.sh" validator
```
