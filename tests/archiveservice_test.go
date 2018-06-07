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
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"

	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/archive/service"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/data/tests"
	. "github.com/etiennelndr/archiveservice/errors"
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

const (
	numberOfRows = 80
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

	// If there are already 80 elements in the Archive table then
	// it's useless to reset and add new elements to the database
	var maxID sql.NullInt64 // Better to use the type sql.NullInt64 to avoid nil error conversion
	err = tx.QueryRow("SELECT MAX(id) FROM " + TABLE).Scan(&maxID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// If maxID's Valid parameter is set to false then it means its value is nil
	if !maxID.Valid || maxID.Int64 != numberOfRows {
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

		// Commit changes
		tx.Commit()
		// Close the connection with the database
		db.Close()

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
		for i := 0; i < numberOfRows/2; i++ {
			// Create the value
			var signe = float64(rand.Int63n(2))
			if signe == 0 {
				elementList.AppendElement(NewValueOfSine(Float(rand.Float64())))
			} else {
				elementList.AppendElement(NewValueOfSine(Float(-rand.Float64())))
			}
			objectType = ObjectType{
				Area:    UShort(2),
				Service: UShort(3),
				Version: UOctet(1),
				Number:  UShort((*elementList)[i].GetTypeShortForm()),
			}
			// Object instance identifier
			var objectInstanceIdentifier = Long(int64(i + 1))
			// Variables for ArchiveDetailsList
			var objectKey = ObjectKey{
				Domain: identifierList,
				InstId: Long(0),
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
		_, errorsList, err := archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
		if errorsList != nil || err != nil {
			if err != nil {
				return err
			} else if errorsList != nil {
				return errors.New(string(*errorsList.ErrorNumber) + ": " + string(*errorsList.ErrorComment))
			}
		}

		// Store fourty new elements (total 80 elements)
		identifierList = IdentifierList([]*Identifier{NewIdentifier("en"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
		for i := 0; i < archiveDetailsList.Size(); i++ {
			var objectInstanceIdentifier = Long(int64(i + 41))
			archiveDetailsList[i].InstId = objectInstanceIdentifier
			archiveDetailsList[i].Details.Source.Key.Domain = identifierList
		}
		_, errorsList, err = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
		if errorsList != nil || err != nil {
			if err != nil {
				return err
			} else if errorsList != nil {
				return errors.New(string(*errorsList.ErrorNumber) + ": " + string(*errorsList.ErrorComment))
			} else {
				return errors.New("UNKNOWN ERROR")
			}
		}
	} else {
		// Commit changes
		tx.Commit()
		// Close the connection with the database
		db.Close()
	}

	return nil
}

// checkAndInitDatabase Checks if the Archive table is initialized or not
// If not, it initializes it and inserts datas in the table Archive
func checkAndInitDatabase() error {
	if !isDatabaseInitialized {
		err := initDabase()
		if err != nil {
			return err
		}
		isDatabaseInitialized = true
	}
	return nil
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
func TestRetrieveOK(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		println(err.Error())
		t.FailNow()
	}

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)
	// Variable that defines the ArchiveService
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	var longList = LongList([]*Long{NewLong(0)})

	// Variables to retrieve the return of this function
	var archiveDetailsList *ArchiveDetailsList
	var elementList ElementList
	var errorsList *ServiceError
	// Start the consumer
	archiveDetailsList, elementList, errorsList, err = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)

	if errorsList != nil || err != nil || archiveDetailsList == nil || elementList == nil {
		println(errorsList)
		t.FailNow()
	}
}

func TestRetrieveKO_3_4_3_2_2(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	var longList = LongList([]*Long{NewLong(0)})
	// Area is equal to 0
	var objectType = ObjectType{
		Area:    UShort(0),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}

	// Service is equal to 0
	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(0),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}

	// Vesrion is equal to 0
	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(0),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}

	// Number is equal to 0
	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(0),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}
}

func TestRetrieveKO_3_4_3_2_4(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

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
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}

	identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("*"), NewIdentifier("archiveservice")})
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}

	identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("*")})
	_, _, errorsList, _ = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}
}

//======================================================================//
//								QUERY									//
//======================================================================//
func TestQueryOK(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	//var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	archiveQuery := &ArchiveQuery{
		Related:   Long(0),
		SortOrder: NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *CompositeFilterSetList

	// Variable to retrieve the responses
	var responses []interface{}

	// Start the consumer
	responses, errorsList, err = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	if errorsList != nil || err != nil || responses == nil {
		t.FailNow()
	}
}

func TestQueryKO_3_4_4_2_9(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	archiveQuery := &ArchiveQuery{
		Related:   Long(0),
		SortOrder: NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList = NewCompositeFilterSetList(1)

	// Start the consumer
	_, errorsList, _ = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}
}

func TestQueryOK_3_4_4_2_14(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	archiveQuery := &ArchiveQuery{
		Related:       Long(0),
		SortOrder:     NewBoolean(true),
		SortFieldName: NewString("domain"),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *CompositeFilterSetList

	// Variable to retrieve the responses
	var responses []interface{}

	// Start the consumer
	responses, errorsList, err = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	if errorsList != nil || err != nil || responses == nil {
		t.FailNow()
	}
	// Now, verify the responses
	for i, resp := range responses {
		if i%4 == 1 {
			if (i == 1 && *(*resp.(*IdentifierList))[0] != "en") || (i == 5 && *(*resp.(*IdentifierList))[0] != "fr") {
				t.FailNow()
			}
		} else if i%4 == 3 {
			if resp.(ElementList).Size() != numberOfRows/2 {
				t.FailNow()
			}
		}
	}
}

func TestQueryKO_3_4_4_2_14(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	archiveQuery := &ArchiveQuery{
		Related:       Long(0),
		SortOrder:     NewBoolean(true),
		SortFieldName: NewString("invalidname"),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *CompositeFilterSetList

	// Start the consumer
	_, errorsList, _ = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}
}

func TestQueryKO_3_4_4_2_16(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	archiveQuery := &ArchiveQuery{
		Related: Long(0),
	}
	archiveQueryList.AppendElement(archiveQuery)
	// Create the query filters list
	var nilValue *Identifier
	var queryFilterList = NewCompositeFilterSetList(0)
	compositeFilter := NewCompositeFilter(String("domain"), COM_EXPRESSIONOPERATOR_CONTAINS, nilValue)
	compositeFilterList := NewCompositeFilterList(0)
	compositeFilterList.AppendElement(compositeFilter)
	queryFilter := NewCompositeFilterSet(compositeFilterList)
	queryFilterList.AppendElement(queryFilter)

	// Start the consumer
	_, errorsList, _ = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) || !strings.Contains(string(*errorsList.ErrorComment), string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR)) {
		fmt.Println(*errorsList.ErrorComment)
		t.FailNow()
	}
}

func TestQueryKO_3_4_4_2_19(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	archiveQuery := &ArchiveQuery{
		Related: Long(0),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *CompositeFilterSetList

	// Start the consumer WITHOUT any wildcard value in the objectType
	resp, _, _ := archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	for i := 0; i < len(resp)/4; i++ {
		objType := resp[i*4].(*ObjectType)
		if objType != nil {
			t.FailNow()
		}
	}

	// Start the consumer WITH a wildcard value in the objectType
	objectType.Area = 0
	resp, _, _ = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	for i := 0; i < len(resp)/4; i++ {
		objType := resp[i*4].(*ObjectType)
		if objType == nil {
			t.FailNow()
		}
	}
}

func TestQueryKO_3_4_4_2_25(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean *Boolean
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	archiveQuery := &ArchiveQuery{
		Related: Long(0),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *CompositeFilterSetList

	// Start the consumer with the initial Boolean set to NIL
	var newBoolean *Boolean
	resp, _, _ := archiveService.Query(consumerURL, providerURL, newBoolean, objectType, *archiveQueryList, queryFilterList)

	fmt.Println(resp)
	for i := 0; i < len(resp)/4; i++ {
		elementList := resp[i*4+3]
		if elementList != ElementList(nil) {
			t.FailNow()
		}
	}

	// Start the consumer with the initial Boolean set to TRUE
	boolean = NewBoolean(true)
	resp, _, _ = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	for i := 0; i < len(resp)/4; i++ {
		elementList := resp[i*4+3]
		if elementList == ElementList(nil) {
			t.FailNow()
		}
	}

	// Start the consumer with the initial Boolean set to FALSE
	boolean = NewBoolean(false)
	resp, _, _ = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	fmt.Println(resp)
	for i := 0; i < len(resp)/4; i++ {
		elementList := resp[i*4+3]
		if elementList != ElementList(nil) {
			t.FailNow()
		}
	}
}

//======================================================================//
//								COUNT									//
//======================================================================//
func TestCountOK(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	archiveQuery := &ArchiveQuery{
		Domain:    &domain,
		Related:   Long(0),
		SortOrder: NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	domain2 := IdentifierList([]*Identifier{NewIdentifier("en"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	archiveQuery2 := &ArchiveQuery{
		Domain:    &domain2,
		Related:   Long(0),
		SortOrder: NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery2)
	var queryFilterList *CompositeFilterSetList

	// Variable to retrieve the return of this function
	var longList *LongList
	// Start the consumer
	longList, errorsList, err = archiveService.Count(consumerURL, providerURL, objectType, archiveQueryList, queryFilterList)

	if errorsList != nil || err != nil || longList == nil {
		t.FailNow()
	}

	if longList.Size() != 2 {
		t.FailNow()
	}
	for i := 0; i < longList.Size(); i++ {
		if *longList.GetElementAt(i).(*Long) != Long(40) {
			t.FailNow()
		}
	}
}

func TestCountKO_3_4_5_2_9(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	archiveQuery := &ArchiveQuery{
		Domain:    &domain,
		Related:   Long(0),
		SortOrder: NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList = NewCompositeFilterSetList(1)

	_, errorsList, _ = archiveService.Count(consumerURL, providerURL, objectType, archiveQueryList, queryFilterList)

	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}
}

func TestCountOK_3_4_5_2_14(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	archiveQuery := &ArchiveQuery{
		Domain:        &domain,
		Related:       Long(0),
		SortOrder:     NewBoolean(true),
		SortFieldName: NewString("domain"),
	}
	archiveQueryList.AppendElement(archiveQuery)

	var queryFilterList *CompositeFilterSetList

	// Variable to retrieve the return of this function
	var longList *LongList
	// Start the consumer
	longList, errorsList, err = archiveService.Count(consumerURL, providerURL, objectType, archiveQueryList, queryFilterList)

	if errorsList != nil || err != nil || longList == nil {
		t.FailNow()
	}

	if longList.Size() != 1 {
		t.FailNow()
	}
	for i := 0; i < longList.Size(); i++ {
		if *longList.GetElementAt(i).(*Long) != Long(40) {
			t.FailNow()
		}
	}
}

func TestCountKO_3_4_5_2_14(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	archiveQuery := &ArchiveQuery{
		Domain:        &domain,
		Related:       Long(0),
		SortOrder:     NewBoolean(true),
		SortFieldName: NewString("invalidname"),
	}
	archiveQueryList.AppendElement(archiveQuery)

	var queryFilterList *CompositeFilterSetList

	// Start the consumer
	_, errorsList, _ = archiveService.Count(consumerURL, providerURL, objectType, archiveQueryList, queryFilterList)

	if errorsList == nil || *errorsList.ErrorNumber != *NewUInteger(uint32(COM_ERROR_INVALID)) {
		t.FailNow()
	}
}

func TestCountKO_3_4_5_2_16(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := NewArchiveQueryList(0)
	var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	archiveQuery := &ArchiveQuery{
		Domain:    &domain,
		Related:   Long(0),
		SortOrder: NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)

	// Create the query filters list
	var nilValue *Identifier
	var queryFilterList = NewCompositeFilterSetList(0)
	compositeFilter := NewCompositeFilter(String("domain"), COM_EXPRESSIONOPERATOR_CONTAINS, nilValue)
	compositeFilterList := NewCompositeFilterList(0)
	compositeFilterList.AppendElement(compositeFilter)
	queryFilter := NewCompositeFilterSet(compositeFilterList)
	queryFilterList.AppendElement(queryFilter)

	// Start the consumer
	_, errorsList, _ = archiveService.Count(consumerURL, providerURL, objectType, archiveQueryList, queryFilterList)

	if errorsList == nil || *errorsList.ErrorNumber != UInteger(uint32(COM_ERROR_INVALID)) || !strings.Contains(string(*errorsList.ErrorComment), string(ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR)) {
		t.FailNow()
	}
}

//======================================================================//
//								STORE									//
//======================================================================//
func TestStoreOK(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(81)
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
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Variable to retrieve the return of this function
	var longList *LongList
	// Start the consumer
	longList, errorsList, err = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList != nil || err != nil || longList == nil {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_1(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean *Boolean
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(0)
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
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Variable to retrieve the return of this function
	var longList *LongList

	// First, start the consumer with the boolean set to NIL
	// Start the consumer
	longList, _, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if longList != nil {
		t.FailNow()
	}

	// Then, start the consumer with the boolean set to FALSE
	boolean = NewBoolean(false)
	// Start the consumer
	longList, _, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if longList != nil {
		t.FailNow()
	}

	// Finally, start the consumer with the boolean set to TRUE
	boolean = NewBoolean(true)
	// Start the consumer
	longList, _, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if longList == nil {
		t.FailNow()
	}
}

func TestStoreOK_3_4_6_2_6(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(0)
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
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Variable to retrieve the return of this function
	var longList *LongList

	// Start the consumer
	longList, _, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if longList == nil {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_6(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(45)
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
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_DUPLICATE {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_8(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = NewValueOfSineList(2)
	(*elementList)[0] = NewValueOfSine(0)
	(*elementList)[1] = NewValueOfSine(0.5)
	var boolean = NewBoolean(true)
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = Long(0)
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
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_STORE_LIST_SIZE_ERROR || *errorsList.ErrorExtra.(*Long) != 1 {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_9(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = NewBoolean(true)
	// Area is equal to 0
	var objectType = ObjectType{
		Area:    UShort(0),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = Long(0)
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
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	// Service is equal to 0
	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(0),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	// Version is equal to 0
	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(0),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	// Number is equal to 0
	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(0),
	}
	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_10(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variables to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = NewBoolean(true)
	// Area is equal to 0
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("*"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = Long(0)
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
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("network")
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_11(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = NewBoolean(true)
	// Area is equal to 0
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = Long(0)
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
		Related: NewLong(1),
		Source:  &objectID,
	}
	// Bad value for the NETWORK
	var network *Identifier
	var timestamp = NewFineTime(time.Now())
	var provider = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR {
		fmt.Println(errorsList)
		t.FailNow()
	}

	// Bad value for the TIMESTAMP
	network = NewIdentifier("network")
	*timestamp = FineTime(time.Unix(int64(0), int64(0)))
	archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR {
		fmt.Println(errorsList)
		t.FailNow()
	}

	// Bad value for the TIMESTAMP
	*timestamp = FineTime(time.Now())
	*provider = URI("*")
	archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

	// Start the consumer
	_, errorsList, _ = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR {
		fmt.Println(errorsList)
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_12(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// TODO: do this test
}

//======================================================================//
//								UPDATE									//
//======================================================================//
func TestUpdateOK(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(1)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = ObjectId{
		Type: &objectType,
		Key:  &objectKey,
	}
	var objectDetails = ObjectDetails{
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("new.network")
	var fineTime = NewFineTime(time.Now())
	var uri = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, fineTime, uri)})

	// Start the consumer
	errorsList, err = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList != nil || err != nil {
		t.FailNow()
	}
}

func TestUpdateKO_3_4_7_2_5(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	// Object instance identifier: UNKNOWN
	var objectInstanceIdentifier = *NewLong(155)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = ObjectId{
		Type: &objectType,
		Key:  &objectKey,
	}
	var objectDetails = ObjectDetails{
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("new.network")
	var fineTime = NewFineTime(time.Now())
	var uri = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, fineTime, uri)})

	// Start the consumer
	errorsList, _ = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != MAL_ERROR_UNKNOWN || *errorsList.ErrorComment != ARCHIVE_SERVICE_UNKNOWN_ELEMENT {
		t.FailNow()
	}
}

func TestUpdateKO_3_4_7_2_8_ObjectType(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = ObjectType{
		Area:    UShort(0),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	// Object instance identifier: UNKNOWN
	var objectInstanceIdentifier = *NewLong(41)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = ObjectId{
		Type: &objectType,
		Key:  &objectKey,
	}
	var objectDetails = ObjectDetails{
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("new.network")
	var fineTime = NewFineTime(time.Now())
	var uri = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, fineTime, uri)})

	// Start the consumer
	errorsList, _ = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(0),
		Version: UOctet(1),
		Number:  UShort((*elementList)[0].GetTypeShortForm()),
	}
	// Start the consumer
	errorsList, _ = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(0),
		Number:  UShort((*elementList)[0].GetTypeShortForm()),
	}
	// Start the consumer
	errorsList, _ = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(0),
	}
	// Start the consumer
	errorsList, _ = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}
}

func TestUpdateKO_3_4_7_2_8_ObjectInstanceIdentifier(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(0)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = ObjectId{
		Type: &objectType,
		Key:  &objectKey,
	}
	var objectDetails = ObjectDetails{
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("new.network")
	var fineTime = NewFineTime(time.Now())
	var uri = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, fineTime, uri)})

	// Start the consumer
	errorsList, _ = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_AREA_OBJECT_INSTANCE_IDENTIFIER_VALUE_ERROR {
		t.FailNow()
	}
}

func TestUpdateKO_3_4_7_2_8_Domain(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("*"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	// Object instance identifier
	var objectInstanceIdentifier = *NewLong(41)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = ObjectId{
		Type: &objectType,
		Key:  &objectKey,
	}
	var objectDetails = ObjectDetails{
		Related: NewLong(1),
		Source:  &objectID,
	}
	var network = NewIdentifier("new.network")
	var fineTime = NewFineTime(time.Now())
	var uri = NewURI("main/start")
	var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, fineTime, uri)})

	// Start the consumer
	errorsList, _ = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR {
		fmt.Println(*errorsList.ErrorComment)
		t.FailNow()
	}
}

//======================================================================//
//								DELETE									//
//======================================================================//
func TestDeleteOK(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	var longList = NewLongList(0)
	longList.AppendElement(NewLong(15))

	// Variable to retrieve the return of this function
	var respLongList *LongList
	// Start the consumer
	respLongList, errorsList, err = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longList)

	if errorsList != nil || err != nil || respLongList == nil {
		t.FailNow()
	}
}

func TestDeleteKO_3_4_8_2_3_ObjectType(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = ObjectType{
		Area:    UShort(0),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	var longList = NewLongList(0)
	longList.AppendElement(NewLong(15))

	// Start the consumer
	_, errorsList, _ = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(0),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, errorsList, _ = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(0),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, errorsList, _ = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}

	objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(0),
	}
	// Start the consumer
	_, errorsList, _ = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR {
		t.FailNow()
	}
}

func TestDeleteKO_3_4_8_2_3_Domain(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("*"), NewIdentifier("test")})
	var longList = NewLongList(0)
	longList.AppendElement(NewLong(15))

	// Start the consumer
	_, errorsList, _ = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longList)

	if errorsList == nil || *errorsList.ErrorNumber != COM_ERROR_INVALID || *errorsList.ErrorComment != ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR {
		t.FailNow()
	}
}

func TestDeleteKO_3_4_8_2_6(t *testing.T) {
	// Check if the Archive table is initialized or not
	err := checkAndInitDatabase()
	if err != nil {
		t.FailNow()
	}

	// Variable to retrieve the return of this function
	var errorsList *ServiceError
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
	var longList = NewLongList(0)
	longList.AppendElement(NewLong(175))

	// Start the consumer
	_, errorsList, _ = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longList)

	if errorsList == nil || *errorsList.ErrorNumber != MAL_ERROR_UNKNOWN || *errorsList.ErrorComment != ARCHIVE_SERVICE_UNKNOWN_ELEMENT {
		t.FailNow()
	}
}
