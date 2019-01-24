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
	. "github.com/CNES/ccsdsmo-malgo/mal"
)

type SineList []*Sine

var (
	NullSineList *SineList = nil
)

const (
	COM_SINE_LIST_TYPE_SHORT_FORM Integer = -0x02
	COM_SINE_LIST_SHORT_FORM      Long    = 0x2000301FFFFFE
)

func NewSineList(size int) *SineList {
	var list SineList = SineList(make([]*Sine, size))
	return &list
}

// ================================================================================
// Defines COM SineList type as an ElementList
// ================================================================================
func (list *SineList) Size() int {
	if list != nil {
		return len(*list)
	}
	return -1
}

func (list *SineList) GetElementAt(i int) Element {
	if list != nil {
		if i < list.Size() {
			return (*list)[i]
		}
		return nil
	}
	return nil
}

func (list *SineList) AppendElement(element Element) {
	if list != nil {
		*list = append(*list, element.(*Sine))
	}
}

func (*SineList) Composite() Composite {
	return new(SineList)
}

// ================================================================================
// Defines COM SineList type as a MAL Element
// ================================================================================
// Registers COM SineList type for polymorpsism handling
func init() {
	RegisterMALElement(COM_SINE_LIST_SHORT_FORM, NullSineList)
}

// Returns the absolute short form of the element type.
func (*SineList) GetShortForm() Long {
	return COM_SINE_LIST_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (*SineList) GetAreaNumber() UShort {
	return 2
}

// Returns the version of the area this element belongs to
func (v *SineList) GetAreaVersion() UOctet {
	return 1
}

func (*SineList) GetServiceNumber() UShort {
	return 3
}

// Returns the relative short form of the element type.
func (*SineList) GetTypeShortForm() Integer {
	//	return MAL_ENTITY_REQUEST_TYPE_SHORT_FORM & 0x01FFFF00
	return COM_SINE_LIST_TYPE_SHORT_FORM
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (list *SineList) Encode(encoder Encoder) error {
	err := encoder.EncodeUInteger(NewUInteger(uint32(len([]*Sine(*list)))))
	if err != nil {
		return err
	}
	for _, e := range []*Sine(*list) {
		encoder.EncodeNullableElement(e)
	}
	return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (list *SineList) Decode(decoder Decoder) (Element, error) {
	return DecodeSineList(decoder)
}

// Decodes an instance of SineList using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded SineList instance.
func DecodeSineList(decoder Decoder) (*SineList, error) {
	size, err := decoder.DecodeUInteger()
	if err != nil {
		return nil, err
	}
	list := SineList(make([]*Sine, int(*size)))
	for i := 0; i < len(list); i++ {
		element, err := decoder.DecodeNullableElement(NullSine)
		if err != nil {
			return nil, err
		}
		list[i] = element.(*Sine)
	}
	return &list, nil
}

// The method allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (list *SineList) CreateElement() Element {
	return NewSineList(0)
}

func (list *SineList) IsNull() bool {
	return list == nil
}

func (*SineList) Null() Element {
	return NullSineList
}
