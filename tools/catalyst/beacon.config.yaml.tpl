datadir: /tmp/prysmcatalyst
force-clear-db: true
min-sync-peers: 0
http-web3provider: http://localhost:8545
bootstrap-node:
chain-config-file: {{ .ChainConfigPath }}
genesis-state: {{ .GenesisStatePath }}