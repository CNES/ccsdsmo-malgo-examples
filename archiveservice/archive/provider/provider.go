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
	"strings"
	"time"

	. "github.com/CNES/ccsdsmo-malgo/com"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/api"
	. "github.com/CNES/ccsdsmo-malgo/mal/encoding/binary"

	. "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/archive/constants"
	arch "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/archive/storage"
	. "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/data"
	. "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/errors"
)

// Define Provider's structure
type Provider struct {
	ctx     *Context
	cctx    *ClientContext
	factory EncodingFactory
}

// Create a provider
func createProvider(url string) (*Provider, error) {
	ctx, err := NewContext(url)
	if err != nil {
		return nil, err
	}

	cctx, err := NewClientContext(ctx, "archiveServiceProvider")
	if err != nil {
		return nil, err
	}

	factory := new(FixedBinaryEncoding)

	provider := &Provider{ctx, cctx, factory}

	return provider, nil
}

// StartProvider : TODO:
func StartProvider(url string) (*Provider, error) {
	// Create the provider
	provider, err := createProvider(url)
	if err != nil {
		return nil, err
	}

	// Create and launch the Retrieve handler
	err = provider.retrieveHandler()
	if err != nil {
		return nil, err
	}

	// Create and launch the Query handler
	err = provider.queryHandler()
	if err != nil {
		return nil, err
	}

	// Create and launch the Count handler
	err = provider.countHandler()
	if err != nil {
		return nil, err
	}

	// Create and launch the Store handler
	err = provider.storeHandler()
	if err != nil {
		return nil, err
	}

	// Create and launch the Update handler
	err = provider.updateHandler()
	if err != nil {
		return nil, err
	}

	// Create and launch the Delete handler
	err = provider.deleteHandler()
	if err != nil {
		return nil, err
	}

	return provider, nil
}

// Close : Allow to close the context of a specific provider
func (provider *Provider) Close() {
	provider.ctx.Close()
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
// Create a handler for the retrieve operation
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
			err = provider.retrieveVerifyParameters(transaction, objectType, identifierList)
			if err != nil {
				return err
			}

			// ----- Call Ack operation -----
			err = provider.retrieveAck(transaction)
			if err != nil {
				provider.retrieveAckError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}

			/*fmt.Println("RetrieveHandler received:\n\t>>>",
			objectType, "\n\t>>>",
			identifierList, "\n\t>>>",
			longList)*/

			// Retrieve these objects in the archive
			archiveDetailsList, elementList, err := arch.RetrieveInArchive(*objectType, *identifierList, *longList)
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

// VERIFY PARAMETERS : TODO:
func (provider *Provider) retrieveVerifyParameters(transaction InvokeTransaction, objectType *ObjectType, identifierList *IdentifierList) error {
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		provider.retrieveAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR))
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(*identifierList)[i] == "*" {
			provider.retrieveAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR))
		}
	}

	return nil
}

// INVOKE : TODO:
func (provider *Provider) retrieveInvoke(msg *Message) (*ObjectType, *IdentifierList, *LongList, error) {
	element, err := msg.DecodeParameter(NullObjectType)
	if err != nil {
		return nil, nil, nil, err
	}
	objectType := element.(*ObjectType)

	element, err = msg.DecodeParameter(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, err
	}
	identifierList := element.(*IdentifierList)

	element, err = msg.DecodeLastParameter(NullLongList, false)
	if err != nil {
		return nil, nil, nil, err
	}
	longList := element.(*LongList)

	return objectType, identifierList, longList, nil
}

// ACK : TODO:
func (provider *Provider) retrieveAck(transaction InvokeTransaction) error {
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

// ACK ERROR : TODO:
func (provider *Provider) retrieveAckError(transaction InvokeTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Ack(body, true)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE : TODO:
func (provider *Provider) retrieveResponse(transaction InvokeTransaction, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := body.EncodeParameter(archiveDetailsList)
	if err != nil {
		return err
	}

	err = body.EncodeLastParameter(elementList, true)
	if err != nil {
		return err
	}

	transaction.Reply(body, false)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE ERROR : TODO:
func (provider *Provider) retrieveResponseError(transaction InvokeTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Response operation with Error status
	err = transaction.Reply(body, true)
	if err != nil {
		return err
	}

	return nil
}

//======================================================================//
//								QUERY									//
//======================================================================//
// Create a handler for the query operation
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
			err = provider.queryVerifyParameters(transaction, archiveQueryList, queryFilterList)
			if err != nil {
				return err
			}

			// ----- Call Ack operation -----
			err = provider.queryAck(transaction)
			if err != nil {
				provider.queryAckError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE, NewLongList(0))
				return err
			}

			// Variables to send to the consumer
			var objType []*ObjectType
			var archDetList []*ArchiveDetailsList
			var idList []*IdentifierList
			var elementList []ElementList

			/*fmt.Println("QueryHandler received:\n\t>>>",
			boolean, "\n\t>>>",
			objectType, "\n\t>>>",
			archiveQueryList, "\n\t>>>",
			queryFilterList)*/

			for i := 0; i < archiveQueryList.Size()-1; i++ {
				// Do a query to the archive
				if queryFilterList != nil {
					objType, archDetList, idList, elementList, err = arch.QueryArchive(boolean, *objectType, *(*archiveQueryList)[i], queryFilterList.GetElementAt(i))
				} else {
					objType, archDetList, idList, elementList, err = arch.QueryArchive(boolean, *objectType, *(*archiveQueryList)[i], nil)
				}
				if err != nil {
					// Send an INVALID error
					if err.Error() == string(ARCHIVE_SERVICE_QUERY_SORT_FIELD_NAME_INVALID_ERROR) ||
						err.Error() == string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR) ||
						strings.Contains(err.Error(), string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR)) {
						provider.queryUpdateError(transaction, COM_ERROR_INVALID, String(err.Error()), NewLongList(0))
					}
					// Otherwise, send an INTERNAL error
					provider.queryUpdateError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
					return err
				}
				for j := 0; j < len(archDetList); j++ {
					// Call Update operation
					err = provider.queryUpdate(transaction, objType[j], idList[j], archDetList[j], elementList[j])
					if err != nil {
						// Send an INTERNAL error
						provider.queryUpdateError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
						return err
					}
				}
			}

			// Do a query to the archive
			if queryFilterList != nil {
				objType, archDetList, idList, elementList, err = arch.QueryArchive(boolean, *objectType, *(*archiveQueryList)[archiveQueryList.Size()-1], queryFilterList.GetElementAt(archiveQueryList.Size()-1))
			} else {
				objType, archDetList, idList, elementList, err = arch.QueryArchive(boolean, *objectType, *(*archiveQueryList)[archiveQueryList.Size()-1], nil)
			}
			if err != nil {
				// Send an INVALID error
				if err.Error() == string(ARCHIVE_SERVICE_QUERY_SORT_FIELD_NAME_INVALID_ERROR) ||
					strings.Contains(err.Error(), string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR)) {
					provider.queryUpdateError(transaction, COM_ERROR_INVALID, String(err.Error()), NewLongList(0))
				}
				// Otherwise, send an INTERNAL error
				provider.queryUpdateError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}
			// ----- Call Response operation -----
			// Unless archive query list size is equal to 1 (we didn't enter in the previous loop)
			for j := 0; j < len(archDetList); j++ {
				if j == len(archDetList)-1 {
					// Call Response operation
					err = provider.queryResponse(transaction, objType[j], idList[j], archDetList[j], elementList[j])
					if err != nil {
						// Send an INTERNAL error
						provider.queryResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
						return err
					}
					break
				}
				err = provider.queryUpdate(transaction, objType[j], idList[j], archDetList[j], elementList[j])
				if err != nil {
					// Send an INTERNAL error
					provider.queryUpdateError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
					return err
				}
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

// VERIFY PARAMETERS : TODO:
func (provider *Provider) queryVerifyParameters(transaction ProgressTransaction, archiveQueryList *ArchiveQueryList, queryFilterList QueryFilterList) error {
	if queryFilterList != nil && archiveQueryList.Size() != queryFilterList.Size() {
		provider.queryAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_QUERY_LISTS_SIZE_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_QUERY_LISTS_SIZE_ERROR))
	}

	return nil
}

// PROGRESS : TODO:
func (provider *Provider) queryProgress(msg *Message) (*Boolean, *ObjectType, *ArchiveQueryList, QueryFilterList, error) {
	// Decode Boolean
	boolean, err := msg.DecodeParameter(NullBoolean)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode ObjectType
	objectType, err := msg.DecodeParameter(NullObjectType)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode ArchiveQueryList
	archiveQueryList, err := msg.DecodeParameter(NullArchiveQueryList)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode QueryFilterList
	queryFilterList, err := msg.DecodeLastParameter(nil, true)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	if queryFilterList == nil {
		return boolean.(*Boolean), objectType.(*ObjectType), archiveQueryList.(*ArchiveQueryList), nil, nil
	}

	return boolean.(*Boolean), objectType.(*ObjectType), archiveQueryList.(*ArchiveQueryList), queryFilterList.(QueryFilterList), nil
}

// ACK : TODO:
func (provider *Provider) queryAck(transaction ProgressTransaction) error {
	// Call Ack operation
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

// ACK ERROR : TODO:
func (provider *Provider) queryAckError(transaction ProgressTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Ack(body, true)
	if err != nil {
		return err
	}

	return nil
}

// UPDATE : TODO:
func (provider *Provider) queryUpdate(transaction ProgressTransaction, objectType *ObjectType, identifierList *IdentifierList, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	// create a body for the operation call
	body := transaction.NewBody()

	// Encode ObjectType
	err := body.EncodeParameter(objectType)
	if err != nil {
		return err
	}

	// Encode IdentifierList
	err = body.EncodeParameter(identifierList)
	if err != nil {
		return err
	}

	// Encode ArchiveDetailsList
	err = body.EncodeParameter(archiveDetailsList)
	if err != nil {
		return err
	}

	// Encode ElementList
	err = body.EncodeLastParameter(elementList, true)
	if err != nil {
		return err
	}

	// Call Update operation
	err = transaction.Update(body, false)
	if err != nil {
		return err
	}

	return nil
}

// UPDATE ERROR : TODO:
func (provider *Provider) queryUpdateError(transaction ProgressTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Update(body, true)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE : TODO:
func (provider *Provider) queryResponse(transaction ProgressTransaction, objectType *ObjectType, identifierList *IdentifierList, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	// create a body for the operation call
	body := transaction.NewBody()

	// Encode ObjectType
	err := body.EncodeParameter(objectType)
	if err != nil {
		return err
	}

	// Encode IdentifierList
	err = body.EncodeParameter(identifierList)
	if err != nil {
		return err
	}

	// Encode ArchiveDetailsList
	err = body.EncodeParameter(archiveDetailsList)
	if err != nil {
		return err
	}

	// Encode ElementList
	err = body.EncodeLastParameter(elementList, true)
	if err != nil {
		return err
	}

	// Call Update operation
	err = transaction.Reply(body, false)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE ERROR : TODO:
func (provider *Provider) queryResponseError(transaction ProgressTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Reply(body, true)
	if err != nil {
		return err
	}

	return nil
}

//======================================================================//
//								COUNT									//
//======================================================================//
// Create a handler for the count operation
func (provider *Provider) countHandler() error {
	countHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(InvokeTransaction)

			// Call Invoke operation
			objectType, archiveQueryList, queryFilterList, err := provider.countInvoke(msg)
			if err != nil {
				provider.countAckError(transaction, MAL_ERROR_BAD_ENCODING, MAL_ERROR_BAD_ENCODING_MESSAGE, NewLongList(0))
				return err
			}

			// ----- Verify the parameters -----
			err = provider.countVerifyParameters(transaction, archiveQueryList, queryFilterList)
			if err != nil {
				return err
			}

			// Call Ack operation
			err = provider.retrieveAck(transaction)
			if err != nil {
				provider.countAckError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				return err
			}

			/*fmt.Println("CountHandler received:\n\t>>>",
			objectType, "\n\t>>>",
			archiveQueryList, "\n\t>>>",
			queryFilterList)*/

			// This variable will be created automatically in the future
			longList, err := arch.CountInArchive(*objectType, *archiveQueryList, queryFilterList)
			if err != nil {
				// Send an INVALID error
				if err.Error() == string(ARCHIVE_SERVICE_QUERY_SORT_FIELD_NAME_INVALID_ERROR) ||
					strings.Contains(err.Error(), string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR)) {
					provider.countResponseError(transaction, COM_ERROR_INVALID, String(err.Error()), NewLongList(0))
				} else {
					// Otherwise, send an INTERNAL error
					provider.countResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				}
				return err
			}
			// Call Response operation
			err = provider.countResponse(transaction, longList)
			if err != nil {
				provider.countResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
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

// VERIFY PARAMETERS : TODO:
func (provider *Provider) countVerifyParameters(transaction InvokeTransaction, archiveQueryList *ArchiveQueryList, queryFilterList QueryFilterList) error {
	if queryFilterList != nil && archiveQueryList.Size() != queryFilterList.Size() {
		provider.countAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_QUERY_LISTS_SIZE_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_QUERY_LISTS_SIZE_ERROR))
	}

	return nil
}

// INVOKE : TODO:
func (provider *Provider) countInvoke(msg *Message) (*ObjectType, *ArchiveQueryList, QueryFilterList, error) {
	// Decode ObjectType
	objectType, err := msg.DecodeParameter(NullObjectType)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode ArchiveQueryList
	archiveQueryList, err := msg.DecodeParameter(NullArchiveQueryList)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode QueryFilterList
	queryFilterList, err := msg.DecodeLastParameter(nil, true)
	if err != nil {
		return nil, nil, nil, err
	}
	if queryFilterList == nil {
		return objectType.(*ObjectType), archiveQueryList.(*ArchiveQueryList), nil, nil
	}

	return objectType.(*ObjectType), archiveQueryList.(*ArchiveQueryList), queryFilterList.(QueryFilterList), nil
}

// ACK : TODO:
func (provider *Provider) countAck(transaction InvokeTransaction) error {
	// Call Ack operation
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

// ACK ERROR : TODO:
func (provider *Provider) countAckError(transaction InvokeTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Ack(body, true)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE : TODO:
func (provider *Provider) countResponse(transaction InvokeTransaction, longList *LongList) error {
	// create a body for the operation call
	body := transaction.NewBody()

	// Encode LongList
	err := body.EncodeLastParameter(longList, false)
	if err != nil {
		return err
	}

	// Call Response operation
	err = transaction.Reply(body, false)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE ERROR : TODO:
func (provider *Provider) countResponseError(transaction InvokeTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Reply(body, true)
	if err != nil {
		return err
	}

	return nil
}

//======================================================================//
//								STORE									//
//======================================================================//
// Create a handler for the store operation
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

			/*fmt.Println("StoreHandler received:\n\t>>>",
			boolean, "\n\t>>>",
			objectType, "\n\t>>>",
			identifierList, "\n\t>>>",
			archiveDetailsList, "\n\t>>>",
			elementList)*/

			// Store these objects in the archive
			var longList *LongList
			longList, err = arch.StoreInArchive(boolean, *objectType, *identifierList, *archiveDetailsList, elementList)
			if err != nil {
				if err.Error() == string(COM_ERROR_DUPLICATE) {
					provider.storeResponseError(transaction, COM_ERROR_DUPLICATE, COM_ERROR_DUPLICATE_MESSAGE, NewLongList(0))
				} else {
					provider.storeResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				}
				return err
			}

			// TODO: for each object stored, an 'ObjectStored' event may be published
			// to the event service

			// Call Response operation
			err = provider.storeResponse(transaction, longList)
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

// VERIFY PARAMETERS : TODO:
func (provider *Provider) storeVerifyParameters(transaction RequestTransaction, boolean *Boolean, objectType *ObjectType, identifierList *IdentifierList, archiveDetailsList *ArchiveDetailsList, elementList ElementList) error {
	// The fourth and fifth lists must be the same size
	if archiveDetailsList.Size() != elementList.Size() {
		if archiveDetailsList.Size() <= elementList.Size() {
			provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_LIST_SIZE_ERROR, NewLong(int64(archiveDetailsList.Size())))
		} else {
			provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_LIST_SIZE_ERROR, NewLong(int64(elementList.Size())))
		}
		return errors.New(string(ARCHIVE_SERVICE_STORE_LIST_SIZE_ERROR))
	}

	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR))
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(*identifierList)[i] == "*" {
			provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR))
		}
	}

	// Verify the parameters network, timestamp and provider of the object ArchiveDetails
	for i := 0; i < archiveDetailsList.Size(); i++ {
		if (*archiveDetailsList)[i].Network == nil || *(*archiveDetailsList)[i].Network == "0" || *(*archiveDetailsList)[i].Network == "*" ||
			(*archiveDetailsList)[i].Timestamp == nil || *(*archiveDetailsList)[i].Timestamp == FineTime(time.Unix(int64(0), int64(0))) ||
			(*archiveDetailsList)[i].Provider == nil || *(*archiveDetailsList)[i].Provider == "0" || *(*archiveDetailsList)[i].Provider == "*" {
			provider.storeResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR))
		}
	}

	// TODO: Raise INVALID error for 3.4.6.2.12

	return nil
}

// REQUEST : TODO:
func (provider *Provider) storeRequest(msg *Message) (*Boolean, *ObjectType, *IdentifierList, *ArchiveDetailsList, ElementList, error) {
	// Decode Boolean
	boolean, err := msg.DecodeParameter(NullBoolean)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode ObjectType
	objectType, err := msg.DecodeParameter(NullObjectType)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode IdentifierList
	identifierList, err := msg.DecodeParameter(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode ArchiveDetailsList
	archiveDetailsList, err := msg.DecodeParameter(NullArchiveDetailsList)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	// Decode ElementList
	elementList, err := msg.DecodeLastParameter(nil, true)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}

	return boolean.(*Boolean), objectType.(*ObjectType), identifierList.(*IdentifierList), archiveDetailsList.(*ArchiveDetailsList), elementList.(ElementList), nil
}

// RESPONSE : TODO:
func (provider *Provider) storeResponse(transaction RequestTransaction, longList *LongList) error {
	// create a body for the operation call
	body := transaction.NewBody()

	// Encode LongList
	err := body.EncodeLastParameter(longList, false)
	if err != nil {
		return err
	}

	// Call Response operation
	err = transaction.Reply(body, false)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE ERROR : TODO:
func (provider *Provider) storeResponseError(transaction RequestTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Response operation with Error status
	err = transaction.Reply(body, true)
	if err != nil {
		return err
	}

	return nil
}

//======================================================================//
//								UPDATE									//
//======================================================================//
// Create a handler for the update operation
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

			/*fmt.Println("UpdateHandler received:\n\t>>>",
			objectType, "\n\t>>>",
			identifierList, "\n\t>>>",
			archiveDetailsList, "\n\t>>>",
			elementList)*/

			// Update these objects
			err = arch.UpdateArchive(*objectType, *identifierList, *archiveDetailsList, elementList)
			if err != nil {
				if err.Error() == string(MAL_ERROR_UNKNOWN_MESSAGE) {
					provider.updateAckError(transaction, MAL_ERROR_UNKNOWN, ARCHIVE_SERVICE_UNKNOWN_ELEMENT, NewLongList(0))
				} else {
					provider.updateAckError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				}
				return err
			}

			// TODO: for each object updated, an 'ObjectUpdated' event may be published
			// to the event service

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

// VERIFY PARAMETERS : TODO:
func (provider *Provider) updateVerifyParameters(transaction SubmitTransaction, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList) error {
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		provider.updateAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR))
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(identifierList)[i] == "*" {
			provider.updateAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR))
		}
	}

	// Verify object instance identifier
	for i := 0; i < archiveDetailsList.Size(); i++ {
		if archiveDetailsList[i].InstId == 0 {
			provider.updateAckError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_AREA_OBJECT_INSTANCE_IDENTIFIER_VALUE_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_AREA_OBJECT_INSTANCE_IDENTIFIER_VALUE_ERROR))
		}
	}

	return nil
}

// SUBMIT : TODO:
func (provider *Provider) updateSubmit(msg *Message) (*ObjectType, *IdentifierList, *ArchiveDetailsList, ElementList, error) {
	// Decode ObjectType
	objectType, err := msg.DecodeParameter(NullObjectType)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode IdentifierList
	identifierList, err := msg.DecodeParameter(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode ArchiveDetailsList
	archiveDetailsList, err := msg.DecodeParameter(NullArchiveDetailsList)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Decode ElementList
	elementList, err := msg.DecodeLastParameter(nil, true)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return objectType.(*ObjectType), identifierList.(*IdentifierList), archiveDetailsList.(*ArchiveDetailsList), elementList.(ElementList), nil
}

// ACK : TODO:
func (provider *Provider) updateAck(transaction SubmitTransaction) error {
	// Call Ack operation
	err := transaction.Ack(nil, false)
	if err != nil {
		return err
	}
	return nil
}

// ACK ERROR : TODO:
func (provider *Provider) updateAckError(transaction SubmitTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Ack operation with Error status
	err = transaction.Ack(body, true)
	if err != nil {
		return err
	}

	return nil
}

//======================================================================//
//								DELETE									//
//======================================================================//
// Create a handler for the delete operation
func (provider *Provider) deleteHandler() error {
	deleteHandler := func(msg *Message, t Transaction) error {
		if msg != nil {
			transaction := t.(RequestTransaction)

			// Call Request operation
			objectType, identifierList, longListRequest, err := provider.deleteRequest(msg)
			if err != nil {
				provider.deleteResponseError(transaction, MAL_ERROR_BAD_ENCODING, MAL_ERROR_BAD_ENCODING_MESSAGE, NewLongList(0))
				return err
			}

			// ----- Verify the parameters -----
			err = provider.deleteVerifyParameters(transaction, *objectType, *identifierList)
			if err != nil {
				return err
			}

			/*fmt.Println("DeleteHandler received:\n\t>>>",
			objectType, "\n\t>>>",
			identifierList, "\n\t>>>",
			longListRequest)*/

			// Delete these objects
			longListResponse, err := arch.DeleteInArchive(*objectType, *identifierList, *longListRequest)
			if err != nil {
				if err.Error() == string(MAL_ERROR_UNKNOWN_MESSAGE) {
					provider.deleteResponseError(transaction, MAL_ERROR_UNKNOWN, ARCHIVE_SERVICE_UNKNOWN_ELEMENT, NewLongList(0))
				} else {
					provider.deleteResponseError(transaction, MAL_ERROR_INTERNAL, MAL_ERROR_INTERNAL_MESSAGE+String(" "+err.Error()), NewLongList(0))
				}
				return err
			}

			// TODO: for each object deleted, an 'ObjectDeleted' event may be published
			// to the event service

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

// VERIFY PARAMETERS : TODO:
func (provider *Provider) deleteVerifyParameters(transaction RequestTransaction, objectType ObjectType, identifierList IdentifierList) error {
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0 {
		provider.deleteResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR, NewLongList(1))
		return errors.New(string(ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR))
	}

	// Verify IdentifierList
	for i := 0; i < identifierList.Size(); i++ {
		if *(identifierList)[i] == "*" {
			provider.deleteResponseError(transaction, COM_ERROR_INVALID, ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR, NewLongList(1))
			return errors.New(string(ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR))
		}
	}

	return nil
}

// REQUEST : TODO:
func (provider *Provider) deleteRequest(msg *Message) (*ObjectType, *IdentifierList, *LongList, error) {
	// Decode ObjectType
	objectType, err := msg.DecodeParameter(NullObjectType)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode IdentifierList
	identifierList, err := msg.DecodeParameter(NullIdentifierList)
	if err != nil {
		return nil, nil, nil, err
	}

	// Decode LongList
	longList, err := msg.DecodeLastParameter(NullLongList, false)
	if err != nil {
		return nil, nil, nil, err
	}

	return objectType.(*ObjectType), identifierList.(*IdentifierList), longList.(*LongList), nil
}

// RESPONSE : TODO:
func (provider *Provider) deleteResponse(transaction RequestTransaction, longList LongList) error {
	// create a body for the operation call
	body := transaction.NewBody()

	// Encode LongList
	err := body.EncodeLastParameter(&longList, false)
	if err != nil {
		return err
	}

	// Call Response operation
	err = transaction.Reply(body, false)
	if err != nil {
		return err
	}

	return nil
}

// RESPONSE ERROR : TODO:
func (provider *Provider) deleteResponseError(transaction RequestTransaction, errorNumber UInteger, errorComment String, errorExtra Element) error {
	// create a body for the operation call
	body := transaction.NewBody()

	err := EncodeError(body, errorNumber, errorComment, errorExtra)
	if err != nil {
		return err
	}

	// Call Response operation with Error status
	err = transaction.Reply(body, true)
	if err != nil {
		return err
	}

	return nil
}
