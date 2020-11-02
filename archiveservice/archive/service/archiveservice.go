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
 /
/**
 * This file has been duplicated from the original source in github.com/etiennelndr/archiveservice/archive/service.
 * It provides the consumer interface to the COM/Archive service.
 * It could be renamed consumer.
 *
 * The implementation has been changed to use a single Context and a single ClientContext from the malgo implementation.
 * The interface has been change accordingly to remove reference to a client URL.
*/
package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/com/archive"
	"github.com/CNES/ccsdsmo-malgo/mal"

	// Init TCP transport
	_ "github.com/CNES/ccsdsmo-malgo/mal/transport/tcp"

	// Blank imports to register all the mal and com elements
	_ "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/testarchivearea/testarchiveservice"

	. "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/archive/provider"
	. "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/service"
)

// ArchiveService : TODO:
type ArchiveService struct {
	AreaIdentifier    mal.Identifier
	ServiceIdentifier mal.Identifier
	AreaNumber        mal.UShort
	ServiceNumber     mal.UShort
	AreaVersion       mal.UOctet

	running bool
	wg      sync.WaitGroup
}

// CreateService : TODO:
func (*ArchiveService) CreateService() Service {
	archiveService := &ArchiveService{
		AreaIdentifier:    com.AREA_NAME,
		ServiceIdentifier: archive.SERVICE_NAME,
		AreaNumber:        com.AREA_NUMBER,
		ServiceNumber:     archive.SERVICE_NUMBER,
		AreaVersion:       com.AREA_VERSION,
		running:           true,
		wg:                *new(sync.WaitGroup),
	}

	return archiveService
}

//======================================================================//
//                          START: Consumer                             //
//======================================================================//

// Retrieve : TODO:
func (archiveService *ArchiveService) Retrieve(providerURL string, objectType com.ObjectType, identifierList mal.IdentifierList, longList mal.LongList) (*archive.ArchiveDetailsList, mal.ElementList, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve operation")

	// IN
	var providerURI = mal.NewURI(providerURL + "/archiveServiceProvider")

	op, err := archive.NewRetrieveOperation(providerURI)
	if err != nil {
		return nil, nil, err
	}
	err = op.Invoke(&objectType, &identifierList, &longList)
	if err != nil {
		return nil, nil, err
	}
	archiveDetailsList, elementList, err := op.GetResponse()
	if err != nil {
		return nil, nil, err
	}

	return archiveDetailsList, elementList, nil
}

// Query : TODO:
func (archiveService *ArchiveService) Query(providerURL string, boolean *mal.Boolean, objectType com.ObjectType, archiveQueryList archive.ArchiveQueryList, queryFilterList archive.QueryFilterList) ([]interface{}, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Query operation")

	// IN
	var providerURI = mal.NewURI(providerURL + "/archiveServiceProvider")
	op, err := archive.NewQueryOperation(providerURI)
	if err != nil {
		return nil, err
	}
	err = op.Progress(boolean, &objectType, &archiveQueryList, queryFilterList)
	if err != nil {
		return nil, err
	}

	// Create the interface that will receive all the responses
	responses := []interface{}{}

	for endUpdateLoop := false; !endUpdateLoop; {
		// Call Update operation
		respObjType, respIDList, respArchDetList, respElemList, err := op.GetUpdate()
		if err != nil {
			return responses, err
		}
		if respArchDetList != nil {
			// Put the objects in the interface
			responses = append(responses, respObjType, respIDList, respArchDetList, respElemList)
		} else {
			endUpdateLoop = true
		}
	}

	// Call Response operation
	respObjType, respIDList, respArchDetList, respElemList, err := op.GetResponse()
	if err != nil {
		return responses, err
	}
	if respArchDetList != nil {
		// Put the objects in the interface
		responses = append(responses, respObjType, respIDList, respArchDetList, respElemList)
	}

	return responses, nil
}

// Count : TODO:
func (archiveService *ArchiveService) Count(providerURL string, objectType *com.ObjectType, archiveQueryList *archive.ArchiveQueryList, queryFilterList archive.QueryFilterList) (*mal.LongList, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Count operation")

	// IN
	var providerURI = mal.NewURI(providerURL + "/archiveServiceProvider")
	op, err := archive.NewCountOperation(providerURI)
	if err != nil {
		return nil, err
	}
	err = op.Invoke(objectType, archiveQueryList, queryFilterList)
	if err != nil {
		return nil, err
	}
	longList, err := op.GetResponse()
	if err != nil {
		return nil, err
	}

	return longList, nil
}

// Store : TODO:
func (archiveService *ArchiveService) Store(providerURL string, boolean *mal.Boolean, objectType com.ObjectType, identifierList mal.IdentifierList, archiveDetailsList archive.ArchiveDetailsList, elementList mal.ElementList) (*mal.LongList, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Store operation")

	// IN
	var providerURI = mal.NewURI(providerURL + "/archiveServiceProvider")
	op, err := archive.NewStoreOperation(providerURI)
	if err != nil {
		return nil, err
	}
	longList, err := op.Request(boolean, &objectType, &identifierList, &archiveDetailsList, elementList)
	if err != nil {
		return nil, err
	}

	return longList, nil
}

// Update : TODO:
func (archiveService *ArchiveService) Update(providerURL string, objectType com.ObjectType, identifierList mal.IdentifierList, archiveDetailsList archive.ArchiveDetailsList, elementList mal.ElementList) error {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Update operation")

	// IN
	var providerURI = mal.NewURI(providerURL + "/archiveServiceProvider")
	op, err := archive.NewUpdateOperation(providerURI)
	if err != nil {
		return err
	}
	err = op.Submit(&objectType, &identifierList, &archiveDetailsList, elementList)
	if err != nil {
		return err
	}

	return nil
}

// Delete : TODO:
func (archiveService *ArchiveService) Delete(providerURL string, objectType com.ObjectType, identifierList mal.IdentifierList, longList mal.LongList) (*mal.LongList, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Delete Consumer")

	// IN
	var providerURI = mal.NewURI(providerURL + "/archiveServiceProvider")
	op, err := archive.NewDeleteOperation(providerURI)
	if err != nil {
		return nil, err
	}
	respLongList, err := op.Request(&objectType, &identifierList, &longList)
	if err != nil {
		return nil, err
	}

	return respLongList, nil
}

//======================================================================//
//                          START: Provider                             //
//======================================================================//

// StartProvider : TODO:
func (archiveService *ArchiveService) StartProvider(providerURL string) error {
	archiveService.wg.Add(2)
	// Start the retrieve provider
	go archiveService.launchProvider(providerURL)
	// Start a simple method to stop the providers
	go archiveService.stopProviders()
	// Wait until the end of the six operations
	archiveService.wg.Wait()

	return nil
}

// stopProviders Stop the providers
func (archiveService *ArchiveService) stopProviders() {
	// Inform the WaitGroup that this goroutine is finished at the end of this function
	defer archiveService.wg.Done()
	// Wait a little bit
	time.Sleep(200 * time.Millisecond)
	for archiveService.running {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Stop providers ? [Yes/No] ")
		text, _ := reader.ReadString('\n')
		stop := strings.TrimRight(text, "\n")
		if len(stop) > 0 && strings.ToLower(stop)[0] == []byte("y")[0] {
			// Good bye
			archiveService.running = false
		}
	}
}

// launchSpecificProvider Start a provider for a specific operation
func (archiveService *ArchiveService) launchProvider(providerURL string) error {
	// Inform the WaitGroup that this goroutine is finished at the end of this function
	defer archiveService.wg.Done()
	// Declare variables
	var provider *archive.Provider
	var err error

	// Start Operation
	provider, err = StartProvider(providerURL)
	if err != nil {
		return err
	}

	// Close the provider at the end of the function
	defer provider.Close()

	// Start communication
	for archiveService.running == true {
		time.Sleep(1 * time.Second)
	}

	return nil
}
