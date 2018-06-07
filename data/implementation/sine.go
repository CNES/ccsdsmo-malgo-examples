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

// Sine :
type Sine struct {
	T Long
	Y Float
}

var (
	NullSine *Sine = nil
)

const (
	COM_SINE_TYPE_SHORT_FORM Integer = 0x01
	COM_SINE_SHORT_FORM      Long    = 0x2000301000001
)

func NewSine(t Long, y Float) *Sine {
	valueOfSine := &Sine{
		T: t,
		Y: y,
	}
	return valueOfSine
}

// ----- Defines COM Sine as a MAL Composite -----
func (v *Sine) Composite() Composite {
	return v
}

// ================================================================================
// Defines COM Sine type as a MAL Element
// ================================================================================
// Registers COM Sine type for polymorpsism handling
func init() {
	RegisterMALElement(COM_SINE_SHORT_FORM, NullSine)
}

// ----- Defines COM Sine as a MAL Element -----
// Returns the absolute short form of the element type
func (*Sine) GetShortForm() Long {
	return COM_SINE_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (*Sine) GetAreaNumber() UShort {
	return 2
}

// Returns the version of the area this element belongs to
func (v *Sine) GetAreaVersion() UOctet {
	return 1
}

func (*Sine) GetServiceNumber() UShort {
	return 3
}

// Returns the relative short form of the element type
func (*Sine) GetTypeShortForm() Integer {
	return COM_SINE_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (v *Sine) Encode(encoder Encoder) error {
	// Encode t
	err := encoder.EncodeLong(&v.T)
	if err != nil {
		return err
	}

	// Encode y
	return encoder.EncodeFloat(&v.Y)
}

// Decodes and instance of Sine using the supplied decoder
func (*Sine) Decode(decoder Decoder) (Element, error) {
	return DecodeSine(decoder)
}

func DecodeSine(decoder Decoder) (*Sine, error) {
	// Decode t
	t, err := decoder.DecodeLong()
	if err != nil {
		return nil, err
	}

	// Decode y
	y, err := decoder.DecodeFloat()
	if err != nil {
		return nil, err
	}

	valueOfSine := &Sine{
		T: *t,
		Y: *y,
	}

	return valueOfSine, nil
}

// The methods allows the creation of an element in a generic way, i.e., using     the MAL Element polymorphism
func (*Sine) CreateElement() Element {
	return new(Sine)
}

func (v *Sine) IsNull() bool {
	return v == nil
}

func (*Sine) Null() Element {
	return NullSine
}
