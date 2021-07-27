package migration

import (
	v1Alpha1 "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	v2 "github.com/prysmaticlabs/prysm/proto/prysm/v2"
)


// V1Alpha1ToV2Attestation --
func V1Alpha1ToV2Attestation(src *v1Alpha1.Attestation) *v2.Attestation {
	if src == nil {
		return &v2.Attestation{}
	}
	return &v2.Attestation{
	
		AggregationBits: src.AggregationBits,
		Data: src.Data,
		Signature: src.Signature,
	}
}

// V2ToV1Alpha1Attestation --
func V2ToV1Alpha1Attestation(src *v2.Attestation) *v1Alpha1.Attestation {
	if src == nil {
		return &v1Alpha1.Attestation{}
	}
	return &v1Alpha1.Attestation{
	
		AggregationBits: src.AggregationBits,
		Data: src.Data,
		Signature: src.Signature,
	}
}

// V1Alpha1ToV2AttestationData --
func V1Alpha1ToV2AttestationData(src *v1Alpha1.AttestationData) *v2.AttestationData {
	if src == nil {
		return &v2.AttestationData{}
	}
	return &v2.AttestationData{
	
		Slot: src.Slot,
		CommitteeIndex: src.CommitteeIndex,
		BeaconBlockRoot: src.BeaconBlockRoot,
		Source: src.Source,
		Target: src.Target,
	}
}

// V2ToV1Alpha1AttestationData --
func V2ToV1Alpha1AttestationData(src *v2.AttestationData) *v1Alpha1.AttestationData {
	if src == nil {
		return &v1Alpha1.AttestationData{}
	}
	return &v1Alpha1.AttestationData{
	
		Slot: src.Slot,
		CommitteeIndex: src.CommitteeIndex,
		BeaconBlockRoot: src.BeaconBlockRoot,
		Source: src.Source,
		Target: src.Target,
	}
}

// V1Alpha1ToV2Checkpoint --
func V1Alpha1ToV2Checkpoint(src *v1Alpha1.Checkpoint) *v2.Checkpoint {
	if src == nil {
		return &v2.Checkpoint{}
	}
	return &v2.Checkpoint{
	
		Epoch: src.Epoch,
		Root: src.Root,
	}
}

// V2ToV1Alpha1Checkpoint --
func V2ToV1Alpha1Checkpoint(src *v2.Checkpoint) *v1Alpha1.Checkpoint {
	if src == nil {
		return &v1Alpha1.Checkpoint{}
	}
	return &v1Alpha1.Checkpoint{
	
		Epoch: src.Epoch,
		Root: src.Root,
	}
}

// V1Alpha1ToV2Genesis --
func V1Alpha1ToV2Genesis(src *v1Alpha1.Genesis) *v2.Genesis {
	if src == nil {
		return &v2.Genesis{}
	}
	return &v2.Genesis{
	
		GenesisTime: src.GenesisTime,
		DepositContractAddress: src.DepositContractAddress,
		GenesisValidatorsRoot: src.GenesisValidatorsRoot,
	}
}

// V2ToV1Alpha1Genesis --
func V2ToV1Alpha1Genesis(src *v2.Genesis) *v1Alpha1.Genesis {
	if src == nil {
		return &v1Alpha1.Genesis{}
	}
	return &v1Alpha1.Genesis{
	
		GenesisTime: src.GenesisTime,
		DepositContractAddress: src.DepositContractAddress,
		GenesisValidatorsRoot: src.GenesisValidatorsRoot,
	}
}

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

