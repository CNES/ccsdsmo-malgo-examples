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
			var related *Long
			var network *Identifier
			var provider *URI

			// We can retrieve this object
			err = tx.QueryRow("SELECT element, timestamp, `details.related`, network, provider, source FROM "+TABLE+" WHERE objectInstanceIdentifier = ? AND area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
				*objectInstanceIdentifierList[i],
				objectType.Area,
				objectType.Service,
				objectType.Version,
				objectType.Number,
				domain).Scan(&encodedElement,
				&timestamp,
				related,
				network,
				provider,
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
			objectDetails := ObjectDetails{related, objectId}
			// Create the ArchiveDetails
			archiveDetails := &ArchiveDetails{
				*objectInstanceIdentifierList[i],
				objectDetails,
				network,
				NewFineTime(timestamp),
				provider,
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
		var related *Long
		var network *Identifier
		var provider *URI

		// Retrieve this object and its archive details in the archive
		rows, err := tx.Query("SELECT objectInstanceIdentifier, element, timestamp, `details.related`, network, provider, source FROM "+TABLE+" WHERE area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
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
			if err = rows.Scan(&objectInstanceIdentifier, &encodedElement,
				&timestamp,
				related,
				network,
				provider,
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
			objectDetails := ObjectDetails{related, objectId}
			// Create the ArchiveDetails
			archiveDetails := &ArchiveDetails{
				objectInstanceIdentifier,
				objectDetails,
				network,
				NewFineTime(timestamp),
				provider,
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
func QueryArchive(boolean *Boolean, objectType ObjectType, archiveQuery ArchiveQuery, queryFilter QueryFilter) (*ObjectType, *ArchiveDetailsList, *IdentifierList, ElementList, error) {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, nil, nil, nil, err
	}
	defer db.Close()

	var isObjectTypeEqualToZero = objectType.Area == 0 || objectType.Number == 0 || objectType.Service == 0 || objectType.Version == 0

	if isObjectTypeEqualToZero {

	} else {

	}

	fmt.Println(tx)

	return nil, nil, nil, nil, nil
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
		// First of all, we need to verify if the object instance identifier, combined with the object type
		// and the domain are in the archive
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
		_, err = tx.Exec("UPDATE "+TABLE+" SET element = ?, timestamp = ?, related = ?, network = ?, provider = ?, source = ? WHERE objectInstanceIdentifier = ? AND area = ? AND service = ? AND version = ? AND number = ? AND domain = ?",
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
//                          GLOBAL FUNCTIONS                            //
//======================================================================//
func createTransaction() (*sql.DB, *sql.Tx, error) {
	// Create the handle
	db, err := sql.Open("mysql", USERNAME+":"+PASSWORD+"@/"+DATABASE)
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

// Transform a list of Identifiers to a domain of this
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

func adaptDomainToIdentifierList(domain string) IdentifierList {
	var identifierList = NewIdentifierList(0)
	var domains = strings.Split(domain, ".")
	for i := 0; i < len(domains); i++ {
		identifierList.AppendElement(NewIdentifier(domains[i]))
	}
	return *identifierList
}

func decodeElements(_objectId []byte, _element []byte) (*ObjectId, Element, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the decoder
	decoder := factory.NewDecoder(_objectId)

	// Decode the ArchiveDetails
	elem, err := decoder.DecodeElement(NullObjectId)
	if err != nil {
		return nil, nil, err
	}
	objectId := elem.(*ObjectId)

	// Reallocate the decoder
	decoder = factory.NewDecoder(_element)

	// Decode the Element
	element, err := decoder.DecodeAbstractElement()
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

func createQuery(boolean *Boolean, isObjectTypeEqualToZero bool, archiveQuery ArchiveQuery, queryFilter QueryFilter) (string, error) {
	var queryBuffer bytes.Buffer
	// Only CompositeFilterSet type should be used
	queryBuffer.WriteString("SELECT timestamp, `details.related`, network, provider, source")
	// Check if we need to retrieve the element and its domain
	if *boolean == true {
		queryBuffer.WriteString(", element, domain")
	}
	// If there's not a wildcard value in one of the object type
	// fiels then we have to retrieve the entire object type
	if isObjectTypeEqualToZero == false {
		queryBuffer.WriteString(", area, service, version, number")
	}

	// Prepare the query for the conditions
	queryBuffer.WriteString(" FROM " + TABLE + " WHERE")

	var isThereAlreadyACondition = false

	// Add archive query conditions
	// Domain
	if archiveQuery.Domain != nil {
		domain := adaptDomainToString(*archiveQuery.Domain)
		queryBuffer.WriteString(fmt.Sprintf(" domain = %s", domain))
		isThereAlreadyACondition = true
	}

	// Network
	if archiveQuery.Network != nil {
		if isThereAlreadyACondition {
			queryBuffer.WriteString(" AND")
		} else {
			isThereAlreadyACondition = true
		}
		queryBuffer.WriteString(fmt.Sprintf(" network = %s", *archiveQuery.Network))
	}

	// Provider
	if archiveQuery.Provider != nil {
		if isThereAlreadyACondition {
			queryBuffer.WriteString(" AND")
		} else {
			isThereAlreadyACondition = true
		}
		queryBuffer.WriteString(fmt.Sprintf(" provider = %s", *archiveQuery.Provider))
	}

	// Related (always have to do a query with this condition)
	if isThereAlreadyACondition {
		queryBuffer.WriteString(" AND")
	} else {
		isThereAlreadyACondition = true
	}
	queryBuffer.WriteString(fmt.Sprintf(" related = %d", archiveQuery.Related))

	// Source
	if archiveQuery.Source != nil {
		if isThereAlreadyACondition {
			queryBuffer.WriteString(" AND")
		} else {
			isThereAlreadyACondition = true
		}

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
		if isThereAlreadyACondition {
			queryBuffer.WriteString(" AND")
		} else {
			isThereAlreadyACondition = true
		}
		queryBuffer.WriteString(fmt.Sprintf(" timestamp >= %s", time.Time(*archiveQuery.StartTime)))
	}

	// EndTime
	if archiveQuery.EndTime != nil {
		if isThereAlreadyACondition {
			queryBuffer.WriteString(" AND")
		} else {
			isThereAlreadyACondition = true
		}
		queryBuffer.WriteString(fmt.Sprintf(" timestamp <= %s", time.Time(*archiveQuery.EndTime)))
	}

	// SortOrder
	if archiveQuery.SortOrder != nil {
		// SortFieldName
		queryBuffer.WriteString(fmt.Sprintf(" GROUP BY %s", *archiveQuery.SortFieldName))
		// If sortOrder is false then we returned values shall be sorted
		// in descending order (ascending order is the default value)
		if *archiveQuery.SortOrder == false {
			queryBuffer.WriteString(" DESC")
		}
	}

	return queryBuffer.String(), nil
}

func whichExpressionOperatorIsIt(expressionOperator ExpressionOperator) string {
	switch expressionOperator {
	case COM_EXPRESSIONOPERATOR_EQUAL:
		return "? = ?"
	case COM_EXPRESSIONOPERATOR_DIFFER:
		return "? != ?"
	case COM_EXPRESSIONOPERATOR_GREATER:
		return "? > ?"
	case COM_EXPRESSIONOPERATOR_GREATER_OR_EQUAL:
		return "? >= ?"
	case COM_EXPRESSIONOPERATOR_LESS:
		return "? < ?"
	case COM_EXPRESSIONOPERATOR_LESS_OR_EQUAL:
		return "? <= ?"
	case COM_EXPRESSIONOPERATOR_CONTAINS:
		return "? LIKE '%?%'"
	case COM_EXPRESSIONOPERATOR_ICONTAINS:
		return "? NOT LIKE '%?%'"
	default:
		return ""
	}
}
