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
	"fmt"
	"time"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/EtienneLndr/MAL_API_Go_Project/data"
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
func createProvider(url string, factory EncodingFactory, typeOfProvider string) (*Provider, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, typeOfProvider)
	if err != nil {
		return nil, err
	}

	provider := &Provider{ctx, cctx, factory}

	return provider, nil
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
// StartRetrieveProvider :
func StartRetrieveProvider(url string, factory EncodingFactory) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, factory, "providerRetrieve")
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
			provider.retrieveResponse(transaction, archiveDetailsList, elementList)
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

//======================================================================//
//								QUERY									//
//======================================================================//
// StartRetrieveProvider :
func StartQueryProvider(url string, factory EncodingFactory) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, factory, "providerQuery")
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
			// TODO
		}

		return nil
	}

	err := provider.cctx.RegisterProgressHandler(SERVICE_AREA_NUMBER,
		SERVICE_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_QUERY,
		queryHandler)
	if err != nil {
		return err
	}

	return nil
}

// PROGRESS
func (provider *Provider) queryProgress(msg *Message) (*Boolean, *ObjectType, *ArchiveQueryList, *QueryFilterList, error) {
	return nil, nil, nil, nil, nil
}

// ACK
func (provider *Provider) queryAck(transaction ProgressTransaction) error {
	return nil
}

// UPDATE
func (provider *Provider) queryUpdate(transaction ProgressTransaction, objectType *ObjectType, identifierList *IdentifierList, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	return nil
}

// RESPONSE
func (provider *Provider) queryResponse(transaction ProgressTransaction, objectType *ObjectType, identifierList *IdentifierList, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	return nil
}

//======================================================================//
//								COUNT									//
//======================================================================//
// StartCountProvider :
func StartCountProvider(url string, factory EncodingFactory) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, factory, "providerCount")
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
			// TODO
		}

		return nil
	}

	err := provider.cctx.RegisterInvokeHandler(SERVICE_AREA_NUMBER,
		SERVICE_AREA_VERSION,
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
	return nil, nil, nil, nil
}

// ACK
func (provider *Provider) countAck(transaction InvokeTransaction) error {
	return nil
}

// RESPONSE
func (provider *Provider) countResponse(transaction InvokeTransaction, longList *LongList) error {
	return nil
}

//======================================================================//
//								STORE									//
//======================================================================//
// StartStoreProvider :
func StartStoreProvider(url string, factory EncodingFactory) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, factory, "providerStore")
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
			// TODO
		}

		return nil
	}

	err := provider.cctx.RegisterRequestHandler(SERVICE_AREA_NUMBER,
		SERVICE_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		OPERATION_IDENTIFIER_STORE,
		storeHandler)
	if err != nil {
		return err
	}

	return nil
}

// REQUEST
func (provider *Provider) storeRequest(msg *Message) (*Boolean, *ObjectType, *IdentifierList, *ArchiveDetailsList, ElementList, error) {
	return nil, nil, nil, nil, nil, nil
}

// RESPONSE
func (provider *Provider) storeResponse(transaction RequestTransaction, longList *LongList) error {
	return nil
}

//======================================================================//
//								UPDATE									//
//======================================================================//
// StartUpdateProvider :
func StartUpdateProvider(url string, factory EncodingFactory) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, factory, "providerUpdate")
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
			// TODO
		}

		return nil
	}

	err := provider.cctx.RegisterSubmitHandler(SERVICE_AREA_NUMBER,
		SERVICE_AREA_VERSION,
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
	return nil, nil, nil, nil, nil
}

// ACK
func (provider *Provider) updateAck(transaction SubmitTransaction) error {
	return nil
}

//======================================================================//
//								DELETE									//
//======================================================================//
// StartDeleteProvider :
func StartDeleteProvider(url string, factory EncodingFactory) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url, factory, "providerDelete")
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
			// TODO
		}

		return nil
	}

	err := provider.cctx.RegisterRequestHandler(SERVICE_AREA_NUMBER,
		SERVICE_AREA_VERSION,
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
	return nil, nil, nil, nil
}

// RESPONSE
func (provider *Provider) deleteResponse(longList *LongList) error {
	return nil
}
