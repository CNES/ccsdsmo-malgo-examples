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

func TestRetrieveOK(t *testing.T) {
	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	// Create parameters
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
	archiveDetailsList, elementList, errorsList, err = archiveService.LaunchRetrieveConsumer(objectType, identifierList, longList)

	if errorsList != nil {
		t.Fatal("ErrorsList must be nil")
	}
	if err != nil {
		t.Fatal("Err must be nil")
	}
	if archiveDetailsList == nil {
		t.Fatal("ArchiveDetailsList must not be nil")
	}
	if elementList == nil {
		t.Fatal("ElementList must not be nil")
	}
}
