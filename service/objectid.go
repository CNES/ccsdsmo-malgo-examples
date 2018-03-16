package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ObjectId struct {
	Type ObjectType
	Key ObjectKey
}

var (
	NullObjectId *ObjectId = nil
)

const (
	MAL_OBJECT_ID_TYPE_SHORT_FORM Integer = 0x03
	MAL_OBJECT_ID_SHORT_FORM	      = 0x1000001000003
)

func NewObjectId(t ObjectType, k ObjectKey) *ObjectId {
	var objectId = &ObjectId {
		t,
		k,
	}
	return objectId
}

// ----- Defines COM ObjectId as a MAL Composite -----
func (objectId *ObjectId) Composite() Composite {
	return objectId
}

// ----- Defines COM ObjectId as a MAL Element -----
// Returns the absolute short form of the element type
func (*ObjectId) GetShortForm() Long {
	return MAL_OBJECT_ID_TYPE_SHORT_FORM
}

// Returns the number of the area this element type belongs to
func (o *ObjectId) GetAreaNumber() UShort {
	return o.Type.area
}

// Returns the version of the area this element belongs to
func (o *ObjectId) GetAreaVersion() UOctet {
	return o.Type.version
}

// Returns the number of the service this element belongs to
func (o *ObjectId) GetServiceNumber() UShort {
	return o.Key.service
}

// Returns the relative short form of the elemennt type
func (*ObjectId) GetTypeShortForm() Integer {
	return MAL_OBJECT_ID_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (o *ObjectId) Encode(encoder Encoder) error {
	err := o.Type.Encode(encoder)
	if err != nil {
		return err
	}
	return o.Key.Encode(err)
}

// Decodes an instance of this eleùent type using the supplied decoder
func (o *ObjectId) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectId(decoder)
}

func DecodeObjectId(decoder Decoder) (*ObjectId, error) {
	// Decode Type
	var Type *ObjectType
	element, err := Type.Decode(decoder)
	if err != nil {
		return nil, err
	}
	Type = element.(*ObjectType)

	// Decode Key
	var Key *ObjectKey
	element, err = Key.Decode(decoder)
	if err != nil {
		return nil, err
	}
	Key = element.(*ObjectType)

	objectId := &ObjectId {
		*Type,
		*Key,
	}

	return objectId, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism
func (*ObjectId) CreateElement() Element {
	return new(ObjectId)
}

func (o *ObjectId) IsNull() bool {
	return o == nil
}

func (*ObjectId) NUll() Element {
	return NullObjectId
}
