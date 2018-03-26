package consumer

import (
	"fmt"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/EtienneLndr/MAL_API_Go_Project/data"
)

type RetrieveConsumer struct {
	ctx     *Context
	cctx    *ClientContext
	op      InvokeOperation
	factory EncodingFactory
}

// Allow to close the context of a specific consumer
func (consumer *RetrieveConsumer) Close() {
	consumer.ctx.Close()
}

// Create a consumer
func createConsumer(url string, factory EncodingFactory, providerURI *URI) (*RetrieveConsumer, error) {
	fmt.Println("Consumer: createConsumer")
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "consumerRetrieve")
	if err != nil {
		return nil, err
	}

	op := cctx.NewInvokeOperation(providerURI,
		SERVICE_AREA_NUMBER,
		SERVICE_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_RETRIEVE)

	consumer := &RetrieveConsumer{ctx, cctx, op, factory}

	return consumer, nil
}

func (consumer *RetrieveConsumer) retrieveInvoke(objectType ObjectType, identifierList IdentifierList, elementList ElementList) error {
	fmt.Println("Consumer: retrieveInvoke")

	fmt.Println(objectType)
	fmt.Println(identifierList)
	fmt.Println(elementList)

	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))
	objectType.Encode(encoder)
	identifierList.Encode(encoder)
	encoder.EncodeAbstractElement(elementList)

	_, err := consumer.op.Invoke(encoder.Body())
	if err != nil {
		return err
	}
	return nil
}

func (consumer *RetrieveConsumer) retrieveResponse() (*ArchiveDetailsList, ElementList, error) {
	fmt.Println("Consumer: retrieveResponse")
	resp, err := consumer.op.GetResponse()
	if err != nil {
		return nil, nil, err
	}
	decoder := consumer.factory.NewDecoder(resp.Body)

	println("archivedetails")
	archiveDetails, err := decoder.DecodeElement(NullArchiveDetailsList)
	if err != nil {
		return nil, nil, err
	}

	println("elementlist")
	elementList, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, err
	}

	return archiveDetails.(*ArchiveDetailsList), elementList.(ElementList), nil
}

func StartConsumer(url string, factory EncodingFactory, providerURI *URI, objectType ObjectType, identifierList IdentifierList, elementList ElementList) (*RetrieveConsumer, *ArchiveDetailsList, ElementList, error) {
	// Create the consumer
	consumer, err := createConsumer(url, factory, providerURI)
	if err != nil {
		return nil, nil, nil, err
	}

	// Call Invoke operation
	err = consumer.retrieveInvoke(objectType, identifierList, elementList)
	if err != nil {
		return nil, nil, nil, err
	}

	// Call Response operation
	archiveDetailsList, elementList, err := consumer.retrieveResponse()
	if err != nil {
		return nil, nil, nil, err
	}

	return consumer, archiveDetailsList, elementList, nil
}
