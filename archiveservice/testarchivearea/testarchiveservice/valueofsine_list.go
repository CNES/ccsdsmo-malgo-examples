/**
 * MIT License
 *
 * Copyright (c) 2020 CNES
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
package testarchiveservice

import (
  "github.com/CNES/ccsdsmo-malgo/mal"
  "github.com/CNES/ccsdsmo-malgo-examples/archiveservice/testarchivearea"
)

// Defines ValueOfSineList type

type ValueOfSineList []*ValueOfSine

var NullValueOfSineList *ValueOfSineList = nil

func NewValueOfSineList(size int) *ValueOfSineList {
  var list ValueOfSineList = ValueOfSineList(make([]*ValueOfSine, size))
  return &list
}

// ================================================================================
// Defines ValueOfSineList type as an ElementList

func (receiver *ValueOfSineList) Size() int {
  if receiver != nil {
    return len(*receiver)
  }
  return -1
}

func (receiver *ValueOfSineList) GetElementAt(i int) mal.Element {
  if receiver == nil || i >= receiver.Size() {
    return nil
  }
  return (*receiver)[i]
}

func (receiver *ValueOfSineList) AppendElement(element mal.Element) {
  if receiver != nil {
    *receiver = append(*receiver, element.(*ValueOfSine))
  }
}

// ================================================================================
// Defines ValueOfSineList type as a MAL Composite

func (receiver *ValueOfSineList) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ValueOfSineList type as a MAL Element

const VALUEOFSINE_LIST_TYPE_SHORT_FORM mal.Integer = -1
const VALUEOFSINE_LIST_SHORT_FORM mal.Long = 0x3ea000301ffffff

// Registers ValueOfSineList type for polymorphism handling
func init() {
  mal.RegisterMALElement(VALUEOFSINE_LIST_SHORT_FORM, NullValueOfSineList)
}

// Returns the absolute short form of the element type.
func (receiver *ValueOfSineList) GetShortForm() mal.Long {
  return VALUEOFSINE_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ValueOfSineList) GetAreaNumber() mal.UShort {
  return testarchivearea.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ValueOfSineList) GetAreaVersion() mal.UOctet {
  return testarchivearea.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ValueOfSineList) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ValueOfSineList) GetTypeShortForm() mal.Integer {
  return VALUEOFSINE_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ValueOfSineList) CreateElement() mal.Element {
  return NewValueOfSineList(0)
}

func (receiver *ValueOfSineList) IsNull() bool {
  return receiver == nil
}

func (receiver *ValueOfSineList) Null() mal.Element {
  return NullValueOfSineList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ValueOfSineList) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(VALUEOFSINE_LIST_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*ValueOfSine(*receiver)))))
  if err != nil {
    return err
  }
  for _, e := range []*ValueOfSine(*receiver) {
    encoder.EncodeNullableElement(e)
  }
  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ValueOfSineList) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(VALUEOFSINE_LIST_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  size, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }
  list := ValueOfSineList(make([]*ValueOfSine, int(*size)))
  for i := 0; i < len(list); i++ {
    elem, err := decoder.DecodeNullableElement(NullValueOfSine)
    if err != nil {
      return nil, err
    }
    list[i] = elem.(*ValueOfSine)
  }
  return &list, nil
}
