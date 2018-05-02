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
package tests

import (
	"database/sql"
	"errors"
	"math/rand"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"

	. "github.com/etiennelndr/archiveservice/archive/service"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/errors"
	. "github.com/etiennelndr/archiveservice/tests/data"
)

// Constants for the providers and consumers
const (
	providerURL = "maltcp://127.0.0.1:12400"
	consumerURL = "maltcp://127.0.0.1:14200"
)

// Database ids
const (
	USERNAME = "archiveService"
	PASSWORD = "1a2B3c4D!@?"
	DATABASE = "archive"
	TABLE    = "Archive"
)

// isDatabaseInitialized attribute is true when the database has been initialized
var isDatabaseInitialized = false

// initDatabase is used to init the database
func initDabase() error {
	rand.Seed(time.Now().UnixNano())

	// Open the database
	db, err := sql.Open("mysql", USERNAME+":"+PASSWORD+"@/"+DATABASE+"?parseTime=true")
	if err != nil {
		return err
	}
	defer db.Close()

	// Validate the connection by pinging it
	err = db.Ping()
	if err != nil {
		return err
	}

	// Create the transaction (we have to use this method to use rollback and commit)
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// Delete all the elements of the table Archive
	_, err = tx.Exec("DELETE FROM " + TABLE)
	if err != nil {
		return err
	}

	// Reset the AUTO_INCREMENT value
	_, err = tx.Exec("ALTER TABLE " + TABLE + " AUTO_INCREMENT=0")
	if err != nil {
		return err
	}

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	archiveService = archiveService.CreateService().(*ArchiveService)

	// Insert elements in the table Archive for future tests
	var elementList = NewValueOfSineList(0)
	var boolean = NewBoolean(false)
	// Variable for the different networks
	var networks = []*Identifier{
		NewIdentifier("tests/network1"),
		NewIdentifier("tests/network2"),
	}
	// Variable for the different providers
	var providers = []*URI{
		NewURI("tests/provider1"),
		NewURI("tests/provider2"),
	}

	var objectType ObjectType
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	var archiveDetailsList = *NewArchiveDetailsList(0)

	// Create elements
	for i := 0; i < 40; i++ {
		// Create the value
		var signe = float64(rand.Int63n(2))
		if signe == 0 {
			elementList.AppendElement(NewValueOfSine(Float(rand.Float64())))
		} else {
			elementList.AppendElement(NewValueOfSine(Float(-rand.Float64())))
		}
		objectType = ObjectType{
			UShort(2),
			UShort(3),
			UOctet(1),
			UShort((*elementList)[i].GetTypeShortForm()),
		}
		// Object instance identifier
		var objectInstanceIdentifier = *NewLong(int64(i))
		// Variables for ArchiveDetailsList
		var objectKey = ObjectKey{
			Domain: identifierList,
			InstId: objectInstanceIdentifier,
		}
		var objectID = ObjectId{
			Type: &objectType,
			Key:  &objectKey,
		}
		var objectDetails = ObjectDetails{
			Related: NewLong(0),
			Source:  &objectID,
		}
		var network = networks[rand.Int63n(int64(len(networks)))]
		var timestamp = NewFineTime(time.Now())
		var provider = providers[rand.Int63n(int64(len(providers)))]
		archiveDetailsList.AppendElement(NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider))
	}
	_, errorsList, err := archiveService.Store(consumerURL, providerURL, *boolean, objectType, identifierList, archiveDetailsList, elementList)
	if errorsList != nil || err != nil {
		if err != nil {
			return err
		} else if errorsList != nil {
			return errors.New(string(*errorsList.ErrorNumber) + ": " + string(*errorsList.ErrorComment))
		} else {
			return errors.New("UNKNOWN ERROR")
		}
	}

	return nil
}

// checkAndInitDatabase Checks if the Archive table is intitializes or not
// If not, it initializes it and inserts datas in the table Archive
func checkAndInitDatabase(t *testing.T) {
	if !isDatabaseInitialized {
		err := initDabase()
		if err != nil {
			t.FailNow()
		}
		isDatabaseInitialized = true
	}
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
func TestRetrieveOK(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)
	println("aight")

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)
	// Variable that defines the ArchiveService
	var objectType = ObjectType{
		UShort(2),
		UShort(3),
		UOctet(1),
		UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	var longList = LongList([]*Long{NewLong(0)})

	// Variables to retrieve the return of this function
	var archiveDetailsList *ArchiveDetailsList
	var elementList ElementList
	var errorsList *ServiceError
	var err error
	// Start the consumer
	archiveDetailsList, elementList, errorsList, err = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)

	if errorsList != nil || err != nil || archiveDetailsList == nil || elementList == nil {
		t.FailNow()
	}
}

func TestRetrieveKO_3_4_3_2_2(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	var longList = LongList([]*Long{NewLong(0)})
	var objectType = ObjectType{
		UShort(0),
		UShort(3),
		UOctet(1),
		UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(70000) {
		t.FailNow()
	}

	objectType = ObjectType{
		UShort(2),
		UShort(0),
		UOctet(1),
		UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(70000) {
		t.FailNow()
	}

	objectType = ObjectType{
		UShort(2),
		UShort(3),
		UOctet(0),
		UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(70000) {
		t.FailNow()
	}

	objectType = ObjectType{
		UShort(2),
		UShort(3),
		UOctet(1),
		UShort(0),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(70000) {
		t.FailNow()
	}
}

func TestRetrieveKO_3_4_3_2_4(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var identifierList = IdentifierList([]*Identifier{NewIdentifier("*"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	var longList = LongList([]*Long{NewLong(0)})
	var objectType = ObjectType{
		UShort(2),
		UShort(3),
		UOctet(1),
		UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(70000) {
		t.FailNow()
	}

	identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("*"), NewIdentifier("archiveservice")})
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(70000) {
		t.FailNow()
	}

	identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("*")})
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(70000) {
		t.FailNow()
	}
}

//======================================================================//
//								QUERY									//
//======================================================================//
func TestQueryOK(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestQueryKO_3_4_4_2_9(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestQueryKO_3_4_4_2_14(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestQueryKO_3_4_4_2_16(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestQueryKO_3_4_4_2_19(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestQueryKO_3_4_4_2_24(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestQueryKO_3_4_4_2_25(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

//======================================================================//
//								COUNT									//
//======================================================================//
func TestCountOK(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestCountKO_3_4_5_2_9(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestCountKO_3_4_5_2_14(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestCountKO_3_4_5_2_16(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestCountKO_3_4_5_2_19(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestCountKO_3_4_5_2_24(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestCountKO_3_4_5_2_25(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

//======================================================================//
//								STORE									//
//======================================================================//
func TestStoreOK(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestStoreKO_3_4_6_2_1(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestStoreKO_3_4_6_2_6(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestStoreKO_3_4_6_2_8(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestStoreKO_3_4_6_2_9(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestStoreKO_3_4_6_2_10(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestStoreKO_3_4_6_2_11(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

func TestStoreKO_3_4_6_2_12(t *testing.T) {
	// Check if the Archive table is intitializes or not
	checkAndInitDatabase(t)

	t.FailNow()
}

//======================================================================//
//								UPDATE									//
//======================================================================//

//======================================================================//
//								DELETE									//
//======================================================================//
