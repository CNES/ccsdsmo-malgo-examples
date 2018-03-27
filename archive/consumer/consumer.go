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
package consumer

import (
	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/api"

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/EtienneLndr/MAL_API_Go_Project/data"
)

type Consumer struct {
	ctx     *Context
	cctx    *ClientContext
	op      InvokeOperation
	factory EncodingFactory
}

// Allow to close the context of a specific consumer
func (consumer *Consumer) Close() {
	consumer.ctx.Close()
}

// Create a consumer
func createConsumer(url string, factory EncodingFactory, providerURI *URI, typeOfConsumer string, operation UShort) (*Consumer, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, typeOfConsumer)
	if err != nil {
		return nil, err
	}

	op := cctx.NewInvokeOperation(providerURI,
		SERVICE_AREA_NUMBER,
		SERVICE_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		operation)

	consumer := &Consumer{ctx, cctx, op, factory}

	return consumer, nil
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
func (consumer *Consumer) retrieveInvoke(objectType ObjectType, identifierList IdentifierList, longList LongList) error {
	println("yooooooooooooooo")
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

func (consumer *Consumer) retrieveResponse() (*ArchiveDetailsList, ElementList, error) {
	resp, err := consumer.op.GetResponse()
	if err != nil {
		return nil, nil, err
	}
	decoder := consumer.factory.NewDecoder(resp.Body)

	archiveDetails, err := decoder.DecodeElement(NullArchiveDetailsList)
	if err != nil {
		return nil, nil, err
	}

	elementList, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, err
	}

	return archiveDetails.(*ArchiveDetailsList), elementList.(ElementList), nil
}

// StartRetrieveConsumer : TODO
func StartRetrieveConsumer(url string, factory EncodingFactory, providerURI *URI, objectType ObjectType, identifierList IdentifierList, longList LongList) (*Consumer, *ArchiveDetailsList, ElementList, error) {
	// Create the consumer
	consumer, err := createConsumer(url, factory, providerURI, "consumerRetrieve", OPERATION_IDENTIFIER_RETRIEVE)
	if err != nil {
		return nil, nil, nil, err
	}

	// Call Invoke operation
	err = consumer.retrieveInvoke(objectType, identifierList, longList)
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
