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
	"testing"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"

	. "github.com/etiennelndr/archiveservice/archive/service"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/errors"
	. "github.com/etiennelndr/archiveservice/main/data"
)

// Constants for the providers and consumers
const (
	providerURL = "maltcp://127.0.0.1:12400"
	consumerURL = "maltcp://127.0.0.1:14200"
)

func TestRetrieveOK(t *testing.T) {
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
