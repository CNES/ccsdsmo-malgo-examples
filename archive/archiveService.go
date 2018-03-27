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
	"fmt"
	"time"

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/consumer"
	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/provider"
	. "github.com/EtienneLndr/MAL_API_Go_Project/service"
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

func (archiveService *ArchiveService) retrieveConsumer(objectType ObjectType, identifierList IdentifierList, elementList ElementList) (*Consumer, error) {
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
		elementList)
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
func (archiveService *ArchiveService) queryProvider() error {
	return nil
}

func (archiveService *ArchiveService) queryConsumer() error {
	return nil
}

//======================================================================//
//								COUNT									//
//======================================================================//
/**
 * Operation        : Count
 * Operation number : 3
 */
func (archiveService *ArchiveService) countProvider() error {
	return nil
}

func (archiveService *ArchiveService) countConsumer() error {
	return nil
}

//======================================================================//
//								STORE									//
//======================================================================//
/**
 * Operation        : Store
 * Operation number : 4
 */
func (archiveService *ArchiveService) storeProvider() error {
	return nil
}

func (archiveService *ArchiveService) storeConsumer() error {
	return nil
}

//======================================================================//
//								UPDATE									//
//======================================================================//
/**
 * Operation        : Update
 * Operation number : 5
 */
func (archiveService *ArchiveService) updateProvider() error {
	return nil
}

func (archiveService *ArchiveService) updateConsumer() error {
	return nil
}

//======================================================================//
//								DELETE									//
//======================================================================//
/**
 * Operation        : Delete
 * Operation number : 6
 */
func (archiveService *ArchiveService) deleteProvider() error {
	return nil
}

func (archiveService *ArchiveService) deleteConsumer() error {
	return nil
}

// StartConsumer : TODO
func (archiveService *ArchiveService) StartConsumer(objectType ObjectType, identifierList IdentifierList, elementList ElementList) error {
	// Start Operations
	consumer, err := archiveService.retrieveConsumer(objectType, identifierList, elementList)
	if err != nil {
		return err
	}

	// Close the consumer
	consumer.Close()

	return nil
}

// StartProvider : TODO
func (archiveService *ArchiveService) StartProvider() error {
	// Start Operations
	provider, err := archiveService.retrieveProvider()
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
