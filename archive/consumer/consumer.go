package consumer

import (
	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/EtienneLndr/MAL_API_Go_Project/data"

	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
)

type ArchiveConsumer struct {
	ctx     *Context
	cctx    *ClientContext
	op      InvokeOperation
	factory EncodingFactory
}

// Allow to close the context of a specific consumer
func (consumer *ArchiveConsumer) Close() {
	consumer.ctx.Close()
}

// Create a consumer
func CreateConsumer(url string, factory EncodingFactory) (*ArchiveConsumer, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "consumer")
	if err != nil {
		return nil, err
	}

	op := cctx.NewInvokeOperation(cctx.Uri,
		ARCHIVE_SERVICE_AREA_NUMBER,
		ARCHIVE_SERVICE_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_RETRIEVE)

	consumer := &ArchiveConsumer{ctx, cctx, op, factory}

	return consumer, nil
}

func (consumer *ArchiveConsumer) retrieveInvoke(objectType ObjectType, identifierList IdentifierList, longList LongList) error {
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))
	objectType.Encode(encoder)
	identifierList.Encode(encoder)
	longList.Encode(encoder)

	_, err := consumer.op.Invoke(encoder.Body())
	if err != nil {
		return err
	}
	return nil
}

func (consumer *ArchiveConsumer) retrieveResponse() (*ArchiveDetailsList, *ElementList, error) {
	return NullArchiveDetailsList, nil, nil
}
