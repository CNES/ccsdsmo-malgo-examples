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
	. "github.com/ccsdsmo/malgo/mal/encoding/binary"

	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/errors"
)

// InvokeConsumer :
type InvokeConsumer struct {
	ctx     *Context
	cctx    *ClientContext
	op      InvokeOperation
	factory EncodingFactory
}

// ProgressConsumer :
type ProgressConsumer struct {
	ctx     *Context
	cctx    *ClientContext
	op      ProgressOperation
	factory EncodingFactory
}

// RequestConsumer :
type RequestConsumer struct {
	ctx     *Context
	cctx    *ClientContext
	op      RequestOperation
	factory EncodingFactory
}

// SubmitConsumer :
type SubmitConsumer struct {
	ctx     *Context
	cctx    *ClientContext
	op      SubmitOperation
	factory EncodingFactory
}

// Close :
func (i *InvokeConsumer) Close() {
	i.ctx.Close()
}

// Close :
func (p *ProgressConsumer) Close() {
	p.ctx.Close()
}

// Close :
func (r *RequestConsumer) Close() {
	r.ctx.Close()
}

// Close :
func (s *SubmitConsumer) Close() {
	s.ctx.Close()
}

//======================================================================//
//								CONSUMERS								//
//======================================================================//
// Create a consumer for an invoke operation
func createInvokeConsumer(url string, providerURI *URI, typeOfConsumer string, operation UShort) (*InvokeConsumer, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, typeOfConsumer)
	if err != nil {
		return nil, err
	}

	op := cctx.NewInvokeOperation(providerURI,
		COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		operation)

	factory := new(FixedBinaryEncoding)

	consumer := &InvokeConsumer{ctx, cctx, op, factory}

	return consumer, nil
}

// Create a consumer for a progress operation
func createProgressConsumer(url string, providerURI *URI, typeOfConsumer string, operation UShort) (*ProgressConsumer, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, typeOfConsumer)
	if err != nil {
		return nil, err
	}

	op := cctx.NewProgressOperation(providerURI,
		COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		operation)

	factory := new(FixedBinaryEncoding)

	consumer := &ProgressConsumer{ctx, cctx, op, factory}

	return consumer, nil
}

// Create a consumer for a request operation
func createRequestConsumer(url string, providerURI *URI, typeOfConsumer string, operation UShort) (*RequestConsumer, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, typeOfConsumer)
	if err != nil {
		return nil, err
	}

	op := cctx.NewRequestOperation(providerURI,
		COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		operation)

	factory := new(FixedBinaryEncoding)

	consumer := &RequestConsumer{ctx, cctx, op, factory}

	return consumer, nil
}

// Create a consumer for a submit operation
func createSubmitConsumer(url string, providerURI *URI, typeOfConsumer string, operation UShort) (*SubmitConsumer, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, typeOfConsumer)
	if err != nil {
		return nil, err
	}

	op := cctx.NewSubmitOperation(providerURI,
		COM_AREA_NUMBER,
		COM_AREA_VERSION,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		operation)

	factory := new(FixedBinaryEncoding)

	consumer := &SubmitConsumer{ctx, cctx, op, factory}

	return consumer, nil
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
// StartRetrieveConsumer : TODO
func StartRetrieveConsumer(url string, providerURI *URI, objectType ObjectType, identifierList IdentifierList, longList LongList) (*InvokeConsumer, *ArchiveDetailsList, ElementList, *ServiceError, error) {
	// Create the consumer
	consumer, err := createInvokeConsumer(url, providerURI, "consumerRetrieve", OPERATION_IDENTIFIER_RETRIEVE)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Call Invoke operation
	errorsList, err := consumer.retrieveInvoke(objectType, identifierList, longList)
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, nil, err
	} else if errorsList != nil {
		// Close consumer
		consumer.Close()
		return nil, nil, nil, errorsList, nil
	}

	// Call Response operation
	archiveDetailsList, elementList, errorsList, err := consumer.retrieveResponse()
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, nil, err
	} else if errorsList != nil {
		// Close consumer
		consumer.Close()
		return nil, nil, nil, errorsList, nil
	}

	return consumer, archiveDetailsList, elementList, nil, nil
}

// Invoke & Ack
func (consumer *InvokeConsumer) retrieveInvoke(objectType ObjectType, identifierList IdentifierList, longList LongList) (*ServiceError, error) {
	// Create the encoder
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))
	// Encode ObjectType
	err := objectType.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// Encode IdentifierList
	err = identifierList.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// Encode LongList
	err = longList.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// Call Invoke operation
	resp, err := consumer.op.Invoke(encoder.Body())
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, err
			}

			return errorsList, nil
		}
		return nil, err
	}

	return nil, nil
}

// Response
func (consumer *InvokeConsumer) retrieveResponse() (*ArchiveDetailsList, ElementList, *ServiceError, error) {
	// Call Response operation
	resp, err := consumer.op.GetResponse()
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, nil, nil, err
			}

			return nil, nil, errorsList, nil
		}
		return nil, nil, nil, err
	}

	// Create the decoder
	decoder := consumer.factory.NewDecoder(resp.Body)

	// Decode ArchiveDetailsList
	archiveDetailsList, err := decoder.DecodeElement(NullArchiveDetailsList)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode ElementList
	elementList, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, nil, err
	}

	return archiveDetailsList.(*ArchiveDetailsList), elementList.(ElementList), nil, nil
}

//======================================================================//
//								QUERY									//
//======================================================================//
// StartQueryConsumer : TODO
func StartQueryConsumer(url string, providerURI *URI, boolean *Boolean, objectType ObjectType, archiveQueryList ArchiveQueryList, queryFilterList QueryFilterList) (*ProgressConsumer, []interface{}, *ServiceError, error) {
	// Create the consumer
	consumer, err := createProgressConsumer(url, providerURI, "consumerQuery", OPERATION_IDENTIFIER_QUERY)
	if err != nil {
		return nil, nil, nil, err
	}

	// Call Progress function
	errorsList, err := consumer.queryProgress(boolean, objectType, archiveQueryList, queryFilterList)
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, err
	} else if errorsList != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, errorsList, nil
	}

	// Create the interface that will receive all the responses
	responses := []interface{}{}
	// Call Update operation
	respObjType, respIDList, respArchDetList, respElemList, errorsList, err := consumer.queryUpdate()
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, err
	} else if errorsList != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, errorsList, nil
	}

	for respArchDetList != nil {
		// Put the objects in the interface
		responses = append(responses, respObjType, respIDList, respArchDetList, respElemList)

		// Call Update operation until it returns nil variables
		respObjType, respIDList, respArchDetList, respElemList, errorsList, err = consumer.queryUpdate()
		if err != nil {
			// Close consummer
			consumer.Close()
			return nil, nil, nil, err
		} else if errorsList != nil {
			// Close consummer
			consumer.Close()
			return nil, nil, errorsList, nil
		}
	}

	// Call Response operation
	respObjType, respIDList, respArchDetList, respElemList, errorsList, err = consumer.queryResponse()
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, err
	} else if errorsList != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, errorsList, nil
	}

	// Put the objects in the interface
	responses = append(responses, respObjType, respIDList, respArchDetList, respElemList)

	return consumer, responses, nil, nil
}

// Progress & Ack
func (consumer *ProgressConsumer) queryProgress(boolean *Boolean, objectType ObjectType, archiveQueryList ArchiveQueryList, queryFilterList QueryFilterList) (*ServiceError, error) {
	// Create the encoder
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode Boolean
	err := encoder.EncodeNullableElement(boolean)
	if err != nil {
		return nil, err
	}

	// Encode ObjectType
	err = objectType.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// Encode ArchiveQueryList
	err = archiveQueryList.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// Encode QueryFilterList
	err = encoder.EncodeNullableAbstractElement(queryFilterList)
	if err != nil {
		return nil, err
	}

	// Call Progress operation
	resp, err := consumer.op.Progress(encoder.Body())
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, err
			}

			return errorsList, nil
		}
		return nil, err
	}

	return nil, nil
}

// Update
func (consumer *ProgressConsumer) queryUpdate() (*ObjectType, *IdentifierList, *ArchiveDetailsList, ElementList, *ServiceError, error) {
	// Call Update operation
	updt, err := consumer.op.GetUpdate()
	if err != nil {
		// Verify if an error occurs during the operation
		if updt.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(updt.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}

			return nil, nil, nil, nil, errorsList, nil
		}
		return nil, nil, nil, nil, nil, err
	}

	if updt != nil {
		// Create the decoder to decode the multiple variables
		decoder := consumer.factory.NewDecoder(updt.Body)

		// Decode ObjectType
		objectType, err := decoder.DecodeNullableElement(NullObjectType)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}

		// Decode IdentifierList
		identifierList, err := decoder.DecodeNullableElement(NullIdentifierList)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}

		// Decode ArchiveDetailsList
		archiveDetailsList, err := decoder.DecodeNullableElement(NullArchiveDetailsList)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}

		// Decode ElementList
		elementList, err := decoder.DecodeNullableAbstractElement()
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}

		return objectType.(*ObjectType), identifierList.(*IdentifierList), archiveDetailsList.(*ArchiveDetailsList), elementList.(ElementList), nil, nil
	}
	return nil, nil, nil, nil, nil, nil
}

// Response
func (consumer *ProgressConsumer) queryResponse() (*ObjectType, *IdentifierList, *ArchiveDetailsList, ElementList, *ServiceError, error) {
	// Call Update operation
	resp, err := consumer.op.GetResponse()
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, nil, nil, nil, nil, err
			}

			return nil, nil, nil, nil, errorsList, nil
		}
		return nil, nil, nil, nil, nil, err
	}

	// Create the decoder to decode the multiple variables
	decoder := consumer.factory.NewDecoder(resp.Body)

	// Decode ObjectType
	objectType, err := decoder.DecodeNullableElement(NullObjectType)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode IdentifierList
	identifierList, err := decoder.DecodeNullableElement(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode ArchiveDetailsList
	archiveDetailsList, err := decoder.DecodeNullableElement(NullArchiveDetailsList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode ElementList
	elementList, err := decoder.DecodeNullableAbstractElement()
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Not a good method but it works...
	if elementList == nil {
		return objectType.(*ObjectType), identifierList.(*IdentifierList), archiveDetailsList.(*ArchiveDetailsList), nil, nil, nil
	}

	return objectType.(*ObjectType), identifierList.(*IdentifierList), archiveDetailsList.(*ArchiveDetailsList), elementList.(ElementList), nil, nil
}

//======================================================================//
//								COUNT									//
//======================================================================//
// StartCountConsumer : TODO
func StartCountConsumer(url string, providerURI *URI, objectType *ObjectType, archiveQueryList *ArchiveQueryList, queryFilterList QueryFilterList) (*InvokeConsumer, *LongList, *ServiceError, error) {
	// Create the consumer
	consumer, err := createInvokeConsumer(url, providerURI, "consumerCount", OPERATION_IDENTIFIER_COUNT)
	if err != nil {
		return nil, nil, nil, err
	}

	// Call Invoke function
	errorsList, err := consumer.countInvoke(objectType, archiveQueryList, queryFilterList)
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, err
	} else if errorsList != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, errorsList, nil
	}

	// Call Response function
	longList, errorsList, err := consumer.countResponse()
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, err
	} else if errorsList != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, errorsList, nil
	}

	return consumer, longList, nil, nil
}

// Invoke & Ack
func (consumer *InvokeConsumer) countInvoke(objectType *ObjectType, archiveQueryList *ArchiveQueryList, queryFilterList QueryFilterList) (*ServiceError, error) {
	// Create the encoder
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode ObjectType
	err := encoder.EncodeNullableElement(objectType)
	if err != nil {
		return nil, err
	}

	// Encode ArchiveQueryList
	err = encoder.EncodeNullableElement(archiveQueryList)
	if err != nil {
		return nil, err
	}

	// Encode QueryFilterList
	err = encoder.EncodeNullableAbstractElement(queryFilterList)
	if err != nil {
		return nil, err
	}

	// Call Invoke operation
	// TODO: we should retrieve the msg to verify if the ack is an error or not
	resp, err := consumer.op.Invoke(encoder.Body())
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, err
			}

			return errorsList, nil
		}
		return nil, err
	}

	return nil, nil
}

// Response
func (consumer *InvokeConsumer) countResponse() (*LongList, *ServiceError, error) {
	// Call Response operation
	resp, err := consumer.op.GetResponse()
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, nil, err
			}

			return nil, errorsList, nil
		}
		return nil, nil, err
	}

	// Create the decoder
	decoder := consumer.factory.NewDecoder(resp.Body)

	// Decode LongList
	longList, err := decoder.DecodeNullableElement(NullLongList)
	if err != nil {
		return nil, nil, err
	}

	return longList.(*LongList), nil, nil
}

//======================================================================//
//								STORE									//
//======================================================================//
// StartStoreConsumer : TODO
func StartStoreConsumer(url string, providerURI *URI, boolean *Boolean, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*RequestConsumer, *LongList, *ServiceError, error) {
	// Create the consumer
	consumer, err := createRequestConsumer(url, providerURI, "consumerStore", OPERATION_IDENTIFIER_STORE)
	if err != nil {
		return nil, nil, nil, err
	}

	// Call Request function and retrieve the Response
	longList, errorsList, err := consumer.storeRequest(boolean, objectType, identifierList, archiveDetailsList, elementList)
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, err
	} else if errorsList != nil {
		// Close consumer
		consumer.Close()
		return nil, nil, errorsList, nil
	}

	return consumer, longList, nil, nil
}

// Request & Response
func (consumer *RequestConsumer) storeRequest(boolean *Boolean, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*LongList, *ServiceError, error) {
	// Create the encoder
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode Boolean
	err := encoder.EncodeNullableElement(boolean)
	if err != nil {
		return nil, nil, err
	}

	// Encode ObjectType
	err = objectType.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}

	// Encode IdentifierList
	err = identifierList.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}

	// Encode ArchiveDetailsList
	err = archiveDetailsList.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}

	// Encode ElementList
	err = encoder.EncodeAbstractElement(elementList)
	if err != nil {
		return nil, nil, err
	}

	// Call Request operation and retrieve the Response
	resp, err := consumer.op.Request(encoder.Body())
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, nil, err
			}

			return nil, errorsList, nil
		}
		return nil, nil, err
	}

	// Create the decoder
	decoder := consumer.factory.NewDecoder(resp.Body)

	// Decode LongList
	longList, err := decoder.DecodeNullableElement(NullLongList)
	if err != nil {
		return nil, nil, err
	}

	return longList.(*LongList), nil, nil
}

//======================================================================//
//								UPDATE									//
//======================================================================//
// StartUpdateConsumer : TODO
func StartUpdateConsumer(url string, providerURI *URI, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*SubmitConsumer, *ServiceError, error) {
	// Create the consumer
	consumer, err := createSubmitConsumer(url, providerURI, "consumerUpdate", OPERATION_IDENTIFIER_UPDATE)
	if err != nil {
		return nil, nil, err
	}

	// Call Submit function
	errorsList, err := consumer.updateSubmit(objectType, identifierList, archiveDetailsList, elementList)
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, err
	} else if errorsList != nil {
		// Close consumer
		consumer.Close()
		return nil, errorsList, nil
	}

	return consumer, nil, nil
}

// Submit & Ack
func (consumer *SubmitConsumer) updateSubmit(objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*ServiceError, error) {
	// Create the encoder
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode ObjectType
	err := objectType.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// Encode IdentifierList
	err = identifierList.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// Encode ArchiveDetailsList
	err = archiveDetailsList.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// Encode ElementList
	err = encoder.EncodeAbstractElement(elementList)
	if err != nil {
		return nil, err
	}

	// Call Submit operation
	resp, err := consumer.op.Submit(encoder.Body())
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, err
			}

			return errorsList, nil
		}
		return nil, err
	}

	return nil, nil
}

//======================================================================//
//								DELETE									//
//======================================================================//
// StartDeleteConsumer : TODO
func StartDeleteConsumer(url string, providerURI *URI, objectType ObjectType, identifierList IdentifierList, longList LongList) (*RequestConsumer, *LongList, *ServiceError, error) {
	// Create the consumer
	consumer, err := createRequestConsumer(url, providerURI, "consumerDelete", OPERATION_IDENTIFIER_DELETE)
	if err != nil {
		return nil, nil, nil, err
	}

	// Call Request function and retrieve the Response
	respLongList, errorsList, err := consumer.deleteRequest(objectType, identifierList, longList)
	if err != nil {
		// Close consummer
		consumer.Close()
		return nil, nil, nil, err
	} else if errorsList != nil {
		// Close consumer
		consumer.Close()
		return nil, nil, errorsList, nil
	}

	return consumer, respLongList, nil, nil
}

// Request & Reponse
func (consumer *RequestConsumer) deleteRequest(objectType ObjectType, identifierList IdentifierList, longList LongList) (*LongList, *ServiceError, error) {
	// Create the encoder
	encoder := consumer.factory.NewEncoder(make([]byte, 0, 8192))

	// Encode ObjectType
	err := objectType.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}

	// Encode IdentifierList
	err = identifierList.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}

	// Encode LongList
	err = longList.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}

	// Call Request operation and retrieve the Response
	resp, err := consumer.op.Request(encoder.Body())
	if err != nil {
		// Verify if an error occurs during the operation
		if resp.IsErrorMessage {
			// Create the decoder
			decoder := consumer.factory.NewDecoder(resp.Body)
			// Decode the error
			errorsList, err := DecodeError(decoder)
			if err != nil {
				return nil, nil, err
			}

			return nil, errorsList, nil
		}
		return nil, nil, err
	}
	// Create the decoder
	decoder := consumer.factory.NewDecoder(resp.Body)

	// Decode LongList
	respLongList, err := decoder.DecodeElement(NullLongList)
	if err != nil {
		return nil, nil, err
	}

	return respLongList.(*LongList), nil, nil
}
