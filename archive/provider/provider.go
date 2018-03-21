package provider

import (
	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/EtienneLndr/MAL_API_Go_Project/data"

	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
)

// Define Provider's structure
type RetrieveProvider struct {
	ctx     *Context
	cctx    *ClientContext
	op      InvokeOperation
	factory EncodingFactory
}

// Allow to close the context of a specific provider
func (provider *RetrieveProvider) Close() {
	provider.ctx.Close()
}

// Create a provider
func CreateRetrieveProvider(url string, factory EncodingFactory) (*RetrieveProvider, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "providerRetrieve")
	if err != nil {
		return nil, err
	}

	op := cctx.NewInvokeOperation(cctx.Uri,
		ARCHIVE_SERVICE_AREA_NUMBER,
		ARCHIVE_SERVICE_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_RETRIEVE)

	provider := &RetrieveProvider{ctx, cctx, op, factory}

	return provider, nil
}

func (provider *RetrieveProvider) retrieveAck() error {
	return nil
}

func (provider *RetrieveProvider) retrieveResponse(archiveDetailsList ArchiveDetailsList, elementList ElementList) error {
	return nil
}

func (provider *RetrieveProvider) retrieveInvoke() (*ObjectType, *IdentifierList, *LongList, error) {
	return NullObjectType, NullIdentifierList, NullLongList, nil
}
