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
package tests

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/com/archive"
	"github.com/CNES/ccsdsmo-malgo/mal"
	malapi "github.com/CNES/ccsdsmo-malgo/mal/api"

	. "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/archive/service"
	//	. "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/errors"
	"github.com/CNES/ccsdsmo-malgo-examples/archiveservice/testarchivearea"
	"github.com/CNES/ccsdsmo-malgo-examples/archiveservice/testarchivearea/testarchiveservice"
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

func NewValueOfSine(f mal.Float) *testarchiveservice.ValueOfSine {
	return &testarchiveservice.ValueOfSine{f}
}

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
		var elementList = testarchiveservice.NewValueOfSineList(0)
		var boolean = mal.NewBoolean(false)
		// Variable for the different networks
		var networks = []*mal.Identifier{
			mal.NewIdentifier("tests/network1"),
			mal.NewIdentifier("tests/network2"),
		}
		// Variable for the different providers
		var providers = []*mal.URI{
			mal.NewURI("tests/provider1"),
			mal.NewURI("tests/provider2"),
		}

		var objectType com.ObjectType
		var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
		var archiveDetailsList = *archive.NewArchiveDetailsList(0)

		// Create elements
		for i := 0; i < numberOfRows/2; i++ {
			// Create the value
			var signe = float64(rand.Int63n(2))
			if signe == 0 {
				elementList.AppendElement(NewValueOfSine(mal.Float(rand.Float64())))
			} else {
				elementList.AppendElement(NewValueOfSine(mal.Float(-rand.Float64())))
			}
			objectType = com.ObjectType{
				Area:    testarchivearea.AREA_NUMBER,
				Service: testarchiveservice.SERVICE_NUMBER,
				Version: testarchivearea.AREA_VERSION,
				Number:  mal.UShort((*elementList)[i].GetTypeShortForm()),
			}
			// Object instance identifier
			var objectInstanceIdentifier = mal.Long(int64(i + 1))
			// Variables for ArchiveDetailsList
			var objectKey = com.ObjectKey{
				Domain: identifierList,
				InstId: mal.Long(0),
			}
			var objectID = com.ObjectId{
				Type: objectType,
				Key:  objectKey,
			}
			var objectDetails = com.ObjectDetails{
				Related: mal.NewLong(0),
				Source:  &objectID,
			}
			var network = networks[rand.Int63n(int64(len(networks)))]
			var timestamp = mal.NewFineTime(time.Now())
			var provider = providers[rand.Int63n(int64(len(providers)))]
			archiveDetailsList.AppendElement(&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider})
		}
		_, err := archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
		if err != nil {
			return err
		}

		// Store fourty new elements (total 80 elements)
		identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("en"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
		for i := 0; i < archiveDetailsList.Size(); i++ {
			var objectInstanceIdentifier = mal.Long(int64(i + 41))
			archiveDetailsList[i].InstId = objectInstanceIdentifier
			archiveDetailsList[i].Details.Source.Key.Domain = identifierList
		}
		_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
		if err != nil {
			return err
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

func TestMain(m *testing.M) {
	err := testSetup()
	if err != nil {
		return
	}
	retCode := m.Run()
	_ = testTeardown()
	os.Exit(retCode)
}

var malContext *mal.Context = nil
var clientContext *malapi.ClientContext = nil

func testSetup() error {
	dfltConsumerURL := "maltcp://127.0.0.1:14200"
	malContext, err := mal.NewContext(dfltConsumerURL)
	if err != nil {
		fmt.Printf("error creating MAL context for URI %s: %s", dfltConsumerURL, err)
		return err
	}
	clientContext, err = malapi.NewClientContext(malContext, "test")
	if err != nil {
		fmt.Printf("error creating client context: %s", err)
		return err
	}
	// InitMalContext(clientContext)
	archive.Init(clientContext)
	return nil
}

func testTeardown() error {
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
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	var longList = mal.LongList([]*mal.Long{mal.NewLong(0)})

	// Variables to retrieve the return of this function
	var archiveDetailsList *archive.ArchiveDetailsList
	var elementList mal.ElementList
	// Start the consumer
	archiveDetailsList, elementList, err = archiveService.Retrieve(providerURL, objectType, identifierList, longList)

	if err != nil || archiveDetailsList == nil || elementList == nil {
		println(err)
		t.FailNow()
	}
}

func TestRetrieveKO_3_4_3_2_2(t *testing.T) {
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

	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	var longList = mal.LongList([]*mal.Long{mal.NewLong(0)})
	// Area is equal to 0
	var objectType = com.ObjectType{
		Area:    mal.UShort(0),
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, err = archiveService.Retrieve(providerURL, objectType, identifierList, longList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	// Service is equal to 0
	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: mal.UShort(0),
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, err = archiveService.Retrieve(providerURL, objectType, identifierList, longList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	// Version is equal to 0
	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: mal.UOctet(0),
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, err = archiveService.Retrieve(providerURL, objectType, identifierList, longList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	// Number is equal to 0
	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(0),
	}
	// Start the consumer
	_, _, err = archiveService.Retrieve(providerURL, objectType, identifierList, longList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestRetrieveKO_3_4_3_2_4(t *testing.T) {
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

	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("*"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	var longList = mal.LongList([]*mal.Long{mal.NewLong(0)})
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, _, err = archiveService.Retrieve(providerURL, objectType, identifierList, longList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("*"), mal.NewIdentifier("archiveservice")})
	_, _, err = archiveService.Retrieve(providerURL, objectType, identifierList, longList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("*")})
	_, _, err = archiveService.Retrieve(providerURL, objectType, identifierList, longList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	//var domain = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	archiveQuery := &archive.ArchiveQuery{
		Related:   mal.Long(0),
		SortOrder: mal.NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	//	var queryFilterList *archive.CompositeFilterSetList

	// Variable to retrieve the responses
	var responses []interface{}

	// Start the consumer
	//	responses, err = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)
	responses, err = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, archive.NullCompositeFilterSetList)

	if err != nil || responses == nil {
		t.FailNow()
	}
}

func TestQueryKO_3_4_4_2_9(t *testing.T) {
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
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	archiveQuery := &archive.ArchiveQuery{
		Related:   mal.Long(0),
		SortOrder: mal.NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList = archive.NewCompositeFilterSetList(1)

	// Start the consumer
	_, err = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestQueryOK_3_4_4_2_14(t *testing.T) {
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
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	archiveQuery := &archive.ArchiveQuery{
		Related:       mal.Long(0),
		SortOrder:     mal.NewBoolean(true),
		SortFieldName: mal.NewString("domain"),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *archive.CompositeFilterSetList

	// Variable to retrieve the responses
	var responses []interface{}

	// Start the consumer
	responses, err = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	if err != nil || responses == nil {
		t.FailNow()
	}
	// Now, verify the responses
	for i, resp := range responses {
		if i%4 == 1 {
			if (i == 1 && *(*resp.(*mal.IdentifierList))[0] != "en") || (i == 5 && *(*resp.(*mal.IdentifierList))[0] != "fr") {
				t.FailNow()
			}
		} else if i%4 == 3 {
			if resp.(mal.ElementList).Size() != numberOfRows/2 {
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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	archiveQuery := &archive.ArchiveQuery{
		Related:       mal.Long(0),
		SortOrder:     mal.NewBoolean(true),
		SortFieldName: mal.NewString("invalidname"),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *archive.CompositeFilterSetList

	// Start the consumer
	_, err = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestQueryKO_3_4_4_2_16(t *testing.T) {
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
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	archiveQuery := &archive.ArchiveQuery{
		Related: mal.Long(0),
	}
	archiveQueryList.AppendElement(archiveQuery)
	// Create the query filters list
	var queryFilterList = archive.NewCompositeFilterSetList(0)
	compositeFilter := &archive.CompositeFilter{mal.String("domain"), archive.EXPRESSIONOPERATOR_CONTAINS, nil}
	compositeFilterList := archive.NewCompositeFilterList(0)
	compositeFilterList.AppendElement(compositeFilter)
	queryFilter := &archive.CompositeFilterSet{*compositeFilterList}
	queryFilterList.AppendElement(queryFilter)

	// Start the consumer
	_, err = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
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
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	archiveQuery := &archive.ArchiveQuery{
		Related: mal.Long(0),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *archive.CompositeFilterSetList

	// Start the consumer WITHOUT any wildcard value in the objectType
	resp, _ := archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	for i := 0; i < len(resp)/4; i++ {
		objType := resp[i*4].(*com.ObjectType)
		if objType != nil {
			t.FailNow()
		}
	}

	// Start the consumer WITH a wildcard value in the objectType
	objectType.Area = 0
	resp, _ = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	for i := 0; i < len(resp)/4; i++ {
		objType := resp[i*4].(*com.ObjectType)
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
	var boolean *mal.Boolean
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	archiveQuery := &archive.ArchiveQuery{
		Related: mal.Long(0),
	}
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList *archive.CompositeFilterSetList

	// Start the consumer with the initial Boolean set to NIL
	var newBoolean *mal.Boolean
	resp, _ := archiveService.Query(providerURL, newBoolean, objectType, *archiveQueryList, queryFilterList)

	fmt.Println(resp)
	for i := 0; i < len(resp)/4; i++ {
		elementList := resp[i*4+3]
		if elementList != mal.ElementList(nil) {
			t.FailNow()
		}
	}

	// Start the consumer with the initial Boolean set to TRUE
	boolean = mal.NewBoolean(true)
	resp, _ = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	for i := 0; i < len(resp)/4; i++ {
		elementList := resp[i*4+3]
		if elementList == mal.ElementList(nil) {
			t.FailNow()
		}
	}

	// Start the consumer with the initial Boolean set to FALSE
	boolean = mal.NewBoolean(false)
	resp, _ = archiveService.Query(providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

	fmt.Println(resp)
	for i := 0; i < len(resp)/4; i++ {
		elementList := resp[i*4+3]
		if elementList != mal.ElementList(nil) {
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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	var domain = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	archiveQuery := &archive.ArchiveQuery{
		Domain:    &domain,
		Related:   mal.Long(0),
		SortOrder: mal.NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	domain2 := mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("en"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	archiveQuery2 := &archive.ArchiveQuery{
		Domain:    &domain2,
		Related:   mal.Long(0),
		SortOrder: mal.NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery2)
	var queryFilterList *archive.CompositeFilterSetList

	// Variable to retrieve the return of this function
	var longList *mal.LongList
	// Start the consumer
	longList, err = archiveService.Count(providerURL, objectType, archiveQueryList, queryFilterList)

	if err != nil || longList == nil {
		t.FailNow()
	}

	if longList.Size() != 2 {
		t.FailNow()
	}
	for i := 0; i < longList.Size(); i++ {
		if *longList.GetElementAt(i).(*mal.Long) != mal.Long(40) {
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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	var domain = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	archiveQuery := &archive.ArchiveQuery{
		Domain:    &domain,
		Related:   mal.Long(0),
		SortOrder: mal.NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)
	archiveQueryList.AppendElement(archiveQuery)
	var queryFilterList = archive.NewCompositeFilterSetList(1)

	_, err = archiveService.Count(providerURL, objectType, archiveQueryList, queryFilterList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestCountOK_3_4_5_2_14(t *testing.T) {
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

	var objectType = &com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	var domain = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	archiveQuery := &archive.ArchiveQuery{
		Domain:        &domain,
		Related:       mal.Long(0),
		SortOrder:     mal.NewBoolean(true),
		SortFieldName: mal.NewString("domain"),
	}
	archiveQueryList.AppendElement(archiveQuery)

	var queryFilterList *archive.CompositeFilterSetList

	// Variable to retrieve the return of this function
	var longList *mal.LongList
	// Start the consumer
	longList, err = archiveService.Count(providerURL, objectType, archiveQueryList, queryFilterList)

	if err != nil || longList == nil {
		t.FailNow()
	}

	if longList.Size() != 1 {
		t.FailNow()
	}
	for i := 0; i < longList.Size(); i++ {
		if *longList.GetElementAt(i).(*mal.Long) != mal.Long(40) {
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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = &com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	var domain = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	archiveQuery := &archive.ArchiveQuery{
		Domain:        &domain,
		Related:       mal.Long(0),
		SortOrder:     mal.NewBoolean(true),
		SortFieldName: mal.NewString("invalidname"),
	}
	archiveQueryList.AppendElement(archiveQuery)

	var queryFilterList *archive.CompositeFilterSetList

	// Start the consumer
	_, err = archiveService.Count(providerURL, objectType, archiveQueryList, queryFilterList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestCountKO_3_4_5_2_16(t *testing.T) {
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

	var objectType = &com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	archiveQueryList := archive.NewArchiveQueryList(0)
	var domain = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	archiveQuery := &archive.ArchiveQuery{
		Domain:    &domain,
		Related:   mal.Long(0),
		SortOrder: mal.NewBoolean(true),
	}
	archiveQueryList.AppendElement(archiveQuery)

	// Create the query filters list
	var queryFilterList = archive.NewCompositeFilterSetList(0)
	compositeFilter := &archive.CompositeFilter{mal.String("domain"), archive.EXPRESSIONOPERATOR_CONTAINS, nil}
	compositeFilterList := archive.NewCompositeFilterList(0)
	compositeFilterList.AppendElement(compositeFilter)
	queryFilter := &archive.CompositeFilterSet{*compositeFilterList}
	queryFilterList.AppendElement(queryFilter)

	// Start the consumer
	_, err = archiveService.Count(providerURL, objectType, archiveQueryList, queryFilterList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *mal.NewLong(81)
	// Variables for ArchiveDetailsList
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("network")
	var timestamp = mal.NewFineTime(time.Now())
	var provider = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Variable to retrieve the return of this function
	var longList *mal.LongList
	// Start the consumer
	longList, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if err != nil || longList == nil {
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
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean *mal.Boolean
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *mal.NewLong(0)
	// Variables for ArchiveDetailsList
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("network")
	var timestamp = mal.NewFineTime(time.Now())
	var provider = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Variable to retrieve the return of this function
	var longList *mal.LongList

	// First, start the consumer with the boolean set to NIL
	// Start the consumer
	longList, _ = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if longList != nil {
		t.FailNow()
	}

	// Then, start the consumer with the boolean set to FALSE
	boolean = mal.NewBoolean(false)
	// Start the consumer
	longList, _ = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

	if longList != nil {
		t.FailNow()
	}

	// Finally, start the consumer with the boolean set to TRUE
	boolean = mal.NewBoolean(true)
	// Start the consumer
	longList, _ = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

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
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *mal.NewLong(0)
	// Variables for ArchiveDetailsList
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("network")
	var timestamp = mal.NewFineTime(time.Now())
	var provider = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Variable to retrieve the return of this function
	var longList *mal.LongList

	// Start the consumer
	longList, _ = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the store consumer
	// Create parameters
	// Object that's going to be stored in the archive
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = *mal.NewLong(45)
	// Variables for ArchiveDetailsList
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("network")
	var timestamp = mal.NewFineTime(time.Now())
	var provider = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_DUPLICATE {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_8(t *testing.T) {
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
	var elementList = testarchiveservice.NewValueOfSineList(2)
	(*elementList)[0] = NewValueOfSine(0)
	(*elementList)[1] = NewValueOfSine(0.5)
	var boolean = mal.NewBoolean(true)
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = mal.Long(0)
	// Variables for ArchiveDetailsList
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("network")
	var timestamp = mal.NewFineTime(time.Now())
	var provider = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
	errid, isUIntegerList := malerr.ExtraInfo.(*mal.UIntegerList)
	if !isUIntegerList {
		fmt.Printf("DEBUG SL no LongList\n")
		t.FailNow()
	} else {
		fmt.Printf("DEBUG SL %v\n", *(*errid)[0])
	}
	if !isUIntegerList || uint32(*(*errid)[0]) != 1 {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_9(t *testing.T) {
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
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = mal.NewBoolean(true)
	// Area is equal to 0
	var objectType = com.ObjectType{
		Area:    mal.UShort(0),
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = mal.Long(0)
	// Variables for ArchiveDetailsList
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("network")
	var timestamp = mal.NewFineTime(time.Now())
	var provider = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	// Service is equal to 0
	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: mal.UShort(0),
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	// Version is equal to 0
	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: mal.UOctet(0),
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	// Number is equal to 0
	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(0),
	}
	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_10(t *testing.T) {
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
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = mal.NewBoolean(true)
	// Area is equal to 0
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("*"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = mal.Long(0)
	// Variables for ArchiveDetailsList
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("network")
	var timestamp = mal.NewFineTime(time.Now())
	var provider = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestStoreKO_3_4_6_2_11(t *testing.T) {
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
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0)
	var boolean = mal.NewBoolean(true)
	// Area is equal to 0
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice")})
	// Object instance identifier
	var objectInstanceIdentifier = mal.Long(0)
	// Variables for ArchiveDetailsList
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	// Bad value for the NETWORK
	var network *mal.Identifier
	var timestamp = mal.NewFineTime(time.Now())
	var provider = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	// Bad value for the TIMESTAMP
	network = mal.NewIdentifier("network")
	*timestamp = mal.FineTime(time.Unix(int64(0), int64(0)))
	archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	// Bad value for the TIMESTAMP
	*timestamp = mal.FineTime(time.Now())
	*provider = mal.URI("*")
	archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, timestamp, provider}})

	// Start the consumer
	_, err = archiveService.Store(providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	// Object instance identifier
	var objectInstanceIdentifier = *mal.NewLong(1)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("new.network")
	var fineTime = mal.NewFineTime(time.Now())
	var uri = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, fineTime, uri}})

	// Start the consumer
	err = archiveService.Update(providerURL, objectType, identifierList, archiveDetailsList, elementList)

	if err != nil {
		t.FailNow()
	}
}

func TestUpdateKO_3_4_7_2_5(t *testing.T) {
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

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	// Object instance identifier: UNKNOWN
	var objectInstanceIdentifier = *mal.NewLong(155)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("new.network")
	var fineTime = mal.NewFineTime(time.Now())
	var uri = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, fineTime, uri}})

	// Start the consumer
	err = archiveService.Update(providerURL, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != mal.ERROR_UNKNOWN {
		t.FailNow()
	}
}

func TestUpdateKO_3_4_7_2_8_ObjectType(t *testing.T) {
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

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = com.ObjectType{
		Area:    mal.UShort(0),
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	// Object instance identifier: UNKNOWN
	var objectInstanceIdentifier = *mal.NewLong(41)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("new.network")
	var fineTime = mal.NewFineTime(time.Now())
	var uri = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, fineTime, uri}})

	// Start the consumer
	err = archiveService.Update(providerURL, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: mal.UShort(0),
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort((*elementList)[0].GetTypeShortForm()),
	}
	// Start the consumer
	err = archiveService.Update(providerURL, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: mal.UOctet(0),
		Number:  mal.UShort((*elementList)[0].GetTypeShortForm()),
	}
	// Start the consumer
	err = archiveService.Update(providerURL, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(0),
	}
	// Start the consumer
	err = archiveService.Update(providerURL, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestUpdateKO_3_4_7_2_8_ObjectInstanceIdentifier(t *testing.T) {
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

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	// Object instance identifier
	var objectInstanceIdentifier = *mal.NewLong(0)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("new.network")
	var fineTime = mal.NewFineTime(time.Now())
	var uri = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, fineTime, uri}})

	// Start the consumer
	err = archiveService.Update(providerURL, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestUpdateKO_3_4_7_2_8_Domain(t *testing.T) {
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

	// Start the update consumer
	// Create parameters
	// ---- ELEMENTLIST ----
	// Object that's going to be updated in the archive
	var elementList = testarchiveservice.NewValueOfSineList(1)
	(*elementList)[0] = NewValueOfSine(0.5)
	// ---- OBJECTTYPE ----
	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort((*elementList)[0].GetTypeShortForm()),
	}
	// ---- IDENTIFIERLIST ----
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("*"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	// Object instance identifier
	var objectInstanceIdentifier = *mal.NewLong(41)
	// Variables for ArchiveDetailsList
	// ---- ARCHIVEDETAILSLIST ----
	var objectKey = com.ObjectKey{
		Domain: identifierList,
		InstId: objectInstanceIdentifier,
	}
	var objectID = com.ObjectId{
		Type: objectType,
		Key:  objectKey,
	}
	var objectDetails = com.ObjectDetails{
		Related: mal.NewLong(1),
		Source:  &objectID,
	}
	var network = mal.NewIdentifier("new.network")
	var fineTime = mal.NewFineTime(time.Now())
	var uri = mal.NewURI("main/start")
	var archiveDetailsList = archive.ArchiveDetailsList([]*archive.ArchiveDetails{&archive.ArchiveDetails{objectInstanceIdentifier, objectDetails, network, fineTime, uri}})

	// Start the consumer
	err = archiveService.Update(providerURL, objectType, identifierList, archiveDetailsList, elementList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
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

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	var longList = mal.NewLongList(0)
	longList.AppendElement(mal.NewLong(15))

	// Variable to retrieve the return of this function
	var respLongList *mal.LongList
	// Start the consumer
	respLongList, err = archiveService.Delete(providerURL, objectType, identifierList, *longList)

	if err != nil || respLongList == nil || respLongList.Size() == 0 {
		t.FailNow()
	}
}

func TestDeleteKO_3_4_8_2_3_ObjectType(t *testing.T) {
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

	var objectType = com.ObjectType{
		Area:    mal.UShort(0),
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	var longList = mal.NewLongList(0)
	longList.AppendElement(mal.NewLong(15))

	// Start the consumer
	_, err = archiveService.Delete(providerURL, objectType, identifierList, *longList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: mal.UShort(0),
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, err = archiveService.Delete(providerURL, objectType, identifierList, *longList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: mal.UOctet(0),
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	// Start the consumer
	_, err = archiveService.Delete(providerURL, objectType, identifierList, *longList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}

	objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(0),
	}
	// Start the consumer
	_, err = archiveService.Delete(providerURL, objectType, identifierList, *longList)
	malerr, ismalerr = err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestDeleteKO_3_4_8_2_3_Domain(t *testing.T) {
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

	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("*"), mal.NewIdentifier("test")})
	var longList = mal.NewLongList(0)
	longList.AppendElement(mal.NewLong(15))

	// Start the consumer
	_, err = archiveService.Delete(providerURL, objectType, identifierList, *longList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != com.ERROR_INVALID {
		t.FailNow()
	}
}

func TestDeleteKO_3_4_8_2_6(t *testing.T) {
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

	var objectType = com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: testarchiveservice.SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(testarchiveservice.VALUEOFSINE_TYPE_SHORT_FORM),
	}
	var identifierList = mal.IdentifierList([]*mal.Identifier{mal.NewIdentifier("fr"), mal.NewIdentifier("cnes"), mal.NewIdentifier("archiveservice"), mal.NewIdentifier("test")})
	var longList = mal.NewLongList(0)
	longList.AppendElement(mal.NewLong(175))

	// Start the consumer
	_, err = archiveService.Delete(providerURL, objectType, identifierList, *longList)
	malerr, ismalerr := err.(*malapi.MalError)
	if err == nil || !ismalerr || malerr.Code != mal.ERROR_UNKNOWN {
		t.FailNow()
	}
}
