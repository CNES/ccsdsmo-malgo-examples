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
package archive

import (
	"errors"
	"fmt"
	"time"

	. "github.com/EtienneLndr/archiveservice/archive/constants"
	. "github.com/EtienneLndr/archiveservice/archive/consumer"
	. "github.com/EtienneLndr/archiveservice/archive/provider"
	. "github.com/EtienneLndr/archiveservice/service"
	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/encoding/binary"
)

type ArchiveService struct {
	areaIdentifier    Identifier
	serviceIdentifier Identifier
	areaNumber        Integer
	serviceNumber     Integer
	areaVersion       Integer
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
		SERVICE_AREA_NUMBER,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		SERVICE_AREA_VERSION,
	}

	return archiveService
}

//======================================================================//
//								RETRIEVE								//
//======================================================================//
/**
 * Operation        : Retrieve
 * Operation number : 1
 */
func (archiveService *ArchiveService) retrieveProvider() (*Provider, error) {
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve Provider")

	transport := new(FixedBinaryEncoding)
	provider, err := StartRetrieveProvider(providerURL, transport)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (archiveService *ArchiveService) retrieveConsumer(objectType ObjectType, identifierList IdentifierList, longList LongList) (*InvokeConsumer, error) {
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve Consumer")

	// IN
	transport := new(FixedBinaryEncoding)
	var providerURI = NewURI(providerURLRetrieve)
	// OUT
	consumer, archiveDetailsList, elementList, err := StartRetrieveConsumer(consumerURL,
		transport,
		providerURI,
		objectType,
		identifierList,
		longList)
	if err != nil {
		return nil, err
	}

	// TODO (AF): do sthg with these objects
	fmt.Println("RetrieveConsumer received:\n\t>>>",
		consumer, "\n\t>>>",
		archiveDetailsList, "\n\t>>>",
		elementList)

	return consumer, nil
}

//======================================================================//
//								QUERY									//
//======================================================================//
/**
 * Operation        : Query
 * Operation number : 2
 */
func (archiveService *ArchiveService) queryProvider() (*Provider, error) {
	// Maybe we should not have to return an error
	fmt.Println("Creation : Query Provider")

	transport := new(FixedBinaryEncoding)
	provider, err := StartQueryProvider(providerURL, transport)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func (archiveService *ArchiveService) queryConsumer() (*ProgressConsumer, error) {
	return nil, nil
}

//======================================================================//
//								COUNT									//
//======================================================================//
/**
 * Operation        : Count
 * Operation number : 3
 */
func (archiveService *ArchiveService) countProvider() (*Provider, error) {
	return nil, nil
}

func (archiveService *ArchiveService) countConsumer() (*InvokeConsumer, error) {
	return nil, nil
}

//======================================================================//
//								STORE									//
//======================================================================//
/**
 * Operation        : Store
 * Operation number : 4
 */
func (archiveService *ArchiveService) storeProvider() (*Provider, error) {
	return nil, nil
}

func (archiveService *ArchiveService) storeConsumer() (*RequestConsumer, error) {
	return nil, nil
}

//======================================================================//
//								UPDATE									//
//======================================================================//
/**
 * Operation        : Update
 * Operation number : 5
 */
func (archiveService *ArchiveService) updateProvider() (*Provider, error) {
	return nil, nil
}

func (archiveService *ArchiveService) updateConsumer() (*SubmitConsumer, error) {
	return nil, nil
}

//======================================================================//
//								DELETE									//
//======================================================================//
/**
 * Operation        : Delete
 * Operation number : 6
 */
func (archiveService *ArchiveService) deleteProvider() (*Provider, error) {
	return nil, nil
}

func (archiveService *ArchiveService) deleteConsumer() (*RequestConsumer, error) {
	return nil, nil
}

//======================================================================//
//							START: Consumer								//
//======================================================================//
// LaunchRetrieveConsumer : TODO
func (archiveService *ArchiveService) LaunchRetrieveConsumer(objectType ObjectType, identifierList IdentifierList, longList LongList) error {
	// Start Operation
	consumer, err := archiveService.retrieveConsumer(objectType, identifierList, longList)

	if err != nil {
		return err
	}

	// Close the consumer
	consumer.Close()

	return nil
}

//======================================================================//
//							START: Provider								//
//======================================================================//
// LaunchProvider : TODO
func (archiveService *ArchiveService) LaunchProvider(operation UShort) error {
	// Declare variables
	var provider *Provider
	var err error

	// Start Operation
	switch operation {
	case OPERATION_IDENTIFIER_RETRIEVE:
		provider, err = archiveService.retrieveProvider()
		break
	case OPERATION_IDENTIFIER_QUERY:
		provider, err = archiveService.queryProvider()
		break
	case OPERATION_IDENTIFIER_COUNT:
		provider, err = archiveService.countProvider()
		break
	case OPERATION_IDENTIFIER_STORE:
		provider, err = archiveService.storeProvider()
		break
	case OPERATION_IDENTIFIER_UPDATE:
		provider, err = archiveService.updateProvider()
		break
	case OPERATION_IDENTIFIER_DELETE:
		provider, err = archiveService.deleteProvider()
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
	var running bool = true
	for running == true {
		time.Sleep(120 * time.Second)
		running = false
	}

	return nil
}
