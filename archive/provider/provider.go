/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package provider

import (
	"errors"
	"fmt"
	"sync"
	"time"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"
	. "github.com/ccsdsmo/malgo/mal/encoding/binary"

	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/archive/storage"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/errors"
)

var (
	ctx    *Context
	locker sync.Mutex
)

// Define Provider's structure
type Provider struct {
	ctx     *Context
	cctx    *ClientContext
	factory EncodingFactory
}

// Allow to close the context of a specific provider
func (provider *Provider) Close() {
	provider.ctx.Close()
}

// Create a provider
func createProvider(url string, typeOfProvider string) (*Provider, error) {
	// Declare variables
	var err error
	locker.Lock()
	if ctx == nil {
		ctx, err = NewContext(url)
		if err != nil {
			return nil, err
		}
	}
	locker.Unlock()

	cctx, err := NewClientContext(ctx, typeOfProvider)
	if err != nil {
		return nil, err
	}

	factory := new(FixedBinaryEncoding)

	provider := &Provider{ctx, cctx, factory}

	return provider, nil
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
// StartRetrieveProvider :
func StartRetrieveProvider(url string) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, "providerRetrieve")
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

// Create retrieve handler
func (provider *Provider) retrieveHandler() error {
	retrieveHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			// ----- Create Invoke Transaction -----
			transaction := t.(InvokeTransaction)

			// ----- Call invoke operation and store objects -----
			objectType, identifierList, longList, err := provider.retrieveInvoke(msg)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// ----- Verify the parameters -----
			// Verify ObjectType values (all of its attributes must not be equal to '0')
			if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
				// TODO: we must send an INVALID error to the consumer
				fmt.Println("ERROR: INVALID")
			}

			// Verify IdentifierList
			for i := 0; i < identifierList.Size(); i++ {
				if *(*identifierList)[i] == "*" {
					// TODO: we must send an INVALID error to the consumer
					fmt.Println("ERROR: INVALID")
				}
			}

			// Verify LongList
			for i := 0; i < identifierList.Size(); i++ {
				if *(*longList)[i] == 0 {
					// TODO: we must send an INVALID error to the consumer
					fmt.Println("ERROR: INVALID")
				}
			}

			// ----- Call Ack operation -----
			err = provider.retrieveAck(transaction)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// Hold on, wait a little
			time.Sleep(SLEEP_TIME * time.Millisecond)

			// TODO (AF): do sthg with these objects
			fmt.Println("RetrieveHandler received:\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				longList)

			var archiveDetailsList = new(ArchiveDetailsList)
			var elementList = new(ArchiveQueryList)
			// ----- Call Response operation -----
			err = provider.retrieveResponse(transaction, archiveDetailsList, elementList)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}
		}

		return nil
	}

	err := provider.cctx.RegisterInvokeHandler(COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_RETRIEVE,
		retrieveHandler)
	if err != nil {
		return err
	}

	return nil
}

// INVOKE
func (provider *Provider) retrieveInvoke(msg *Message) (*ObjectType, *IdentifierList, *LongList, error) {
	decoder := provider.factory.NewDecoder(msg.Body)

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

// ACK
func (provider *Provider) retrieveAck(transaction InvokeTransaction) error {
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

// RESPONSE
func (provider *Provider) retrieveResponse(transaction InvokeTransaction, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	err := archiveDetailsList.Encode(encoder)
	if err != nil {
		return err
	}

	err = encoder.EncodeAbstractElement(elementList)
	if err != nil {
		return err
	}

	transaction.Reply(encoder.Body(), false)

	return nil
}

//======================================================================//
//								QUERY									//
//======================================================================//
// StartRetrieveProvider :
func StartQueryProvider(url string) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, "providerQuery")
	if err != nil {
		return nil, err
	}

	// Create and launch the handler
	err = provider.queryHandler()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Create query handler
func (provider *Provider) queryHandler() error {
	queryHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(ProgressTransaction)

			// ----- Retrieve the objects thanks to the progress operation -----
			boolean, objectType, archiveQueryList, queryFilter, err := provider.queryProgress(msg)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// ----- Verify the parameters -----
			// TODO: form a single query from ArchiveQueryList and QueryFilterList

			// ----- Call Ack operation -----
			err = provider.queryAck(transaction)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// Hold on buddy, wait a little
			time.Sleep(SLEEP_TIME * time.Millisecond)

			// TODO: do sthg with these objects
			fmt.Println("QueryHandler received:\n\t>>>",
				boolean, "\n\t>>>",
				objectType, "\n\t>>>",
				archiveQueryList, "\n\t>>>",
				queryFilter)

			// TODO: This value will depend in the future of the number of objects to send to the consumer
			var nbObjects = 10
			// TODO: These variables will be created automatically in the future
			var objType = new(ObjectType)
			var archDetList = new(ArchiveDetailsList)
			var idList = new(IdentifierList)
			var elementList ElementList
			if *boolean != true {
				elementList = nil
				idList = nil
			} else {
				idList = new(IdentifierList)
				elementList = new(LongList)
			}

			for i := 0; i < nbObjects; i++ {
				// Call Update operation
				err = provider.queryUpdate(transaction, objType, idList, archDetList, elementList)
				if err != nil {
					// TODO: we're (maybe) supposed to say to the consumer that an error occured
					return err
				}
			}

			// ----- Call Response operation -----
			err = provider.queryResponse(transaction, objType, idList, archDetList, elementList)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}
		}

		return nil
	}

	err := provider.cctx.RegisterProgressHandler(COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_QUERY,
		queryHandler)
	if err != nil {
		return err
	}

	return nil
}

// PROGRESS
func (provider *Provider) queryProgress(msg *Message) (*Boolean, *ObjectType, *ArchiveQueryList, QueryFilterList, error) {
	// Create the decoder
	decoder := provider.factory.NewDecoder(msg.Body)

	// Decode Boolean
	boolean, err := decoder.DecodeElement(NullBoolean)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode ObjectType
	objectType, err := decoder.DecodeElement(NullObjectType)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode ArchiveQueryList
	archiveQueryList, err := decoder.DecodeElement(NullArchiveQueryList)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode QueryFilterList
	queryFilterList, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return boolean.(*Boolean), objectType.(*ObjectType), archiveQueryList.(*ArchiveQueryList), queryFilterList.(QueryFilterList), nil
}

// ACK
func (provider *Provider) queryAck(transaction ProgressTransaction) error {
	// Call Ack operation
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

// UPDATE
func (provider *Provider) queryUpdate(transaction ProgressTransaction, objectType *ObjectType, identifierList *IdentifierList, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode ObjectType
	err := objectType.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode IdentifierList
	err = identifierList.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode ArchiveDetailsList
	err = archiveDetailsList.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode ElementList
	err = encoder.EncodeAbstractElement(elementList)
	if err != nil {
		return err
	}

	// Call Update operation
	err = transaction.Update(encoder.Body(), false)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE
func (provider *Provider) queryResponse(transaction ProgressTransaction, objectType *ObjectType, identifierList *IdentifierList, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode ObjectType
	err := objectType.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode IdentifierList
	err = identifierList.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode ArchiveDetailsList
	err = archiveDetailsList.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode ElementList
	err = encoder.EncodeAbstractElement(elementList)
	if err != nil {
		return err
	}

	// Call Update operation
	err = transaction.Reply(encoder.Body(), false)
	if err != nil {
		return err
	}

	return nil
}

//======================================================================//
//								COUNT									//
//======================================================================//
// StartCountProvider :
func StartCountProvider(url string) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, "providerCount")
	if err != nil {
		return nil, err
	}

	// Create and launch the handler
	err = provider.countHandler()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Create count handler
func (provider *Provider) countHandler() error {
	countHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(InvokeTransaction)

			// Call Invoke operation
			objectType, archiveQueryList, queryFilterList, err := provider.countInvoke(msg)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// ----- Verify the parameters -----

			// Call Ack operation
			err = provider.retrieveAck(transaction)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// Hold on, wait a little
			time.Sleep(SLEEP_TIME * time.Millisecond)

			// TODO (AF): do sthg with these objects
			fmt.Println("CountHandler received:\n\t>>>",
				objectType, "\n\t>>>",
				archiveQueryList, "\n\t>>>",
				queryFilterList)

			// This variable will be created automatically in the future
			var longList = new(LongList)
			// Call Response operation
			err = provider.countResponse(transaction, longList)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}
		}

		return nil
	}

	err := provider.cctx.RegisterInvokeHandler(COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_COUNT,
		countHandler)
	if err != nil {
		return err
	}

	return nil
}

// INVOKE
func (provider *Provider) countInvoke(msg *Message) (*ObjectType, *ArchiveQueryList, QueryFilterList, error) {
	// Create the decoder
	decoder := provider.factory.NewDecoder(msg.Body)

	// Decode ObjectType
	objectType, err := decoder.DecodeElement(NullObjectType)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode ArchiveQueryList
	archiveQueryList, err := decoder.DecodeElement(NullArchiveQueryList)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode QueryFilterList
	queryFilterList, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, nil, err
	}

	return objectType.(*ObjectType), archiveQueryList.(*ArchiveQueryList), queryFilterList.(QueryFilterList), nil
}

// ACK
func (provider *Provider) countAck(transaction InvokeTransaction) error {
	// Call Ack operation
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

// RESPONSE
func (provider *Provider) countResponse(transaction InvokeTransaction, longList *LongList) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode LongList
	err := longList.Encode(encoder)
	if err != nil {
		return err
	}

	// Call Response operation
	err = transaction.Reply(encoder.Body(), false)
	if err != nil {
		return err
	}

	return nil
}

//======================================================================//
//								STORE									//
//======================================================================//
// StartStoreProvider :
func StartStoreProvider(url string) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, "providerStore")
	if err != nil {
		return nil, err
	}

	// Create and launch the handler
	err = provider.storeHandler()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Create store handler
func (provider *Provider) storeHandler() error {
	storeHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(RequestTransaction)

			fmt.Println("Store Handler")

			// Call Request operation
			boolean, objectType, identifierList, archiveDetailsList, elementList, err := provider.storeRequest(msg)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				fmt.Println(err)
				return err
			}

			// ----- Verify the parameters -----
			err = provider.storeVerifyParameters(transaction, boolean, objectType, identifierList, archiveDetailsList, elementList)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// Hold on, wait a little
			time.Sleep(SLEEP_TIME * time.Millisecond)

			// TODO (AF): do sthg with these objects
			fmt.Println("StoreHandler received:\n\t>>>",
				boolean, "\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				archiveDetailsList, "\n\t>>>",
				elementList)

			err = StoreInArchive(*objectType, *identifierList, *archiveDetailsList, elementList)
			if err != nil {
				return err
			}

			// This variable will be created automatically in the future
			var longList LongList
			if *boolean == true {
				longList = LongList([]*Long{NewLong(1), NewLong(2), NewLong(3)})
			} else {
				longList = *NewLongList(0)
			}
			// Call Response operation
			err = provider.storeResponse(transaction, &longList)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}
		}

		return nil
	}

	err := provider.cctx.RegisterRequestHandler(COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_STORE,
		storeHandler)
	if err != nil {
		return err
	}

	return nil
}

// VERIFY PARAMETERS
func (provider *Provider) storeVerifyParameters(transaction RequestTransaction, boolean *Boolean, objectType *ObjectType, identifierList *IdentifierList, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	// The fourth and fifth lists must be the same size
	if archiveDetailsList.Size() != elementList.Size() {
		fmt.Println("ArchiveDetailsList and ElementList must have the same size")
		provider.storeResponseError(transaction, COM_ERROR_INVALID, "ArchiveDetailsList and ElementList must have the same size", NewLongList(1))
		return errors.New("ArchiveDetailsList and ElementList must have the same size")
	}

	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		fmt.Println("ObjectType's attributes must not be equal to 'O'")
		provider.storeResponseError(transaction, COM_ERROR_INVALID, "ObjectType's attributes must not be equal to 'O'", NewLongList(1))
		return errors.New("ObjectType's attributes must not be equal to 'O'")
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(*identifierList)[i] == "*" {
			fmt.Println("IdenfierList elements must not be equal to '*'")
			provider.storeResponseError(transaction, COM_ERROR_INVALID, "IdenfierList elements must not be equal to '*'", NewLongList(1))
			return errors.New("IdenfierList elements must not be equal to '*'")
		}
	}

	// Verify the parameters network, timestamp and provider of the object ArchiveDetails
	mapNetwork := map[*Identifier]bool{
		NewIdentifier("0"): true,
		NewIdentifier("*"): true,
		nil:                true,
	}
	mapTimestamp := map[*FineTime]bool{
		NewFineTime(time.Unix(int64(0), int64(0))): true,
		nil: true,
	}
	mapProvider := map[*URI]bool{
		NewURI("0"): true,
		NewURI("*"): true,
		nil:         true,
	}
	for i := 0; i < archiveDetailsList.Size(); i++ {
		if mapNetwork[(*archiveDetailsList)[i].Network] || mapTimestamp[(*archiveDetailsList)[i].Timestamp] || mapProvider[(*archiveDetailsList)[i].Provider] {
			fmt.Println("ArchiveDetailsList elements must not be equal to '0', '*' or NULL")
			provider.storeResponseError(transaction, COM_ERROR_INVALID, "ArchiveDetailsList elements must not be equal to '0', '*' or NULL", NewLongList(1))
			return errors.New("ArchiveDetailsList elements must not be equal to '0', '*' or NULL")
		}
	}

	return nil
}

// REQUEST
func (provider *Provider) storeRequest(msg *Message) (*Boolean, *ObjectType, *IdentifierList, *ArchiveDetailsList, ElementList, error) {
	// Create the decoder
	decoder := provider.factory.NewDecoder(msg.Body)

	fmt.Println(msg.Body)

	// Decode Boolean
	boolean, err := decoder.DecodeElement(NullBoolean)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode ObjectType
	objectType, err := decoder.DecodeElement(NullObjectType)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode IdentifierList
	identifierList, err := decoder.DecodeElement(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode ArchiveDetailsList
	archiveDetailsList, err := decoder.DecodeElement(NullArchiveDetailsList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode ElementList
	elementList, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return boolean.(*Boolean), objectType.(*ObjectType), identifierList.(*IdentifierList), archiveDetailsList.(*ArchiveDetailsList), elementList.(ElementList), nil
}

// RESPONSE
func (provider *Provider) storeResponse(transaction RequestTransaction, longList *LongList) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode LongList
	err := longList.Encode(encoder)
	if err != nil {
		return err
	}

	// Call Response operation
	err = transaction.Reply(encoder.Body(), false)
	if err != nil {
		return err
	}

	return nil
}

func (provider *Provider) storeResponseError(transaction RequestTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode UInteger
	err := errorNumber.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode String
	err = errorComment.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode Element
	err = errorExtra.Encode(encoder)
	if err != nil {
		return err
	}

	// Call Response operation with Error status
	err = transaction.Reply(encoder.Body(), true)
	if err != nil {
		return err
	}

	return nil
}

//======================================================================//
//								UPDATE									//
//======================================================================//
// StartUpdateProvider :
func StartUpdateProvider(url string) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, "providerUpdate")
	if err != nil {
		return nil, err
	}

	// Create and launch the handler
	err = provider.updateHandler()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Create update handler
func (provider *Provider) updateHandler() error {
	updateHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(SubmitTransaction)

			// Call Submit operation
			objectType, identifierList, archiveDetailsList, elementList, err := provider.updateSubmit(msg)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// ----- Verify the parameters -----

			// Call Ack operation
			err = provider.updateAck(transaction)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// TODO (AF): do sthg with these objects
			fmt.Println("UpdateHandler received:\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				archiveDetailsList, "\n\t>>>",
				elementList)
		}

		return nil
	}

	err := provider.cctx.RegisterSubmitHandler(COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_UPDATE,
		updateHandler)
	if err != nil {
		return err
	}

	return nil
}

// SUBMIT
func (provider *Provider) updateSubmit(msg *Message) (*ObjectType, *IdentifierList, *ArchiveDetailsList, ElementList, error) {
	// Create the decoder
	decoder := provider.factory.NewDecoder(msg.Body)

	// Decode ObjectType
	objectType, err := decoder.DecodeElement(NullObjectType)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode IdentifierList
	identifierList, err := decoder.DecodeElement(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode ArchiveDetailsList
	archiveDetailsList, err := decoder.DecodeElement(NullArchiveDetailsList)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode ElementList
	elementList, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return objectType.(*ObjectType), identifierList.(*IdentifierList), archiveDetailsList.(*ArchiveDetailsList), elementList.(ElementList), nil
}

// ACK
func (provider *Provider) updateAck(transaction SubmitTransaction) error {
	// Call Ack operation
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

//======================================================================//
//								DELETE									//
//======================================================================//
// StartDeleteProvider :
func StartDeleteProvider(url string) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, "providerDelete")
	if err != nil {
		return nil, err
	}

	// Create and launch the handler
	err = provider.deleteHandler()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Create delete handler
func (provider *Provider) deleteHandler() error {
	deleteHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(RequestTransaction)

			// Call Request operation
			objectType, identifierList, longListRequest, err := provider.deleteRequest(msg)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}

			// ----- Verify the parameters -----

			// TODO (AF): do sthg with these objects
			fmt.Println("DeleteHandler received:\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				longListRequest)

			// Hold on dude, wait a little
			time.Sleep(SLEEP_TIME * time.Millisecond)

			// This variable will be created automatically in the future
			var longListResponse = new(LongList)
			// Call Response operation
			err = provider.deleteResponse(transaction, longListResponse)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				return err
			}
		}

		return nil
	}

	err := provider.cctx.RegisterRequestHandler(COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_DELETE,
		deleteHandler)
	if err != nil {
		return err
	}

	return nil
}

// REQUEST
func (provider *Provider) deleteRequest(msg *Message) (*ObjectType, *IdentifierList, *LongList, error) {
	// Create the decoder
	decoder := provider.factory.NewDecoder(msg.Body)

	// Decode ObjectType
	objectType, err := decoder.DecodeElement(NullObjectType)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode IdentifierList
	identifierList, err := decoder.DecodeElement(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode LongList
	longList, err := decoder.DecodeElement(NullLongList)
	if err != nil {
		return nil, nil, nil, err
	}

	return objectType.(*ObjectType), identifierList.(*IdentifierList), longList.(*LongList), nil
}

// RESPONSE
func (provider *Provider) deleteResponse(transaction RequestTransaction, longList *LongList) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode LongList
	err := longList.Encode(encoder)
	if err != nil {
		return err
	}

	// Call Response operation
	err = transaction.Reply(encoder.Body(), false)
	if err != nil {
		return err
	}

	return nil
}
