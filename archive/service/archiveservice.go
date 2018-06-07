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
package service

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"

	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/archive/consumer"
	. "github.com/etiennelndr/archiveservice/archive/provider"
	. "github.com/etiennelndr/archiveservice/data"
	_ "github.com/etiennelndr/archiveservice/data"
	_ "github.com/etiennelndr/archiveservice/data/implementation"
	_ "github.com/etiennelndr/archiveservice/data/tests"
	. "github.com/etiennelndr/archiveservice/errors"
	. "github.com/etiennelndr/archiveservice/service"
)

// ArchiveService : TODO:
type ArchiveService struct {
	AreaIdentifier    Identifier
	ServiceIdentifier Identifier
	AreaNumber        UShort
	ServiceNumber     Integer
	AreaVersion       UOctet

	running bool
	wg      sync.WaitGroup
}

// CreateService : TODO:
func (*ArchiveService) CreateService() Service {
	archiveService := &ArchiveService{
		AreaIdentifier:    ARCHIVE_SERVICE_AREA_IDENTIFIER,
		ServiceIdentifier: ARCHIVE_SERVICE_SERVICE_IDENTIFIER,
		AreaNumber:        COM_AREA_NUMBER,
		ServiceNumber:     ARCHIVE_SERVICE_SERVICE_NUMBER,
		AreaVersion:       COM_AREA_VERSION,
		running:           true,
		wg:                *new(sync.WaitGroup),
	}

	return archiveService
}

//======================================================================//
//                          START: Consumer                             //
//======================================================================//

// Retrieve : TODO:
func (archiveService *ArchiveService) Retrieve(consumerURL string, providerURL string, objectType ObjectType, identifierList IdentifierList, longList LongList) (*ArchiveDetailsList, ElementList, *ServiceError, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve Consumer")

	// IN
	var providerURI = NewURI(providerURL + "/archiveServiceProvider")
	// OUT
	consumer, archiveDetailsList, elementList, errorsList, err := StartRetrieveConsumer(consumerURL,
		providerURI,
		objectType,
		identifierList,
		longList)
	if err != nil {
		return nil, nil, nil, err
	} else if errorsList != nil {
		return nil, nil, errorsList, nil
	}

	// Close the consumer
	consumer.Close()

	return archiveDetailsList, elementList, nil, nil
}

// Query : TODO:
func (archiveService *ArchiveService) Query(consumerURL string, providerURL string, boolean *Boolean, objectType ObjectType, archiveQueryList ArchiveQueryList, queryFilterList QueryFilterList) ([]interface{}, *ServiceError, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Query Consumer")

	// IN
	var providerURI = NewURI(providerURL + "/archiveServiceProvider")
	// OUT
	consumer, responses, errorsList, err := StartQueryConsumer(consumerURL,
		providerURI,
		boolean,
		objectType,
		archiveQueryList,
		queryFilterList)
	if err != nil {
		return nil, nil, err
	} else if errorsList != nil {
		return nil, errorsList, nil
	}

	// Close the consumer
	consumer.Close()

	return responses, nil, nil
}

// Count : TODO:
func (archiveService *ArchiveService) Count(consumerURL string, providerURL string, objectType *ObjectType, archiveQueryList *ArchiveQueryList, queryFilterList QueryFilterList) (*LongList, *ServiceError, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Count Consumer")

	// IN
	var providerURI = NewURI(providerURL + "/archiveServiceProvider")
	// OUT
	consumer, longList, errorsList, err := StartCountConsumer(consumerURL,
		providerURI,
		objectType,
		archiveQueryList,
		queryFilterList)
	if err != nil {
		return nil, nil, err
	} else if errorsList != nil {
		return nil, errorsList, nil
	}

	// Close the consumer
	consumer.Close()

	return longList, nil, nil
}

// Store : TODO:
func (archiveService *ArchiveService) Store(consumerURL string, providerURL string, boolean *Boolean, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*LongList, *ServiceError, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Store Consumer")

	// IN
	var providerURI = NewURI(providerURL + "/archiveServiceProvider")
	// OUT
	consumer, longList, errorsList, err := StartStoreConsumer(consumerURL,
		providerURI,
		boolean,
		objectType,
		identifierList,
		archiveDetailsList,
		elementList)
	if err != nil {
		return nil, nil, err
	} else if errorsList != nil {
		return nil, errorsList, nil
	}

	// Close the consumer
	consumer.Close()

	return longList, nil, nil
}

// Update : TODO:
func (archiveService *ArchiveService) Update(consumerURL string, providerURL string, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*ServiceError, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Update Consumer")

	// IN
	var providerURI = NewURI(providerURL + "/archiveServiceProvider")
	// OUT
	consumer, errorsList, err := StartUpdateConsumer(consumerURL,
		providerURI,
		objectType,
		identifierList,
		archiveDetailsList,
		elementList)
	if err != nil {
		return nil, err
	} else if errorsList != nil {
		return errorsList, nil
	}

	// Close the consumer
	consumer.Close()

	return nil, nil
}

// Delete : TODO:
func (archiveService *ArchiveService) Delete(consumerURL string, providerURL string, objectType ObjectType, identifierList IdentifierList, longList LongList) (*LongList, *ServiceError, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Delete Consumer")

	// IN
	var providerURI = NewURI(providerURL + "/archiveServiceProvider")
	// OUT
	consumer, respLongList, errorsList, err := StartDeleteConsumer(consumerURL,
		providerURI,
		objectType,
		identifierList,
		longList)
	if err != nil {
		return nil, nil, err
	} else if errorsList != nil {
		return nil, errorsList, nil
	}

	// Close the consumer
	consumer.Close()

	return respLongList, nil, nil
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
	var provider *Provider
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
