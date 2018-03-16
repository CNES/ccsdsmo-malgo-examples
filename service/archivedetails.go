package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ArchiveDetails struct {
	instId Long
	details ObjectDetails
	network *Identifier
	timestamp *FineTime
	provider *URI
}

var (
	NullArchiveDetails *ArchiveDetails = nil
)

const (
        MAL_ARCHIVE_DETAILS_TYPE_SHORT_FORM Integer = 0x01
	MAL_ARCHIVE_DETAILS_SHORT_FORM      Long    = 0x10000010000001

)

func NewArchiveDetails(instId Long, details ObjectDetails, network *Identifier, timestamp *FineTime, provider *URI) *ArchiveDetails {
	archiveDetails := &ArchiveDetails {
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
func (*ArchiveDetails) GetAreaNumber() UShort {
	return o.details.GetAreaNumber()
}

// Returns the version of the area this element belongs to
func (*ArchiveDetails) GetAreaVersion() UOctet {
	return o.details.GetAreaVersion()
}

func (*ArchiveDetails) GetServiceNumber() UShort {
	return o.details.GetServiceNumber()
}

// Returns the relative short form of the element type
func (*ArchiveDetails) GetTypeShortForm() Integer {
	return MAL_ARCHIVE_DETAILS_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (a *ArchiveDetails) Encode(encoder Encoder) error {
	// Encode instId
	err := encoder.EncodeLong(&a.instId)
	if err != nil {
		return err
	}

	// Encode details
	err = a.details.Encode(encoder)
	if err != nil {
		return err
	}

	// Encode network
	err = encoder.EncodeNullableIdentifier(a.network)
	if err != nil {
		return err
	}

	// Encode timestamp
	err = encoder.EncodeNullableFineTime(a.timestamp)
	if err != nil {
		return err
	}

	// Encode provider
	return encoder.EncodeNullableURI(a.provider)
}

// Decodes and instance of ArchiveDetails using the supplied decoder
func (*ArchiveDetails) Decode(decoder Decoder) (Element, error) {
	return DecodeArchiveDetails(decoder)
}

func DecodeArchiveDetails(decoder Decoder) (*ArchiveDetails, error) {
	// Decode instId
	instId, err := decoder.DecodeLong()
	if err != nil {
		return nil, err
	}

	// Decode details
	var details *ObjectDetails
	element, err := details.Decode(decoder)
	if err != nil {
		return nil, err
	}
	details = element.(*ObjectDetails)

	// Decode network
	network, err := decoder.DecodeNullableIdentifier()
	if err != nil {
		return nil, err
	}

	// Decode timestamp
	timestamp, err := decoder.DecodeNullableFineTime()
	if err != nil {
		return nil, err
	}

	// Decode provider
	provider, err := decoder.DecodeNullableURI()
	if err != nil {
		return nil, err
	}

	archiveDetails := &ArchiveDetails {
		*instId,
		*details,
		network,
		timestamp,
		provider,
	}

	return archiveDetails
}

// The methods allows the creation of an element in a generic way, i.e., using     the MAL Element polymorphism
func (*ArchiveDetails) CreateElement() Element {
	return new(ArchiveDetails)
}

func (a *ArchiveDetails) IsNull() bool {
	return a == nil
}

func (*ArchiveDetails) NUll() Element {
	return NullArchiveDetails
}
