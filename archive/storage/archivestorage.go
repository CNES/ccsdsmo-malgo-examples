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
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
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
	fmt.Println("IN: RetrieveInArchive")
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	// Convert domain
	domain := adaptDomain(identifierList)

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

	// Create variables to return the elements and information
	var archiveDetailsList = *NewArchiveDetailsList(0)
	var elementList = element.(ElementList)
	elementList = elementList.CreateElement().(ElementList)
	// Then, retrieve these elements and their information
	if !isAll {
		for i := 0; i < objectInstanceIdentifierList.Size(); i++ {
			// Variables to store the different elements present in the database
			var encodedObjectDetails []byte
			var encodedElement []byte

			// We can retrieve this object
			err = tx.QueryRow("SELECT element, objectDetails FROM "+TABLE+" WHERE objectInstanceIdentifier = ? ",
				*objectInstanceIdentifierList[i]).Scan(&encodedElement, &encodedObjectDetails)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					return nil, nil, errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
				}
				return nil, nil, err
			}

			archiveDetails, element, err := decodeRetrieveElements(*objectInstanceIdentifierList[i], encodedObjectDetails, encodedElement)
			if err != nil {
				return nil, nil, err
			}

			archiveDetailsList.AppendElement(archiveDetails)
			elementList.AppendElement(element)
		}
	} else {
		// Retrieve all these elements (no particular object instance iedentifiers)
		// Variables to store the different elements present in the database
		var objectInstanceIdentifier Long
		var encodedObjectDetails []byte
		var encodedElement []byte

		// Retrieve this object and its archive details in the archive
		rows, err := tx.Query("SELECT objectInstanceIdentifier, element, objectDetails FROM "+TABLE+" WHERE objectTypeArea = ? AND objectTypeService = ? AND objectTypeVersion = ? AND objectTypeNumber = ? AND domain = ?",
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
			if err = rows.Scan(&objectInstanceIdentifier, &encodedElement, &encodedObjectDetails); err != nil {
				return nil, nil, err
			}

			archiveDetails, element, err := decodeRetrieveElements(objectInstanceIdentifier, encodedObjectDetails, encodedElement)
			if err != nil {
				return nil, nil, err
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
func QueryArchive() error {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return err
	}
	defer db.Close()

	// Commit changes
	tx.Commit()

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

	// Commit changes
	tx.Commit()

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
	domain := adaptDomain(identifierList)

	for i := 0; i < archiveDetailsList.Size(); i++ {
		// First of all, we need to encode the element and the objectDetails
		factory := new(FixedBinaryEncoding)

		// Create the encoder
		encoder := factory.NewEncoder(make([]byte, 0, 8192))

		// Encode the Element
		err := encoder.EncodeAbstractElement(elementList.GetElementAt(i))
		if err != nil {
			return nil, err
		}
		encodedElement := encoder.Body()

		// Reallocate the encoder
		encoder = factory.NewEncoder(make([]byte, 0, 8192))

		// Encode the ObjectDetails
		err = archiveDetailsList[i].Details.Encode(encoder)
		if err != nil {
			return nil, err
		}
		encodedObjectDetails := encoder.Body()

		//encoder := factory.
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
					err := insertInDatabase(tx, objectInstanceIdentifier, encodedElement, objectType, domain, encodedObjectDetails)
					if err != nil {
						// An error occurred, do a rollback
						tx.Rollback()
						return nil, err
					}

					// Insert this new object instance identifier in the returned list
					longList = append(longList, NewLong(objectInstanceIdentifier))

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
			err = insertInDatabase(tx, int64(archiveDetailsList[i].InstId), encodedElement, objectType, domain, encodedObjectDetails)
			if err != nil {
				// An error occurred, do a rollback
				tx.Rollback()
				return nil, err
			}

			// Insert this new object instance identifier in the returned list
			longList = append(longList, &archiveDetailsList[i].InstId)
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
	domain := adaptDomain(identifierList)

	for i := 0; i < elementList.Size(); i++ {
		// First of all, we need to verify if the object instance identifier, combined with the object type
		// and the domain are in the archive
		var queryReturn int
		err := tx.QueryRow("SELECT objectInstanceIdentifier FROM "+TABLE+" WHERE objectInstanceIdentifier = ? AND objectTypeArea = ? AND objectTypeService = ? AND objectTypeVersion = ? AND objectTypeNumber = ? AND domain = ?",
			archiveDetailsList[i].InstId,
			objectType.Area,
			objectType.Service,
			objectType.Version,
			objectType.Number,
			domain).Scan(&queryReturn)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
			}
			return err
		}

		// TODO: Encode element and objectDetails
		element, objectDetails, err := encodeUpdateElements(elementList.GetElementAt(i), archiveDetailsList[i].Details)

		// If no error, the object is in the archive and we can update it
		_, err = tx.Exec("UPDATE "+TABLE+" SET element = ? AND objectDetails = ? WHERE objectInstanceIdentifier = ? AND objectTypeArea = ? AND objectTypeService = ? AND objectTypeVersion = ? AND objectTypeNumber = ? AND domain = ?",
			element,
			objectDetails,
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
func DeleteInArchive(objectType ObjectType, identifierList IdentifierList, longListRequest LongList) (*LongList, error) {
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Variable to return
	var longList = NewLongList(0)

	// Create the domain (It might change in the future)
	domain := adaptDomain(identifierList)

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
		rows, err := tx.Query("SELECT objectInstanceIdentifier FROM "+TABLE+" WHERE objectTypeArea = ? AND objectTypeService = ? AND objectTypeVersion = ? AND objectTypeNumber = ? AND domain = ?",
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
		_, err = tx.Exec("DELETE FROM "+TABLE+" WHERE objectTypeArea = ? AND objectTypeService = ? AND objectTypeVersion = ? AND objectTypeNumber = ? AND domain = ?",
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
			err := tx.QueryRow("SELECT objectInstanceIdentifier FROM "+TABLE+" WHERE objectInstanceIdentifier = ? AND objectTypeArea = ? AND objectTypeService = ? AND objectTypeVersion = ? AND objectTypeNumber = ? AND domain = ?",
				*longListRequest[i],
				objectType.Area,
				objectType.Service,
				objectType.Version,
				objectType.Number,
				domain).Scan(&objInstID)
			if err != nil {
				if err.Error() == "sql: no rows in result set" {
					return nil, errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
				}
				return nil, err
			}

			_, err = tx.Exec("DELETE FROM "+TABLE+" WHERE objectInstanceIdentifier = ? AND objectTypeArea = ? AND objectTypeService = ? AND objectTypeVersion = ? AND objectTypeNumber = ? AND domain = ?",
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
func insertInDatabase(tx *sql.Tx, objectInstanceIdentifier int64, element []byte, objectType ObjectType, domain String, objectDetails []byte) error {
	_, err := tx.Exec("INSERT INTO "+TABLE+" VALUES ( NULL , ? , ? , ? , ? , ? , ? , ? , ? )",
		objectInstanceIdentifier,
		element,
		objectType.Area,
		objectType.Service,
		objectType.Version,
		objectType.Number,
		domain,
		objectDetails)
	if err != nil {
		return err
	}
	return nil
}

// resetAutoIncrement take the maximum id in the database and set the
// AUTO_INCREMENT at this value (actually it's this value to which we added 1)
func resetAutoIncrement(tx *sql.Tx) error {
	_, err := tx.Exec("SELECT @max := max(id)+1 FROM " + TABLE)
	if err != nil {
		return err
	}
	_, err = tx.Exec("SET @alter_statement = CONCAT('ALTER TABLE " + TABLE + " AUTO_INCREMENT = ', @max)")
	if err != nil {
		return err
	}
	_, err = tx.Exec("PREPARE stmt1 FROM @alter_statement;")
	if err != nil {
		return err
	}
	_, err = tx.Exec("EXECUTE stmt1")
	if err != nil {
		return err
	}
	_, err = tx.Exec("DEALLOCATE PREPARE stmt1;")
	if err != nil {
		return err
	}
	return nil
}

func adaptDomain(identifierList IdentifierList) String {
	var domain String
	for i := 0; i < identifierList.Size(); i++ {
		domain += String(*identifierList.GetElementAt(i).(*Identifier))
		if i+1 < identifierList.Size() {
			domain += "."
		}
	}
	return domain
}

func decodeRetrieveElements(objectInstanceIdentifier Long, _objectDetails []byte, _element []byte) (*ArchiveDetails, Element, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the decoder
	decoder := factory.NewDecoder(_objectDetails)

	// Decode the ArchiveDetails
	elem, err := decoder.DecodeElement(NullObjectDetails)
	if err != nil {
		return nil, nil, err
	}
	objectDetails := elem.(*ObjectDetails)

	// Reallocate the decoder
	decoder = factory.NewDecoder(_element)

	// Decode the Element
	element, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, nil, err
	}

	// Create the new elements
	archiveDetails := &ArchiveDetails{
		InstId:  objectInstanceIdentifier,
		Details: *objectDetails,
	}

	return archiveDetails, element, nil
}

func encodeUpdateElements(_element Element, _objectDetails ObjectDetails) ([]byte, []byte, error) {
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

	// Encode ObjectDetails
	err = _objectDetails.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}
	objectDetails := encoder.Body()

	return element, objectDetails, nil
}

// This part is usefull for type short form conversion (from typeShortForm to listShortForm)
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
