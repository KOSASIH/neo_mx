package interceptors

import (
	blockData "github.com/ElrondNetwork/elrond-go-sandbox/data/block"
	"github.com/ElrondNetwork/elrond-go-sandbox/hashing"
	"github.com/ElrondNetwork/elrond-go-sandbox/marshal"
	"github.com/ElrondNetwork/elrond-go-sandbox/p2p"
	"github.com/ElrondNetwork/elrond-go-sandbox/process/block"
	"github.com/ElrondNetwork/elrond-go-sandbox/sharding"
	"github.com/ElrondNetwork/elrond-go-sandbox/storage"
)

// TxBlockBodyInterceptor represents an interceptor used for transaction block bodies
type TxBlockBodyInterceptor struct {
	*GenericBlockBodyInterceptor
}

// NewMiniBlocksInterceptor creates a new instance of a TxBlockBodyInterceptor
func NewMiniBlocksInterceptor(
	marshalizer marshal.Marshalizer,
	cache storage.Cacher,
	storer storage.Storer,
	hasher hashing.Hasher,
	shardCoordinator sharding.Coordinator,
) (*TxBlockBodyInterceptor, error) {
	gbbi, err := NewGenericBlockBodyInterceptor(
		marshalizer,
		cache,
		storer,
		hasher,
		shardCoordinator)

	if err != nil {
		return nil, err
	}

	return &TxBlockBodyInterceptor{GenericBlockBodyInterceptor: gbbi}, nil
}

// ProcessReceivedMessage will be the callback func from the p2p.Messenger and will be called each time a new message was received
// (for the topic this validator was registered to)
func (tbbi *TxBlockBodyInterceptor) ProcessReceivedMessage(message p2p.MessageP2P) error {
	err := tbbi.checkMessage(message)
	if err != nil {
		return err
	}

	txBlockBody := block.NewInterceptedTxBlockBody()
	miniBlocks := make([]*blockData.MiniBlock, 0)
	x := message.Data()
	err = tbbi.marshalizer.Unmarshal(&miniBlocks, x)
	if err != nil {
		return err
	}
	txBlockBody.TxBlockBody = miniBlocks

	return tbbi.processBlockBody(message.Data(), txBlockBody)
}