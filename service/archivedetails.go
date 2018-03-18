package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ArchiveDetails struct {
	instId    Long
	details   ObjectDetails
	network   *Identifier
	timestamp *FineTime
	provider  *URI
}

var (
	NullArchiveDetails *ArchiveDetails = nil
)

const (
	MAL_ARCHIVE_DETAILS_TYPE_SHORT_FORM Integer = 0x01
	MAL_ARCHIVE_DETAILS_SHORT_FORM      Long    = 0x10000010000001
)

func NewArchiveDetails(instId Long, details ObjectDetails, network *Identifier, timestamp *FineTime, provider *URI) *ArchiveDetails {
	archiveDetails := &ArchiveDetails{
		instId,
		details,
		network,
		timestamp,
		provider,
	}
	return archiveDetails
}

// ----- Defines COM ArchiveDetails as a MAL Composite -----
func (a *ArchiveDetails) Composite() Composite {
	return a
}

// ----- Defines COM ArchiveDetails as a MAL Element -----
// Returns the absolute short form of the element type
func (*ArchiveDetails) GetShortForm() Long {
	return MAL_ARCHIVE_DETAILS_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (a *ArchiveDetails) GetAreaNumber() UShort {
	return a.details.GetAreaNumber()
}

// Returns the version of the area this element belongs to
func (a *ArchiveDetails) GetAreaVersion() UOctet {
	return a.details.GetAreaVersion()
}

func (a *ArchiveDetails) GetServiceNumber() UShort {
	return a.details.GetServiceNumber()
}

// Returns the relative short form of the element type
func (*ArchiveDetails) GetTypeShortForm() Integer {
	return MAL_ARCHIVE_DETAILS_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (a *ArchiveDetails) Encode(encoder Encoder) error {
	// Encode instId (Long)
	err := encoder.EncodeElement(&a.instId)
	if err != nil {
		return err
	}

	// Encode details (ObjectDetails)
	err = encoder.EncodeElement(&a.details)
	if err != nil {
		return err
	}

	// Encode network (NullableIdentifier)
	err = encoder.EncodeNullableElement(a.network)
	if err != nil {
		return err
	}

	// Encode timestamp (NullableFineTime)
	err = encoder.EncodeNullableElement(a.timestamp)
	if err != nil {
		return err
	}

	// Encode provider (NullableURI)
	return encoder.EncodeNullableElement(a.provider)
}

// Decodes and instance of ArchiveDetails using the supplied decoder
func (*ArchiveDetails) Decode(decoder Decoder) (Element, error) {
	return DecodeArchiveDetails(decoder)
}

func DecodeArchiveDetails(decoder Decoder) (*ArchiveDetails, error) {
	// Decode instId (Long)
	element, err := decoder.DecodeElement(NullLong)
	if err != nil {
		return nil, err
	}
	instId := element.(*Long)

	// Decode details (ObjectDetails)
	element, err = decoder.DecodeElement(NullObjectDetails)
	if err != nil {
		return nil, err
	}
	details := element.(*ObjectDetails)

	// Decode network (NullableIdentifier)
	element, err = decoder.DecodeNullableElement(NullIdentifier)
	if err != nil {
		return nil, err
	}
	network := element.(*Identifier)

	// Decode timestamp (NullableFineTime)
	element, err = decoder.DecodeNullableElement(NullFineTime)
	if err != nil {
		return nil, err
	}
	timestamp := element.(*FineTime)

	// Decode provider (NullableURI)
	element, err = decoder.DecodeNullableElement(NullURI)
	if err != nil {
		return nil, err
	}
	provider := element.(*URI)

	archiveDetails := &ArchiveDetails{
		*instId,
		*details,
		network,
		timestamp,
		provider,
	}

	return archiveDetails, nil
}

// The methods allows the creation of an element in a generic way, i.e., using     the MAL Element polymorphism
func (*ArchiveDetails) CreateElement() Element {
	return new(ArchiveDetails)
}

func (a *ArchiveDetails) IsNull() bool {
	return a == nil
}

func (*ArchiveDetails) Null() Element {
	return NullArchiveDetails
}
