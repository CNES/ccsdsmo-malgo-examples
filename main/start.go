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

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive"
	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"
)

func main() {
	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("ERROR: You must use this program like this: go run start.go [provider|consumer]")
		return
	}

	var archiveService *ArchiveService
	// Create the Archive Service
	element := archiveService.CreateService()
	archiveService = element.(*ArchiveService)

	if args[0] == "provider" {
		// Start the provider
		archiveService.StartProvider(OPERATION_IDENTIFIER_RETRIEVE)
	} else if args[0] == "consumer" {
		// Create parameters
		var objectType *ObjectType
		var identifierList *IdentifierList
		var longList = NewLongList(10)

		// Start the consumer
		archiveService.StartConsumer(OPERATION_IDENTIFIER_RETRIEVE, objectType, identifierList, longList)

	} else {
		var bidule Integer
		if bidule == nil {
			println("yoloooooooooooo")
		}
		fmt.Println("ERROR: You must use this program like this: go run start.go [provider|consumer]")
		return
	}
}
