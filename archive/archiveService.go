package archive

import (
	"fmt"
	"time"

	. "github.com/EtienneLndr/MAL_API_Go_Project/service"
	. "github.com/EtienneLndr/MAL_API_Go_Project/service/provider"
	. "github.com/ccsdsmo/malgo/mal"
)

type ArchiveService struct {
	areaIdentifier    Identifier
	serviceIdentifier Identifier
	areaNumber        Integer
	serviceNumber     Integer
	areaVersion       Integer
}

// Constants for the Archive Service
const (
	ARCHIVE_SERVICE_AREA_IDENTIFIER    = "COM"
	ARCHIVE_SERVICE_SERVICE_IDENTIFIER = "Archive"
	ARCHIVE_SERVICE_AREA_NUMBER        = 2
	ARCHIVE_SERVICE_SERVICE_NUMBER     = 2
	ARCHIVE_SERVICE_AREA_VERSION       = 1
)

// Constants for the operations
const (
	OPERATION_IDENTIFIER_RETRIEVE = iota + 1
	OPERATION_IDENTIFIER_QUERY
	OPERATION_IDENTIFIER_COUNT
	OPERATION_IDENTIFIER_STORE
	OPERATION_IDENTIFIER_UPDATE
	OPERATION_IDENTIFIER_DELETE
)

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
func (archiveService *ArchiveService) Retrieve(provider *Provider) error {
	// Maybe we should not have to return an error
	fmt.Println("Creation : Retrieve")

	return nil
}

/**
 * Operation        : Query
 * Operation number : 2
 */
func (archiveService *ArchiveService) Query(provider *Provider) error {
	fmt.Println("Creation : Query")

	return nil
}

/**
 * Operation        : Count
 * Operation number : 3
 */
func (archiveService *ArchiveService) Count(provider *Provider) error {
	fmt.Println("Creation : Count")

	return nil
}

/**
 * Operation        : Store
 * Operation number : 4
 */
func (archiveService *ArchiveService) Store(provider *Provider) error {
	fmt.Println("Creation : Store")

	return nil
}

/**
 * Operation        : Update
 * Operation number : 5
 */
func (archiveService *ArchiveService) Update(provider *Provider) error {
	fmt.Println("Creation : Update")

	return nil
}

/**
 * Operation        : Delete
 * Operation number : 6
 */
func (archiveService *ArchiveService) Delete(provider *Provider) error {
	fmt.Println("Creation : Delete")

	return nil
}

func (archiveService *ArchiveService) Start() error {
	provider, err := CreateProvider(providerURL)
	if err != nil {
		return err
	}
	defer provider.Close()

	// Start Operations
	go archiveService.Retrieve(provider)
	go archiveService.Query(provider)
	go archiveService.Count(provider)
	go archiveService.Store(provider)
	go archiveService.Update(provider)
	go archiveService.Delete(provider)

	// Start communication
	var running bool = true
	for running == true {
		time.Sleep(10 * time.Second)
		running = false
	}

	return nil
}
