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
	"fmt"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
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
	// Create the handle
	archiveDatabase, err := createArchiveDatabase(USERNAME, PASSWORD, DATABASE)
	if err != nil {
		return err
	}
	defer archiveDatabase.db.Close()

	for i := 0; i < archiveDetailsList.Size(); i++ {
		if archiveDetailsList[i].InstId == 0 {
			// We have to create a new and unused object instance identifier
		} else {
			// We must verify if the object instance identifier is not already present in the table
			statementVerify, err := archiveDatabase.db.Prepare("SELECT objectInstanceIdentifier FROM " + TABLE + " WHERE objectInstanceIdentifier = ? ")
			if err != nil {
				return err
			}
			defer statementVerify.Close()

			// Execute the query
			var queryReturn int
			err = statementVerify.QueryRow(archiveDetailsList[i].InstId).Scan(&queryReturn)
			if err != nil {
				return err
			}

			fmt.Println(queryReturn)
		}
	}

	/*statementStore, err := archiveDatabase.db.Prepare("INSERT INTO " + TABLE + " VALUES ( NULL , ? , ? )")
	if err != nil {
		return err
	}
	defer statementStore.Close()

	for i := 0; i < archiveDetailsList.Size(); i++ {
		_, err = statementStore.Exec(archiveDetailsList[i].InstId, elementList[i])
		if err != nil {
			return err
		}
	}*/

	return nil
}
