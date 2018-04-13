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

	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/archive/service"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/errors"
)

func main() {
	args := os.Args[1:]

	var _list = MAL_SHORT_LIST_SHORT_FORM
	var _type = MAL_SHORT_SHORT_FORM
	fmt.Println(_list)
	fmt.Println(_type)
	fmt.Println(_list - _type)

	fmt.Printf("BEFORE:\n%08b\n", _type)
	quatuor5 := (MAL_SHORT_SHORT_FORM & 0x0000F0) >> 4
	fmt.Println(quatuor5)
	var listByte []byte
	listByte = append(listByte, 1, 0, 0, 1)
	if quatuor5 == 0x0 {
		var b byte
		for i := 2; i >= 0; i-- {
			b = byte(_type>>uint(i*8)) ^ 255
			if i == 0 {
				b++
			}
			listByte = append(listByte, b)
		}
	} else {

	}

	var finalListType = uint64(listByte[0]) << 42
	finalListType &= uint64(listByte[0]) << 36
	fmt.Printf("%08b\n", finalListType)
	fmt.Printf("%08b\n", _list)

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
				UShort(MAL_LONG_TYPE_SHORT_FORM),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
			var longList = LongList([]*Long{NewLong(29)})

			// Variables to retrieve the return of this function
			var archiveDetailsList *ArchiveDetailsList
			var elementList ElementList
			// Start the consumer
			archiveDetailsList, elementList, errorsList, err = archiveService.LaunchRetrieveConsumer(objectType, identifierList, longList)

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
				UShort(MAL_LONG_TYPE_SHORT_FORM),
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
				UShort((*elementList)[0].GetTypeShortForm()),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
			// Object instance identifier
			var objectInstanceIdentifier = *NewLong(10)
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
			// ---- ELEMENTLIST ----
			// Object that's going to be updated in the archive
			var elementList = NewLongList(1)
			(*elementList)[0] = NewLong(29)
			// ---- OBJECTTYPE ----
			var objectType = ObjectType{
				UShort(archiveService.AreaNumber),
				UShort(archiveService.ServiceNumber),
				UOctet(archiveService.AreaVersion),
				UShort((*elementList)[0].GetTypeShortForm()),
			}
			// ---- IDENTIFIERLIST ----
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
			// Object instance identifier
			var objectInstanceIdentifier = *NewLong(13)
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
			var network = NewIdentifier("network")
			var fineTime = NewFineTime(time.Now())
			var uri = NewURI("main/start")
			var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, fineTime, uri)})

			// Start the consumer
			errorsList, err = archiveService.LaunchUpdateConsumer(objectType, identifierList, archiveDetailsList, elementList)

			break
		case "delete":
			// Start the delete consumer
			// Create parameters
			var objectType = ObjectType{
				UShort(archiveService.AreaNumber),
				UShort(archiveService.ServiceNumber),
				UOctet(archiveService.AreaVersion),
				UShort(MAL_LONG_TYPE_SHORT_FORM),
			}
			var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
			var longList = NewLongList(10)

			// Variable to retrieve the return of this function
			var respLongList *LongList
			// Start the consumer
			respLongList, errorsList, err = archiveService.LaunchDeleteConsumer(objectType, identifierList, *longList)

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
