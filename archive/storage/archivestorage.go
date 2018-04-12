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
func RetrieveInArchive(objectType ObjectType, domain IdentifierList, objectInstanceIdentifiers LongList) (ArchiveDetailsList, ElementList, error) {
	fmt.Println("IN: RetrieveInArchive")
	// Create the transaction to execute future queries
	db, tx, err := createTransaction()
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()

	// Create variables to return the elements and information
	//var archiveDetailsList ArchiveDetailsList
	//var elementList ElementList
	// Then, retrieve these elements and their information
	for i := 0; i < objectInstanceIdentifiers.Size(); i++ {
		// Variables to store the different elements present in the database
		var objectInstanceIdentifier Long
		var encodedObjectDetails []byte
		var encodedElement []byte

		// We can retrieve this object
		err = tx.QueryRow("SELECT objectInstanceIdentifier, element, objectDetails FROM "+TABLE+" WHERE objectInstanceIdentifier = ? ", objectInstanceIdentifier).Scan(&objectInstanceIdentifier, &encodedElement, &encodedObjectDetails)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return nil, nil, errors.New(string(MAL_ERROR_UNKNOWN_MESSAGE))
			}
			return nil, nil, err
		}

		/*archiveDetails, element, err := decodeRetrieveElements(objectInstanceIdentifier, encodedObjectDetails, encodedElement)
		if err != nil {
			return nil, nil, err
		}

		archiveDetailsList = append(archiveDetailsList, archiveDetails)*/
		//elementList = append(elementList, element)

	}

	return nil, nil, nil
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
	fmt.Println(tx)

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
	var longList *LongList

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
		for rows.Next() {
			var instID *Long
			if err = rows.Scan(instID); err != nil {
				return nil, err
			}

			*longList = append(*longList, instID)
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
	} else {

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

// TODO: might change in the future
func adaptDomain(identifierList IdentifierList) String {
	// Create the domain (It might change in the future)
	var domain String
	domain += "/"
	for i := 0; i < identifierList.Size(); i++ {
		domain += String(*identifierList.GetElementAt(i).(*Identifier))
		domain += "/"
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
