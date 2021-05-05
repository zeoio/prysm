package shard

import "github.com/prysmaticlabs/prysm/shared/params"

// UpdatedGasPrice returns the updated gas price based on the EIP 1599 formulas.
//
// Spec code:
// def compute_updated_gasprice(prev_gasprice: Gwei, shard_block_length: uint64, adjustment_quotient: uint64) -> Gwei:
//    if shard_block_length > TARGET_SAMPLES_PER_BLOCK:
//        delta = max(1, prev_gasprice * (shard_block_length - TARGET_SAMPLES_PER_BLOCK)
//                       // TARGET_SAMPLES_PER_BLOCK // adjustment_quotient)
//        return min(prev_gasprice + delta, MAX_GASPRICE)
//    else:
//        delta = max(1, prev_gasprice * (TARGET_SAMPLES_PER_BLOCK - shard_block_length)
//                       // TARGET_SAMPLES_PER_BLOCK // adjustment_quotient)
//        return max(prev_gasprice, MIN_GASPRICE + delta) - delta

func UpdatedGasPrice(prevGasPrice uint64, shardBlockLength uint64, adjustmentQuotient uint64) uint64 {
	targetBlockSize := params.BeaconConfig().TargetShardBlockSize
	maxGasPrice := params.BeaconConfig().MaxGasPrice
	minGasPrice := params.BeaconConfig().MinGasPrice
	// Delta can't be more than 1.
	delta := uint64(1)
	if shardBlockLength > targetBlockSize {
		if delta > prevGasPrice*(shardBlockLength-targetBlockSize)/targetBlockSize/adjustmentQuotient {
			delta = prevGasPrice * (shardBlockLength - targetBlockSize) / targetBlockSize / adjustmentQuotient
		}
		// Max gas price is the upper bound.
		if prevGasPrice+delta > maxGasPrice {
			return maxGasPrice
		}
		return prevGasPrice + delta
	}
	if delta > prevGasPrice*(targetBlockSize-shardBlockLength)/targetBlockSize/adjustmentQuotient {
		delta = prevGasPrice * (targetBlockSize - shardBlockLength) / targetBlockSize / adjustmentQuotient
	}

	// Min gas price is the lower bound.
	if prevGasPrice < minGasPrice+delta {
		return minGasPrice
	}
	return prevGasPrice - delta
}
