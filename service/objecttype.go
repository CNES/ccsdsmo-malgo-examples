package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ObjectType struct {
	area    UShort
	service UShort
	version UOctet
	number  UShort
}

var (
	NullObjectType *ObjectType = nil
)

const (
	MAL_OBJECT_TYPE_TYPE_SHORT_FORM Integer = 0x01
	MAL_OBJECT_TYPE_SHORT_FORM      Long    = 0x1000001000001
)

func NewObjectType(area UShort, service UShort, version UOctet, number UShort) *ObjectType {
	objectType := &ObjectType{
		area,
		service,
		version,
		number,
	}
	return objectType
}

// ----- Defines MAL ObjectType as a MAL Composite -----
func (objectType *ObjectType) Composite() Composite {
	return objectType
}

// ----- Defines MAL ObjectType as a MAL Element -----
// Returns the absolute short form of the element type.
func (*ObjectType) GetShortForm() Long {
	return MAL_OBJECT_TYPE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (o *ObjectType) GetAreaNumber() UShort {
	return o.area
}

// Returns the version of the area this element type belongs to.
func (o *ObjectType) GetAreaVersion() UOctet {
	return o.version
}

// Returns the number of the service this element type belongs to.
func (o *ObjectType) GetServiceNumber() UShort {
	return o.service
}

// Returns the relative short form of the element type.
func (*ObjectType) GetTypeShortForm() Integer {
	return MAL_OBJECT_TYPE_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder.
func (o *ObjectType) Encode(encoder Encoder) error {
	// Encode area (UShort)
	err := encoder.EncodeElement(&o.area)
	if err != nil {
		return err
	}

	// Encode service (UShort)
	err = encoder.EncodeElement(&o.service)
	if err != nil {
		return err
	}

	// Encode version (UOctet)
	err = encoder.EncodeElement(&o.version)
	if err != nil {
		return err
	}

	// Encode number (UShort)
	return encoder.EncodeElement(&o.number)
}

// Decodes an instance of this element type using the supplied decoder.
func (*ObjectType) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectType(decoder)
}

// Decodes an instance of ObjectType using the supplied decoder
func DecodeObjectType(decoder Decoder) (*ObjectType, error) {
	// Decode area (UShort)
	element, err := decoder.DecodeElement(NullUShort)
	if err != nil {
		return nil, err
	}
	area := element.(*UShort)

	// Decode service (UShort)
	element, err = decoder.DecodeElement(NullUShort)
	if err != nil {
		return nil, err
	}
	service := element.(*UShort)

	// Decode version (UOctet)
	element, err = decoder.DecodeElement(NullOctet)
	if err != nil {
		return nil, err
	}
	version := element.(*UOctet)

	// Decode number (UShort)
	element, err = decoder.DecodeElement(NullShort)
	if err != nil {
		return nil, err
	}
	number := element.(*UShort)

	objectType := &ObjectType{
		*area,
		*service,
		*version,
		*number,
	}

	return objectType, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (*ObjectType) CreateElement() Element {
	return new(ObjectType)
}

func (o *ObjectType) IsNull() bool {
	return o == nil
}

func (*ObjectType) Null() Element {
	return NullObjectType
}
