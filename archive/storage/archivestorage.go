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

// Package storage : TODO:
package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"time"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/encoding/binary"

	. "github.com/etiennelndr/archiveservice/archive/constants"
	"github.com/etiennelndr/archiveservice/archive/utils"
	. "github.com/etiennelndr/archiveservice/data"

	// Init mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Database ids
const (
	USERNAME = "archiveService"
	PASSWORD = "1a2B3c4D!@?"
	DATABASE = "archive"
	TABLE    = "Archive"
)

// Database columns
var databaseFields = []string{
	"id",
	"objectInstanceIdentifier",
	"element",
	"area",
	"service",
	"version",
	"number",
	"domain",
	"timestamp",
	"`details.related`",
	"network",
	"provider",
	"`details.source`",
}

//======================================================================//
//                            RETRIEVE                                  //
//======================================================================//

// RetrieveInArchive : TODO:
func RetrieveInArchive(objectType ObjectType, identifierList IdentifierList, objectInstanceIdentifierList LongList) (ArchiveDetailsList, ElementList, error) {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	// Convert domain
	domain := utils.AdaptDomainToString(identifierList)

	// First of all, we need to verity the object instance identifiers values
	var isAll = false
	for i := 0; i < objectInstanceIdentifierList.Size(); i++ {
		if *objectInstanceIdentifierList[i] == 0 {
			isAll = true
			break
		}
	}

	// Transform Type Short Form to List Short Form
	listShortForm := utils.ConvertToListShortForm(objectType)
	// Get Element in the MAL Registry
	element, err := LookupMALElement(listShortForm)
	if err != nil {
		return nil, nil, err
	}

	// select a.objectInstanceIdentifier, t, y, timestamp, `details.related`, network,
	// provider, `details.source` from Archive a INNER JOIN Sine s ON
	// a.objectInstanceIdentifier = s.objectInstanceIdentifier;

	// Create variables to return the elements and information
	var archiveDetailsList = *NewArchiveDetailsList(0)
	var elementList = element.(ElementList)
	elementList = elementList.CreateElement().(ElementList)
	// Then, retrieve these elements and their information
	if !isAll {
		for i := 0; i < objectInstanceIdentifierList.Size(); i++ {
			// Variables to store the different elements present in the database
			var encodedObjectId []byte
			var encodedElement []byte
			var timestamp time.Time
			var related Long
			var network Identifier
			var provider URI

			// We can retrieve this object
			err = tx.QueryRow("SELECT element, timestamp, `details.related`, network, provider, `details.source` FROM "+TABLE+" WHERE objectInstanceIdentifier = ? AND area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
				*objectInstanceIdentifierList[i],
				objectType.Area,
				objectType.Service,
				objectType.Version,
				objectType.Number,
				domain).Scan(&encodedElement,
				&timestamp,
				&related,
				&network,
				&provider,
				&encodedObjectId)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					return nil, nil, errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
				}
				return nil, nil, err
			}

			// Decode the Element and the ObjectId for the ArchiveDetails
			objectId, element, err := utils.DecodeElements(encodedObjectId, encodedElement)
			if err != nil {
				return nil, nil, err
			}

			// Create the ArchiveDetails
			// First, create the ObjectDetails
			objectDetails := ObjectDetails{&related, objectId}
			// Create the ArchiveDetails
			archiveDetails := &ArchiveDetails{
				*objectInstanceIdentifierList[i],
				objectDetails,
				&network,
				NewFineTime(timestamp),
				&provider,
			}

			archiveDetailsList.AppendElement(archiveDetails)
			elementList.AppendElement(element)
		}
	} else {
		// Retrieve all these elements (no particular object instance identifiers)
		// Variables to store the different elements present in the database
		var objectInstanceIdentifier Long
		var encodedObjectId []byte
		var encodedElement []byte
		var timestamp time.Time
		var related Long
		var network Identifier
		var provider URI

		// Retrieve this object and its archive details in the archive
		rows, err := tx.Query("SELECT objectInstanceIdentifier, element, timestamp, `details.related`, network, provider, `details.source` FROM "+TABLE+" WHERE area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
			objectType.Area,
			objectType.Service,
			objectType.Version,
			objectType.Number,
			domain)
		if err != nil {
			return nil, nil, err
		}

		var countElements int
		for rows.Next() {
			if err = rows.Scan(&objectInstanceIdentifier,
				&encodedElement,
				&timestamp,
				&related,
				&network,
				&provider,
				&encodedObjectId); err != nil {
				return nil, nil, err
			}

			// Decode the Element and the ObjectId for the ArchiveDetails
			objectId, element, err := utils.DecodeElements(encodedObjectId, encodedElement)
			if err != nil {
				return nil, nil, err
			}

			// Create the ArchiveDetails
			// First, create the ObjectDetails
			objectDetails := ObjectDetails{&related, objectId}
			// Create the ArchiveDetails
			archiveDetails := &ArchiveDetails{
				objectInstanceIdentifier,
				objectDetails,
				&network,
				NewFineTime(timestamp),
				&provider,
			}

			archiveDetailsList.AppendElement(archiveDetails)
			elementList.AppendElement(element)
			countElements++
		}

		if countElements == 0 {
			return nil, nil, errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
		}
	}

	// Commit changes
	tx.Commit()

	return archiveDetailsList, elementList, nil
}

//======================================================================//
//                              QUERY                                   //
//======================================================================//

// QueryArchive : TODO:
func QueryArchive(boolean *Boolean, objectType ObjectType, archiveQuery ArchiveQuery, queryFilter QueryFilter) ([]*ObjectType, []*ArchiveDetailsList, []*IdentifierList, []ElementList, error) {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer db.Close()

	// Verify the parameters
	err = verifyParameters(archiveQuery, queryFilter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	var isObjectTypeEqualToZero = objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0

	// First of all we have to create the query
	query, err := createQuery(boolean, objectType, isObjectTypeEqualToZero, archiveQuery, queryFilter)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// Variables to return
	var objectTypeToReturn []*ObjectType
	var identifierListToReturn []*IdentifierList
	var archiveDetailsListToReturn []*ArchiveDetailsList
	var elementListToReturn []ElementList

	if boolean != nil && *boolean == true && isObjectTypeEqualToZero == true {
		// Retrieve all of the elements
		// Variables to store the different elements present in the database
		var objectInstanceIdentifier Long
		var encodedObjectId []byte
		var encodedElement []byte
		var timestamp time.Time
		var related Long
		var network Identifier
		var provider URI
		var area UShort
		var service UShort
		var version UOctet
		var number UShort
		var domain string

		rows, err := tx.Query(query)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		// Define a slice to sort on the differents values to return
		// TODO: find a name for this attribute
		var findASmartName []interface{}
		// Map for the different object types
		var objectTypeMap = make(map[ObjectType]uint)
		var countObjectType uint
		// Map for the different domains
		var domainMap = make(map[string]uint)
		var countDomain uint

		for rows.Next() {
			if err = rows.Scan(&objectInstanceIdentifier, &timestamp, &related, &network, &provider, &encodedObjectId, &encodedElement, &domain, &area, &service, &version, &number); err != nil {
				return nil, nil, nil, nil, err
			}

			//
			var isAlreadyUsed = true
			// Verify the object type value
			objectTypeFromDB := ObjectType{area, service, version, number}
			if _, ok := objectTypeMap[objectTypeFromDB]; !ok {
				objectTypeMap[objectTypeFromDB] = countObjectType
				countObjectType++
				isAlreadyUsed = false
			}
			// Verify the domain value
			if _, ok := domainMap[domain]; !ok {
				domainMap[domain] = countDomain
				countDomain++
				isAlreadyUsed = false
			}

			if isAlreadyUsed {
				var index uint
				// Retrieve the index in the general slice
				duo := []interface{}{objectTypeFromDB, domain}
				for i := 0; i < len(findASmartName)/2; i++ {
					if reflect.DeepEqual(duo, findASmartName[i]) {
						index = uint(i)
						break
					}
				}

				// ArchiveDetailsList
				// Decode the object id
				objId, err := utils.DecodeObjectID(encodedObjectId)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				objectDet := ObjectDetails{
					&related,
					objId,
				}
				archDetails := &ArchiveDetails{objectInstanceIdentifier, objectDet, &network, NewFineTime(timestamp), &provider}
				// Append this ArchiveDetails to the desired ArchiveDetailsList
				archiveDetailsListToReturn[index].AppendElement(archDetails)

				// ElementList
				// Decode the element
				elem, err := utils.DecodeElement(encodedElement)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				elementListToReturn[index].AppendElement(elem)
			} else {
				// First, append the objectType and the domain in the general slice
				findASmartName = append(findASmartName, objectTypeFromDB, domain)

				// IdentifierList
				idList := utils.AdaptDomainToIdentifierList(domain)
				identifierListToReturn = append(identifierListToReturn, &idList)

				// ObjectType
				objectTypeToReturn = append(objectTypeToReturn, &objectTypeFromDB)

				// ArchiveDetailsList
				archDetailsList := NewArchiveDetailsList(0)
				// Decode the object id
				objId, err := utils.DecodeObjectID(encodedObjectId)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				objectDet := ObjectDetails{
					&related,
					objId,
				}
				archDetails := &ArchiveDetails{objectInstanceIdentifier, objectDet, &network, NewFineTime(timestamp), &provider}
				archDetailsList.AppendElement(archDetails)
				archiveDetailsListToReturn = append(archiveDetailsListToReturn, archDetailsList)

				// ElementList
				// Decode the element
				elem, err := utils.DecodeElement(encodedElement)
				if err != nil {
					return nil, nil, nil, nil, err
				}

				// Transform Type Short Form to List Short Form
				listShortForm := utils.ConvertToListShortForm(objectTypeFromDB)
				// Get Element in the MAL Registry
				element, err := LookupMALElement(listShortForm)
				var elementList = element.(ElementList)
				elementList = elementList.CreateElement().(ElementList)

				elementList.AppendElement(elem)
				elementListToReturn = append(elementListToReturn, elementList)
			}
		}
	} else if boolean != nil && *boolean == true && isObjectTypeEqualToZero == false {
		// Retrieve all of the elements unless the object type
		// Variables to store the different elements present in the database
		var objectInstanceIdentifier Long
		var encodedObjectId []byte
		var encodedElement []byte
		var timestamp time.Time
		var related Long
		var network Identifier
		var provider URI
		var domain string
		var area UShort
		var service UShort
		var version UOctet
		var number UShort

		rows, err := tx.Query(query)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		// Map for the different domains
		var domainMap = make(map[string]uint)
		var countDomain uint

		for rows.Next() {
			if err = rows.Scan(&objectInstanceIdentifier, &timestamp, &related, &network, &provider, &encodedObjectId, &encodedElement, &domain, &area, &service, &version, &number); err != nil {
				return nil, nil, nil, nil, err
			}

			//
			var isAlreadyUsed = true
			// Verify the domain value
			if _, ok := domainMap[domain]; !ok {
				domainMap[domain] = countDomain
				countDomain++
				isAlreadyUsed = false
			}

			if isAlreadyUsed {
				// ArchiveDetailsList
				// Decode the object id
				objId, err := utils.DecodeObjectID(encodedObjectId)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				objectDet := ObjectDetails{
					&related,
					objId,
				}
				archDetails := &ArchiveDetails{objectInstanceIdentifier, objectDet, &network, NewFineTime(timestamp), &provider}
				// Append this ArchiveDetails to the desired ArchiveDetailsList
				archiveDetailsListToReturn[domainMap[domain]].AppendElement(archDetails)

				// ElementList
				// Decode the element
				elem, err := utils.DecodeElement(encodedElement)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				elementListToReturn[domainMap[domain]].AppendElement(elem)
			} else {
				// IdentifierList
				idList := utils.AdaptDomainToIdentifierList(domain)
				identifierListToReturn = append(identifierListToReturn, &idList)

				// ObjectType
				objectTypeToReturn = append(objectTypeToReturn, nil)

				// ArchiveDetailsList
				archDetailsList := NewArchiveDetailsList(0)
				// Decode the object id
				objId, err := utils.DecodeObjectID(encodedObjectId)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				objectDet := ObjectDetails{
					&related,
					objId,
				}
				archDetails := &ArchiveDetails{objectInstanceIdentifier, objectDet, &network, NewFineTime(timestamp), &provider}
				archDetailsList.AppendElement(archDetails)
				archiveDetailsListToReturn = append(archiveDetailsListToReturn, archDetailsList)

				// ElementList
				// Decode the element
				elem, err := utils.DecodeElement(encodedElement)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				// Transform Type Short Form to List Short Form
				listShortForm := utils.ConvertToListShortForm(ObjectType{area, service, version, number})
				// Get Element in the MAL Registry
				element, err := LookupMALElement(listShortForm)
				var elementList = element.(ElementList)
				elementList = elementList.CreateElement().(ElementList)

				elementList.AppendElement(elem)
				elementListToReturn = append(elementListToReturn, elementList)
			}
		}
	} else if (boolean == nil || *boolean == false) && isObjectTypeEqualToZero == true {
		// Retrieve only the object type and the archive details
		// Variables to store the different elements present in the database
		var objectInstanceIdentifier Long
		var encodedObjectId []byte
		var timestamp time.Time
		var related Long
		var network Identifier
		var provider URI
		var area UShort
		var service UShort
		var version UOctet
		var number UShort

		rows, err := tx.Query(query)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		// Map for the different object types
		var objectTypeMap = make(map[ObjectType]uint)
		var countObjectType uint

		for rows.Next() {
			if err = rows.Scan(&objectInstanceIdentifier, &timestamp, &related, &network, &provider, &encodedObjectId, &area, &service, &version, &number); err != nil {
				return nil, nil, nil, nil, err
			}

			//
			var isAlreadyUsed = true
			// Verify the object type value
			objectTypeFromDB := ObjectType{area, service, version, number}
			if _, ok := objectTypeMap[objectTypeFromDB]; !ok {
				objectTypeMap[objectTypeFromDB] = countObjectType
				countObjectType++
				isAlreadyUsed = false
			}

			if isAlreadyUsed {
				// ArchiveDetailsList
				// Decode the object id
				objId, err := utils.DecodeObjectID(encodedObjectId)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				objectDet := ObjectDetails{
					&related,
					objId,
				}
				archDetails := &ArchiveDetails{objectInstanceIdentifier, objectDet, &network, NewFineTime(timestamp), &provider}
				// Append this ArchiveDetails to the desired ArchiveDetailsList
				archiveDetailsListToReturn[objectTypeMap[objectTypeFromDB]].AppendElement(archDetails)

			} else {
				// IdentifierList
				identifierListToReturn = append(identifierListToReturn, nil)

				// ObjectType
				objectTypeToReturn = append(objectTypeToReturn, &objectTypeFromDB)

				// ArchiveDetailsList
				archDetailsList := NewArchiveDetailsList(0)
				// Decode the object id
				objId, err := utils.DecodeObjectID(encodedObjectId)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				objectDet := ObjectDetails{
					&related,
					objId,
				}
				archDetails := &ArchiveDetails{objectInstanceIdentifier, objectDet, &network, NewFineTime(timestamp), &provider}
				archDetailsList.AppendElement(archDetails)
				archiveDetailsListToReturn = append(archiveDetailsListToReturn, archDetailsList)

				// ElementList
				var longList *LongList
				elementListToReturn = append(elementListToReturn, longList)
			}
		}
	} else { // (*boolean == false or boolean == nil) and isObjectTypeEqualToZero == false
		// Retrieve only the archive details
		// Variables to store the different elements present in the database
		var objectInstanceIdentifier Long
		var encodedObjectId []byte
		var timestamp time.Time
		var related Long
		var network Identifier
		var provider URI

		rows, err := tx.Query(query)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		// Set the identifierListToReturn to null
		identifierListToReturn = append(identifierListToReturn, nil)
		// Set the objectTypeToReturn to null
		objectTypeToReturn = append(objectTypeToReturn, nil)
		// Set the ElementList to null
		var longList *LongList
		elementListToReturn = append(elementListToReturn, longList)

		var isAlreadyUsed = false
		for rows.Next() {
			if err = rows.Scan(&objectInstanceIdentifier, &timestamp, &related, &network, &provider, &encodedObjectId); err != nil {
				return nil, nil, nil, nil, err
			}

			// ArchiveDetailsList
			// Decode the object id
			objId, err := utils.DecodeObjectID(encodedObjectId)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			objectDet := ObjectDetails{
				&related,
				objId,
			}
			archDetails := &ArchiveDetails{objectInstanceIdentifier, objectDet, &network, NewFineTime(timestamp), &provider}
			// Append this ArchiveDetails to the desired ArchiveDetailsList
			if !isAlreadyUsed {
				// ArchiveDetailsList
				archDetailsList := NewArchiveDetailsList(0)
				archDetailsList.AppendElement(archDetails)
				archiveDetailsListToReturn = append(archiveDetailsListToReturn, archDetailsList)
				isAlreadyUsed = true
			} else {
				archiveDetailsListToReturn[0].AppendElement(archDetails)
			}
		}
	}

	// Finally, it is useful to verify the size of the archive details
	// list (if it is equal to 0 we have to append a nil element to
	// each list)
	if len(archiveDetailsListToReturn) == 0 {
		var longList *LongList
		objectTypeToReturn = append(objectTypeToReturn, nil)
		archiveDetailsListToReturn = append(archiveDetailsListToReturn, nil)
		identifierListToReturn = append(identifierListToReturn, nil)
		elementListToReturn = append(elementListToReturn, longList)
		// TODO: maybe find another way to initialize elementListToReturn
	}

	// Commit changes
	tx.Commit()

	return objectTypeToReturn, archiveDetailsListToReturn, identifierListToReturn, elementListToReturn, nil
}

// verifyParameters : TODO:
func verifyParameters(archiveQuery ArchiveQuery, queryFilter QueryFilter) error {
	// Check sortFieldName value
	var isSortFieldNameADefinedField = false
	for i := 0; i < len(databaseFields); i++ {
		if (archiveQuery.SortFieldName != nil && string(*archiveQuery.SortFieldName) == databaseFields[i]) ||
			archiveQuery.SortFieldName == nil {
			isSortFieldNameADefinedField = true
			break
		}
	}
	if !isSortFieldNameADefinedField {
		// Return a new error
		return errors.New(string(ARCHIVE_SERVICE_QUERY_SORT_FIELD_NAME_INVALID_ERROR))
	}

	// Check if QueryFilter doesn't contain an error
	if queryFilter != nil {
		compositerFilterSet := queryFilter.(*CompositeFilterSet)

		for i := 0; i < compositerFilterSet.Filters.Size(); i++ {
			var filter = compositerFilterSet.Filters.GetElementAt(i).(*CompositeFilter)
			if (filter.Type == COM_EXPRESSIONOPERATOR_CONTAINS || filter.Type == COM_EXPRESSIONOPERATOR_ICONTAINS || filter.Type == COM_EXPRESSIONOPERATOR_GREATER || filter.Type == COM_EXPRESSIONOPERATOR_GREATER_OR_EQUAL || filter.Type == COM_EXPRESSIONOPERATOR_LESS || filter.Type == COM_EXPRESSIONOPERATOR_LESS_OR_EQUAL) && filter.FieldValue == nil {
				return errors.New(string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR) + ": must not contain NULL value")
			} else if _, ok := filter.FieldValue.(*Blob); ok {
				if filter.Type != COM_EXPRESSIONOPERATOR_EQUAL && filter.Type != COM_EXPRESSIONOPERATOR_DIFFER {
					return errors.New(string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR) + ": must not use this expression operator for a blob")
				}
			} else if filter.Type == COM_EXPRESSIONOPERATOR_CONTAINS || filter.Type == COM_EXPRESSIONOPERATOR_ICONTAINS {
				if _, ok := filter.FieldValue.(*String); !ok {
					return errors.New(string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR) + ": must not use this expression operator for a non-String")
				}
			}
		}
	}

	return nil
}

//======================================================================//
//                              COUNT                                   //
//======================================================================//

// CountInArchive : TODO:
func CountInArchive(objectType ObjectType, archiveQueryList ArchiveQueryList, queryFilterList QueryFilterList) (*LongList, error) {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//
	var longList = NewLongList(0)

	for i := 0; i < archiveQueryList.Size(); i++ {
		// Verify the parameters
		if queryFilterList != nil {
			err = verifyParameters(*archiveQueryList[i], queryFilterList.GetElementAt(i))
		} else {
			err = verifyParameters(*archiveQueryList[i], nil)
		}
		if err != nil {
			return nil, err
		}
		// Create the query
		var query string
		if queryFilterList != nil {
			query, err = createCountQuery(objectType, *archiveQueryList[i], queryFilterList.GetElementAt(i))
		} else {
			query, err = createCountQuery(objectType, *archiveQueryList[i], nil)
		}
		if err != nil {
			return nil, err
		}

		// Create a variable to Store the response
		var response int64
		// Execute the query
		err = tx.QueryRow(query).Scan(&response)
		if err != nil {
			return nil, err
		}

		// Add this response in the long list
		longList.AppendElement(NewLong(response))
	}

	// Commit changes
	tx.Commit()

	return longList, nil
}

//======================================================================//
//                              STORE                                   //
//======================================================================//

// StoreInArchive : Use this function to store objects in an COM archive
func StoreInArchive(boolean *Boolean, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*LongList, error) {
	rand.Seed(time.Now().UnixNano())

	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Variable to return all the object instance identifiers
	var longList *LongList

	// Create the domain (It might change in the future)
	domain := utils.AdaptDomainToString(identifierList)

	// Init the list to return (if boolean is not equal to false)
	if boolean != nil && *boolean {
		longList = NewLongList(0)
	}
	for i := 0; i < archiveDetailsList.Size(); i++ {
		if archiveDetailsList[i].InstId == 0 {
			// We have to create a new and unused object instance identifier
			for {
				var objectInstanceIdentifier = rand.Int63n(int64(LONG_MAX))
				isObjInstIDInDB, err := isObjectInstanceIdentifierInDatabase(tx, objectInstanceIdentifier)
				if err != nil {
					// An error occurred, do a rollback
					tx.Rollback()
					return nil, err
				}
				if !isObjInstIDInDB {
					// OK, we can insert the object with this instance identifier
					err := insertInDatabase(tx, objectInstanceIdentifier, elementList.GetElementAt(i), objectType, domain, *archiveDetailsList[i])
					if err != nil {
						// An error occurred, do a rollback
						tx.Rollback()
						return nil, err
					}

					if boolean != nil && *boolean {
						// Insert this new object instance identifier in the returned list
						longList.AppendElement(NewLong(objectInstanceIdentifier))
					}

					break
				}
			}
		} else {
			// We must verify if the object instance identifier is not already present in the table
			isObjInstIDInDB, err := isObjectInstanceIdentifierInDatabase(tx, int64(archiveDetailsList[i].InstId))
			if err != nil {
				// An error occurred, do a rollback
				tx.Rollback()
				return nil, err
			}
			if isObjInstIDInDB {
				// This object is already in the database, do a rollback and raise a DUPLICATE error
				tx.Rollback()
				return nil, errors.New(string(COM_ERROR_DUPLICATE))
			}

			// This object is not present in the archive
			err = insertInDatabase(tx, int64(archiveDetailsList[i].InstId), elementList.GetElementAt(i), objectType, domain, *archiveDetailsList[i])
			if err != nil {
				// An error occurred, do a rollback
				tx.Rollback()
				return nil, err
			}

			if boolean != nil && *boolean {
				// Insert this new object instance identifier in the returned list
				longList.AppendElement(NewLong(int64(archiveDetailsList[i].InstId)))
			}
		}
	}

	// Commit changes
	tx.Commit()

	return longList, nil
}

//======================================================================//
//                              UPDATE                                  //
//======================================================================//

// UpdateArchive : TODO:
func UpdateArchive(objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) error {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the domain (It might change in the future)
	domain := utils.AdaptDomainToString(identifierList)

	for i := 0; i < elementList.Size(); i++ {
		// First of all, we need to verify if the object instance identifier, combined
		// with the object type and the domain which are in the archive
		var queryReturn int
		err := tx.QueryRow("SELECT objectInstanceIdentifier FROM "+TABLE+" WHERE objectInstanceIdentifier = ? AND area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
			archiveDetailsList[i].InstId,
			objectType.Area,
			objectType.Service,
			objectType.Version,
			objectType.Number,
			domain).Scan(&queryReturn)
		if err != nil {
			tx.Rollback()
			if err.Error() == "sql: no rows in result set" {
				return errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
			}
			return err
		}

		encodedElement, encodedObjectId, err := utils.EncodeElements(elementList.GetElementAt(i), *archiveDetailsList[i].Details.Source)
		if err != nil {
			tx.Rollback()
			return err
		}
		// If no error, the object is in the archive and we can update it
		_, err = tx.Exec("UPDATE "+TABLE+" SET element = ?, timestamp = ?, `details.related` = ?, network = ?, provider = ?, `details.source` = ? WHERE objectInstanceIdentifier = ? AND area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
			encodedElement,
			time.Time(*archiveDetailsList[i].Timestamp),
			*archiveDetailsList[i].Details.Related,
			*archiveDetailsList[i].Network,
			*archiveDetailsList[i].Provider,
			encodedObjectId,
			archiveDetailsList[i].InstId,
			objectType.Area,
			objectType.Service,
			objectType.Version,
			objectType.Number,
			domain)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit changes
	tx.Commit()

	return nil
}

//======================================================================//
//                              DELETE                                  //
//======================================================================//

// DeleteInArchive : TODO:
func DeleteInArchive(objectType ObjectType, identifierList IdentifierList, longListRequest LongList) (LongList, error) {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Variable to return
	var longList LongList

	// Create the domain (It might change in the future)
	domain := utils.AdaptDomainToString(identifierList)

	// Variable to say if we have to delete all of the objects or not
	var isAll = false
	for i := 0; i < longListRequest.Size(); i++ {
		if *longListRequest[i] == 0 {
			isAll = true
			break
		}
	}

	if isAll {
		// Retrieve the objectInstanceIdentifier
		rows, err := tx.Query("SELECT objectInstanceIdentifier FROM "+TABLE+" WHERE area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
			objectType.Area,
			objectType.Service,
			objectType.Version,
			objectType.Number,
			domain)
		if err != nil {
			return nil, err
		}

		var countElements int
		for rows.Next() {
			var instID Long
			if err = rows.Scan(&instID); err != nil {
				return nil, err
			}

			longList.AppendElement(&instID)
			countElements++
		}

		if countElements == 0 {
			return nil, errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
		}

		// Delete all these objects
		_, err = tx.Exec("DELETE FROM "+TABLE+" WHERE area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
			objectType.Area,
			objectType.Service,
			objectType.Version,
			objectType.Number,
			domain)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Set AUTO_INCREMENT to max(id)+1
		err = resetAutoIncrement(tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	} else {
		for i := 0; i < longListRequest.Size(); i++ {
			// Check if the object is in the archive
			var objInstID int
			err := tx.QueryRow("SELECT objectInstanceIdentifier FROM "+TABLE+" WHERE objectInstanceIdentifier = ? AND area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
				*longListRequest[i],
				objectType.Area,
				objectType.Service,
				objectType.Version,
				objectType.Number,
				domain).Scan(&objInstID)
			if err != nil {
				tx.Rollback()
				if err.Error() == "sql: no rows in result set" {
					return nil, errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
				}
				return nil, err
			}

			_, err = tx.Exec("DELETE FROM "+TABLE+" WHERE objectInstanceIdentifier = ? AND area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
				*longListRequest[i],
				objectType.Area,
				objectType.Service,
				objectType.Version,
				objectType.Number,
				domain)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			longList.AppendElement(longListRequest.GetElementAt(i))
		}

		// Set AUTO_INCREMENT to max(id)+1
		err = resetAutoIncrement(tx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit changes
	tx.Commit()

	return longList, nil
}

//======================================================================//
//                           LOCAL FUNCTIONS                            //
//======================================================================//
// createTransaction : TODO:
func createTransaction() (*sql.DB, *sql.Tx, error) {
	// Open the database
	db, err := sql.Open("mysql", USERNAME+":"+PASSWORD+"@/"+DATABASE+"?parseTime=true")
	if err != nil {
		return nil, nil, err
	}

	// Validate the connection by pinging it
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	// Create the transaction (we have to use this method to use rollback and commit)
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}

	return db, tx, nil
}

// isObjectInstanceIdentifierInDatabase: This function allows to verify if an instance of
// an object is already in the archive
func isObjectInstanceIdentifierInDatabase(tx *sql.Tx, objectInstanceIdentifier int64) (bool, error) {
	// Execute the query
	// Before, create a variable to retrieve the result
	var queryReturn int
	// Then, execute the query
	err := tx.QueryRow("SELECT objectInstanceIdentifier FROM "+TABLE+" WHERE objectInstanceIdentifier = ? ", objectInstanceIdentifier).Scan(&queryReturn)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

// insertInDatabase: This function allows to insert an element in the archive
func insertInDatabase(tx *sql.Tx, objectInstanceIdentifier int64, element Element, objectType ObjectType, domain String, archiveDetails ArchiveDetails) error {
	// Encode the Element and the ObjectId from the ArchiveDetails
	encodedElement, encodedObjectID, err := utils.EncodeElements(element, *archiveDetails.Details.Source)
	if err != nil {
		return err
	}

	// Execute the query to insert all the values in the database
	_, err = tx.Exec("INSERT INTO "+TABLE+" VALUES ( NULL , ? , ? , ? , ? , ? , ? , ? , ? , ? , ? , ? , ? )",
		objectInstanceIdentifier,
		encodedElement,
		objectType.Area,
		objectType.Service,
		objectType.Version,
		objectType.Number,
		domain,
		time.Time(*archiveDetails.Timestamp),
		*archiveDetails.Details.Related,
		*archiveDetails.Network,
		*archiveDetails.Provider,
		encodedObjectID)
	if err != nil {
		return err
	}
	return nil
}

// resetAutoIncrement takes the maximum id in the database and set the
// AUTO_INCREMENT at this value (actually it's this value to which we added 1)
func resetAutoIncrement(tx *sql.Tx) error {
	// Create a value to store the maximum id (to which we added 1)
	_, err := tx.Exec("SELECT @max := max(id)+1 FROM " + TABLE)
	if err != nil {
		return err
	}
	// Create the statement (in our case, we use an "alter table" statement)
	_, err = tx.Exec("SET @alter_statement = CONCAT('ALTER TABLE " + TABLE + " AUTO_INCREMENT = ', @max)")
	if err != nil {
		return err
	}
	// Prepare the statement
	_, err = tx.Exec("PREPARE stmt1 FROM @alter_statement;")
	if err != nil {
		return err
	}
	// Then execute the statement
	_, err = tx.Exec("EXECUTE stmt1")
	if err != nil {
		return err
	}
	// Finally, deallocate the statement
	_, err = tx.Exec("DEALLOCATE PREPARE stmt1;")
	if err != nil {
		return err
	}
	return nil
}

// createCountQuery allows the provider to create automatically a query for the Count operation
func createCountQuery(objectType ObjectType, archiveQuery ArchiveQuery, queryFilter QueryFilter) (string, error) {
	var queryBuffer bytes.Buffer
	// Only CompositeFilterSet type should be used
	queryBuffer.WriteString("SELECT COUNT(id)")

	err := createCommonQuery(&queryBuffer, objectType, archiveQuery, queryFilter)
	if err != nil {
		return "", err
	}

	return queryBuffer.String(), nil
}

// createQuery allows the provider to create automatically a query for the Query operation
func createQuery(boolean *Boolean, objectType ObjectType, isObjectTypeEqualToZero bool, archiveQuery ArchiveQuery, queryFilter QueryFilter) (string, error) {
	var queryBuffer bytes.Buffer
	// Only CompositeFilterSet type should be used
	queryBuffer.WriteString("SELECT objectInstanceIdentifier, timestamp, `details.related`, network, provider, `details.source`")
	// Check if we need to retrieve the element and its domain
	if boolean != nil && *boolean == true {
		queryBuffer.WriteString(", element, domain")
	}
	// If there's a wildcard value in one of the object type
	// fields then we have to retrieve the entire object type
	if isObjectTypeEqualToZero == true || (isObjectTypeEqualToZero == false && (boolean != nil && *boolean == true)) {
		queryBuffer.WriteString(", area, service, version, number")
	}

	err := createCommonQuery(&queryBuffer, objectType, archiveQuery, queryFilter)
	if err != nil {
		return "", err
	}

	return queryBuffer.String(), nil
}

// createCommonQuery is a common way of generating a part of a query
func createCommonQuery(queryBuffer *bytes.Buffer, objectType ObjectType, archiveQuery ArchiveQuery, queryFilter QueryFilter) error {
	// Prepare the query for the conditions
	queryBuffer.WriteString(" FROM " + TABLE + " WHERE")

	// Attribute to check if there is already a condition before
	var isThereAlreadyACondition = false

	// Conditions on the object type attributes
	// Area
	if objectType.Area != 0 {
		queryBuffer.WriteString(fmt.Sprintf(" area = %d", objectType.Area))
		isThereAlreadyACondition = true
	}
	// Service
	if objectType.Service != 0 {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" service = %d", objectType.Service))
	}
	// Version
	if objectType.Version != 0 {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" version = %d", objectType.Version))
	}
	// Number
	if objectType.Number != 0 {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" number = %d", objectType.Number))
	}

	// Add archive query conditions
	// Domain
	if archiveQuery.Domain != nil {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
		domain := utils.AdaptDomainToString(*archiveQuery.Domain)
		queryBuffer.WriteString(fmt.Sprintf(" domain = '%s'", domain))
	}

	// Network
	if archiveQuery.Network != nil {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" network = '%s'", *archiveQuery.Network))
	}

	// Provider
	if archiveQuery.Provider != nil {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" provider = '%s'", *archiveQuery.Provider))
	}

	// Related (always have to do a query with this condition)
	utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
	queryBuffer.WriteString(fmt.Sprintf(" `details.related` = %d", archiveQuery.Related))

	// Source
	if archiveQuery.Source != nil {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)

		// Encode the ObjectId
		// Create the factory
		factory := new(FixedBinaryEncoding)

		// Create the encoder
		encoder := factory.NewEncoder(make([]byte, 0, 8192))

		// Encode it
		err := archiveQuery.Source.Encode(encoder)
		if err != nil {
			return err
		}
		queryBuffer.WriteString(fmt.Sprintf(" `details.source` = %s", encoder.Body()))
	}

	// StartTime
	if archiveQuery.StartTime != nil {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" timestamp >= %s", time.Time(*archiveQuery.StartTime)))
	}

	// EndTime
	if archiveQuery.EndTime != nil {
		utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" timestamp <= %s", time.Time(*archiveQuery.EndTime)))
	}

	// Add query filter conditions
	if queryFilter != nil {
		compositerFilterSet := queryFilter.(*CompositeFilterSet)

		for i := 0; i < compositerFilterSet.Filters.Size(); i++ {
			utils.CheckCondition(&isThereAlreadyACondition, queryBuffer)
			var fieldValue = (*compositerFilterSet.Filters)[i].FieldValue
			// Transform the expresion operator
			expressionOperator := (*compositerFilterSet.Filters)[i].Type.TransformOperator()
			if (*compositerFilterSet.Filters)[i].Type == COM_EXPRESSIONOPERATOR_CONTAINS || (*compositerFilterSet.Filters)[i].Type == COM_EXPRESSIONOPERATOR_ICONTAINS {
				queryBuffer.WriteString(fmt.Sprintf(" %s %s", (*compositerFilterSet.Filters)[i].FieldName,
					expressionOperator))
				queryBuffer.WriteString(fmt.Sprintf("%v", reflect.ValueOf(fieldValue).Elem().Interface()))
				queryBuffer.WriteString("%'")
			} else {
				if fieldValue == nil {
					queryBuffer.WriteString(fmt.Sprintf(" %s %s %s", (*compositerFilterSet.Filters)[i].FieldName,
						expressionOperator,
						"NULL"))
				} else {
					queryBuffer.WriteString(fmt.Sprintf(" %s %s", (*compositerFilterSet.Filters)[i].FieldName,
						expressionOperator))
					queryBuffer.WriteString(fmt.Sprintf(" %v", reflect.ValueOf(fieldValue).Elem().Interface()))
				}
			}
		}
	}

	// SortOrder
	if archiveQuery.SortOrder != nil {
		// SortFieldName
		if archiveQuery.SortFieldName != nil {
			queryBuffer.WriteString(fmt.Sprintf(" ORDER BY %s", *archiveQuery.SortFieldName))
		} else {
			queryBuffer.WriteString(" ORDER BY timestamp")
		}
		// If sortOrder is false then returned values shall be sorted
		// in descending order (ascending order is the default value)
		if *archiveQuery.SortOrder == false {
			queryBuffer.WriteString(" DESC")
		}
	}

	return nil
}
