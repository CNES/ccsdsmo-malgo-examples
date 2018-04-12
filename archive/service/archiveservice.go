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
	"errors"
	"fmt"
	"sync"
	"time"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	_ "github.com/ccsdsmo/malgo/mal/transport/tcp"

	. "github.com/etiennelndr/archiveservice/archive/constants"
	. "github.com/etiennelndr/archiveservice/archive/consumer"
	. "github.com/etiennelndr/archiveservice/archive/provider"
	. "github.com/etiennelndr/archiveservice/data"
	. "github.com/etiennelndr/archiveservice/errors"
	. "github.com/etiennelndr/archiveservice/service"
)

type ArchiveService struct {
	AreaIdentifier    Identifier
	ServiceIdentifier Identifier
	AreaNumber        UShort
	ServiceNumber     Integer
	AreaVersion       UOctet
}

// Constants for the providers and consumers
const (
	providerURL         = "maltcp://127.0.0.1:12400"
	providerURLRetrieve = providerURL + "/providerRetrieve"
	providerURLQuery    = providerURL + "/providerQuery"
	providerURLCount    = providerURL + "/providerCount"
	providerURLStore    = providerURL + "/providerStore"
	providerURLUpdate   = providerURL + "/providerUpdate"
	providerURLDelete   = providerURL + "/providerDelete"
	consumerURL         = "maltcp://127.0.0.1:14200"
)

// CreateService : TODO
func (*ArchiveService) CreateService() Service {
	archiveService := &ArchiveService{
		ARCHIVE_SERVICE_AREA_IDENTIFIER,
		ARCHIVE_SERVICE_SERVICE_IDENTIFIER,
		COM_AREA_NUMBER,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		COM_AREA_VERSION,
	}

	return archiveService
}

//======================================================================//
//                          START: Consumer                             //
//======================================================================//
// LaunchRetrieveConsumer : TODO
func (archiveService *ArchiveService) LaunchRetrieveConsumer(objectType ObjectType, identifierList IdentifierList, longList LongList) (*ArchiveDetailsList, ElementList, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve Consumer")

	// IN
	var providerURI = NewURI(providerURLRetrieve)
	// OUT
	consumer, archiveDetailsList, elementList, err := StartRetrieveConsumer(consumerURL,
		providerURI,
		objectType,
		identifierList,
		longList)
	if err != nil {
		return nil, nil, err
	}
	// Close the consumer
	consumer.Close()

	return archiveDetailsList, elementList, nil
}

// LaunchQueryConsumer : TODO
func (archiveService *ArchiveService) LaunchQueryConsumer(boolean Boolean, objectType ObjectType, archiveQueryList ArchiveQueryList, queryFilterList QueryFilterList) ([]interface{}, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Query Consumer")

	// IN
	var providerURI = NewURI(providerURLQuery)
	// OUT
	consumer, responses, err := StartQueryConsumer(consumerURL,
		providerURI,
		boolean,
		objectType,
		archiveQueryList,
		queryFilterList)
	if err != nil {
		return nil, err
	}

	// Close the consumer
	consumer.Close()

	return responses, nil
}

// LaunchCountConsumer : TODO
func (archiveService *ArchiveService) LaunchCountConsumer(objectType ObjectType, archiveQueryList ArchiveQueryList, queryFilterList QueryFilterList) (*LongList, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Count Consumer")

	// IN
	var providerURI = NewURI(providerURLCount)
	// OUT
	consumer, longList, err := StartCountConsumer(consumerURL,
		providerURI,
		objectType,
		archiveQueryList,
		queryFilterList)
	if err != nil {
		return nil, err
	}

	// Close the consumer
	consumer.Close()

	return longList, nil
}

// LaunchStoreConsumer : TODO
func (archiveService *ArchiveService) LaunchStoreConsumer(boolean Boolean, objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*LongList, *ServiceError, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Store Consumer")

	// IN
	var providerURI = NewURI(providerURLStore)
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

// LaunchUpdateConsumer : TODO
func (archiveService *ArchiveService) LaunchUpdateConsumer(objectType ObjectType, identifierList IdentifierList, archiveDetailsList ArchiveDetailsList, elementList ElementList) (*ServiceError, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Update Consumer")

	// IN
	var providerURI = NewURI(providerURLUpdate)
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

// LaunchDeleteConsumer : TODO
func (archiveService *ArchiveService) LaunchDeleteConsumer(objectType ObjectType, identifierList IdentifierList, longList LongList) (*LongList, error) {
	// Start Operation
	// Maybe we should not have to return an error
	fmt.Println("Creation : Delete Consumer")

	// IN
	var providerURI = NewURI(providerURLDelete)
	// OUT
	consumer, respLongList, err := StartDeleteConsumer(consumerURL,
		providerURI,
		objectType,
		identifierList,
		longList)
	if err != nil {
		return nil, err
	}

	// Close the consumer
	consumer.Close()

	return respLongList, nil
}

//======================================================================//
//                          START: Provider                             //
//======================================================================//
// LaunchProvider : TODO
func (archiveService *ArchiveService) LaunchProvider(operation UShort, wg sync.WaitGroup) error {
	// Inform the WaitGroup that this goroutine is finished at the end of this function
	defer wg.Done()
	// Declare variables
	var provider *Provider
	var err error

	// Start Operation
	switch operation {
	case OPERATION_IDENTIFIER_RETRIEVE:
		fmt.Println("Creation : Retrieve Provider")
		provider, err = StartRetrieveProvider(providerURL)
		break
	case OPERATION_IDENTIFIER_QUERY:
		fmt.Println("Creation : Query Provider")
		provider, err = StartQueryProvider(providerURL)
		break
	case OPERATION_IDENTIFIER_COUNT:
		fmt.Println("Creation : Count Provider")
		provider, err = StartCountProvider(providerURL)
		break
	case OPERATION_IDENTIFIER_STORE:
		fmt.Println("Creation : Store Provider")
		provider, err = StartStoreProvider(providerURL)
		break
	case OPERATION_IDENTIFIER_UPDATE:
		fmt.Println("Creation : Update Provider")
		provider, err = StartUpdateProvider(providerURL)
		break
	case OPERATION_IDENTIFIER_DELETE:
		fmt.Println("Creation : Delete Provider")
		provider, err = StartDeleteProvider(providerURL)
		break
	default:
		return errors.New("Unknown operation")
	}

	if err != nil {
		return err
	}

	// Close the provider at the end of the function
	defer provider.Close()

	// Start communication
	var running = true
	for running == true {
		time.Sleep(10 * time.Second)
		//running = false
	}

	return nil
}
