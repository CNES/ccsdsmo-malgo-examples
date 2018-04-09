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
// Defines COM CompositeFilterSetList type
// ################################################################################

type CompositeFilterSetList []*CompositeFilterSet

var (
	NullCompositeFilterSetList *CompositeFilterSetList = nil
)

const (
	COM_COMPOSITE_FILTER_SET_LIST_TYPE_SHORT_FORM Integer = -0x04
	COM_COMPOSITE_FILTER_SET_LIST_SHORT_FORM      Long    = 0x2000002FFFFFC
)

func NewCompositeFilterSetList(size int) *CompositeFilterSetList {
	var list CompositeFilterSetList = CompositeFilterSetList(make([]*CompositeFilterSet, size))
	return &list
}

// ================================================================================
// Defines COM CompositeFilterSetList type as a QueryFilterList
// ================================================================================
func (list *CompositeFilterSetList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (*CompositeFilterSetList) Composite() Composite {
	return new(CompositeFilterSetList)
}

func (list *CompositeFilterSetList) QueryFilterList() QueryFilterList {
	return list
}

func (list *CompositeFilterSetList) GetElementAt(i int) QueryFilter {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

// ================================================================================
// Defines COM CompositeFilterSetList type as a QueryFilter
// ================================================================================
func (list *CompositeFilterSetList) QueryFilter() QueryFilter {
	return list
}

// ================================================================================
// Defines COM CompositeFilterSetList type as a MAL Element
// ================================================================================
// Registers COM CompositeFilterSetList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_COMPOSITE_FILTER_SET_LIST_SHORT_FORM, NullCompositeFilterSetList)
}

// Returns the absolute short form of the element type.
func (*CompositeFilterSetList) GetShortForm() Long {
	return COM_COMPOSITE_FILTER_SET_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (*CompositeFilterSetList) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (*CompositeFilterSetList) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (*CompositeFilterSetList) GetServiceNumber() UShort {
	return DEFAULT_SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (*CompositeFilterSetList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_COMPOSITE_FILTER_SET_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *CompositeFilterSetList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*CompositeFilterSet(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*CompositeFilterSet(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *CompositeFilterSetList) Decode(decoder Decoder) (Element, error) {
	return DecodeCompositeFilterSetList(decoder)
}

// Decodes an instance of CompositeFilterSetList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded CompositeFilterSetList instance.
func DecodeCompositeFilterSetList(decoder Decoder) (*CompositeFilterSetList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := CompositeFilterSetList(make([]*CompositeFilterSet, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullCompositeFilterSet)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*CompositeFilterSet)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *CompositeFilterSetList) CreateElement() Element {
	return NewCompositeFilterList(0)
}

func (list *CompositeFilterSetList) IsNull() bool {
	return list == nil
}

func (*CompositeFilterSetList) Null() Element {
	return NullCompositeFilterSetList
}
