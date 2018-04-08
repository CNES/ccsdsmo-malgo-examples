/**
 * MIT License
 *
 * Copyright (c) 2018 CNES
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
package data

import (
	"errors"

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/etiennelndr/archiveservice/archive/constants"
)

// ################################################################################
// Defines COM ArchiveQueryList type
// ################################################################################

type ArchiveQueryList []*ArchiveQuery

var (
	NullArchiveQueryList *ArchiveQueryList = nil
)

const (
	COM_ARCHIVE_QUERY_LIST_TYPE_SHORT_FORM Integer = -0x02
	COM_ARCHIVE_QUERY_LIST_SHORT_FORM      Long    = 0x2000002FFFFFE
)

func NewArchiveQueryList(size int) *ArchiveQueryList {
	var list ArchiveQueryList = ArchiveQueryList(make([]*ArchiveQuery, size))
	return &list
}

// ================================================================================
// Defines COM ArchiveQueryList type as an ElementList
// ================================================================================
func (list *ArchiveQueryList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *ArchiveQueryList) GetElementAt(i int) (Element, error) {
	if list != nil {
		if i <= list.Size() {
			return (*list)[i], nil
		}
		return nil, errors.New("Index must not be upper or equal to the list size")
	}
	return nil, errors.New("List must not be null")
}

func (*ArchiveQueryList) Composite() Composite {
	return new(ArchiveQueryList)
}

// ================================================================================
// Defines COM ArchiveQueryList type as a MAL Element
// ================================================================================
// Registers COM ArchiveQueryList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_ARCHIVE_QUERY_LIST_SHORT_FORM, NullArchiveQueryList)
}

// Returns the absolute short form of the element type.
func (*ArchiveQueryList) GetShortForm() Long {
	return COM_ARCHIVE_QUERY_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*ArchiveQueryList) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*ArchiveQueryList) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*ArchiveQueryList) GetServiceNumber() UShort {
	return DEFAULT_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*ArchiveQueryList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_ARCHIVE_QUERY_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ArchiveQueryList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*ArchiveQuery(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*ArchiveQuery(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *ArchiveQueryList) Decode(decoder Decoder) (Element, error) {
	return DecodeArchiveQueryList(decoder)
}

// Decodes an instance of ArchiveQueryList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded ArchiveQueryList instance.
func DecodeArchiveQueryList(decoder Decoder) (*ArchiveQueryList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ArchiveQueryList(make([]*ArchiveQuery, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullArchiveQuery)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*ArchiveQuery)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *ArchiveQueryList) CreateElement() Element {
	return NewArchiveQueryList(0)
}

func (list *ArchiveQueryList) IsNull() bool {
	return list == nil
}

func (*ArchiveQueryList) Null() Element {
	return NullArchiveQueryList
}
