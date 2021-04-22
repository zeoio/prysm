package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/minio/sha256-simd"
	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/validator/keymanager/derived"
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip39"
	util "github.com/wealdtech/go-eth2-util"
)

const (
	mnemonic       = "lumber kind orange gold firm achieve tree robust peasant april very word ordinary before treat way ivory jazz cereal debate juice evil flame sadness"
	validatorCount = 64
	eth1Config     = "tools/catalyst/eth1_config.yaml"
	eth2Config     = "tools/catalyst/eth2_config.yaml"
)

var log = logrus.WithField("prefix", "catalyst-tool")

func main() {
	genesisTime := uint64(time.Now().Unix())
	eth1Genesis, err := loadEth1GenesisConf(eth1Config)
	if err != nil {
		panic(err)
	}
	eth1Genesis.Timestamp = genesisTime

	eth1Db := rawdb.NewMemoryDatabase()
	eth1GenesisBlock := eth1Genesis.ToBlock(eth1Db)
	params.LoadChainConfigFile(eth2Config)

	validators, err := loadValidatorKeys()
	if err != nil {
		panic(err)
	}

	if uint64(len(validators)) < params.BeaconConfig().MinGenesisActiveValidatorCount {
		log.Warnf(
			"Not enough validators for genesis - have %d total, but need %d",
			len(validators),
			params.BeaconConfig().MinGenesisActiveValidatorCount,
		)
	}

	eth1BlockHash := eth1GenesisBlock.Hash()

	beaconState, err := state.EmptyGenesisState()
	if err != nil {
		panic(err)
	}
	if err := beaconState.SetValidators(validators); err != nil {
		panic(err)
	}
	eth1Data := &ethpb.Eth1Data{
		DepositRoot:  make([]byte, 32),
		DepositCount: 0,
		BlockHash:    eth1BlockHash[:],
	}
	beaconState, err = state.OptimizedGenesisBeaconState(genesisTime, beaconState, eth1Data)
	if err != nil {
		panic(err)
	}
	if err := beaconState.SetGenesisTime(genesisTime + params.BeaconConfig().GenesisDelay); err != nil {
		panic(err)
	}
	if err := beaconState.SetLatestExecutionPayloadHeader(&pbp2p.ExecutionPayloadHeader{
		BlockHash:        eth1BlockHash[:],
		ParentHash:       eth1GenesisBlock.ParentHash().Bytes(),
		Coinbase:         eth1GenesisBlock.Coinbase().Bytes(),
		StateRoot:        eth1GenesisBlock.Root().Bytes(),
		Number:           eth1GenesisBlock.NumberU64(),
		GasLimit:         eth1GenesisBlock.GasLimit(),
		GasUsed:          eth1GenesisBlock.GasUsed(),
		Timestamp:        eth1GenesisBlock.Time(),
		ReceiptRoot:      eth1GenesisBlock.ReceiptHash().Bytes(),
		LogsBloom:        eth1GenesisBlock.Bloom().Bytes(),
		TransactionsRoot: make([]byte, 32),
	}); err != nil {
		panic(err)
	}

	t := beaconState.GenesisTime()
	log.Infof(
		"eth2 genesis at %d + %d = %d (%v)",
		genesisTime,
		params.BeaconConfig().GenesisDelay,
		t,
		time.Unix(int64(t), 0),
	)

	//fmt.Println("done preparing state, serializing SSZ now...")
	//f, err := os.OpenFile(g.StateOutputPath, os.O_CREATE|os.O_WRONLY, 0777)
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//buf := bufio.NewWriter(f)
	//defer buf.Flush()
	//w := codec.NewEncodingWriter(f)
	//if err := state.Serialize(w); err != nil {
	//	return err
	//}
	//fmt.Println("done!")

	// Generates a beacon.config.yaml for Prysm's beacon node and a validator.config.yaml for the validator
}

func loadEth1GenesisConf(configPath string) (*core.Genesis, error) {
	eth1ConfData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read eth1 config file: %v", err)
	}
	var eth1Genesis core.Genesis
	if err := json.NewDecoder(bytes.NewReader(eth1ConfData)).Decode(&eth1Genesis); err != nil {
		return nil, fmt.Errorf("failed to decode eth1 config file: %v", err)
	}
	return &eth1Genesis, nil
}

func loadValidatorKeys() ([]*ethpb.Validator, error) {
	// Uses the provided mnemonic seed phrase to generate the
	// appropriate seed file for recovering a derived wallets.
	if ok := bip39.IsMnemonicValid(mnemonic); !ok {
		panic(bip39.ErrInvalidMnemonic)
	}
	seed := bip39.NewSeed(strings.TrimSpace(mnemonic), "")
	validators := make([]*ethpb.Validator, validatorCount)
	for i := 0; i < validatorCount; i++ {
		signingKey, err := util.PrivateKeyFromSeedAndPath(
			seed, fmt.Sprintf(derived.ValidatingKeyDerivationPathTemplate, i),
		)
		if err != nil {
			return nil, errors.Wrap(err, "got bad withdrawal")
		}
		withdrawalKey, err := util.PrivateKeyFromSeedAndPath(
			seed, fmt.Sprintf(derived.WithdrawalKeyDerivationPathTemplate, i),
		)
		if err != nil {
			return nil, err
		}
		validators[i] = &ethpb.Validator{
			PublicKey:                  signingKey.PublicKey().Marshal(),
			ActivationEligibilityEpoch: 0,
			ActivationEpoch:            0,
			ExitEpoch:                  params.BeaconConfig().FarFutureEpoch,
			WithdrawableEpoch:          params.BeaconConfig().FarFutureEpoch,
		}
		h := sha256.New()
		h.Write(withdrawalKey.PublicKey().Marshal())
		validators[i].WithdrawalCredentials = h.Sum(nil)
		validators[i].WithdrawalCredentials[0] = params.BeaconConfig().BLSWithdrawalPrefixByte
		validators[i].EffectiveBalance = params.BeaconConfig().MaxEffectiveBalance
	}
	return validators, nil
}
