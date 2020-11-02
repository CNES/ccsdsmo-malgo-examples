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

// Defines ValueOfSine type

type ValueOfSine struct {
  Value mal.Float
}

var (
  NullValueOfSine *ValueOfSine = nil
)
func NewValueOfSine() *ValueOfSine {
  return new(ValueOfSine)
}

// ================================================================================
// Defines ValueOfSine type as a MAL Composite

func (receiver *ValueOfSine) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines ValueOfSine type as a MAL Element

const VALUEOFSINE_TYPE_SHORT_FORM mal.Integer = 1
const VALUEOFSINE_SHORT_FORM mal.Long = 0x3ea000301000001

// Registers ValueOfSine type for polymorphism handling
func init() {
  mal.RegisterMALElement(VALUEOFSINE_SHORT_FORM, NullValueOfSine)
}

// Returns the absolute short form of the element type.
func (receiver *ValueOfSine) GetShortForm() mal.Long {
  return VALUEOFSINE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *ValueOfSine) GetAreaNumber() mal.UShort {
  return testarchivearea.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *ValueOfSine) GetAreaVersion() mal.UOctet {
  return testarchivearea.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *ValueOfSine) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *ValueOfSine) GetTypeShortForm() mal.Integer {
  return VALUEOFSINE_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *ValueOfSine) CreateElement() mal.Element {
  return new(ValueOfSine)
}

func (receiver *ValueOfSine) IsNull() bool {
  return receiver == nil
}

func (receiver *ValueOfSine) Null() mal.Element {
  return NullValueOfSine
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *ValueOfSine) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(VALUEOFSINE_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeFloat(&receiver.Value)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *ValueOfSine) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(VALUEOFSINE_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  Value, err := decoder.DecodeFloat()
  if err != nil {
    return nil, err
  }

  var composite = ValueOfSine {
    Value: *Value,
  }
  return &composite, nil
}
