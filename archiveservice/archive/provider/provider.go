/**
 * MIT License
 *
 * Copyright (c) 2018-2020 CNES
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
	//	"errors"
	"strings"
	"time"

	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/com/archive"
	"github.com/CNES/ccsdsmo-malgo/mal"
	malapi "github.com/CNES/ccsdsmo-malgo/mal/api"

	. "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/archive/constants"
	arch "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/archive/storage"
)

// Define Provider's implementation structure
type ProviderImpl struct {
	uri string
}

// StartProvider : TODO:
func StartProvider(url string) (*archive.Provider, error) {
	ctx, err := mal.NewContext(url)
	if err != nil {
		return nil, err
	}
	return archive.NewProvider(ctx, "archiveServiceProvider", &ProviderImpl{"archiveServiceProvider"})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
func (*ProviderImpl) Retrieve(opHelper *archive.RetrieveHelper, objType *com.ObjectType, domain *mal.IdentifierList, objInstIds *mal.LongList) error {
	// ----- Verify the parameters -----
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objType.Area == 0 || objType.Number == 0 || objType.Service == 0 || objType.Version == 0 {
		extraInfo := mal.NewUIntegerList(1)
		(*extraInfo)[0] = mal.NewUInteger(1)
		return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
	}
	// Verify IdentifierList
	for i := 0; i < domain.Size(); i++ {
		if *(*domain)[i] == "*" {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(uint32(i + 1))
			return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
		}
	}

	// ----- Call Ack operation -----
	err := opHelper.Ack()
	if err != nil {
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	// Retrieve these objects in the archive
	archiveDetailsList, elementList, err := arch.RetrieveInArchive(*objType, *domain, *objInstIds)
	if err != nil {
		if err.Error() == string(mal.ERROR_UNKNOWN_MESSAGE) {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(0)
			return malapi.NewMalError(mal.ERROR_UNKNOWN, extraInfo)
		}
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	return opHelper.Reply(&archiveDetailsList, elementList)
}

//======================================================================//
//								QUERY									//
//======================================================================//
func (*ProviderImpl) Query(opHelper *archive.QueryHelper, returnBody *mal.Boolean, objType *com.ObjectType, archiveQuery *archive.ArchiveQueryList, queryFilter archive.QueryFilterList) error {
	// ----- Verify the parameters -----
	if queryFilter != nil && archiveQuery.Size() != queryFilter.Size() {
		extraInfo := mal.NewUIntegerList(1)
		(*extraInfo)[0] = mal.NewUInteger(uint32(min(archiveQuery.Size(), queryFilter.Size())))
		return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
	}

	// ----- Call Ack operation -----
	err := opHelper.Ack()
	if err != nil {
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	// Variables to send to the consumer
	var objTypes []*com.ObjectType
	var archDetList []*archive.ArchiveDetailsList
	var idList []*mal.IdentifierList
	var elementList []mal.ElementList

	for i := 0; i < archiveQuery.Size(); i++ {
		// Do a query to the archive
		if queryFilter != nil {
			objTypes, archDetList, idList, elementList, err = arch.QueryArchive(returnBody, *objType, *(*archiveQuery)[i], queryFilter.GetElementAt(i).(archive.QueryFilter))
		} else {
			objTypes, archDetList, idList, elementList, err = arch.QueryArchive(returnBody, *objType, *(*archiveQuery)[i], nil)
		}
		if err != nil {
			// Send an INVALID error
			if err.Error() == string(ARCHIVE_SERVICE_QUERY_SORT_FIELD_NAME_INVALID_ERROR) ||
				err.Error() == string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR) ||
				strings.Contains(err.Error(), string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR)) {
				extraInfo := mal.NewUIntegerList(1)
				(*extraInfo)[0] = mal.NewUInteger(uint32(i + 1))
				return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
			}
			// Otherwise, send an INTERNAL error
			return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
		}
		for j := 0; j < len(archDetList); j++ {
			if i == archiveQuery.Size()-1 && j == len(archDetList)-1 {
				// Call Response operation
				err = opHelper.Reply(objTypes[j], idList[j], archDetList[j], elementList[j])
			} else {
				// Call Update operation
				err = opHelper.Update(objTypes[j], idList[j], archDetList[j], elementList[j])
			}
			if err != nil {
				// Send an INTERNAL error
				return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
			}
		}
	}

	return nil
}

//======================================================================//
//								COUNT									//
//======================================================================//
func (*ProviderImpl) Count(opHelper *archive.CountHelper, objType *com.ObjectType, archiveQuery *archive.ArchiveQueryList, queryFilter archive.QueryFilterList) error {
	// ----- Verify the parameters -----
	if queryFilter != nil && archiveQuery.Size() != queryFilter.Size() {
		extraInfo := mal.NewUIntegerList(1)
		(*extraInfo)[0] = mal.NewUInteger(uint32(min(archiveQuery.Size(), queryFilter.Size())))
		return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
	}

	// ----- Call Ack operation -----
	err := opHelper.Ack()
	if err != nil {
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	// This variable will be created automatically in the future
	longList, err := arch.CountInArchive(*objType, *archiveQuery, queryFilter)
	if err != nil {
		// Send an INVALID error
		if err.Error() == string(ARCHIVE_SERVICE_QUERY_SORT_FIELD_NAME_INVALID_ERROR) ||
			strings.Contains(err.Error(), string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR)) {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(0)
			return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
		}
		// Otherwise, send an INTERNAL error
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}
	// Call Response operation
	err = opHelper.Reply(longList)
	if err != nil {
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	return nil
}

//======================================================================//
//								STORE									//
//======================================================================//
func (*ProviderImpl) Store(opHelper *archive.StoreHelper, returnObjInstIds *mal.Boolean, objType *com.ObjectType, domain *mal.IdentifierList, objDetails *archive.ArchiveDetailsList, objBodies mal.ElementList) error {
	// ----- Verify the parameters -----
	// The fourth and fifth lists must be the same size
	if objBodies != nil && objBodies.Size() != objDetails.Size() {
		extraInfo := mal.NewUIntegerList(1)
		(*extraInfo)[0] = mal.NewUInteger(1)
		return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
	}
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objType.Area == 0 || objType.Number == 0 || objType.Service == 0 || objType.Version == 0 {
		extraInfo := mal.NewUIntegerList(1)
		(*extraInfo)[0] = mal.NewUInteger(1)
		return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
	}
	// Verify IdentifierList
	for i := 0; i < domain.Size(); i++ {
		if *(*domain)[i] == "*" {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(uint32(i + 1))
			return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
		}
	}
	// Verify the parameters network, timestamp and provider of the object ArchiveDetails
	for i := 0; i < objDetails.Size(); i++ {
		if (*objDetails)[i].Network == nil || *(*objDetails)[i].Network == "0" || *(*objDetails)[i].Network == "*" ||
			(*objDetails)[i].Timestamp == nil || *(*objDetails)[i].Timestamp == mal.FineTime(time.Unix(int64(0), int64(0))) ||
			(*objDetails)[i].Provider == nil || *(*objDetails)[i].Provider == "0" || *(*objDetails)[i].Provider == "*" {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(uint32(i + 1))
			return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
		}
	}
	// TODO: Raise INVALID error for 3.4.6.2.12

	// Store these objects in the archive
	var longList *mal.LongList
	longList, err := arch.StoreInArchive(returnObjInstIds, *objType, *domain, *objDetails, objBodies)
	if err != nil {
		if err.Error() == string(com.ERROR_DUPLICATE) {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(0)
			return malapi.NewMalError(com.ERROR_DUPLICATE, extraInfo)
		}
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	// TODO: for each object stored, an 'ObjectStored' event may be published
	// to the event service

	// Call Response operation
	err = opHelper.Reply(longList)
	if err != nil {
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	return nil
}

//======================================================================//
//								UPDATE									//
//======================================================================//
func (*ProviderImpl) Update(opHelper *archive.UpdateHelper, objType *com.ObjectType, domain *mal.IdentifierList, objDetails *archive.ArchiveDetailsList, objBodies mal.ElementList) error {

	// ----- Verify the parameters -----
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objType.Area == 0 || objType.Number == 0 || objType.Service == 0 || objType.Version == 0 {
		extraInfo := mal.NewUIntegerList(1)
		(*extraInfo)[0] = mal.NewUInteger(1)
		return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
	}
	// Verify IdentifierList
	for i := 0; i < domain.Size(); i++ {
		if *(*domain)[i] == "*" {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(1)
			return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
		}
	}
	// Verify object instance identifier
	for i := 0; i < objDetails.Size(); i++ {
		if (*objDetails)[i].InstId == 0 {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(uint32(i + 1))
			return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
		}
	}

	// Update these objects
	err := arch.UpdateArchive(*objType, *domain, *objDetails, objBodies)
	if err != nil {
		if err.Error() == string(mal.ERROR_UNKNOWN_MESSAGE) {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(0)
			return malapi.NewMalError(mal.ERROR_UNKNOWN, extraInfo)
		}
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	// TODO: for each object updated, an 'ObjectUpdated' event may be published
	// to the event service

	// Call Ack operation
	err = opHelper.Ack()
	if err != nil {
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	return nil
}

//======================================================================//
//								DELETE									//
//======================================================================//
func (*ProviderImpl) Delete(opHelper *archive.DeleteHelper, objType *com.ObjectType, domain *mal.IdentifierList, objInstIds *mal.LongList) error {

	// ----- Verify the parameters -----
	// Verify ObjectType values (all of its attributes must not be equal to '0')
	if objType.Area == 0 || objType.Number == 0 || objType.Service == 0 || objType.Version == 0 {
		extraInfo := mal.NewUIntegerList(1)
		(*extraInfo)[0] = mal.NewUInteger(1)
		return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
	}
	// Verify IdentifierList
	for i := 0; i < domain.Size(); i++ {
		if *(*domain)[i] == "*" {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(1)
			return malapi.NewMalError(com.ERROR_INVALID, extraInfo)
		}
	}

	// Delete these objects
	longListResponse, err := arch.DeleteInArchive(*objType, *domain, *objInstIds)
	if err != nil {
		if err.Error() == string(mal.ERROR_UNKNOWN_MESSAGE) {
			extraInfo := mal.NewUIntegerList(1)
			(*extraInfo)[0] = mal.NewUInteger(0)
			return malapi.NewMalError(mal.ERROR_UNKNOWN, extraInfo)
		}
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	// TODO: for each object deleted, an 'ObjectDeleted' event may be published
	// to the event service

	// Call Response operation
	err = opHelper.Reply(&longListResponse)
	if err != nil {
		return malapi.NewMalError(mal.ERROR_INTERNAL, mal.NewString(err.Error()))
	}

	return nil
}
