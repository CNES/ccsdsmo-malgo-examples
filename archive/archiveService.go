package archive

import (
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
 * Operation		: Retrieve
 * Operation number : 1
 */
func (*ArchiveService) retrieve() error {
	return nil
}

/**
 * Operation 		: Query
 * Operation number : 2
 */
func (*ArchiveService) query() error {
	return nil
}

/**
 * Operation 		: Count
 * Operation number : 3
 */
func (*ArchiveService) count() error {
	return nil
}

/**
 * Operation 		: Store
 * Operation number : 4
 */
func (*ArchiveService) store() error {
	return nil
}

/**
 * Operation 		: Update
 * Operation number : 5
 */
func (*ArchiveService) update() error {
	return nil
}

/**
 * Operation 		: Delete
 * Operation number : 6
 */
func (*ArchiveService) delete() error {
	return nil
}
