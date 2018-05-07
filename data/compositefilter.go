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

type CompositeFilter struct {
	FieldName  String
	Type       ExpressionOperator
	FieldValue Attribute
}

var (
	NullCompositeFilter *CompositeFilter = nil
)

const (
	COM_COMPOSITE_FILTER_TYPE_SHORT_FORM Integer = 0x03
	COM_COMPOSITE_FILTER_SHORT_FORM      Long    = 0x2000201000003
)

func NewCompositeFilter(fieldName String, _type ExpressionOperator, fieldValue Attribute) *CompositeFilter {
	compositeFilter := &CompositeFilter{
		fieldName,
		_type,
		fieldValue,
	}
	return compositeFilter
}

// ----- Defines COM CompositeFilter as a MAL Composite -----
func (a *CompositeFilter) Composite() Composite {
	return a
}

// ================================================================================
// Defines COM CompositeFilter type as a MAL Element
// ================================================================================
// Registers COM CompositeFilter type for polymorpsism handling
func init() {
	RegisterMALElement(COM_COMPOSITE_FILTER_SHORT_FORM, NullCompositeFilter)
}

// ----- Defines COM CompositeFilter as a MAL Element -----
// Returns the absolute short form of the element type
func (*CompositeFilter) GetShortForm() Long {
	return COM_COMPOSITE_FILTER_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (*CompositeFilter) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element belongs to
func (*CompositeFilter) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

func (*CompositeFilter) GetServiceNumber() UShort {
	return ARCHIVE_SERVICE_SERVICE_NUMBER
}

// Returns the relative short form of the element type
func (*CompositeFilter) GetTypeShortForm() Integer {
	return COM_COMPOSITE_FILTER_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (c *CompositeFilter) Encode(encoder Encoder) error {
	// FieldName  String
	err := encoder.EncodeNullableElement(&c.FieldName)
	if err != nil {
		return err
	}

	// Type       ExpressionOperator
	err = encoder.EncodeSmallEnum(uint8(c.Type))
	if err != nil {
		return err
	}

	// FieldValue Attribute
	return encoder.EncodeNullableAttribute(c.FieldValue)
}

// Decodes and instance of CompositeFilter using the supplied decoder
func (*CompositeFilter) Decode(decoder Decoder) (Element, error) {
	return DecodeCompositeFilter(decoder)
}

func DecodeCompositeFilter(decoder Decoder) (*CompositeFilter, error) {
	// FieldName  String
	element, err := decoder.DecodeNullableElement(NullString)
	if err != nil {
		return nil, err
	}
	fieldName := element.(*String)

	// Type       ExpressionOperator
	elementType, err := decoder.DecodeSmallEnum()
	if err != nil {
		return nil, err
	}
	_type := ExpressionOperator(elementType)

	// FieldValue Attribute Nullable
	fieldValue, err := decoder.DecodeNullableAttribute()
	if err != nil {
		return nil, err
	}

	// Create CompositeFilter
	compositeFilter := &CompositeFilter{
		*fieldName,
		_type,
		fieldValue,
	}

	return compositeFilter, nil
}

// The methods allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism
func (*CompositeFilter) CreateElement() Element {
	return new(CompositeFilter)
}

func (c *CompositeFilter) IsNull() bool {
	return c == nil
}

func (*CompositeFilter) Null() Element {
	return NullCompositeFilter
}
