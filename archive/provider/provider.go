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
func (provider *Provider) retrieveAck(transaction InvokeTransaction) error {
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

func (provider *Provider) retrieveResponse(archiveDetailsList *ArchiveDetailsList, elementList ElementList, transaction InvokeTransaction) error {
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

func (provider *Provider) retrieveInvoke(msg *Message) (*ObjectType, *IdentifierList, ElementList, error) {
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

	elementList, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, nil, err
	}

	return objectType, identifierList, elementList.(ElementList), nil
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

//======================================================================//
//								QUERY									//
//======================================================================//

//======================================================================//
//								COUNT									//
//======================================================================//

//======================================================================//
//								STORE									//
//======================================================================//

//======================================================================//
//								UPDATE									//
//======================================================================//

//======================================================================//
//								DELETE									//
//======================================================================//
