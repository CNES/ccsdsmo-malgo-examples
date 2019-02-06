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
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal/debug"
)

var (
	logger debug.Logger = debug.GetLogger("archive.data")
)

type ValueOfSine struct {
	Value Float
}

var (
	NullValueOfSine *ValueOfSine = nil
)

const (
	COM_VALUE_OF_SINE_TYPE_SHORT_FORM Integer = 0x01
	COM_VALUE_OF_SINE_SHORT_FORM      Long    = 0x2000301000001
)

func NewValueOfSine(value Float) *ValueOfSine {
	valueOfSine := &ValueOfSine{
		value,
	}
	return valueOfSine
}

// ----- Defines COM ValueOfSine as a MAL Composite -----
func (v *ValueOfSine) Composite() Composite {
	return v
}

// ================================================================================
// Defines COM ValueOfSine type as a MAL Element
// ================================================================================
// Registers COM ValueOfSine type for polymorpsism handling
func init() {
	RegisterMALElement(COM_VALUE_OF_SINE_SHORT_FORM, NullValueOfSine)
	// In the tests the short form is also used as the COM type number
	// the fields values are those used in the test
	comObjType := com.ObjectType{
		Area:    UShort(2),
		Service: UShort(3),
		Version: UOctet(1),
		Number:  UShort(COM_VALUE_OF_SINE_TYPE_SHORT_FORM),
	}
	err := comObjType.RegisterMALBodyType(COM_VALUE_OF_SINE_SHORT_FORM)
	if err!=nil {
		logger.Errorf("ValueOfSine.init, cannot register COM object: %s", err.Error())
	}
}

// ----- Defines COM ValueOfSine as a MAL Element -----
// Returns the absolute short form of the element type
func (*ValueOfSine) GetShortForm() Long {
	return COM_VALUE_OF_SINE_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (*ValueOfSine) GetAreaNumber() UShort {
	return 2
}

// Returns the version of the area this element belongs to
func (v *ValueOfSine) GetAreaVersion() UOctet {
	return 1
}

func (*ValueOfSine) GetServiceNumber() UShort {
	return 3
}

// Returns the relative short form of the element type
func (*ValueOfSine) GetTypeShortForm() Integer {
	return COM_VALUE_OF_SINE_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (v *ValueOfSine) Encode(encoder Encoder) error {
	// Encode value
	return encoder.EncodeFloat(&v.Value)
}

// Decodes and instance of ValueOfSine using the supplied decoder
func (*ValueOfSine) Decode(decoder Decoder) (Element, error) {
	return DecodeValueOfSine(decoder)
}

func DecodeValueOfSine(decoder Decoder) (*ValueOfSine, error) {
	// Decode value
	value, err := decoder.DecodeFloat()
	if err != nil {
		return nil, err
	}

	valueOfSine := &ValueOfSine{
		*value,
	}

	return valueOfSine, nil
}

// The methods allows the creation of an element in a generic way, i.e., using     the MAL Element polymorphism
func (*ValueOfSine) CreateElement() Element {
	return new(ValueOfSine)
}

func (v *ValueOfSine) IsNull() bool {
	return v == nil
}

func (*ValueOfSine) Null() Element {
	return NullValueOfSine
}
