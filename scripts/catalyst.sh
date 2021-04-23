#!/bin/bash

function color() {
    # Usage: color "31;5" "string"
    # Some valid values for color:
    # - 5 blink, 1 strong, 4 underlined
    # - fg: 31 red,  32 green, 33 yellow, 34 blue, 35 purple, 36 cyan, 37 white
    # - bg: 40 black, 41 red, 44 blue, 45 purple
    printf '\033[%sm%s\033[0m\n' "$@"
}

if [[ $1 == init ]]; then
    color "37" "Initializing catalyst"
    bazel run //tools/catalyst --define=ssz=minimal -- -base-path=$(pwd) -state-output=$(pwd)/tools/catalyst/genesis.ssz
fi

if [[ $1 == beacon-chain ]]; then
    color "37" "Launching beacon-chain"
    bazel run //beacon-chain --define=ssz=minimal -- --config-file=$(pwd)/tools/catalyst/beacon.config.yaml
fi

if [[ $1 == validator ]]; then
    color "37" "Launching validator"
    bazel run //validator --define=ssz=minimal -- --config-file=$(pwd)/tools/catalyst/validator.config.yaml
fi
