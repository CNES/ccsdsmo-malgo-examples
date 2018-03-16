package service

import (
	. "github.com/ccsdsmo/malgo/mal"
)

type ObjectKey struct {
	domain IdentifierList
	instId Long
}

var (
	NullObjectKey *ObjectKey = nil
)

const (
	MAL_OBJECT_KEY_TYPE_SHORT_FORM Integer = 0x02
	MAL_OBJECT_KEY_SHORT_FORM      Long    = 0x1000001000002
)

func NewObjectKey(domain IdentifierList, instId Long) *ObjectKey {
	var objectKey = &ObjectKey{
		domain,
		instId,
	}
	return objectKey
}

// ----- Defines COM ObjectKey as a MAL Composite -----
func (objectKey *ObjectKey) Composite() Composite {
	return objectKey
}

// ----- Defines COM ObjectKey as a MAL Element -----
// Returns the absolute short form of the element type.
func (*ObjectKey) GetShortForm() Long {
	return MAL_OBJECT_KEY_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ObjectKey) GetAreaNumber() UShort {
	return MAL_ATTRIBUTE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ObjectKey) GetAreaVersion() UOctet {
	return MAL_ATTRIBUTE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ObjectKey) GetServiceNumber() UShort {
	return MAL_ATTRIBUTE_AREA_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ObjectKey) GetTypeShortForm() Integer {
	return MAL_OBJECT_KEY_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder.
func (o *ObjectKey) Encode(encoder Encoder) error {
	err := o.domain.Encode(encoder)
	if err != nil {
		return err
	}
	return encoder.EncodeLong(&o.instId)
}

// Decodes an instance of this element type using the supplied decoder.
func (*ObjectKey) Decode(decoder Decoder) (Element, error) {
	return DecodeObjectKey(decoder)
}

func DecodeObjectKey(decoder Decoder) (*ObjectKey, error) {
	var domain *IdentifierList
	element, err := domain.Decode(decoder)
	if err != nil {
		return nil, err
	}
	domain = element.(*IdentifierList)

	instId, err := decoder.DecodeLong()
	if err != nil {
		return nil, err
	}

	objectKey := &ObjectKey{
		*domain,
		*instId,
	}
	return objectKey, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (*ObjectKey) CreateElement() Element {
	return new(ObjectKey)
}

func (o *ObjectKey) IsNull() bool {
	return o == nil
}

func (*ObjectKey) Null() Element {
	return NullObjectKey
}
