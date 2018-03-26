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

// Constant for the provider url
const (
	providerURL = "maltcp://127.0.0.1:12400"
	consumerURL = "maltcp://127.0.0.1:15400"
)

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

/**
 * Operation        : Retrieve
 * Operation number : 1
 */
func (archiveService *ArchiveService) RetrieveProvider() (*RetrieveProvider, error) {
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve Provider")

	transport := new(FixedBinaryEncoding)
	provider, err := StartProvider(providerURL, transport)
	if err != nil {
		return nil, err
	}

	// TODO (AF): do sthg with these objects
	fmt.Println("RetrieveProvider received:\n\t>>>",
		provider)

	return provider, nil
}

func (archiveService *ArchiveService) RetrieveConsumer() (*RetrieveConsumer, error) {
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve Consumer")

	transport := new(FixedBinaryEncoding)
	// IN
	var objectType ObjectType
	var identifierList IdentifierList
	var longList LongList
	var providerURI = NewURI(providerURL + "/providerRetrieve")
	// OUT
	consumer, archiveDetailsList, elementList, err := StartConsumer(consumerURL, transport, providerURI, objectType, identifierList, longList)
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

/**
 * Operation        : Query
 * Operation number : 2
 */
func (archiveService *ArchiveService) Query() error {
	fmt.Println("Creation : Query")

	return nil
}

/**
 * Operation        : Count
 * Operation number : 3
 */
func (archiveService *ArchiveService) Count() error {
	fmt.Println("Creation : Count")

	return nil
}

/**
 * Operation        : Store
 * Operation number : 4
 */
func (archiveService *ArchiveService) Store() error {
	fmt.Println("Creation : Store")

	return nil
}

/**
 * Operation        : Update
 * Operation number : 5
 */
func (archiveService *ArchiveService) Update() error {
	fmt.Println("Creation : Update")

	return nil
}

/**
 * Operation        : Delete
 * Operation number : 6
 */
func (archiveService *ArchiveService) Delete() error {
	fmt.Println("Creation : Delete")

	return nil
}

func (archiveService *ArchiveService) StartConsumer() error {
	// Start Operations
	consumer, err := archiveService.RetrieveConsumer()
	if err != nil {
		return err
	}
	defer consumer.Close()
	/*archiveService.QueryProvider()
	archiveService.CountProvider()
	archiveService.StoreProvider()
	archiveService.UpdateProvider()
	archiveService.DeleteProvider()*/

	// Start communication
	var running bool = true
	for running == true {
		time.Sleep(10 * time.Second)
		running = false
	}

	return nil
}

func (archiveService *ArchiveService) StartProvider() error {
	// Start Operations
	provider, err := archiveService.RetrieveProvider()
	if err != nil {
		return err
	}
	defer provider.Close()
	/*archiveService.QueryConsumer()
	archiveService.CountConsumer()
	archiveService.StoreConsumer()
	archiveService.UpdateConsumer()
	archiveService.DeleteConsumer()*/

	// Start communication
	var running bool = true
	for running == true {
		time.Sleep(120 * time.Second)
		running = false
	}

	return nil
}
