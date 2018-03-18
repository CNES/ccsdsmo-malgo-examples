package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ObjectId struct {
	Type ObjectType
	Key  ObjectKey
}

var (
	NullObjectId *ObjectId = nil
)

const (
	MAL_OBJECT_ID_TYPE_SHORT_FORM Integer = 0x03
	MAL_OBJECT_ID_SHORT_FORM              = 0x1000001000003
)

func NewObjectId(t ObjectType, k ObjectKey) *ObjectId {
	var objectId = &ObjectId{
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
	return MAL_OBJECT_ID_SHORT_FORM
}

// Returns the number of the area this element type belongs to
func (o *ObjectId) GetAreaNumber() UShort {
	return o.Type.GetAreaNumber()
}

// Returns the version of the area this element belongs to
func (o *ObjectId) GetAreaVersion() UOctet {
	return o.Type.GetAreaVersion()
}

// Returns the number of the service this element belongs to
func (o *ObjectId) GetServiceNumber() UShort {
	return o.Key.GetServiceNumber()
}

// Returns the relative short form of the elemennt type
func (*ObjectId) GetTypeShortForm() Integer {
	return MAL_OBJECT_ID_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (o *ObjectId) Encode(encoder Encoder) error {
	// Encode Type (ObjectType)
	err := encoder.EncodeElement(&o.Type)
	if err != nil {
		return err
	}

	// Encode Key (ObjectKey)
	return encoder.EncodeElement(&o.Key)
}

// Decodes an instance of this eleùent type using the supplied decoder
func (o *ObjectId) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectId(decoder)
}

func DecodeObjectId(decoder Decoder) (*ObjectId, error) {
	// Decode Type (ObjectType)
	element, err := decoder.DecodeElement(NullObjectType)
	if err != nil {
		return nil, err
	}
	Type := element.(*ObjectType)

	// Decode Key (ObjectKey)
	element, err = decoder.DecodeElement(NullObjectKey)
	if err != nil {
		return nil, err
	}
	Key := element.(*ObjectKey)

	objectId := &ObjectId{
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

func (*ObjectId) Null() Element {
	return NullObjectId
}
