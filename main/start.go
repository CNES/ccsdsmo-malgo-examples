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
	"sync"
	"time"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/archive/service"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/errors"
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
	element := archiveService.CreateService()
	archiveService = element.(*ArchiveService)

	if args[0] == "provider" {
		var wg sync.WaitGroup
		wg.Add(6)
		// Start the retrieve provider
		go archiveService.LaunchProvider(OPERATION_IDENTIFIER_RETRIEVE, wg)
		// Start the query provider
		go archiveService.LaunchProvider(OPERATION_IDENTIFIER_QUERY, wg)
		// Start the count provider
		go archiveService.LaunchProvider(OPERATION_IDENTIFIER_COUNT, wg)
		// Start the store provider
		go archiveService.LaunchProvider(OPERATION_IDENTIFIER_STORE, wg)
		// Start the update provider
		go archiveService.LaunchProvider(OPERATION_IDENTIFIER_UPDATE, wg)
		// Start the delete provider
		go archiveService.LaunchProvider(OPERATION_IDENTIFIER_DELETE, wg)
		wg.Wait()
	} else if args[0] == "consumer" {
		switch args[1] {
		case "retrieve":
			// Start the retrieve consumer
			// Create parameters
			var objectType = ObjectType{
				UShort(archiveService.AreaNumber),
				UShort(archiveService.ServiceNumber),
				UOctet(archiveService.AreaVersion),
				UShort(archiveService.ServiceNumber),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("test"), NewIdentifier("archiveService")})
			var longList = LongList([]*Long{NewLong(29), NewLong(31)})

			// // Variables to retrieve the return of this function
			var archiveDetailsList *ArchiveDetailsList
			var elementList ElementList
			// Start the consumer
			archiveDetailsList, elementList, err = archiveService.LaunchRetrieveConsumer(objectType, identifierList, longList)

			fmt.Println("Retrieve Consumer received:\n\t>>>", archiveDetailsList, "\n\t>>>", elementList)

			break
		case "query":
			// Start the query consumer
			// Create parameters
			var boolean = NewBoolean(true)
			var objectType = ObjectType{
				UShort(archiveService.AreaNumber),
				UShort(archiveService.ServiceNumber),
				UOctet(archiveService.AreaVersion),
				UShort(archiveService.ServiceNumber),
			}
			var archiveQueryList = NewArchiveQueryList(10)
			var queryFilterList = NewCompositeFilterSetList(10)

			// Variable to retrieve the responses
			var responses []interface{}
			// Start the consumer
			responses, err = archiveService.LaunchQueryConsumer(*boolean, objectType, *archiveQueryList, queryFilterList)

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
				UShort(archiveService.AreaNumber),
				UShort(archiveService.ServiceNumber),
				UOctet(archiveService.AreaVersion),
				UShort(archiveService.ServiceNumber),
			}
			var archiveQueryList = NewArchiveQueryList(10)
			var queryFilterList = NewCompositeFilterSetList(10)

			// Variable to retrieve the return of this function
			var longList *LongList
			// Start the consumer
			longList, err = archiveService.LaunchCountConsumer(objectType, *archiveQueryList, queryFilterList)

			fmt.Println("Count Consumer received:\n\t>>>", longList)

			break
		case "store":
			// Start the store consumer
			// Create parameters
			// Object that's going to be stored in the archive
			var elementList = NewLongList(1)
			(*elementList)[0] = NewLong(29)

			var boolean = NewBoolean(true)
			var objectType = ObjectType{
				UShort(archiveService.AreaNumber),
				UShort(archiveService.ServiceNumber),
				UOctet(archiveService.AreaVersion),
				UShort(elementList.GetShortForm()),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("test"), NewIdentifier("archiveService")})
			// Object instance identifier
			var objectInstanceIdentifier = *NewLong(0)
			// Variables for ArchiveDetailsList
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
			var network = NewIdentifier("network")
			var fineTime = NewFineTime(time.Now())
			var uri = NewURI("main/start")
			var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, fineTime, uri)})

			// Variable to retrieve the return of this function
			var longList *LongList
			// Start the consumer
			longList, errorsList, err = archiveService.LaunchStoreConsumer(*boolean, objectType, identifierList, archiveDetailsList, elementList)

			fmt.Println("Store Consumer received:\n\t>>>", longList,
				"\n\t>>>", errorsList)

			break
		case "update":
			// Start the update consumer
			// Create parameters
			var objectType = ObjectType{
				UShort(archiveService.AreaNumber),
				UShort(archiveService.ServiceNumber),
				UOctet(archiveService.AreaVersion),
				UShort(archiveService.ServiceNumber),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("test"), NewIdentifier("archiveService")})
			var archiveDetailsList = NewArchiveDetailsList(10)
			var elementList = NewLongList(10)

			// Start the consumer
			err = archiveService.LaunchUpdateConsumer(objectType, identifierList, *archiveDetailsList, elementList)

			break
		case "delete":
			// Start the delete consumer
			// Create parameters
			var objectType = ObjectType{
				UShort(archiveService.AreaNumber),
				UShort(archiveService.ServiceNumber),
				UOctet(archiveService.AreaVersion),
				UShort(archiveService.ServiceNumber),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("test"), NewIdentifier("archiveService")})
			var longList = NewLongList(10)

			// Variable to retrieve the return of this function
			var respLongList *LongList
			// Start the consumer
			respLongList, err = archiveService.LaunchDeleteConsumer(objectType, identifierList, *longList)

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
	}
}
