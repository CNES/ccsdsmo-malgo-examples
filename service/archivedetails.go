package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ArchiveDetails struct {
	domain IdentifierList
	instId Long
}

var (
	NullArchiveDetails *ArchiveDetails = nil
)

const (
	MAL_ARCHIVE_DETAILS_TYPE_SHORT_FORM Integer = 0x02
	MAL_ARCHIVE_DETAILS_SHORT_FORM      Long    = 0x1000001000002
)

func NewArchiveDetails(domain IdentifierList, instId Long) *ArchiveDetails {
	var archiveDetails = &ArchiveDetails{
		domain,
		instId,
	}
	return archiveDetails
}

// ----- Defines MAL ArchiveDetails as a MAL Composite -----
func (archiveDetails *ArchiveDetails) Composite() Composite {
	return archiveDetails
}

// ----- Defines MAL ArchiveDetails as a MAL Element -----
// Returns the absolute short form of the element type.
func (*ArchiveDetails) GetShortForm() Long {
	return MAL_ARCHIVE_DETAILS_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ArchiveDetails) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ArchiveDetails) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ArchiveDetails) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ArchiveDetails) GetTypeShortForm() Integer {
	return MAL_ARCHIVE_DETAILS_TYPE_SHORT_FORM
}

func (a *ArchiveDetails) Encode(encoder Encoder) error {
	err := a.domain.Encode(encoder)
	if err != nil {
		return err
	}
	return encoder.EncodeLong(&a.instId)
}

// Decodes an instance of this element type using the supplied decoder.
func (o *ArchiveDetails) Decode(decoder Decoder) (Element, error) {
	return DecodeArchiveDetails(decoder)
}

func DecodeArchiveDetails(decoder Decoder) (*ArchiveDetails, error) {
	var domain *IdentifierList
	element, err := decoder.DecodeElement(domain)
	if err != nil {
		return nil, err
	}
	domain = element.(*IdentifierList)

	instId, err := decoder.DecodeLong()
	if err != nil {
		return nil, err
	}

	archiveDetails := &ArchiveDetails{
		*domain,
		*instId,
	}
	return archiveDetails, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (request *ArchiveDetails) CreateElement() Element {
	// TODO (AF):
	//	return new(EntityRequest)
	return NewEntityRequest()
}

func (o *ArchiveDetails) IsNull() bool {
	return o == nil
}

func (*ArchiveDetails) Null() Element {
	return NullArchiveDetails
}
