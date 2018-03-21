package archive

import (
	"fmt"
	"time"

	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/consumer"
	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/provider"
	. "github.com/EtienneLndr/MAL_API_Go_Project/service"
	. "github.com/ccsdsmo/malgo/mal"
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
)

func (*ArchiveService) CreateService() Service {
	archiveService := &ArchiveService{
		ARCHIVE_SERVICE_AREA_IDENTIFIER,
		ARCHIVE_SERVICE_SERVICE_IDENTIFIER,
		ARCHIVE_SERVICE_AREA_NUMBER,
		ARCHIVE_SERVICE_SERVICE_NUMBER,
		ARCHIVE_SERVICE_AREA_VERSION,
	}

	return archiveService
}

/**
 * Operation        : Retrieve
 * Operation number : 1
 */
func (archiveService *ArchiveService) Retrieve() error {
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve")

	provider, err := CreateRetrieveProvider(providerURL)
	if err != nil {
		return err
	}
	defer provider.Close()

	return nil
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

func (archiveService *ArchiveService) Start() error {
	// Start Operations
	go archiveService.Retrieve()
	go archiveService.Query()
	go archiveService.Count()
	go archiveService.Store()
	go archiveService.Update()
	go archiveService.Delete()

	// Start communication
	var running bool = true
	for running == true {
		time.Sleep(10 * time.Second)
		running = false
	}

	return nil
}
