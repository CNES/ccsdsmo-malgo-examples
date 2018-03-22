package provider

import (
	"fmt"
	"time"

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
	factory EncodingFactory
}

// Allow to close the context of a specific provider
func (provider *RetrieveProvider) Close() {
	provider.ctx.Close()
}

func (provider *RetrieveProvider) retrieveAck(transaction InvokeTransaction) error {
	fmt.Println("Provider: retrieveAck")
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

func (provider *RetrieveProvider) retrieveResponse(archiveDetailsList *ArchiveDetailsList, elementList ElementList, transaction InvokeTransaction) error {
	fmt.Println("Provider: retrieveResponse")
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	err := archiveDetailsList.Encode(encoder)
	if err != nil {
		return err
	}

	err = elementList.Encode(encoder)
	if err != nil {
		return err
	}

	transaction.Reply(encoder.Body(), false)

	return nil
}

func (provider *RetrieveProvider) retrieveInvoke(msg *Message) (*ObjectType, *IdentifierList, *LongList, error) {
	fmt.Println("Provider: retrieveInvoke")
	decoder := provider.factory.NewDecoder(msg.Body)

	fmt.Println("Provider: retrieveInvoke -> ok")
	element, err := decoder.DecodeElement(NullObjectType)
	if err != nil {
		return nil, nil, nil, err
	}
	objectType := element.(*ObjectType)

	element, err = decoder.DecodeElement(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, err
	}
	identifierList := element.(*IdentifierList)

	element, err = decoder.DecodeElement(NullLongList)
	if err != nil {
		return nil, nil, nil, err
	}
	longList := element.(*LongList)

	return objectType, identifierList, longList, nil
}

// Create a provider
func createRetrieveProvider(url string, factory EncodingFactory) (*RetrieveProvider, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "providerRetrieve")
	if err != nil {
		return nil, err
	}

	provider := &RetrieveProvider{ctx, cctx, factory}

	return provider, nil
}

// Create retrieve handler
func (provider *RetrieveProvider) retrieveHandler() error {
	fmt.Println("Provider: retrieveHandler")
	retrieveHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			// Create Invoke Transaction
			transaction := t.(InvokeTransaction)

			// Call invoke operation and store objects
			objectType, identifierList, longList, err := provider.retrieveInvoke(msg)
			if err != nil {
				return err
			}

			// Hold on, wait a little
			time.Sleep(250 * time.Millisecond)

			// Call Ack operation
			provider.retrieveAck(transaction)

			// TODO (AF): do sthg with these objects
			fmt.Println("RetrieveHandler received:\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				longList)

			var archiveDetailsList = new(ArchiveDetailsList)
			var elementList = new(ArchiveQueryList)
			// Call Response operation
			provider.retrieveResponse(archiveDetailsList, elementList, transaction)
		}

		return nil
	}

	err := provider.cctx.RegisterInvokeHandler(SERVICE_AREA_NUMBER,
		SERVICE_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_RETRIEVE,
		retrieveHandler)
	if err != nil {
		return err
	}

	return nil
}

// Start :
func StartProvider(url string, factory EncodingFactory) (*RetrieveProvider, error) {
	// Create the provider
	provider, err := createRetrieveProvider(url, factory)
	if err != nil {
		return nil, err
	}

	// Create and launch the handler
	err = provider.retrieveHandler()
	if err != nil {
		return nil, err
	}

	return provider, nil
}
