package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ObjectType struct {
	area *UShort
	service *UShort
	version *UOctet
	number *UShort
}

var (
	NullObjectType *ObjectType = nil
)

const (
	MAL_OBJECT_TYPE_TYPE_SHORT_FORM Integer = 0x01
	MAL_OBJECT_TYPE_SHORT_FORM Long = 0x1000001000001
)

func NewObjectType(area UShort, service UShort, version UOctet, number UShort) *ObjectType {
	objectType := &ObjectType{
		&area,
		&service,
		&version,
		&number,
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
	return *o.area
}

// Returns the version of the area this element type belongs to.
func (o *ObjectType) GetAreaVersion() UOctet {
	return *o.version
}

// Returns the number of the service this element type belongs to.
func (o *ObjectType) GetServiceNumber() UShort {
	return *o.service
}

// Returns the relative short form of the element type.
func (*ObjectType) GetTypeShortForm() Integer {
	return MAL_OBJECT_TYPE_TYPE_SHORT_FORM
}


// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder.
func (o *ObjectType) Encode(encoder Encoder) error {
	err := encoder.EncodeUShort(o.area)
	if err != nil {
		return err
	}

	err = encoder.EncodeUShort(o.service)
	if err != nil {
		return err
	}

	err = encoder.EncodeUOctet(o.version)
	if err != nil {
		return err
	}

	return encoder.EncodeUShort(o.number)
}

// Decodes an instance of this element type using the supplied decoder.
func (o *ObjectType) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectType(decoder)
}

// Decodes an instance of ObjectType using the supplied decoder
func DecodeObjectType(decoder Decoder) (*ObjectType, error) {
	area, err := decoder.DecodeUShort()
	if err != nil {
		return nil, err
	}

	service, err := decoder.DecodeUShort()
	if err != nil {
		return nil, err
	}

	version, err := decoder.DecodeUOctet()
	if err != nil {
		return nil, err
	}

	number, err := decoder.DecodeUShort()
	if err != nil {
		return nil, err
	}

	objectType := &ObjectType{
		area,
		service,
		version,
		number,
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
