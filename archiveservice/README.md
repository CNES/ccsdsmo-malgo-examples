MAL GO ARCHIVE SERVICE
======================


Introduction
============

This service is an implementation of the **Archive Service** described in the CCSDS Recommendation **Mission Operations Common Object Model CCSDS 521.1-B-1**. It uses the High level **MAL GO API** available [here](https://github.com/CNES/ccsdsmo-malgo/) and the **MAL GO generator** available [here](https://github.com/CNES/ccsdsmo-malgo-stubgenerator/).


Download this repository
========================

```
go get github.com/juju/loggo
go get github.com/CNES/ccsdsmo-malgo
go get github.com/CNES/ccsdsmo-malgo-examples/archiveservice
```


Using the Archive Service
=========================

If you want to see the Archive Service in action, click [here](https://github.com/CNES/ccsdsmo-malgo-examples/archiveservice_implementation/) and run this project.

Use of the provider
-------------------

```go
// Variable that defines the ArchiveService
var archiveService *ArchiveService
// Variable to retrieve the error
var err error
// Create the Archive Service
archiveService = archiveService.CreateService().(*ArchiveService)

// Start the providers
err = archiveService.StartProvider("maltcp://127.0.0.1:12400")

if err != nil {
    fmt.Println("Error:", err)
}
```

Use of the consumer
-------------------

### Retrieve

The **retrieve operation** retrieves a set of objects identified by their object instance identifier. In our service it can be used in that way:

```go
// Variable that defines the ArchiveService
var archiveService *ArchiveService
// Create the Archive Service
service := archiveService.CreateService()
archiveService = service.(*ArchiveService)
// Variable that defines the ArchiveService
var objectType = ObjectType{
	Area:    UShort(2),
	Service: UShort(3),
	Version: UOctet(1),
	Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
}
var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
var longList = LongList([]*Long{NewLong(0)})

// Variables to retrieve the return of this function
var archiveDetailsList *ArchiveDetailsList
var elementList ElementList
var errorsList *ServiceError
var err error

// Start the consumer
archiveDetailsList, elementList, errorsList, err = archiveService.Retrieve(consumerURL, providerURL, objectType, identifierList, longList)

// Check errors
if err != nil {
    // Do Something
} else if errorsList != nil {
    // Do something else
}
```

### Query

The **query operation** retrieves a set of object instance identifiers, and optionally the object
bodies, from a list of supplied queries. The **PROGRESS interaction** pattern is used as the
returned set of data may be quite large and this allows it to be split over several MAL
messages. A simple way to use this operation:

```go
// Variable that defines the ArchiveService
var archiveService *ArchiveService
// Create the Archive Service
service := archiveService.CreateService()
archiveService = service.(*ArchiveService)

// Create parameters
var boolean = NewBoolean(true)
var objectType = ObjectType{
	Area:    UShort(2),
	Service: UShort(3),
	Version: UOctet(1),
	Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
}
archiveQueryList := NewArchiveQueryList(0)
//var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
archiveQuery := &ArchiveQuery{
	Related:   Long(0),
	SortOrder: NewBoolean(true),
}
archiveQueryList.AppendElement(archiveQuery)
var queryFilterList *CompositeFilterSetList

// Variables to retrieve the return of this function
var responses []interface{}
var errorsList *ServiceError
var err error
// Start the consumer
responses, errorsList, err = archiveService.Query(consumerURL, providerURL, boolean, objectType, *archiveQueryList, queryFilterList)

// Check errors
if err != nil {
    // Do Something
} else if errorsList != nil {
    // Do something else
}
```

### Count

The **count operation** counts the set of objects based on a supplied query.

```go
// Variable that defines the ArchiveService
var archiveService *ArchiveService
// Create the Archive Service
service := archiveService.CreateService()
archiveService = service.(*ArchiveService)

var objectType = &ObjectType{
	Area:    UShort(2),
	Service: UShort(3),
	Version: UOctet(1),
	Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
}
archiveQueryList := NewArchiveQueryList(0)
var domain = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
archiveQuery := &ArchiveQuery{
	Domain:    &domain,
	Related:   Long(0),
	SortOrder: NewBoolean(true),
}
archiveQueryList.AppendElement(archiveQuery)
var queryFilterList *CompositeFilterSetList

// Variable to retrieve the return of this function
var longList *LongList
var errorsList *ServiceError
var err error
// Start the consumer
longList, errorsList, err = archiveService.Count(consumerURL, providerURL, objectType, archiveQueryList, queryFilterList)

// Check errors
if err != nil {
    // Do Something
} else if errorsList != nil {
    // Do something else
}
```

### Store

The **store operation** stores new objects in the archive and causes an ObjectStored event to be
published by the archive.

When new objects are being stored in an archive by a service provider the archive service
provider is **capable of allocating an unused object instance identifier** for the objects being
stored. The returned object instance identifier should be used by the service provider for
identifying the object instances to its consumer to ensure that **only a single object instance
identifier** is used for each object instance.

In our case, this operation can be used in that way:

```go
// Variable that defines the ArchiveService
var archiveService *ArchiveService
// Create the Archive Service
service := archiveService.CreateService()
archiveService = service.(*ArchiveService)

// Start the store consumer
// Create parameters
// Object that's going to be stored in the archive
var elementList = NewValueOfSineList(1)
(*elementList)[0] = NewValueOfSine(0)
var boolean = NewBoolean(true)
var objectType = ObjectType{
	Area:    UShort(2),
	Service: UShort(3),
	Version: UOctet(1),
	Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
}
var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice")})
// Object instance identifier
var objectInstanceIdentifier = *NewLong(81)
// Variables for ArchiveDetailsList
var objectKey = ObjectKey{
	Domain: identifierList,
	InstId: objectInstanceIdentifier,
}
var objectID = ObjectId{
	Type: &objectType,
	Key:  &objectKey,
}
var objectDetails = ObjectDetails{
	Related: NewLong(1),
	Source:  &objectID,
}
var network = NewIdentifier("network")
var timestamp = NewFineTime(time.Now())
var provider = NewURI("main/start")
var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, timestamp, provider)})

// Variable to retrieve the return of this function
var longList *LongList
var err error
var errorsList *ServiceError
// Start the consumer
longList, errorsList, err = archiveService.Store(consumerURL, providerURL, boolean, objectType, identifierList, archiveDetailsList, elementList)

// Check errors
if err != nil {
    // Do Something
} else if errorsList != nil {
    // Do something else
}
```

### Update

The **update operation** updates an object (or set of objects) and causes an ObjectUpdated event
to be published by the archive.

```go
// Variable that defines the ArchiveService
var archiveService *ArchiveService
// Create the Archive Service
service := archiveService.CreateService()
archiveService = service.(*ArchiveService)

// Start the update consumer
// Create parameters
// ---- ELEMENTLIST ----
// Object that's going to be updated in the archive
var elementList = NewValueOfSineList(1)
(*elementList)[0] = NewValueOfSine(0.5)
// ---- OBJECTTYPE ----
var objectType = ObjectType{
	Area:    UShort(2),
	Service: UShort(3),
	Version: UOctet(1),
	Number:  UShort((*elementList)[0].GetTypeShortForm()),
}
// ---- IDENTIFIERLIST ----
var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
// Object instance identifier
var objectInstanceIdentifier = *NewLong(1)
// Variables for ArchiveDetailsList
// ---- ARCHIVEDETAILSLIST ----
var objectKey = ObjectKey{
	Domain: identifierList,
	InstId: objectInstanceIdentifier,
}
var objectID = ObjectId{
	Type: &objectType,
	Key:  &objectKey,
}
var objectDetails = ObjectDetails{
	Related: NewLong(1),
	Source:  &objectID,
}
var network = NewIdentifier("new.network")
var fineTime = NewFineTime(time.Now())
var uri = NewURI("main/start")
var archiveDetailsList = ArchiveDetailsList([]*ArchiveDetails{NewArchiveDetails(objectInstanceIdentifier, objectDetails, network, fineTime, uri)})
    
// Variable to retrieve the return of this function
var errorsList *ServiceError
var err error
// Start the consumer
errorsList, err = archiveService.Update(consumerURL, providerURL, objectType, identifierList, archiveDetailsList, elementList)

// Check errors
if err != nil {
    // Do Something
} else if errorsList != nil {
    // Do something else
}
```

### Delete

The **delete operation** deletes an object (or set of objects) and causes an ObjectDeleted event
to be published by the archive. A simple way to use this operation:

```go
// Variable that defines the ArchiveService
var archiveService *ArchiveService
// Create the Archive Service
service := archiveService.CreateService()
archiveService = service.(*ArchiveService)

var objectType = ObjectType{
	Area:    UShort(2),
	Service: UShort(3),
	Version: UOctet(1),
	Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
}
var identifierList = IdentifierList([]*Identifier{NewIdentifier("fr"), NewIdentifier("cnes"), NewIdentifier("archiveservice"), NewIdentifier("test")})
var longList = NewLongList(0)
longList.AppendElement(NewLong(15))

// Variable to retrieve the return of this function
var respLongList *LongList
var errorsList *ServiceError
var err error
// Start the consumer
respLongList, errorsList, err = archiveService.Delete(consumerURL, providerURL, objectType, identifierList, *longL

// Check errors
if err != nil {
    // Do Something
} else if errorsList != nil {
    // Do something else
}
```

Implementation details
======================

This Archive service has originally been implemented using the MAL/Go API. It has then be refactored to use the new go generator available [here](https://github.com/CNES/ccsdsmo-malgo-stubgenerator/). This new implementation can then be used as an example for using the go generator.

The archive/storage directory holds the core implementation of the service. It has been left quite unchanged.
The archive/provider directory has been considerably simplified by using the provider generated stubs. It implements the MAL standard API of the service and makes use of the storage part.
The archive/consumer directory has been removed and completely replaced by the consumer generated stubs.
The archive/service directory defines a higher level API, which is a choice of the developer of the service. This higher level API is no rule for any other service.

The testarchivearea directory replaces the former data directory. It holds the definition of COM data types used in the tests. As those data types should be defined as MAL types, they have been declared in a standard XML file TestArchive.xml. The testarchivearea holds the result of the generator processing this area definition. An additional file named valueofsine_impl.go has been created to hold code that could not be generated.

Examples for using the go generated code can then be found in the provider directory for the provider side, and in the service directory for the consumer side.
