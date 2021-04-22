package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/google/uuid"
	"github.com/minio/sha256-simd"
	ethpb "github.com/prysmaticlabs/ethereumapis/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/state"
	pbp2p "github.com/prysmaticlabs/prysm/proto/beacon/p2p/v1"
	"github.com/prysmaticlabs/prysm/shared/fileutil"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/validator/keymanager/derived"
	"github.com/prysmaticlabs/prysm/validator/keymanager/imported"
	"github.com/sirupsen/logrus"
	"github.com/tyler-smith/go-bip39"
	types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
	keystorev4 "github.com/wealdtech/go-eth2-wallet-encryptor-keystorev4"
)

const (
	mnemonic       = "lumber kind orange gold firm achieve tree robust peasant april very word ordinary before treat way ivory jazz cereal debate juice evil flame sadness"
	validatorCount = 64
	eth1Config     = "tools/catalyst/eth1_config.yaml"
	eth2Config     = "tools/catalyst/eth2_config.yaml"
)

var (
	log             = logrus.WithField("prefix", "catalyst-tool")
	basePathFlag    = flag.String("base-path", "", "Base path for Prysm")
	stateOutputFlag = flag.String("state-output", "", "State output path")
)

type accountStore struct {
	PrivateKeys [][]byte `json:"private_keys"`
	PublicKeys  [][]byte `json:"public_keys"`
}

func main() {
	flag.Parse()
	base, err := fileutil.ExpandPath(*basePathFlag)
	if err != nil {
		panic(err)
	}
	genesisTime := uint64(time.Now().Unix())
	eth1Genesis, err := loadEth1GenesisConf(filepath.Join(base, eth1Config))
	if err != nil {
		panic(err)
	}
	eth1Genesis.Timestamp = genesisTime

	eth1Db := rawdb.NewMemoryDatabase()
	eth1GenesisBlock := eth1Genesis.ToBlock(eth1Db)
	params.LoadChainConfigFile(filepath.Join(base, eth2Config))

	validators, privKeys, pubKeys, err := loadValidatorKeys()
	if err != nil {
		panic(err)
	}
	accounts := &accountStore{
		PrivateKeys: privKeys,
		PublicKeys:  pubKeys,
	}
	encodedStore, err := json.MarshalIndent(accounts, "", "\t")
	if err != nil {
		panic(err)
	}
	encryptor := keystorev4.New()
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	cryptoFields, err := encryptor.Encrypt(encodedStore, "foobar")
	if err != nil {
		panic(err)
	}
	ks := &imported.AccountsKeystoreRepresentation{
		Crypto:  cryptoFields,
		ID:      id.String(),
		Version: encryptor.Version(),
		Name:    encryptor.Name(),
	}
	encJSON, err := json.MarshalIndent(ks, "", "\t")
	if err != nil {
		panic(err)
	}
	if err := fileutil.WriteFile(filepath.Join(base, "tools", "catalyst", "wallet", "direct", "accounts", "all-accounts.keystore.json"), encJSON); err != nil {
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

	f, err := os.OpenFile(*stateOutputFlag, os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			panic(err)
		}
	}()
	enc, err := beaconState.MarshalSSZ()
	if err != nil {
		panic(err)
	}
	buf := bufio.NewWriter(f)
	defer func() {
		if err := buf.Flush(); err != nil {
			panic(err)
		}
	}()
	n, err := buf.Write(enc)
	if err != nil {
		panic(err)
	}
	if n != len(enc) {
		panic("Not equal length")
	}
	log.Info("Done!")
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

func loadValidatorKeys() (validators []*ethpb.Validator, privKeys, pubKeys [][]byte, err error) {
	// Uses the provided mnemonic seed phrase to generate the
	// appropriate seed file for recovering a derived wallets.
	if ok := bip39.IsMnemonicValid(mnemonic); !ok {
		panic(bip39.ErrInvalidMnemonic)
	}
	seed := bip39.NewSeed(strings.TrimSpace(mnemonic), "")
	validators = make([]*ethpb.Validator, validatorCount)
	privKeys = make([][]byte, validatorCount)
	pubKeys = make([][]byte, validatorCount)
	for i := 0; i < validatorCount; i++ {
		var signingKey *types.BLSPrivateKey
		signingKey, err = util.PrivateKeyFromSeedAndPath(
			seed, fmt.Sprintf(derived.ValidatingKeyDerivationPathTemplate, i),
		)
		if err != nil {
			return
		}
		var withdrawalKey *types.BLSPrivateKey
		withdrawalKey, err = util.PrivateKeyFromSeedAndPath(
			seed, fmt.Sprintf(derived.WithdrawalKeyDerivationPathTemplate, i),
		)
		if err != nil {
			return
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

		privKeys[i] = signingKey.Marshal()
		pubKeys[i] = signingKey.PublicKey().Marshal()
	}
	return
}
