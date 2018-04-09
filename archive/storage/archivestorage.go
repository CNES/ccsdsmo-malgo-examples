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
	"math/rand"
	"time"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/data"

	_ "github.com/go-sql-driver/mysql"
)

// ArchiveDatabase :
type ArchiveDatabase struct {
	db       *sql.DB
	username string
	password string
	database string
}

const (
	USERNAME = "archiveService"
	PASSWORD = "1a2B3c4D!@?"
	DATABASE = "archive"
	TABLE    = "Archive"
)

func createArchiveDatabase(username string, password string, database string) (*ArchiveDatabase, error) {
	// Get a handle for our database
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)
	if err != nil {
		return nil, err
	}

	// Create the bdd for the ArchiveService
	archiveDatabase := &ArchiveDatabase{
		db,
		username,
		password,
		database,
	}

	return archiveDatabase, nil
}

// StoreInArchive : Use this function to store objects in an COM archive
func StoreInArchive(objectType ObjectType, identifier IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) error {
	rand.Seed(time.Now().UnixNano())

	// Create the handle
	archiveDatabase, err := createArchiveDatabase(USERNAME, PASSWORD, DATABASE)
	if err != nil {
		return err
	}
	defer archiveDatabase.db.Close()

	// TODO : WARNING -> We MUST do a roll back if a "DUPLICATE" error occurs (none
	// of the objects have to be stored)
	for i := 0; i < archiveDetailsList.Size(); i++ {
		if archiveDetailsList[i].InstId != 0 {
			// We must verify if the object instance identifier is not already present in the table
			boolean, err := isObjectInstanceIdentifierInDatabase(archiveDatabase, int64(archiveDetailsList[i].InstId))
			if err != nil {
				return err
			}
			if boolean {
				return errors.New(string(COM_ERROR_DUPLICATE))
			}
		}
	}

	for i := 0; i < archiveDetailsList.Size(); i++ {
		if archiveDetailsList[i].InstId == 0 {
			// We have to create a new and unused object instance identifier
			for {
				var objectInstanceIdentifier = rand.Int63n(int64(LONG_MAX))
				boolean, err := isObjectInstanceIdentifierInDatabase(archiveDatabase, objectInstanceIdentifier)
				if err != nil {
					return err
				}
				if !boolean {
					// OK, we can insert the object with this instance identifier
					// This object is not present in the archive
					err := insertInDatabase(archiveDatabase, objectInstanceIdentifier, elementList.GetElementAt(i))
					if err != nil {
						return err
					}
					break
				}
			}
		} else {
			// This object is not present in the archive
			err := insertInDatabase(archiveDatabase, int64(archiveDetailsList[i].InstId), elementList.GetElementAt(i))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func isObjectInstanceIdentifierInDatabase(archiveDatabase *ArchiveDatabase, objectInstanceIdentifier int64) (bool, error) {
	statementVerify, err := archiveDatabase.db.Prepare("SELECT objectInstanceIdentifier FROM " + TABLE + " WHERE objectInstanceIdentifier = ? ")
	if err != nil {
		return false, err
	}
	defer statementVerify.Close()

	// Execute the query
	// Before, create a variable to retrieve the result
	var queryReturn int
	// Then, execute the query
	err = statementVerify.QueryRow(objectInstanceIdentifier).Scan(&queryReturn)
	if err != nil {
		if err.Error() != "sql: no rows in result set" {
			return false, err
		}
		return false, nil
	}
	return true, nil
}

func insertInDatabase(archiveDatabase *ArchiveDatabase, objectInstanceIdentifier int64, element Element) error {
	statementStore, err := archiveDatabase.db.Prepare("INSERT INTO " + TABLE + " VALUES ( NULL , ? , ? )")
	if err != nil {
		return err
	}
	defer statementStore.Close()

	_, err = statementStore.Exec(objectInstanceIdentifier, element)
	if err != nil {
		return err
	}
	return nil
}
