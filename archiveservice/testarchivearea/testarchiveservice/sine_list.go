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

// Defines SineList type

type SineList []*Sine

var NullSineList *SineList = nil

func NewSineList(size int) *SineList {
  var list SineList = SineList(make([]*Sine, size))
  return &list
}

// ================================================================================
// Defines SineList type as an ElementList

func (receiver *SineList) Size() int {
  if receiver != nil {
    return len(*receiver)
  }
  return -1
}

func (receiver *SineList) GetElementAt(i int) mal.Element {
  if receiver == nil || i >= receiver.Size() {
    return nil
  }
  return (*receiver)[i]
}

func (receiver *SineList) AppendElement(element mal.Element) {
  if receiver != nil {
    *receiver = append(*receiver, element.(*Sine))
  }
}

// ================================================================================
// Defines SineList type as a MAL Composite

func (receiver *SineList) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines SineList type as a MAL Element

const SINE_LIST_TYPE_SHORT_FORM mal.Integer = -2
const SINE_LIST_SHORT_FORM mal.Long = 0x3ea000301fffffe

// Registers SineList type for polymorphism handling
func init() {
  mal.RegisterMALElement(SINE_LIST_SHORT_FORM, NullSineList)
}

// Returns the absolute short form of the element type.
func (receiver *SineList) GetShortForm() mal.Long {
  return SINE_LIST_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *SineList) GetAreaNumber() mal.UShort {
  return testarchivearea.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *SineList) GetAreaVersion() mal.UOctet {
  return testarchivearea.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *SineList) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *SineList) GetTypeShortForm() mal.Integer {
  return SINE_LIST_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *SineList) CreateElement() mal.Element {
  return NewSineList(0)
}

func (receiver *SineList) IsNull() bool {
  return receiver == nil
}

func (receiver *SineList) Null() mal.Element {
  return NullSineList
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *SineList) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(SINE_LIST_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeUInteger(mal.NewUInteger(uint32(len([]*Sine(*receiver)))))
  if err != nil {
    return err
  }
  for _, e := range []*Sine(*receiver) {
    encoder.EncodeNullableElement(e)
  }
  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *SineList) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(SINE_LIST_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  size, err := decoder.DecodeUInteger()
  if err != nil {
    return nil, err
  }
  list := SineList(make([]*Sine, int(*size)))
  for i := 0; i < len(list); i++ {
    elem, err := decoder.DecodeNullableElement(NullSine)
    if err != nil {
      return nil, err
    }
    list[i] = elem.(*Sine)
  }
  return &list, nil
}
