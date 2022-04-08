package staking

import (
	"math/big"

	"github.com/ElrondNetwork/elrond-go-core/core"
	"github.com/ElrondNetwork/elrond-go-core/data"
	"github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-go/dataRetriever"
	"github.com/ElrondNetwork/elrond-go/epochStart"
	"github.com/ElrondNetwork/elrond-go/epochStart/metachain"
	factory2 "github.com/ElrondNetwork/elrond-go/factory"
	mock2 "github.com/ElrondNetwork/elrond-go/integrationTests/mock"
	"github.com/ElrondNetwork/elrond-go/process"
	blproc "github.com/ElrondNetwork/elrond-go/process/block"
	"github.com/ElrondNetwork/elrond-go/process/block/bootstrapStorage"
	"github.com/ElrondNetwork/elrond-go/process/block/postprocess"
	"github.com/ElrondNetwork/elrond-go/process/mock"
	"github.com/ElrondNetwork/elrond-go/sharding"
	"github.com/ElrondNetwork/elrond-go/sharding/nodesCoordinator"
	"github.com/ElrondNetwork/elrond-go/state"
	"github.com/ElrondNetwork/elrond-go/testscommon"
	"github.com/ElrondNetwork/elrond-go/testscommon/dblookupext"
)

func createMetaBlockProcessor(
	nc nodesCoordinator.NodesCoordinator,
	systemSCProcessor process.EpochStartSystemSCProcessor,
	coreComponents factory2.CoreComponentsHolder,
	dataComponents factory2.DataComponentsHolder,
	bootstrapComponents factory2.BootstrapComponentsHolder,
	statusComponents factory2.StatusComponentsHolder,
	stateComponents factory2.StateComponentsHandler,
	validatorsInfoCreator process.ValidatorStatisticsProcessor,
	blockChainHook process.BlockChainHookHandler,
	metaVMFactory process.VirtualMachinesContainerFactory,
	epochStartHandler process.EpochStartTriggerHandler,
) process.BlockProcessor {
	arguments := createMockMetaArguments(coreComponents, dataComponents, bootstrapComponents, statusComponents, nc, systemSCProcessor, stateComponents, validatorsInfoCreator, blockChainHook, metaVMFactory, epochStartHandler)

	metaProc, _ := blproc.NewMetaProcessor(arguments)
	return metaProc
}

func createMockMetaArguments(
	coreComponents factory2.CoreComponentsHolder,
	dataComponents factory2.DataComponentsHolder,
	bootstrapComponents factory2.BootstrapComponentsHolder,
	statusComponents factory2.StatusComponentsHolder,
	nodesCoord nodesCoordinator.NodesCoordinator,
	systemSCProcessor process.EpochStartSystemSCProcessor,
	stateComponents factory2.StateComponentsHandler,
	validatorsInfoCreator process.ValidatorStatisticsProcessor,
	blockChainHook process.BlockChainHookHandler,
	metaVMFactory process.VirtualMachinesContainerFactory,
	epochStartHandler process.EpochStartTriggerHandler,
) blproc.ArgMetaProcessor {
	shardCoordiantor := bootstrapComponents.ShardCoordinator()
	valInfoCreator := createValidatorInfoCreator(coreComponents, dataComponents, shardCoordiantor)
	blockTracker := createBlockTracker(dataComponents.Blockchain().GetGenesisHeader(), shardCoordiantor)
	epochStartDataCreator := createEpochStartDataCreator(coreComponents, dataComponents, shardCoordiantor, epochStartHandler, blockTracker)

	accountsDb := make(map[state.AccountsDbIdentifier]state.AccountsAdapter)
	accountsDb[state.UserAccountsState] = stateComponents.AccountsAdapter()
	accountsDb[state.PeerAccountsState] = stateComponents.PeerAccounts()

	bootStorer, _ := bootstrapStorage.NewBootstrapStorer(coreComponents.InternalMarshalizer(), dataComponents.StorageService().GetStorer(dataRetriever.BootstrapUnit))
	headerValidator := createHeaderValidator(coreComponents)
	vmContainer, _ := metaVMFactory.Create()
	return blproc.ArgMetaProcessor{
		ArgBaseProcessor: blproc.ArgBaseProcessor{
			CoreComponents:                 coreComponents,
			DataComponents:                 dataComponents,
			BootstrapComponents:            bootstrapComponents,
			StatusComponents:               statusComponents,
			AccountsDB:                     accountsDb,
			ForkDetector:                   &mock2.ForkDetectorStub{},
			NodesCoordinator:               nodesCoord,
			FeeHandler:                     postprocess.NewFeeAccumulator(),
			RequestHandler:                 &testscommon.RequestHandlerStub{},
			BlockChainHook:                 blockChainHook,
			TxCoordinator:                  &mock.TransactionCoordinatorMock{},
			EpochStartTrigger:              epochStartHandler,
			HeaderValidator:                headerValidator,
			GasHandler:                     &mock.GasHandlerMock{},
			BootStorer:                     bootStorer,
			BlockTracker:                   blockTracker,
			BlockSizeThrottler:             &mock.BlockSizeThrottlerStub{},
			HistoryRepository:              &dblookupext.HistoryRepositoryStub{},
			EpochNotifier:                  coreComponents.EpochNotifier(),
			RoundNotifier:                  &mock.RoundNotifierStub{},
			ScheduledTxsExecutionHandler:   &testscommon.ScheduledTxsExecutionStub{},
			ScheduledMiniBlocksEnableEpoch: 10000,
			VMContainersFactory:            metaVMFactory,
			VmContainer:                    vmContainer,
		},
		SCToProtocol:                 &mock.SCToProtocolStub{},
		PendingMiniBlocksHandler:     &mock.PendingMiniBlocksHandlerStub{},
		EpochStartDataCreator:        epochStartDataCreator,
		EpochEconomics:               &mock.EpochEconomicsStub{},
		EpochRewardsCreator:          &testscommon.RewardsCreatorStub{},
		EpochValidatorInfoCreator:    valInfoCreator,
		ValidatorStatisticsProcessor: validatorsInfoCreator,
		EpochSystemSCProcessor:       systemSCProcessor,
	}
}

func createValidatorInfoCreator(
	coreComponents factory2.CoreComponentsHolder,
	dataComponents factory2.DataComponentsHolder,
	shardCoordinator sharding.Coordinator,
) process.EpochStartValidatorInfoCreator {
	args := metachain.ArgsNewValidatorInfoCreator{
		ShardCoordinator: shardCoordinator,
		MiniBlockStorage: dataComponents.StorageService().GetStorer(dataRetriever.MiniBlockUnit),
		Hasher:           coreComponents.Hasher(),
		Marshalizer:      coreComponents.InternalMarshalizer(),
		DataPool:         dataComponents.Datapool(),
	}

	valInfoCreator, _ := metachain.NewValidatorInfoCreator(args)
	return valInfoCreator
}

func createEpochStartDataCreator(
	coreComponents factory2.CoreComponentsHolder,
	dataComponents factory2.DataComponentsHolder,
	shardCoordinator sharding.Coordinator,
	epochStartTrigger process.EpochStartTriggerHandler,
	blockTracker process.BlockTracker,
) process.EpochStartDataCreator {
	argsEpochStartDataCreator := metachain.ArgsNewEpochStartData{
		Marshalizer:       coreComponents.InternalMarshalizer(),
		Hasher:            coreComponents.Hasher(),
		Store:             dataComponents.StorageService(),
		DataPool:          dataComponents.Datapool(),
		BlockTracker:      blockTracker,
		ShardCoordinator:  shardCoordinator,
		EpochStartTrigger: epochStartTrigger,
		RequestHandler:    &testscommon.RequestHandlerStub{},
		GenesisEpoch:      0,
	}
	epochStartDataCreator, _ := metachain.NewEpochStartData(argsEpochStartDataCreator)
	return epochStartDataCreator
}

func createBlockTracker(genesisMetaHeader data.HeaderHandler, shardCoordinator sharding.Coordinator) process.BlockTracker {
	genesisBlocks := make(map[uint32]data.HeaderHandler)
	for ShardID := uint32(0); ShardID < shardCoordinator.NumberOfShards(); ShardID++ {
		genesisBlocks[ShardID] = createGenesisBlock(ShardID)
	}

	genesisBlocks[core.MetachainShardId] = genesisMetaHeader
	return mock.NewBlockTrackerMock(shardCoordinator, genesisBlocks)
}

func createGenesisBlock(ShardID uint32) *block.Header {
	rootHash := []byte("roothash")
	return &block.Header{
		Nonce:           0,
		Round:           0,
		Signature:       rootHash,
		RandSeed:        rootHash,
		PrevRandSeed:    rootHash,
		ShardID:         ShardID,
		PubKeysBitmap:   rootHash,
		RootHash:        rootHash,
		PrevHash:        rootHash,
		AccumulatedFees: big.NewInt(0),
		DeveloperFees:   big.NewInt(0),
	}
}

func createGenesisMetaBlock() *block.MetaBlock {
	rootHash := []byte("roothash")
	return &block.MetaBlock{
		Nonce:                  0,
		Round:                  0,
		Signature:              rootHash,
		RandSeed:               rootHash,
		PrevRandSeed:           rootHash,
		PubKeysBitmap:          rootHash,
		RootHash:               rootHash,
		PrevHash:               rootHash,
		AccumulatedFees:        big.NewInt(0),
		DeveloperFees:          big.NewInt(0),
		AccumulatedFeesInEpoch: big.NewInt(0),
		DevFeesInEpoch:         big.NewInt(0),
	}
}

func createHeaderValidator(coreComponents factory2.CoreComponentsHolder) epochStart.HeaderValidator {
	argsHeaderValidator := blproc.ArgsHeaderValidator{
		Hasher:      coreComponents.Hasher(),
		Marshalizer: coreComponents.InternalMarshalizer(),
	}
	headerValidator, _ := blproc.NewHeaderValidator(argsHeaderValidator)
	return headerValidator
}
