package factory

import (
	"time"

	wsDriverFactory "github.com/multiversx/mx-chain-core-go/webSockets/factory"
	indexerFactory "github.com/multiversx/mx-chain-es-indexer-go/process/factory"
	"github.com/multiversx/mx-chain-go/outport"
)

// WrappedOutportDriverWebSocketsSenderFactoryArgs extends the wsDriverFactory.OutportDriverWebSocketSenderFactoryArgs structure with the Enabled field
type WrappedOutportDriverWebSocketsSenderFactoryArgs struct {
	Enabled bool
	wsDriverFactory.ArgsWebSocketsDriverFactory
}

// OutportFactoryArgs holds the factory arguments of different outport drivers
type OutportFactoryArgs struct {
	RetrialInterval                   time.Duration
	ElasticIndexerFactoryArgs         indexerFactory.ArgsIndexerFactory
	EventNotifierFactoryArgs          *EventNotifierFactoryArgs
	WebSocketsSenderDriverFactoryArgs WrappedOutportDriverWebSocketsSenderFactoryArgs
}

// CreateOutport will create a new instance of OutportHandler
func CreateOutport(args *OutportFactoryArgs) (outport.OutportHandler, error) {
	err := checkArguments(args)
	if err != nil {
		return nil, err
	}

	outportHandler, err := outport.NewOutport(args.RetrialInterval)
	if err != nil {
		return nil, err
	}

	err = createAndSubscribeDrivers(outportHandler, args)
	if err != nil {
		return nil, err
	}

	return outportHandler, nil
}

func createAndSubscribeDrivers(outport outport.OutportHandler, args *OutportFactoryArgs) error {
	err := createAndSubscribeElasticDriverIfNeeded(outport, args.ElasticIndexerFactoryArgs)
	if err != nil {
		return err
	}

	err = createAndSubscribeEventNotifierIfNeeded(outport, args.EventNotifierFactoryArgs)
	if err != nil {
		return err
	}

	return createAndSubscribeWebSocketDriver(outport, args.WebSocketsSenderDriverFactoryArgs)
}

func createAndSubscribeElasticDriverIfNeeded(
	outport outport.OutportHandler,
	args indexerFactory.ArgsIndexerFactory,
) error {
	if !args.Enabled {
		return nil
	}

	elasticDriver, err := indexerFactory.NewIndexer(args)
	if err != nil {
		return err
	}

	return outport.SubscribeDriver(elasticDriver)
}

func createAndSubscribeEventNotifierIfNeeded(
	outport outport.OutportHandler,
	args *EventNotifierFactoryArgs,
) error {
	if !args.Enabled {
		return nil
	}

	eventNotifier, err := CreateEventNotifier(args)
	if err != nil {
		return err
	}

	return outport.SubscribeDriver(eventNotifier)
}

func checkArguments(args *OutportFactoryArgs) error {
	if args == nil {
		return outport.ErrNilArgsOutportFactory
	}

	return nil
}

func createAndSubscribeWebSocketDriver(
	outport outport.OutportHandler,
	args WrappedOutportDriverWebSocketsSenderFactoryArgs,
) error {
	if !args.Enabled {
		return nil
	}

	wsFactory, err := wsDriverFactory.NewWebSocketsDriverFactory(args.ArgsWebSocketsDriverFactory)
	if err != nil {
		return err
	}

	wsDriver, err := wsFactory.Create()
	if err != nil {
		return err
	}

	return outport.SubscribeDriver(wsDriver)
}
