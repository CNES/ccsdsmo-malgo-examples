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

// Defines Sine type

type Sine struct {
  T mal.Long
  Y mal.Float
}

var (
  NullSine *Sine = nil
)
func NewSine() *Sine {
  return new(Sine)
}

// ================================================================================
// Defines Sine type as a MAL Composite

func (receiver *Sine) Composite() mal.Composite {
  return receiver
}

// ================================================================================
// Defines Sine type as a MAL Element

const SINE_TYPE_SHORT_FORM mal.Integer = 2
const SINE_SHORT_FORM mal.Long = 0x3ea000301000002

// Registers Sine type for polymorphism handling
func init() {
  mal.RegisterMALElement(SINE_SHORT_FORM, NullSine)
}

// Returns the absolute short form of the element type.
func (receiver *Sine) GetShortForm() mal.Long {
  return SINE_SHORT_FORM
}

// Returns the number of the area this element type belongs to.
func (receiver *Sine) GetAreaNumber() mal.UShort {
  return testarchivearea.AREA_NUMBER
}

// Returns the version of the area this element type belongs to.
func (receiver *Sine) GetAreaVersion() mal.UOctet {
  return testarchivearea.AREA_VERSION
}

// Returns the number of the service this element type belongs to.
func (receiver *Sine) GetServiceNumber() mal.UShort {
    return SERVICE_NUMBER
}

// Returns the relative short form of the element type.
func (receiver *Sine) GetTypeShortForm() mal.Integer {
  return SINE_TYPE_SHORT_FORM
}

// Allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism.
func (receiver *Sine) CreateElement() mal.Element {
  return new(Sine)
}

func (receiver *Sine) IsNull() bool {
  return receiver == nil
}

func (receiver *Sine) Null() mal.Element {
  return NullSine
}

// Encodes this element using the supplied encoder.
// @param encoder The encoder to use, must not be null.
func (receiver *Sine) Encode(encoder mal.Encoder) error {
  specific := encoder.LookupSpecific(SINE_SHORT_FORM)
  if specific != nil {
    return specific(receiver, encoder)
  }

  err := encoder.EncodeLong(&receiver.T)
  if err != nil {
    return err
  }
  err = encoder.EncodeFloat(&receiver.Y)
  if err != nil {
    return err
  }

  return nil
}

// Decodes an instance of this element type using the supplied decoder.
// @param decoder The decoder to use, must not be null.
// @return the decoded instance, may be not the same instance as this Element.
func (receiver *Sine) Decode(decoder mal.Decoder) (mal.Element, error) {
  specific := decoder.LookupSpecific(SINE_SHORT_FORM)
  if specific != nil {
    return specific(decoder)
  }

  T, err := decoder.DecodeLong()
  if err != nil {
    return nil, err
  }
  Y, err := decoder.DecodeFloat()
  if err != nil {
    return nil, err
  }

  var composite = Sine {
    T: *T,
    Y: *Y,
  }
  return &composite, nil
}
