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

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
	. "github.com/etiennelndr/archiveservice/archive"
	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/data"
)

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("ERROR: You must use this program like this:\n\tgo run start.go [provider|consumer] [retrieve|query|count|store|update|delete]")
		return
	}

	var archiveService *ArchiveService
	// Create the Archive Service
	element := archiveService.CreateService()
	archiveService = element.(*ArchiveService)

	if args[0] == "provider" {
		switch args[1] {
		case "retrieve":
			// Start the retrieve provider
			archiveService.LaunchProvider(OPERATION_IDENTIFIER_RETRIEVE)
			break
		case "query":
			// Start the query provider
			archiveService.LaunchProvider(OPERATION_IDENTIFIER_QUERY)
			break
		case "count":
			// Start the count provider
			archiveService.LaunchProvider(OPERATION_IDENTIFIER_COUNT)
			break
		case "store":
			// Start the store provider
			archiveService.LaunchProvider(OPERATION_IDENTIFIER_STORE)
			break
		case "update":
			// Start the update provider
			archiveService.LaunchProvider(OPERATION_IDENTIFIER_UPDATE)
			break
		case "delete":
			// Start the delete provider
			archiveService.LaunchProvider(OPERATION_IDENTIFIER_DELETE)
			break
		default:
			fmt.Println("ERROR: You must use this program like this:\n\tgo run start.go [provider|consumer] [retrieve|query|count|store|update|delete]")
			return
		}
	} else if args[0] == "consumer" {
		switch args[1] {
		case "retrieve":
			// Start the retrieve consumer
			// Create parameters
			var objectType ObjectType
			var identifierList IdentifierList
			var longList = NewLongList(10)

			// Start the consumer
			archiveService.LaunchRetrieveConsumer(objectType, identifierList, *longList)
			break
		case "query":
			// Start the query consumer
			/*var boolean Boolean
			var objectType ObjectType
			var archiveQueryList ArchiveQueryList
			var queryFilterList QueryFilterList*/

			break
		case "count":
			// Start the count consumer

			break
		case "store":
			// Start the store consumer

			break
		case "update":
			// Start the update consumer

			break
		case "delete":
			// Start the delete consumer

			break
		default:
			fmt.Println("ERROR: You must use this program like this:\n\tgo run start.go [provider|consumer] [retrieve|query|count|store|update|delete]")
			return
		}
	} else {
		fmt.Println("ERROR: You must use this program like this:\n\tgo run start.go [provider|consumer] [retrieve|query|count|store|update|delete]")
		return
	}
}
