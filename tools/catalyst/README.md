## Prysm Catalyst Setup


Clone bazel-go-ethereum to your folder of choice
`git clone -b merge git@github.com:prysmaticlabs/bazel-go-ethereum`

**Inside of the Prysm repository**, setup Catalyst configuration for a devnet
```text
./scripts/catalyst.sh init
```

**Inside of the bazel-go-ethereum repository**, setup go-ethereum for Catalyst mode
```text
bazel run //cmd/geth -- \
 --datadir=/tmp/catalyst \
 init \
 ./tools/catalyst/eth1_config.json
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

**Inside of the Prysm repository**, run the Prysm beacon for Catalyst mode
```text
./scripts/catalyst.sh beacon-chain
```

**Inside of the Prysm repository**, run the Prysm validator for Catalyst mode
```text
./scripts/catalyst.sh validator
```
