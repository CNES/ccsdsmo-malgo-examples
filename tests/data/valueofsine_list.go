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
	. "github.com/ccsdsmo/malgo/mal"
)

type ValueOfSineList []*ValueOfSine

var (
	NullValueOfSineList *ValueOfSineList = nil
)

const (
	COM_VALUE_OF_SINE_LIST_TYPE_SHORT_FORM Integer = -0x01
	COM_VALUE_OF_SINE_LIST_SHORT_FORM      Long    = 0x2000301FFFFFF
)

func NewValueOfSineList(size int) *ValueOfSineList {
	var list ValueOfSineList = ValueOfSineList(make([]*ValueOfSine, size))
	return &list
}

// ================================================================================
// Defines COM ValueOfSineList type as an ElementList
// ================================================================================
func (list *ValueOfSineList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *ValueOfSineList) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *ValueOfSineList) AppendElement(element Element) {
	if list != nil {
		*list = append(*list, element.(*ValueOfSine))
	}
}

func (*ValueOfSineList) Composite() Composite {
	return new(ValueOfSineList)
}

// ================================================================================
// Defines COM ValueOfSineList type as a MAL Element
// ================================================================================
// Registers COM ValueOfSineList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_VALUE_OF_SINE_LIST_SHORT_FORM, NullValueOfSineList)
}

// Returns the absolute short form of the element type.
func (*ValueOfSineList) GetShortForm() Long {
	return COM_VALUE_OF_SINE_LIST_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (*ValueOfSineList) GetAreaNumber() UShort {
	return 2
}

// Returns the version of the area this element belongs to
func (v *ValueOfSineList) GetAreaVersion() UOctet {
	return 1
}

func (*ValueOfSineList) GetServiceNumber() UShort {
	return 3
}

// Returns the relative short form of the element type.
func (*ValueOfSineList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_VALUE_OF_SINE_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *ValueOfSineList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*ValueOfSine(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*ValueOfSine(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *ValueOfSineList) Decode(decoder Decoder) (Element, error) {
	return DecodeValueOfSineList(decoder)
}

// Decodes an instance of ValueOfSineList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded ValueOfSineList instance.
func DecodeValueOfSineList(decoder Decoder) (*ValueOfSineList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := ValueOfSineList(make([]*ValueOfSine, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullValueOfSine)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*ValueOfSine)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *ValueOfSineList) CreateElement() Element {
	return NewValueOfSineList(0)
}

func (list *ValueOfSineList) IsNull() bool {
	return list == nil
}

func (*ValueOfSineList) Null() Element {
	return NullValueOfSineList
}
