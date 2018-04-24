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
				provider.retrieveAckError(transaction, MAL_ERROR_BAD_ENCODING, MAL_ERROR_BAD_ENCODING_MESSAGE, NewLongList(0))
				return err
			}

			// ----- Verify the parameters -----
			err = provider.retrieveVerifyParameters(transaction, objectType, identifierList, longList)
			if err != nil {
				return err
			}

			// ----- Call Ack operation -----
			err = provider.retrieveAck(transaction)
			if err != nil {
				provider.retrieveAckError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}

			// Hold on, wait a little
			time.Sleep(SLEEP_TIME * time.Millisecond)

			// TODO: do sthg with these objects
			fmt.Println("RetrieveHandler received:\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				longList)

			// Retrieve these objects in the archive
			archiveDetailsList, elementList, err := RetrieveInArchive(*objectType, *identifierList, *longList)
			if err != nil {
				if err.Error() == string(MAL_ERROR_UNKNOWN_MESSAGE) {
					provider.retrieveResponseError(transaction, MAL_ERROR_UNKNOWN, MAL_ERROR_UNKNOWN_MESSAGE, NewLongList(0))
				} else {
					provider.retrieveResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				}
				return err
			}

			// ----- Call Response operation -----
			err = provider.retrieveResponse(transaction, &archiveDetailsList, elementList)
			if err != nil {
				provider.retrieveResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE, NewLongList(0))
				return err
			}
		}

		return nil
	}

	// Register the handler
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

// VERIFY PARAMETERS
func (provider *Provider) retrieveVerifyParameters(transaction InvokeTransaction, objectType *ObjectType, identifierList *IdentifierList, longList *LongList) error {
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		fmt.Println(ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR)
		provider.retrieveAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR))
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(*identifierList)[i] == "*" {
			fmt.Println(ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR)
			provider.retrieveAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR))
		}
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

// ACK ERROR
func (provider *Provider) retrieveAckError(transaction InvokeTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	encoder, err := EncodeError(encoder, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Ack(encoder.Body(), true)
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
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE ERROR
func (provider *Provider) retrieveResponseError(transaction InvokeTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	encoder, err := EncodeError(encoder, errorNumber, errorComment, errorExtra)
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
			boolean, objectType, archiveQueryList, queryFilterList, err := provider.queryProgress(msg)
			if err != nil {
				provider.queryAckError(transaction, MAL_ERROR_BAD_ENCODING, MAL_ERROR_BAD_ENCODING_MESSAGE, NewLongList(0))
				return err
			}

			// ----- Verify the parameters -----
			// TODO: form a single query by combining ArchiveQueryList and QueryFilterList

			// ----- Call Ack operation -----
			err = provider.queryAck(transaction)
			if err != nil {
				provider.queryAckError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE, NewLongList(0))
				return err
			}

			// Hold on buddy, wait a little
			time.Sleep(SLEEP_TIME * time.Millisecond)

			// Variables to send to the consumer
			var objType *ObjectType
			var archDetList *ArchiveDetailsList
			var idList *IdentifierList
			var elementList ElementList

			// TODO: do sthg with these objects
			fmt.Println("QueryHandler received:\n\t>>>",
				boolean, "\n\t>>>",
				objectType, "\n\t>>>",
				archiveQueryList, "\n\t>>>",
				queryFilterList)

			for i := 0; i < archiveQueryList.Size()-1; i++ {
				// TODO: we'll have to change all of the following lines
				// Do a query to the archive
				if queryFilterList != nil {
					objType, archDetList, idList, elementList, err = QueryArchive(boolean, *objectType, *(*archiveQueryList)[i], queryFilterList.GetElementAt(i))
				} else {
					objType, archDetList, idList, elementList, err = QueryArchive(boolean, *objectType, *(*archiveQueryList)[i], nil)
				}
				if err != nil {
					// TODO: we may have to check if err is not an "UNKNOWN" error
					provider.queryUpdateError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
					return err
				}
				// Call Update operation
				err = provider.queryUpdate(transaction, objType, idList, archDetList, elementList)
				if err != nil {
					// TODO: we're (maybe) supposed to say to the consumer that an error occured
					provider.queryUpdateError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
					return err
				}
			}

			// Do a query to the archive
			if queryFilterList != nil {
				objType, archDetList, idList, elementList, err = QueryArchive(boolean, *objectType, *(*archiveQueryList)[archiveQueryList.Size()-1], queryFilterList.GetElementAt(archiveQueryList.Size()-1))
			} else {
				objType, archDetList, idList, elementList, err = QueryArchive(boolean, *objectType, *(*archiveQueryList)[archiveQueryList.Size()-1], nil)
			}
			if err != nil {
				// TODO: we may have to check if err is not an "UNKNOWN" error
				provider.queryUpdateError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}
			// ----- Call Response operation -----
			// Unless archive query list size is equal to 1 (we didn't enter in the for loop)
			err = provider.queryResponse(transaction, objType, idList, archDetList, elementList)
			if err != nil {
				// TODO: we're (maybe) supposed to say to the consumer that an error occured
				provider.queryResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}
		}

		return nil
	}

	// Register the handler
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

// VERIFY PARAMETERS

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
	queryFilterList, err := decoder.DecodeNullableAbstractElement()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if queryFilterList == nil {
		return boolean.(*Boolean), objectType.(*ObjectType), archiveQueryList.(*ArchiveQueryList), nil, nil
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

// ACK ERROR
func (provider *Provider) queryAckError(transaction ProgressTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	encoder, err := EncodeError(encoder, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Ack(encoder.Body(), true)
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

// UPDATE ERROR
func (provider *Provider) queryUpdateError(transaction ProgressTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	encoder, err := EncodeError(encoder, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Update(encoder.Body(), true)
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
	err := encoder.EncodeNullableElement(objectType)
	if err != nil {
		return err
	}

	// Encode IdentifierList
	err = encoder.EncodeNullableElement(identifierList)
	if err != nil {
		return err
	}

	// Encode ArchiveDetailsList
	err = encoder.EncodeNullableElement(archiveDetailsList)
	if err != nil {
		return err
	}

	// Encode ElementList
	err = encoder.EncodeNullableAbstractElement(elementList)
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

// RESPONSE ERROR
func (provider *Provider) queryResponseError(transaction ProgressTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	encoder, err := EncodeError(encoder, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Reply(encoder.Body(), true)
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

			// TODO: do sthg with these objects
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

	// Register the handler
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

// VERIFY PARAMETERS

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

// ACK ERROR

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

// RESPONSE ERROR

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

			// Call Request operation
			boolean, objectType, identifierList, archiveDetailsList, elementList, err := provider.storeRequest(msg)
			if err != nil {
				provider.storeResponseError(transaction, MAL_ERROR_BAD_ENCODING, MAL_ERROR_BAD_ENCODING_MESSAGE, NewLong(0))
				return err
			}

			// ----- Verify the parameters -----
			err = provider.storeVerifyParameters(transaction, boolean, objectType, identifierList, archiveDetailsList, elementList)
			if err != nil {
				return err
			}

			// Hold on, wait a little
			time.Sleep(SLEEP_TIME * time.Millisecond)

			// TODO: do sthg with these objects
			fmt.Println("StoreHandler received:\n\t>>>",
				boolean, "\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				archiveDetailsList, "\n\t>>>",
				elementList)

			// Store these objects in the archive
			var longList LongList
			longList, err = StoreInArchive(*objectType, *identifierList, *archiveDetailsList, elementList)
			if err != nil {
				if err.Error() == string(COM_ERROR_DUPLICATE) {
					provider.storeResponseError(transaction, COM_ERROR_DUPLICATE, COM_ERROR_DUPLICATE_MESSAGE, NewLongList(0))
				} else {
					provider.storeResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				}
				return err
			}

			// TODO: for each object stored, and 'ObjectStored' event may be published
			// to the event service

			// If boolean is false then we must send an empty LongList
			if !(*boolean) {
				longList = *NewLongList(0)
			}
			// Call Response operation
			err = provider.storeResponse(transaction, &longList)
			if err != nil {
				provider.storeResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}
		}

		return nil
	}

	// Register the handler
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
		fmt.Println(ARCHIVE_SERVICE_STORE_LIST_SIZE_ERROR)
		provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_LIST_SIZE_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_STORE_LIST_SIZE_ERROR))
	}

	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		fmt.Println(ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR)
		provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR))
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(*identifierList)[i] == "*" {
			fmt.Println(ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR)
			provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR))
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
			fmt.Println(ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR)
			provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR))
		}
	}

	// TODO: Raise INVALID error for 3.4.6.2.12

	return nil
}

// REQUEST
func (provider *Provider) storeRequest(msg *Message) (*Boolean, *ObjectType, *IdentifierList, *ArchiveDetailsList, ElementList, error) {
	// Create the decoder
	decoder := provider.factory.NewDecoder(msg.Body)

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

// RESPONSE ERROR
func (provider *Provider) storeResponseError(transaction RequestTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	encoder, err := EncodeError(encoder, errorNumber, errorComment, errorExtra)
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
				provider.updateAckError(transaction, MAL_ERROR_BAD_ENCODING, MAL_ERROR_BAD_ENCODING_MESSAGE, NewLong(0))
				return err
			}

			// ----- Verify the parameters -----
			err = provider.updateVerifyParameters(transaction, *objectType, *identifierList, *archiveDetailsList)
			if err != nil {
				return err
			}

			// TODO: do sthg with these objects
			fmt.Println("UpdateHandler received:\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				archiveDetailsList, "\n\t>>>",
				elementList)

			// Update these objects
			err = UpdateArchive(*objectType, *identifierList, *archiveDetailsList, elementList)
			if err != nil {
				if err.Error() == string(MAL_ERROR_UNKNOWN_MESSAGE) {
					provider.updateAckError(transaction, MAL_ERROR_UNKNOWN, MAL_ERROR_UNKNOWN_MESSAGE, NewLongList(0))
				} else {
					provider.updateAckError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				}
				return err
			}

			// Call Ack operation
			err = provider.updateAck(transaction)
			if err != nil {
				provider.updateAckError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}
		}

		return nil
	}

	// Register the handler
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

// VERIFY PARAMETERS
func (provider *Provider) updateVerifyParameters(transaction SubmitTransaction, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList) error {
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		fmt.Println(ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR)
		provider.updateAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR))
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(identifierList)[i] == "*" {
			fmt.Println(ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR)
			provider.updateAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR))
		}
	}

	// Verify object instance identifier
	for i := 0; i < archiveDetailsList.Size(); i++ {
		if archiveDetailsList[i].InstId == 0 {
			fmt.Println(ARCHIVE_SERVICE_AREA_OBJECT_INSTANCE_IDENTIFIER_VALUE_ERROR)
			provider.updateAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_AREA_OBJECT_INSTANCE_IDENTIFIER_VALUE_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_AREA_OBJECT_INSTANCE_IDENTIFIER_VALUE_ERROR))
		}
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

// ACK ERROR
func (provider *Provider) updateAckError(transaction SubmitTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	encoder, err := EncodeError(encoder, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Ack(encoder.Body(), true)
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
			err = provider.deleteVerifyParameters(transaction, *objectType, *identifierList)
			if err != nil {
				return err
			}

			// TODO: do sthg with these objects
			fmt.Println("DeleteHandler received:\n\t>>>",
				objectType, "\n\t>>>",
				identifierList, "\n\t>>>",
				longListRequest)

			// Delete these objects
			longListResponse, err := DeleteInArchive(*objectType, *identifierList, *longListRequest)
			if err != nil {
				if err.Error() == string(MAL_ERROR_UNKNOWN_MESSAGE) {
					provider.deleteResponseError(transaction, MAL_ERROR_UNKNOWN, MAL_ERROR_UNKNOWN_MESSAGE, NewLongList(0))
				} else {
					provider.deleteResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				}
				return err
			}

			// Call Response operation
			err = provider.deleteResponse(transaction, longListResponse)
			if err != nil {
				provider.deleteResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}
		}

		return nil
	}

	// Register the handler
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

// VERIFY PARAMETERS
func (provider *Provider) deleteVerifyParameters(transaction RequestTransaction, objectType ObjectType, identifierList IdentifierList) error {
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		fmt.Println(ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR)
		provider.deleteResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_STORE_OBJECTTYPE_VALUES_ERROR))
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(identifierList)[i] == "*" {
			fmt.Println(ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR)
			provider.deleteResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_STORE_IDENTIFIERLIST_VALUES_ERROR))
		}
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
func (provider *Provider) deleteResponse(transaction RequestTransaction, longList LongList) error {
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

// RESPONSE ERROR
func (provider *Provider) deleteResponseError(transaction RequestTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// Create the encoder
	encoder := provider.factory.NewEncoder(make([]byte, 0, 8192))

	encoder, err := EncodeError(encoder, errorNumber, errorComment, errorExtra)
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
