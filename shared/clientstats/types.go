package clientstats

const (
	ClientName            = "prysm"
	ProcessNameBeaconNode = "beaconnode"
	ProcessNameValidator  = "validator"
	ProcessNameSystem     = "system"
	APIVersion            = 1
)

// APIMessage are common to all requests to the client-stats API
// Note that there is a "system" type that we do not currently
// support -- if we did APIMessage would be present on the system
// messages as well as validator and beaconnode, whereas
// CommonStats would only be part of beaconnode and validator.
type APIMessage struct {
	APIVersion  int    `json:"version"`
	Timestamp   int64  `json:"timestamp"` // unix timestamp in milliseconds
	ProcessName string `json:"process"`   // validator, beaconnode, system
}

// CommonStats represent generic metrics that are expected on both
// beaconnode and validator metric types. This type is used for
// marshaling metrics to the POST body sent to the metrics collcetor.
// Note that some metrics are labeled NA because they are expected
// to be present with their zero-value when not supported by a client.
type CommonStats struct {
	CPUProcessSecondsTotal int64  `json:"cpu_process_seconds_total"`
	MemoryProcessBytes     int64  `json:"memory_process_bytes"`
	ClientName             string `json:"client_name"`
	ClientVersion          string `json:"client_version"`
	ClientBuild            int64  `json:"client_build"`
	// TODO(#8849): parse the grpc connection string to determine
	// if multiple addresses are present
	SyncEth2FallbackConfigured bool `json:"sync_eth2_fallback_configured"`
	// N/A -- when multiple addresses are provided to grpc, requests are
	// load-balanced between the provided endpoints.
	// This is different from a "fallback" configuration where
	// the second address is treated as a failover.
	SyncEth2FallbackConnected bool `json:"sync_eth2_fallback_connected"`
	APIMessage                `json:",inline"`
}

// BeaconNodeStats embeds CommonStats and represents metrics specific to
// the beacon-node process. This type is used to marshal metrics data
// to the POST body sent to the metrics collcetor. To make the connection
// to client-stats clear, BeaconNodeStats is also used by prometheus
// collection code introduced to support client-stats.
// Note that some metrics are labeled NA because they are expected
// to be present with their zero-value when not supported by a client.
type BeaconNodeStats struct {
	// TODO(#8850): add support for this after slasher refactor is merged
	SlasherActive              bool  `json:"slasher_active"`
	SyncEth1FallbackConfigured bool  `json:"sync_eth1_fallback_configured"`
	SyncEth1FallbackConnected  bool  `json:"sync_eth1_fallback_connected"`
	SyncEth1Connected          bool  `json:"sync_eth1_connected"`
	SyncEth2Synced             bool  `json:"sync_eth2_synced"`
	DiskBeaconchainBytesTotal  int64 `json:"disk_beaconchain_bytes_total"`
	// N/A -- would require significant network code changes at this time
	NetworkLibp2pBytesTotalReceive int64 `json:"network_libp2p_bytes_total_receive"`
	// N/A -- would require significant network code changes at this time
	NetworkLibp2pBytesTotalTransmit int64 `json:"network_libp2p_bytes_total_transmit"`
	// p2p_peer_count where label "state" == "Connected"
	NetworkPeersConnected int64 `json:"network_peers_connected"`
	// beacon_head_slot
	SyncBeaconHeadSlot int64 `json:"sync_beacon_head_slot"`
	CommonStats        `json:",inline"`
}

// ValidatorStats embeds CommonStats and represents metrics specific to
// the validator process. This type is used to marshal metrics data
// to the POST body sent to the metrics collcetor.
// Note that some metrics are labeled NA because they are expected
// to be present with their zero-value when not supported by a client.
type ValidatorStats struct {
	// N/A -- TODO(#8848): verify whether we can obtain this metric from the validator process
	ValidatorTotal int64 `json:"validator_total"`
	// N/A -- TODO(#8848): verify whether we can obtain this metric from the validator process
	ValidatorActive int64 `json:"validator_active"`
	CommonStats     `json:",inline"`
}

// SystemStats scrapes the prometheus node-exporter to
// collect generic system-level metrics.
type SystemStats struct {
	// sort of a hack, but i noticed node_cpu_core_throttles_total is the
	// only metric that corresponds to the number of cores, ie matches
	// cpu cores	: 4
	// in /proc/cpuinfo on my system
	CPUCores int64 `json:"cpu_cores"`
	// for this one we can count the distinct 'cpu=i' labels on
	// node_cpu_seconds_total (which we already need to loop through to find
	// the following cpu timing metrics
	CPUThreads int64 `json:"cpu_threads"`

	// node_cpu_seconds_total{cpu="7",mode="system"} 1995.43
	CPUNodeSystemSecondsTotal int64 `json:"cpu_node_system_seconds_total"`
	// node_cpu_seconds_total{cpu="7",mode="user"} 7925.8
	CPUNodeUserSecondsTotal int64 `json:"cpu_node_user_seconds_total"`
	// node_cpu_seconds_total{cpu="7",mode="iowait"} 105.01
	CPUNodeIOWaitSecondsTotal int64 `json:"cpu_node_iowait_seconds_total"`
	// node_cpu_seconds_total{cpu="7",mode="idle"} 63785.17
	CPUNodeIdleSecondsTotal int64 `json:"cpu_node_idle_seconds_total"`

	// node_memory_MemTotal_bytes
	MemoryNodeBytesTotal int64 `json:"memory_node_bytes_total"`
	// node_memory_MemFree_bytes
	MemoryNodeBytesFree int64 `json:"memory_node_bytes_free"`
	// node_memory_Cached_bytes
	MemoryNodeBytesCached int64 `json:"memory_node_bytes_cached"`
	// node_memory_Buffers_bytes
	MemoryNodeBytesBuffers int64 `json:"memory_node_bytes_buffers"`

	// these are tricky because they can contain multiple filesystems and often some virtual ones
	// node_filesystem_size_bytes
	DiskNodeBytesTotal int64 `json:"disk_node_bytes_total"`
	// node_filesystem_free_bytes ?
	DiskNodeBytesFree int64 `json:"disk_node_bytes_free"`

	// tricky because multiple filesystems
	// node_disk_io_time_seconds_total
	DiskNodeIOSeconds int64 `json:"disk_node_io_seconds"`
	// node_disk_reads_completed_total
	DiskNodeReadsTotal int64 `json:"disk_node_reads_total"`
	// node_disk_writes_completed_total
	DiskNodeWritesTotal int64 `json:"disk_node_writes_total"`

	// counter node_network_receive_bytes_total
	// eg node_network_receive_bytes_total{device="wlp2s0"} 1.6432525497e+10
	// note: should we aggregate across devices? this is tricky because docker is a network device
	NetworkNodeBytesTotalReceive int64 `json:"network_node_bytes_total_receive"`
	// counter: node_network_transmit_bytes_total{device="wlp2s0"} 1.0545479739e+10
	NetworkNodeBytesTotalTransmit int64 `json:"network_node_bytes_total_transmit"`

	// node_boot_time_seconds?
	MiscNodeBootTimeTsSeconds int64 `json:"misc_node_boot_ts_seconds"`
	// pull 'sysname' label from node_uname_info gauge, eg
	// node_uname_info{sysname="Linux"} 1
	// note this looks to be uname -s not uname -o
	// Linux=lin, Windows=win, Mac=macos or 'unk' for unknown
	MiscOS string `json:"misc_os"`

	APIMessage `json:",inline"`
}
