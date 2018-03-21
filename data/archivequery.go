package data

import (
	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
)

type ArchiveQuery struct {
	domain        IdentifierList
	network       Identifier
	provider      URI
	related       Long
	source        ObjectId
	startTime     FineTime
	endTime       FineTime
	sortOrder     Boolean
	sortFieldName String
}

var (
	NullArchiveQuery *ArchiveQuery = nil
)

const (
	MAL_ARCHIVE_QUERY_TYPE_SHORT_FORM Integer = 0x02
	MAL_ARCHIVE_QUERY_SHORT_FORM      Long    = 0x1000001000002
)

func NewArchiveQuery(domain IdentifierList,
	network Identifier,
	provider URI,
	related Long,
	source ObjectId,
	startTime FineTime,
	endTime FineTime,
	sortOrder Boolean,
	sortFieldName String) *ArchiveQuery {
	archiveQuery := &ArchiveQuery{
		domain,
		network,
		provider,
		related,
		source,
		startTime,
		endTime,
		sortOrder,
		sortFieldName,
	}
	return archiveQuery
}

// ----- Defines COM ArchiveQuery as a MAL Composite -----
func (a *ArchiveQuery) Composite() Composite {
	return a
}

// ----- Defines COM ArchiveQuery as a MAL Element -----
// Returns the absolute short form of the element type
func (*ArchiveQuery) GetShortForm() Long {
	return MAL_ARCHIVE_QUERY_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (a *ArchiveQuery) GetAreaNumber() UShort {
	return a.source.GetAreaNumber()
}

// Returns the version of the area this element belongs to
func (a *ArchiveQuery) GetAreaVersion() UOctet {
	return a.source.GetAreaVersion()
}

// Returns the number of the service this element belongs to
func (a *ArchiveQuery) GetServiceNumber() UShort {
	return a.source.GetServiceNumber()
}

// Returns the relative short form of the element type
func (*ArchiveQuery) GetTypeShortForm() Integer {
	return MAL_ARCHIVE_QUERY_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (a *ArchiveQuery) Encode(encoder Encoder) error {
	// Encode domain (NullableIdentifierList)
	err := a.domain.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode network (NullableIdentifier)
	err = encoder.EncodeNullableIdentifier(&a.network)
	if err != nil {
		return err
	}

	// Encode provider (NullableURI)
	err = encoder.EncodeNullableURI(&a.provider)
	if err != nil {
		return err
	}

	// Encode related (Long)
	err = encoder.EncodeLong(&a.related)
	if err != nil {
		return err
	}

	// Encode source (NullableObjectId)
	err = encoder.EncodeNullableElement(&a.source)
	if err != nil {
		return err
	}

	// Encode startTime (NullableFineTime)
	err = encoder.EncodeNullableFineTime(&a.startTime)
	if err != nil {
		return err
	}

	// Encode endTime (NullableFineTime)
	err = encoder.EncodeNullableFineTime(&a.endTime)
	if err != nil {
		return err
	}

	// Encode sortOrder (NullableBoolean)
	err = encoder.EncodeNullableBoolean(&a.sortOrder)
	if err != nil {
		return err
	}

	// Encode sortFieldName (NullableString)
	return encoder.EncodeNullableBoolean(&a.sortOrder)
}

// Decodes an instance of ObjectDetails using the supplied decoder
func (*ArchiveQuery) Decode(decoder Decoder) (Element, error) {
	return DecodeArchiveQuery(decoder)
}

func DecodeArchiveQuery(decoder Decoder) (*ArchiveQuery, error) {
	// Encode domain (NullableIdentifierList)
	element, err := decoder.DecodeNullableElement(NullIdentifierList)
	if err != nil {
		return nil, err
	}
	domain := element.(*IdentifierList)

	// Encode network (NullableIdentifier)
	element, err = decoder.DecodeNullableElement(NullIdentifier)
	if err != nil {
		return nil, err
	}
	network := element.(*Identifier)

	// Encode provider (NullableURI)
	element, err = decoder.DecodeNullableElement(NullURI)
	if err != nil {
		return nil, err
	}
	provider := element.(*URI)

	// Encode related (Long)
	element, err = decoder.DecodeElement(NullLong)
	if err != nil {
		return nil, err
	}
	related := element.(*Long)

	// Encode source (NullableObjectId)
	element, err = decoder.DecodeNullableElement(NullObjectId)
	if err != nil {
		return nil, err
	}
	source := element.(*ObjectId)

	// Encode startTime (NullableFineTime)
	element, err = decoder.DecodeNullableElement(NullFineTime)
	if err != nil {
		return nil, err
	}
	startTime := element.(*FineTime)

	// Encode endTime (NullableFineTime)
	element, err = decoder.DecodeNullableElement(NullFineTime)
	if err != nil {
		return nil, err
	}
	endTime := element.(*FineTime)

	// Encode sortOrder (NullableBoolean)
	element, err = decoder.DecodeNullableElement(NullBoolean)
	if err != nil {
		return nil, err
	}
	sortOrder := element.(*Boolean)

	// Encode sortFieldName (NullableString)
	element, err = decoder.DecodeNullableElement(NullString)
	if err != nil {
		return nil, err
	}
	sortFieldName := element.(*String)

	archiveQuery := &ArchiveQuery{
		*domain,
		*network,
		*provider,
		*related,
		*source,
		*startTime,
		*endTime,
		*sortOrder,
		*sortFieldName,
	}

	return archiveQuery, nil
}

// The methods allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism
func (*ArchiveQuery) CreateElement() Element {
	return new(ArchiveQuery)
}

func (a *ArchiveQuery) IsNull() bool {
	return a == nil
}

func (*ArchiveQuery) Null() Element {
	return NullArchiveQuery
}
