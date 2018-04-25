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
package main

import (
	"fmt"
	"os"
	"time"

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

func main() {
	args := os.Args[1:]

	if len(args) < 1 || len(args) > 2 {
		fmt.Println("ERROR: You must use this program like this:\n\tgo run start.go [provider|[consumer] [retrieve|query|count|store|update|delete]]")
		return
	}

	// Variable that defines the ArchiveService
	var archiveService *ArchiveService
	// Variables to retrieve the different errors
	var err error
	var errorsList *ServiceError
	// Create the Archive Service
	service := archiveService.CreateService()
	archiveService = service.(*ArchiveService)

	if args[0] == "provider" {
		// Start the providers
		archiveService.StartProviders(providerURL)
	} else if args[0] == "consumer" {
		switch args[1] {
		case "retrieve":
			// Start the retrieve consumer
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
			// Start the consumer
			archiveDetailsList, elementList, errorsList, err = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)
			if elementList != nil && err == nil {
				for i := 0; i < elementList.Size(); i++ {
					element := elementList.GetElementAt(i).(*ValueOfSine)
					fmt.Println(*element)
				}
			}
			fmt.Println("Retrieve Consumer received:\n\t>>>", archiveDetailsList, "\n\t>>>", elementList)

			break
		case "query":
			// Start the query consumer
			// Create parameters
			var boolean = NewBoolean(true)
			var objectType = ObjectType{
				UShort(2),
				UShort(3),
				UOctet(1),
				UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
			}
			archiveQueryList := NewArchiveQueryList(0)
			archiveQuery := &ArchiveQuery{
				//Domain: IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")}),
				nil,
				nil,
				nil,
				*NewLong(1),
				nil,
				nil,
				nil,
				nil,
				nil,
			}
			archiveQueryList.AppendElement(archiveQuery)
			var queryFilterList *CompositeFilterSetList

			// Variable to retrieve the responses
			var responses []interface{}
			// Start the consumer
			responses, errorsList, err = archiveService.Query(consumerURL, providerURL, *boolean, objectType, *archiveQueryList, queryFilterList)

			for i := 0; i < len(responses)/4; i++ {
				fmt.Printf("Responses.#%d\n", i)
				fmt.Println("\t> ObjectType        :", responses[i*4])
				fmt.Println("\t> IdentifierList    :", responses[i*4+1])
				fmt.Println("\t> ArchiveDetailsList:", responses[i*4+2])
				fmt.Println("\t> ElementList       :", responses[i*4+3])
			}

			break
		case "count":
			// Start the count consumer
			// Create parameters
			// objectType ObjectType, archiveQueryList ArchiveQueryList, queryFilterList QueryFilterList
			var objectType = ObjectType{
				UShort(2),
				UShort(3),
				UOctet(1),
				UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
			}
			var archiveQueryList = NewArchiveQueryList(10)
			var queryFilterList = NewCompositeFilterSetList(10)

			// Variable to retrieve the return of this function
			var longList *LongList
			// Start the consumer
			longList, errorsList, err = archiveService.Count(consumerURL, providerURL, objectType, *archiveQueryList, queryFilterList)

			fmt.Println("Count Consumer received:\n\t>>>", longList)

			break
		case "store":
			// Start the store consumer
			// Create parameters
			// Object that's going to be stored in the archive
			var elementList = NewValueOfSineList(1)
			(*elementList)[0] = NewValueOfSine(0)
			var boolean = NewBoolean(true)
			var objectType = ObjectType{
				UShort(2),
				UShort(3),
				UOctet(1),
				UShort((*elementList)[0].GetTypeShortForm()),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
			// Object instance identifier
			var objectInstanceIdentifier = *NewLong(2)
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
			longList, errorsList, err = archiveService.Store(consumerURL, providerURL, *boolean, objectType, identifierList, archiveDetailsList, elementList)

			fmt.Println("Store Consumer received:\n\t>>>", longList,
				"\n\t>>>", errorsList)

			break
		case "update":
			// Start the update consumer
			// Create parameters
			// ---- ELEMENTLIST ----
			// Object that's going to be updated in the archive
			var elementList = NewValueOfSineList(1)
			(*elementList)[0] = NewValueOfSine(0.5)
			// ---- OBJECTTYPE ----
			var objectType = ObjectType{
				UShort(2),
				UShort(3),
				UOctet(1),
				UShort((*elementList)[0].GetTypeShortForm()),
			}
			// ---- IDENTIFIERLIST ----
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
			// Object instance identifier
			var objectInstanceIdentifier = *NewLong(1)
			// Variables for ArchiveDetailsList
			// ---- ARCHIVEDETAILSLIST ----
			var objectKey = ObjectKey{
				identifierList,
				objectInstanceIdentifier,
			}
			var objectID = ObjectId{
				&objectType,
				&objectKey,
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

			break
		case "delete":
			// Start the delete consumer
			// Create parameters
			var objectType = ObjectType{
				UShort(2),
				UShort(3),
				UOctet(1),
				UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
			var longList = NewLongList(0)
			longList.AppendElement(NewLong(3))

			// Variable to retrieve the return of this function
			var respLongList *LongList
			// Start the consumer
			respLongList, errorsList, err = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longList)

			fmt.Println("Delete Consumer received:\n\t>>>", respLongList)

			break
		default:
			fmt.Println("ERROR: You must use this program like this:\n\tgo run start.go [provider|[consumer] [retrieve|query|count|store|update|delete]]")
			return
		}
	} else {
		fmt.Println("ERROR: You must use this program like this:\n\tgo run start.go [provider|[consumer] [retrieve|query|count|store|update|delete]]")
		return
	}

	if err != nil {
		fmt.Println("ERROR: sthg unwanted happened,", err)
	} else if errorsList != nil {
		fmt.Println(*errorsList.ErrorNumber)
		fmt.Println(*errorsList.ErrorComment)
	}
}
