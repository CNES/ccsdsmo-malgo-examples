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
	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/etiennelndr/archiveservice/archive/constants"
)

// ################################################################################
// Defines COM CompositeFilterList type
// ################################################################################

type CompositeFilterList []*CompositeFilter

var (
	NullCompositeFilterList *CompositeFilterList = nil
)

const (
	COM_COMPOSITE_FILTER_LIST_TYPE_SHORT_FORM Integer = -0x03
	COM_COMPOSITE_FILTER_LIST_SHORT_FORM      Long    = 0x2000002FFFFFD
)

func NewCompositeFilterList(size int) *CompositeFilterList {
	var list CompositeFilterList = CompositeFilterList(make([]*CompositeFilter, size))
	return &list
}

// ================================================================================
// Defines COM CompositeFilterList type as an ElementList
// ================================================================================
func (list *CompositeFilterList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *CompositeFilterList) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *CompositeFilterList) AppendElement(element Element) {
	if list != nil {
		*list = append(*list, element.(*CompositeFilter))
	}
}

func (*CompositeFilterList) Composite() Composite {
	return new(CompositeFilterList)
}

// ================================================================================
// Defines COM CompositeFilterList type as a MAL Element
// ================================================================================
// Registers COM CompositeFilterList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_COMPOSITE_FILTER_LIST_SHORT_FORM, NullCompositeFilterList)
}

// Returns the absolute short form of the element type.
func (*CompositeFilterList) GetShortForm() Long {
	return COM_COMPOSITE_FILTER_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*CompositeFilterList) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*CompositeFilterList) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*CompositeFilterList) GetServiceNumber() UShort {
	return DEFAULT_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*CompositeFilterList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_COMPOSITE_FILTER_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *CompositeFilterList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*CompositeFilter(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*CompositeFilter(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *CompositeFilterList) Decode(decoder Decoder) (Element, error) {
	return DecodeCompositeFilterList(decoder)
}

// Decodes an instance of CompositeFilterList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded CompositeFilterList instance.
func DecodeCompositeFilterList(decoder Decoder) (*CompositeFilterList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := CompositeFilterList(make([]*CompositeFilter, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullCompositeFilter)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*CompositeFilter)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *CompositeFilterList) CreateElement() Element {
	return NewCompositeFilterList(0)
}

func (list *CompositeFilterList) IsNull() bool {
	return list == nil
}

func (*CompositeFilterList) Null() Element {
	return NullCompositeFilterList
}
