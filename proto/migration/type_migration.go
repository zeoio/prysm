package migration

import (
	v1Alpha1 "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	v2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
)

// V1Alpha1ToV2SyncStatus --
func V1Alpha1ToV2SyncStatus(src *v1Alpha1.SyncStatus) *v2.SyncStatus {
	if src == nil {
		return &v2.SyncStatus{}
	}
	return &v2.SyncStatus{

		Syncing: src.Syncing,
	}
}

// V2ToV1Alpha1SyncStatus --
func V2ToV1Alpha1SyncStatus(src *v2.SyncStatus) *v1Alpha1.SyncStatus {
	if src == nil {
		return &v1Alpha1.SyncStatus{}
	}
	return &v1Alpha1.SyncStatus{

		Syncing: src.Syncing,
	}
}

// V1Alpha1ToV2Genesis --
func V1Alpha1ToV2Genesis(src *v1Alpha1.Genesis) *v2.Genesis {
	if src == nil {
		return &v2.Genesis{}
	}
	return &v2.Genesis{

		GenesisTime:            src.GenesisTime,
		DepositContractAddress: src.DepositContractAddress,
		GenesisValidatorsRoot:  src.GenesisValidatorsRoot,
	}
}

// V2ToV1Alpha1Genesis --
func V2ToV1Alpha1Genesis(src *v2.Genesis) *v1Alpha1.Genesis {
	if src == nil {
		return &v1Alpha1.Genesis{}
	}
	return &v1Alpha1.Genesis{

		GenesisTime:            src.GenesisTime,
		DepositContractAddress: src.DepositContractAddress,
		GenesisValidatorsRoot:  src.GenesisValidatorsRoot,
	}
}
