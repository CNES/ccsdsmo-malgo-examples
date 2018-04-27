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
package storage

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"strings"
	"time"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/encoding/binary"

	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/data"

	_ "github.com/go-sql-driver/mysql"
)

// Database ids
const (
	USERNAME = "archiveService"
	PASSWORD = "1a2B3c4D!@?"
	DATABASE = "archive"
	TABLE    = "Archive"
)

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
	"details.related",
	"network",
	"provider",
	"details.source",
}

//======================================================================//
//                            RETRIEVE                                  //
//======================================================================//
func RetrieveInArchive(objectType ObjectType, identifierList IdentifierList, objectInstanceIdentifierList LongList) (ArchiveDetailsList, ElementList, error) {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	// Convert domain
	domain := adaptDomainToString(identifierList)

	// First of all, we need to verity the object instance identifiers values
	var isAll = false
	for i := 0; i < objectInstanceIdentifierList.Size(); i++ {
		if *objectInstanceIdentifierList[i] == 0 {
			isAll = true
			break
		}
	}

	// Transform Type Short Form to List Short Form
	listShortForm := convertToListShortForm(objectType)
	// Get Element in the MAL Registry
	element, err := LookupMALElement(listShortForm)
	if err != nil {
		return nil, nil, err
	}

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
			objectId, element, err := decodeElements(encodedObjectId, encodedElement)
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
		// Retrieve all these elements (no particular object instance iedentifiers)
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
			objectId, element, err := decodeElements(encodedObjectId, encodedElement)
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

	return archiveDetailsList, elementList, nil
}

//======================================================================//
//                              QUERY                                   //
//======================================================================//
func QueryArchive(boolean *Boolean, objectType ObjectType, archiveQuery ArchiveQuery, queryFilter QueryFilter) ([]*ObjectType, []*ArchiveDetailsList, []*IdentifierList, []ElementList, error) {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer db.Close()

	// Verify the parameters
	err = queryVerifyParameters(archiveQuery, queryFilter)
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

	if *boolean == true && isObjectTypeEqualToZero == true {
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
				for i := 0; i < len(findASmartName)/2; i++ {
					duo := []interface{}{objectTypeFromDB, domain}
					if reflect.DeepEqual(duo, findASmartName[i]) {
						index = uint(i)
						break
					}
				}

				// ArchiveDetailsList
				// Decode the object id
				objId, err := decodeObjectId(encodedObjectId)
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
				elem, err := decodeElement(encodedElement)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				elementListToReturn[index].AppendElement(elem)
			} else {
				// First, append the objectType and the domain in the general slice
				findASmartName = append(findASmartName, objectTypeFromDB, domain)

				// IdentifierList
				idList := adaptDomainToIdentifierList(domain)
				identifierListToReturn = append(identifierListToReturn, &idList)

				// ObjectType
				objectTypeToReturn = append(objectTypeToReturn, &objectTypeFromDB)

				// ArchiveDetailsList
				archDetailsList := NewArchiveDetailsList(0)
				// Decode the object id
				objId, err := decodeObjectId(encodedObjectId)
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
				elem, err := decodeElement(encodedElement)
				if err != nil {
					return nil, nil, nil, nil, err
				}

				// Transform Type Short Form to List Short Form
				listShortForm := convertToListShortForm(objectTypeFromDB)
				// Get Element in the MAL Registry
				element, err := LookupMALElement(listShortForm)
				var elementList = element.(ElementList)
				elementList = elementList.CreateElement().(ElementList)

				elementList.AppendElement(elem)
				elementListToReturn = append(elementListToReturn, elementList)
			}
		}
	} else if *boolean == true && isObjectTypeEqualToZero == false {
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
				objId, err := decodeObjectId(encodedObjectId)
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
				elem, err := decodeElement(encodedElement)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				elementListToReturn[domainMap[domain]].AppendElement(elem)
			} else {
				// IdentifierList
				idList := adaptDomainToIdentifierList(domain)
				identifierListToReturn = append(identifierListToReturn, &idList)

				// ObjectType
				objectTypeToReturn = append(objectTypeToReturn, new(ObjectType))

				// ArchiveDetailsList
				archDetailsList := NewArchiveDetailsList(0)
				// Decode the object id
				objId, err := decodeObjectId(encodedObjectId)
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
				elem, err := decodeElement(encodedElement)
				if err != nil {
					return nil, nil, nil, nil, err
				}
				// Transform Type Short Form to List Short Form
				listShortForm := convertToListShortForm(ObjectType{area, service, version, number})
				// Get Element in the MAL Registry
				element, err := LookupMALElement(listShortForm)
				var elementList = element.(ElementList)
				elementList = elementList.CreateElement().(ElementList)

				elementList.AppendElement(elem)
				elementListToReturn = append(elementListToReturn, elementList)
			}
		}
	} else if *boolean == false && isObjectTypeEqualToZero == true {
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
				objId, err := decodeObjectId(encodedObjectId)
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
				identifierListToReturn = append(identifierListToReturn, new(IdentifierList))

				// ObjectType
				objectTypeToReturn = append(objectTypeToReturn, &objectTypeFromDB)

				// ArchiveDetailsList
				archDetailsList := NewArchiveDetailsList(0)
				// Decode the object id
				objId, err := decodeObjectId(encodedObjectId)
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
				elementListToReturn = append(elementListToReturn, new(LongList))
			}
		}
	} else { // boolean == false and isObjectTypeEqualToZero == false
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

		// Set the identifierListToReturn to nul
		identifierListToReturn = append(identifierListToReturn, new(IdentifierList))
		// Set the objectTypeToReturn to nul
		objectTypeToReturn = append(objectTypeToReturn, new(ObjectType))
		// Set the ElementList to nul
		elementListToReturn = append(elementListToReturn, new(LongList))

		for rows.Next() {
			if err = rows.Scan(&objectInstanceIdentifier, &timestamp, &related, &network, &provider, &encodedObjectId); err != nil {
				return nil, nil, nil, nil, err
			}

			// ArchiveDetailsList
			// Decode the object id
			objId, err := decodeObjectId(encodedObjectId)
			if err != nil {
				return nil, nil, nil, nil, err
			}
			objectDet := ObjectDetails{
				&related,
				objId,
			}
			archDetails := &ArchiveDetails{objectInstanceIdentifier, objectDet, &network, NewFineTime(timestamp), &provider}
			// Append this ArchiveDetails to the desired ArchiveDetailsList
			archiveDetailsListToReturn[0].AppendElement(archDetails)
		}
	}

	// Finally, it is useful to verify the size of the archive details
	// list (if it is equal to 0 we have to append a nil element to
	// each list)
	if len(archiveDetailsListToReturn) == 0 {
		objectTypeToReturn = append(objectTypeToReturn, new(ObjectType))
		archiveDetailsListToReturn = append(archiveDetailsListToReturn, new(ArchiveDetailsList))
		identifierListToReturn = append(identifierListToReturn, new(IdentifierList))
		elementListToReturn = append(elementListToReturn, NewLongList(0))
		// TODO: maybe find another way to initialize elementListToReturn
	}

	return objectTypeToReturn, archiveDetailsListToReturn, identifierListToReturn, elementListToReturn, nil
}

func queryVerifyParameters(archiveQuery ArchiveQuery, queryFilter QueryFilter) error {
	// Check sortFieldName value
	var isSortFieldNameADefinedField = false
	for i := 0; i < len(databaseFields); i++ {
		if archiveQuery.SortFieldName != nil && string(*archiveQuery.SortFieldName) == databaseFields[i] {
			isSortFieldNameADefinedField = true
			break
		}
	}
	if !isSortFieldNameADefinedField {
		// Return a new error
		return errors.New(string(ARCHIVE_SERVICE_QUERY_SORT_FIELD_NAME_INVALID_ERROR))
	}

	// Check if QueryFilter doesn't contain an error
	compositerFilterSet := queryFilter.(*CompositeFilterSet)

	for i := 0; i < compositerFilterSet.Filters.Size(); i++ {
		var filter = compositerFilterSet.Filters.GetElementAt(i).(*CompositeFilter)
		if (filter.Type == COM_EXPRESSIONOPERATOR_CONTAINS ||
			filter.Type == COM_EXPRESSIONOPERATOR_ICONTAINS ||
			filter.Type == COM_EXPRESSIONOPERATOR_GREATER ||
			filter.Type == COM_EXPRESSIONOPERATOR_GREATER_OR_EQUAL ||
			filter.Type == COM_EXPRESSIONOPERATOR_LESS ||
			filter.Type == COM_EXPRESSIONOPERATOR_LESS_OR_EQUAL) && filter.FieldValue == Attribute(nil) {
			return errors.New(string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR) + ": must not contain NULL value")
		}
		if _, ok := filter.FieldValue.(*Blob); ok {
			if filter.Type != COM_EXPRESSIONOPERATOR_EQUAL && filter.Type != COM_EXPRESSIONOPERATOR_DIFFER {
				return errors.New(string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR) + ": must not use this expression operator for a blob")
			}
		}
		if filter.Type == COM_EXPRESSIONOPERATOR_CONTAINS || filter.Type == COM_EXPRESSIONOPERATOR_ICONTAINS {
			if _, ok := filter.FieldValue.(*String); !ok {
				return errors.New(string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR) + ": must not use this expression operator for a non-String")
			}
		}
	}

	return nil
}

//======================================================================//
//                              COUNT                                   //
//======================================================================//
func CountInArchive() error {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return err
	}
	defer db.Close()

	fmt.Println(tx)

	return nil
}

//======================================================================//
//                              STORE                                   //
//======================================================================//
// StoreInArchive : Use this function to store objects in an COM archive
func StoreInArchive(objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (LongList, error) {
	rand.Seed(time.Now().UnixNano())

	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Variable to return all the object instance identifiers
	var longList LongList

	// Create the domain (It might change in the future)
	domain := adaptDomainToString(identifierList)

	for i := 0; i < archiveDetailsList.Size(); i++ {
		if archiveDetailsList[i].InstId == 0 {
			// We have to create a new and unused object instance identifier
			for {
				var objectInstanceIdentifier = rand.Int63n(int64(LONG_MAX))
				boolean, err := isObjectInstanceIdentifierInDatabase(tx, objectInstanceIdentifier)
				if err != nil {
					// An error occurred, do a rollback
					tx.Rollback()
					return nil, err
				}
				if !boolean {
					// OK, we can insert the object with this instance identifier
					err := insertInDatabase(tx, objectInstanceIdentifier, elementList.GetElementAt(i), objectType, domain, *archiveDetailsList[i])
					if err != nil {
						// An error occurred, do a rollback
						tx.Rollback()
						return nil, err
					}

					// Insert this new object instance identifier in the returned list
					longList.AppendElement(NewLong(objectInstanceIdentifier))

					break
				}
			}
		} else {
			// We must verify if the object instance identifier is not already present in the table
			boolean, err := isObjectInstanceIdentifierInDatabase(tx, int64(archiveDetailsList[i].InstId))
			if err != nil {
				// An error occurred, do a rollback
				tx.Rollback()
				return nil, err
			}
			if boolean {
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

			// Insert this new object instance identifier in the returned list
			longList.AppendElement(&archiveDetailsList[i].InstId)
		}
	}

	// Commit changes
	tx.Commit()

	return longList, nil
}

//======================================================================//
//                              UPDATE                                  //
//======================================================================//
func UpdateArchive(objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) error {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return err
	}
	defer db.Close()

	// Create the domain (It might change in the future)
	domain := adaptDomainToString(identifierList)

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

		encodedElement, encodedObjectId, err := encodeElements(elementList.GetElementAt(i), *archiveDetailsList[i].Details.Source)
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
	domain := adaptDomainToString(identifierList)

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
func createTransaction() (*sql.DB, *sql.Tx, error) {
	// Create the handle
	db, err := sql.Open("mysql", USERNAME+":"+PASSWORD+"@/"+DATABASE+"?parseTime=true")
	if err != nil {
		return nil, nil, err
	}

	// Validate DSN data
	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}

	// Create the transaction (me have to use this method to use rollback and commit)
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}

	return db, tx, nil
}

// This function allows to verify if an instance of an object is already in the archive
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

// This function allows insert an element in the archive
func insertInDatabase(tx *sql.Tx, objectInstanceIdentifier int64, element Element, objectType ObjectType, domain String, archiveDetails ArchiveDetails) error {
	// Encode the Element and the ObjectId from the ArchiveDetails
	encodedElement, encodedObjectID, err := encodeElements(element, *archiveDetails.Details.Source)
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

// resetAutoIncrement take the maximum id in the database and set the
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

// adaptDomainToString transforms a list of Identifiers to a domain of this
// type: first.second.third.[...]
func adaptDomainToString(identifierList IdentifierList) String {
	var domain String
	for i := 0; i < identifierList.Size(); i++ {
		domain += String(*identifierList.GetElementAt(i).(*Identifier))
		if i+1 < identifierList.Size() {
			domain += "."
		}
	}
	return domain
}

// adaptDomainToIdentifierList transforms a domain of this
// type: first.second.third.[...] to a list of Identifiers
func adaptDomainToIdentifierList(domain string) IdentifierList {
	var identifierList = NewIdentifierList(0)
	var domains = strings.Split(domain, ".")
	for i := 0; i < len(domains); i++ {
		identifierList.AppendElement(NewIdentifier(domains[i]))
	}
	return *identifierList
}

func decodeObjectId(encodedObjectId []byte) (*ObjectId, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the decoder
	decoder := factory.NewDecoder(encodedObjectId)

	// Decode the ObjectId
	elem, err := decoder.DecodeElement(NullObjectId)
	if err != nil {
		return nil, err
	}
	objectId := elem.(*ObjectId)

	return objectId, nil
}

func decodeElement(encodedObjectElement []byte) (Element, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the decoder
	decoder := factory.NewDecoder(encodedObjectElement)

	// Decode the Element
	element, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, err
	}

	return element, nil
}

func decodeElements(_objectId []byte, _element []byte) (*ObjectId, Element, error) {
	// Decode the ObjectId
	objectId, err := decodeObjectId(_objectId)
	if err != nil {
		return nil, nil, err
	}

	// Decode the Element
	element, err := decodeElement(_element)
	if err != nil {
		return nil, nil, err
	}

	return objectId, element, nil
}

func encodeElements(_element Element, _objectId ObjectId) ([]byte, []byte, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the encoder
	encoder := factory.NewEncoder(make([]byte, 0, 8192))

	// Encode Element
	err := encoder.EncodeAbstractElement(_element)
	if err != nil {
		return nil, nil, err
	}
	element := encoder.Body()

	// Reallocate the encoder
	encoder = factory.NewEncoder(make([]byte, 0, 8192))

	// Encode ObjectId
	err = _objectId.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}
	objectId := encoder.Body()

	return element, objectId, nil
}

// This part is useful for type short form conversion (from typeShortForm to listShortForm)
func typeShortFormToShortForm(objectType ObjectType) Long {
	var typeShortForm = Long(objectType.Number) | 0xFFFFFFF000000
	var areaVersion = (Long(objectType.Version) << 24) | 0xFFFFF00FFFFFF
	var serviceNumber = (Long(objectType.Service) << 32) | 0xF0000FFFFFFFF
	var areaNumber = (Long(objectType.Area) << 48) | 0x0FFFFFFFFFFFF

	return areaNumber & serviceNumber & areaVersion & typeShortForm
}

// convertToListShortForm converts an ObjectType to a Long (which
// will be used for a List Short Form)
func convertToListShortForm(objectType ObjectType) Long {
	var listByte []byte
	listByte = append(listByte, byte(objectType.Area), byte(objectType.Service>>8), byte(objectType.Service), byte(objectType.Version))
	typeShort := typeShortFormToShortForm(objectType)
	quatuor5 := (typeShort & 0x0000F0) >> 4
	if quatuor5 == 0x0 {
		var b byte
		for i := 2; i >= 0; i-- {
			b = byte(typeShort>>uint(i*8)) ^ 255
			if i == 0 {
				b++
			}
			listByte = append(listByte, b)
		}

		var byte0 = Long(listByte[6]) | 0xFFFFFFFFFFF00
		var byte1 = (Long(listByte[5]) << 8) | 0xFFFFFFFFF00FF
		var byte2 = (Long(listByte[4]) << 16) | 0xFFFFFFF00FFFF
		var byte3 = (Long(listByte[3]) << 24) | 0xFFFFF00FFFFFF
		var byte4 = (Long(listByte[2]) << 32) | 0xFFF00FFFFFFFF
		var byte5 = (Long(listByte[1]) << 40) | 0xF00FFFFFFFFFF
		var byte6 = (Long(listByte[0]) << 48) | 0x0FFFFFFFFFFFF

		return byte6 & byte5 & byte4 & byte3 & byte2 & byte1 & byte0
	}

	// Force bytes 2, 3 to 1
	return typeShort | 0x0000000FFFF00
}

func checkCondition(cond *bool, buffer *bytes.Buffer) {
	if *cond {
		buffer.WriteString(" AND")
	} else {
		*cond = true
	}
}

// createQuery allows the provider to create automatically a query
func createQuery(boolean *Boolean, objectType ObjectType, isObjectTypeEqualToZero bool, archiveQuery ArchiveQuery, queryFilter QueryFilter) (string, error) {
	var queryBuffer bytes.Buffer
	// Only CompositeFilterSet type should be used
	queryBuffer.WriteString("SELECT objectInstanceIdentifier, timestamp, `details.related`, network, provider, `details.source`")
	// Check if we need to retrieve the element and its domain
	if *boolean == true {
		queryBuffer.WriteString(", element, domain")
	}
	// If there's a wildcard value in one of the object type
	// fields then we have to retrieve the entire object type
	if isObjectTypeEqualToZero == true || (isObjectTypeEqualToZero == false && *boolean == true) {
		queryBuffer.WriteString(", area, service, version, number")
	}

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
		checkCondition(&isThereAlreadyACondition, &queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" service = %d", objectType.Service))
	}
	// Version
	if objectType.Version != 0 {
		checkCondition(&isThereAlreadyACondition, &queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" version = %d", objectType.Version))
	}
	// Number
	if objectType.Number != 0 {
		checkCondition(&isThereAlreadyACondition, &queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" number = %d", objectType.Number))
	}

	// Add archive query conditions
	// Domain
	if archiveQuery.Domain != nil {
		checkCondition(&isThereAlreadyACondition, &queryBuffer)
		domain := adaptDomainToString(*archiveQuery.Domain)
		queryBuffer.WriteString(fmt.Sprintf(" domain = %s", domain))
	}

	// Network
	if archiveQuery.Network != nil {
		checkCondition(&isThereAlreadyACondition, &queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" network = %s", *archiveQuery.Network))
	}

	// Provider
	if archiveQuery.Provider != nil {
		checkCondition(&isThereAlreadyACondition, &queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" provider = %s", *archiveQuery.Provider))
	}

	// Related (always have to do a query with this condition)
	checkCondition(&isThereAlreadyACondition, &queryBuffer)
	queryBuffer.WriteString(fmt.Sprintf(" `details.related` = %d", archiveQuery.Related))

	// Source
	if archiveQuery.Source != nil {
		checkCondition(&isThereAlreadyACondition, &queryBuffer)

		// Encode the ObjectId
		// Create the factory
		factory := new(FixedBinaryEncoding)

		// Create the encoder
		encoder := factory.NewEncoder(make([]byte, 0, 8192))

		// Encode it
		err := archiveQuery.Source.Encode(encoder)
		if err != nil {
			return "", err
		}
		queryBuffer.WriteString(fmt.Sprintf(" provider = %s", encoder.Body()))
	}

	// StartTime
	if archiveQuery.StartTime != nil {
		checkCondition(&isThereAlreadyACondition, &queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" timestamp >= %s", time.Time(*archiveQuery.StartTime)))
	}

	// EndTime
	if archiveQuery.EndTime != nil {
		checkCondition(&isThereAlreadyACondition, &queryBuffer)
		queryBuffer.WriteString(fmt.Sprintf(" timestamp <= %s", time.Time(*archiveQuery.EndTime)))
	}

	// Add query filter conditions
	if queryFilter != nil {
		compositerFilterSet := queryFilter.(*CompositeFilterSet)

		for i := 0; i < compositerFilterSet.Filters.Size(); i++ {
			checkCondition(&isThereAlreadyACondition, &queryBuffer)
			// Transform the expresion operator
			expressionOperator := whichExpressionOperatorIsIt((*compositerFilterSet.Filters)[i].Type)

			if (*compositerFilterSet.Filters)[i].Type == COM_EXPRESSIONOPERATOR_CONTAINS || (*compositerFilterSet.Filters)[i].Type == COM_EXPRESSIONOPERATOR_ICONTAINS {
				queryBuffer.WriteString(fmt.Sprintf(" %s %s", *(*compositerFilterSet.Filters)[i].FieldName,
					expressionOperator))
				queryBuffer.WriteString("%'")
				queryBuffer.WriteString(fmt.Sprintf(" %s", (*compositerFilterSet.Filters)[i].FieldValue))
			} else {
				queryBuffer.WriteString(fmt.Sprintf(" %s %s %s", *(*compositerFilterSet.Filters)[i].FieldName,
					expressionOperator,
					(*compositerFilterSet.Filters)[i].FieldValue))
			}
		}
	}

	// SortOrder
	if archiveQuery.SortOrder != nil {
		// SortFieldName
		if archiveQuery.SortFieldName != nil {
			queryBuffer.WriteString(fmt.Sprintf(" GROUP BY %s", *archiveQuery.SortFieldName))
		} else {
			queryBuffer.WriteString(" GROUP BY timestamp")
		}
		// If sortOrder is false then returned values shall be sorted
		// in descending order (ascending order is the default value)
		if *archiveQuery.SortOrder == false {
			queryBuffer.WriteString(" DESC")
		}
	}

	return queryBuffer.String(), nil
}

// whichExpressionOperatorIsIt transforms an ExpressionOperator to a string
func whichExpressionOperatorIsIt(expressionOperator ExpressionOperator) string {
	switch expressionOperator {
	case COM_EXPRESSIONOPERATOR_EQUAL:
		return "="
	case COM_EXPRESSIONOPERATOR_DIFFER:
		return "!="
	case COM_EXPRESSIONOPERATOR_GREATER:
		return ">"
	case COM_EXPRESSIONOPERATOR_GREATER_OR_EQUAL:
		return ">="
	case COM_EXPRESSIONOPERATOR_LESS:
		return "<"
	case COM_EXPRESSIONOPERATOR_LESS_OR_EQUAL:
		return "<="
	case COM_EXPRESSIONOPERATOR_CONTAINS:
		return "LIKE '%"
	case COM_EXPRESSIONOPERATOR_ICONTAINS:
		return "NOT LIKE '%"
	default:
		return ""
	}
}
