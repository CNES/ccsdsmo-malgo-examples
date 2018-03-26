package data

import (
	. "github.com/EtienneLndr/MAL_API_Go_Project/archive/constants"
	. "github.com/ccsdsmo/malgo/mal"
)

// ################################################################################
// Defines COM ArchiveDetailsList type
// ################################################################################

type ArchiveDetailsList []*ArchiveDetails

var (
	NullArchiveDetailsList *ArchiveDetailsList = nil
)

const (
	COM_ARCHIVE_DETAILS_LIST_TYPE_SHORT_FORM Integer = -0x01
	COM_ARCHIVE_DETAILS_LIST_SHORT_FORM      Long    = 0x2000002FFFFFF
)

func NewArchiveDetailsList(size int) *ArchiveDetailsList {
	var list ArchiveDetailsList = ArchiveDetailsList(make([]*ArchiveDetails, size))
	return &list
}

// ================================================================================
// Defines COM ArchiveDetailsList type as an ElementList
// ================================================================================
func (list *ArchiveDetailsList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (*ArchiveDetailsList) Composite() Composite {
	return new(ArchiveDetailsList)
}

// ================================================================================
// Defines COM ArchiveDetailsList type as a MAL Element
// ================================================================================
// Registers COM ArchiveDetailsList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_ARCHIVE_DETAILS_LIST_SHORT_FORM, NullArchiveDetailsList)
}

// Returns the absolute short form of the element type.
func (*ArchiveDetailsList) GetShortForm() Long {
	return COM_ARCHIVE_DETAILS_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ArchiveDetailsList) GetAreaNumber() UShort {
	return SERVICE_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ArchiveDetailsList) GetAreaVersion() UOctet {
	return SERVICE_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ArchiveDetailsList) GetServiceNumber() UShort {
	return DEFAULT_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ArchiveDetailsList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_ARCHIVE_QUERY_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ArchiveDetailsList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*ArchiveDetails(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*ArchiveDetails(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *ArchiveDetailsList) Decode(decoder Decoder) (Element, error) {
	return DecodeArchiveDetailsList(decoder)
}

// Decodes an instance of ArchiveDetailsList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded ArchiveDetailsList instance.
func DecodeArchiveDetailsList(decoder Decoder) (*ArchiveDetailsList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ArchiveDetailsList(make([]*ArchiveDetails, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullArchiveDetails)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*ArchiveDetails)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *ArchiveDetailsList) CreateElement() Element {
	return NewArchiveDetailsList(0)
}

func (list *ArchiveDetailsList) IsNull() bool {
	return list == nil
}

func (*ArchiveDetailsList) Null() Element {
	return NullArchiveDetailsList
}
