// Package slasher implements slashing detection for eth2, able to catch slashable attestations
// and proposals that it receives via two event feeds, respectively. Any found slashings
// are then submitted to the beacon node's slashing operations pool. See the design document
// here https://hackmd.io/@prysmaticlabs/slasher.
package slasher

import (
	"context"
	"time"

	"github.com/prysmaticlabs/prysm/beacon-chain/blockchain"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/feed"
	statefeed "github.com/prysmaticlabs/prysm/beacon-chain/core/feed/state"
	"github.com/prysmaticlabs/prysm/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/beacon-chain/db"
	"github.com/prysmaticlabs/prysm/beacon-chain/operations/slashings"
	slashertypes "github.com/prysmaticlabs/prysm/beacon-chain/slasher/types"
	"github.com/prysmaticlabs/prysm/beacon-chain/state/stategen"
	"github.com/prysmaticlabs/prysm/beacon-chain/sync"
	ethpb "github.com/prysmaticlabs/prysm/proto/eth/v1alpha1"
	"github.com/prysmaticlabs/prysm/proto/interfaces"
	"github.com/prysmaticlabs/prysm/shared/attestationutil"
	"github.com/prysmaticlabs/prysm/shared/bytesutil"
	"github.com/prysmaticlabs/prysm/shared/event"
	"github.com/prysmaticlabs/prysm/shared/params"
	"github.com/prysmaticlabs/prysm/shared/slotutil"
)

// ServiceConfig for the slasher service in the beacon node.
// This struct allows us to specify required dependencies and
// parameters for slasher to function as needed.
type ServiceConfig struct {
	IndexedAttestationsFeed *event.Feed
	BeaconBlockHeadersFeed  *event.Feed
	Database                db.SlasherDatabase
	StateNotifier           statefeed.Notifier
	AttestationStateFetcher blockchain.AttestationStateFetcher
	StateGen                stategen.StateManager
	SlashingPoolInserter    slashings.PoolInserter
	HeadStateFetcher        blockchain.HeadFetcher
	SyncChecker             sync.Checker
}

// SlashingChecker is an interface for defining services that the beacon node may interact with to provide slashing data.
type SlashingChecker interface {
	IsSlashableBlock(ctx context.Context, proposal *ethpb.SignedBeaconBlockHeader) (*ethpb.ProposerSlashing, error)
	IsSlashableAttestation(ctx context.Context, attestation *ethpb.IndexedAttestation) ([]*ethpb.AttesterSlashing, error)
}

// Service defining a slasher implementation as part of
// the beacon node, able to detect eth2 slashable offenses.
type Service struct {
	params                 *Parameters
	serviceCfg             *ServiceConfig
	indexedAttsChan        chan *ethpb.IndexedAttestation
	beaconBlockHeadersChan chan *ethpb.SignedBeaconBlockHeader
	attsQueue              *attestationsQueue
	blksQueue              *blocksQueue
	ctx                    context.Context
	cancel                 context.CancelFunc
	slotTicker             *slotutil.SlotTicker
	genesisTime            time.Time
}

// New instantiates a new slasher from configuration values.
func New(ctx context.Context, srvCfg *ServiceConfig) (*Service, error) {
	ctx, cancel := context.WithCancel(ctx)
	return &Service{
		params:                 DefaultParams(),
		serviceCfg:             srvCfg,
		indexedAttsChan:        make(chan *ethpb.IndexedAttestation, 1),
		beaconBlockHeadersChan: make(chan *ethpb.SignedBeaconBlockHeader, 1),
		attsQueue:              newAttestationsQueue(),
		blksQueue:              newBlocksQueue(),
		ctx:                    ctx,
		cancel:                 cancel,
	}, nil
}

// Start listening for received indexed attestations and blocks
// and perform slashing detection on them.
func (s *Service) Start() {
	go s.run()
}

func (s *Service) run() {
	stateChannel := make(chan *feed.Event, 1)
	stateSub := s.serviceCfg.StateNotifier.StateFeed().Subscribe(stateChannel)
	stateEvent := <-stateChannel

	// Wait for us to receive the genesis time via a chain started notification.
	if stateEvent.Type == statefeed.ChainStarted {
		data, ok := stateEvent.Data.(*statefeed.ChainStartedData)
		if !ok {
			log.Error("Could not receive chain start notification, want *statefeed.ChainStartedData")
			return
		}
		s.genesisTime = data.StartTime
		log.WithField("genesisTime", s.genesisTime).Info("Starting slasher, received chain start event")
	} else if stateEvent.Type == statefeed.Initialized {
		// Alternatively, if the chain has already started, we then read the genesis
		// time value from this data.
		data, ok := stateEvent.Data.(*statefeed.InitializedData)
		if !ok {
			log.Error("Could not receive chain start notification, want *statefeed.ChainStartedData")
			return
		}
		s.genesisTime = data.StartTime
		log.WithField("genesisTime", s.genesisTime).Info("Starting slasher, chain already initialized")
	} else {
		// This should not happen.
		log.Error("Could start slasher, could not receive chain start event")
		return
	}

	stateSub.Unsubscribe()
	secondsPerSlot := params.BeaconConfig().SecondsPerSlot
	s.slotTicker = slotutil.NewSlotTicker(s.genesisTime, secondsPerSlot)

	// Wait until the beacon node is synced (short cirtuits if genesis epoch).
	s.waitForSync(s.ctx, s.genesisTime)

	// Perform a backfilling process for data from the synced epoch, N, down to N - WEAK_SUBJECTIVITY_PERIOD.
	s.backfillSlasherHistory(s.ctx)

	log.Info("Completed chain sync, starting slashing detection")
	go s.processQueuedAttestations(s.ctx, s.slotTicker.C())
	go s.processQueuedBlocks(s.ctx, s.slotTicker.C())
	go s.receiveAttestations(s.ctx)
	go s.receiveBlocks(s.ctx)
	go s.pruneSlasherData(s.ctx, s.slotTicker.C())
}

// Stop the slasher service.
func (s *Service) Stop() error {
	s.cancel()
	if s.slotTicker != nil {
		s.slotTicker.Done()
	}
	return nil
}

// Status of the slasher service.
func (s *Service) Status() error {
	return nil
}

func (s *Service) backfillSlasherHistory(ctx context.Context) {
	// Perform backfilling in chunks of N blocks at a time
	blocksOnDisk := make([]interfaces.BeaconBlock, 0)
	blockWrappers := make([]*slashertypes.SignedBlockHeaderWrapper, 0)
	attWrappers := make([]*slashertypes.IndexedAttestationWrapper, 0)
	for _, blk := range blocksOnDisk {
		preState, err := s.serviceCfg.StateGen.StateByRoot(ctx, bytesutil.ToBytes32(blk.ParentRoot()))
		if err != nil {
			return
		}
		if preState == nil || preState.IsNil() {
			return
		}
		for _, att := range blk.Body().Attestations() {
			committee, err := helpers.BeaconCommitteeFromState(preState, att.Data.Slot, att.Data.CommitteeIndex)
			if err != nil {
				log.WithError(err).Error("Could not get attestation committee")
				return
			}
			// Using a different context to prevent timeouts as this operation can be expensive
			// and we want to avoid affecting the critical code path.
			indexedAtt, err := attestationutil.ConvertToIndexed(ctx, att, committee)
			if err != nil {
				log.WithError(err).Error("Could not convert to indexed attestation")
				return
			}
			signingRoot, err := att.Data.HashTreeRoot()
			if err != nil {
				log.WithError(err).Error("Could not get hash tree root of attestation")
				continue
			}
			attWrapper := &slashertypes.IndexedAttestationWrapper{
				IndexedAttestation: indexedAtt,
				SigningRoot:        signingRoot,
			}
			attWrappers = append(attWrappers, attWrapper)
		}
	}

	// TODO: Need to also save the blocks and attestations in slasher's database.

	propSlashings, err := s.detectProposerSlashings(ctx, blockWrappers)
	if err != nil {
		log.WithError(err).Error("Could not detect proposer slashings on backfilled data")
		return
	}
	attSlashings, err := s.checkSlashableAttestations(ctx, attWrappers)
	if err != nil {
		log.WithError(err).Error("Could not detect proposer slashings on backfilled data")
		return
	}
	if err := s.processProposerSlashings(ctx, propSlashings); err != nil {
		log.WithError(err).Error("Could not process proposer slashings")
		return
	}
	if err := s.processAttesterSlashings(ctx, attSlashings); err != nil {
		log.WithError(err).Error("Could not process proposer slashings")
		return
	}
}

func (s *Service) waitForSync(ctx context.Context, genesisTime time.Time) {
	if slotutil.SlotsSinceGenesis(genesisTime) == 0 || !s.serviceCfg.SyncChecker.Syncing() {
		return
	}
	for {
		select {
		case <-s.slotTicker.C():
			// If node is still syncing, do not operate slasher.
			if s.serviceCfg.SyncChecker.Syncing() {
				continue
			}
			return
		case <-ctx.Done():
			return
		}
	}
}
