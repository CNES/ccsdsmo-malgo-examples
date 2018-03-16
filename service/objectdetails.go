package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ObjectDetails struct {
	related Long
	source ObjectId
}

var (
	NullObjectDetails *ObjectDetails = nil
)

const (
	MAL_OBJECT_DETAILS_TYPE_SHORT_FORM Integer = 0x04
	MAL_OBJECT_DETAILS_SHORT_FORM	   Long    = 0x1000001000004
)

func NewObjectDetails(related Long, source ObjectId) *ObjectDetails {
	objectDetails := &ObjectDetails {
		related,
		source,
	}
	return objectDetails
}

// ----- Defines COM ObjectDetails as a MAL Composite -----
func (o *ObjectDetails) Composite() Composite {
	return o
}

// ----- Defines COM ObjectDetails as a MAL Element -----
// Returns the absolute short form of the element type
func (*ObjectDetails) GetShortForm() Long {
	return MAL_OBJECT_DETAILS_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (o *ObjectDetails) GetAreaNumber() UShort {
	return o.source.GetAreaNumber()
}

// Returns the version of the area this element belongs to
func (o *ObjectDetails) GetAreaVersion() UOctet {
	return o.source.GetAreaVersion()
}

// Returns the number of the service this element belongs to
func (o *ObjectDetails) GetServiceNumber() UShort {
	return o.source.GetServiceNumber()
}

// Returns the relative short form of the element type
func (*ObjectDetails) GetTypeShortForm() Integer {
	return MAL_OBJECT_DETAILS_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (o *ObjectDetails) Encode(encoder Encoder) error {
	err := encoder.EncodeLong(&o.related)
	if err != nil {
		return err
	}
	return o.source.Encode(encoder)
}

// Decodes an instance of ObjectDetails using the supplied decoder
func (*ObjectDetails) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectDetails(decoder)
}

func DecodeObjectDetails(decoder Decoder) (*ObjectDetails, error) {
	related, err := decoder.DecodeLong()
	if err != nil {
		return nil, err
	}

	var source *ObjectId
	element, err := source.Decode(decoder)
	if err != nil {
		return nil, err
	}
	source = element.(*ObjectId)

	objectDetails := &ObjectDetails {
		*related,
		*source,
	}
}

// The methods allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism
func (*ObjectDetails) CreateElement() Element {
	return new(ObjectDetails)
}

func (o *ObjectDetails) IsNull() bool {
	return o == nil
}

func (*ObjectDetails) Null() Element {
	return NullObjectDetails
}
