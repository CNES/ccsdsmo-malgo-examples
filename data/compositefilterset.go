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

type CompositeFilterSet struct {
	Filters *CompositeFilterList
}

var (
	NullCompositeFilterSet *CompositeFilterSet = nil
)

const (
	COM_COMPOSITE_FILTER_SET_TYPE_SHORT_FORM Integer = 0x04
	COM_COMPOSITE_FILTER_SET_FORM            Long    = 0x2000201000004
)

func NewCompositeFilterSet(filters *CompositeFilterList) *CompositeFilterSet {
	compositeFilterSet := &CompositeFilterSet{
		filters,
	}
	return compositeFilterSet
}

// ----- Defines COM CompositeFilterSet as a MAL Composite -----
func (c *CompositeFilterSet) Composite() Composite {
	return c
}

// ----- Defines COM CompositeFilterSet as a COM QueryFilter -----
func (c *CompositeFilterSet) QueryFilter() QueryFilter {
	return c
}

// ================================================================================
// Defines COM CompositeFilterSet type as a MAL Element
// ================================================================================
// Registers COM CompositeFilterSet type for polymorpsism handling
func init() {
	RegisterMALElement(COM_COMPOSITE_FILTER_SET_FORM, NullCompositeFilterSet)
}

// ----- Defines COM ArchiveDetails as a MAL Element -----
// Returns the absolute short form of the element type
func (*CompositeFilterSet) GetShortForm() Long {
	return COM_COMPOSITE_FILTER_SET_FORM
}

// Returns the number of the area this element belongs to
func (*CompositeFilterSet) GetAreaNumber() UShort {
	return COM_AREA_NUMBER
}

// Returns the version of the area this element belongs to
func (*CompositeFilterSet) GetAreaVersion() UOctet {
	return COM_AREA_VERSION
}

func (*CompositeFilterSet) GetServiceNumber() UShort {
	return ARCHIVE_SERVICE_SERVICE_NUMBER
}

// Returns the relative short form of the element type
func (*CompositeFilterSet) GetTypeShortForm() Integer {
	return COM_COMPOSITE_FILTER_SET_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (c *CompositeFilterSet) Encode(encoder Encoder) error {
	return encoder.EncodeNullableElement(c.Filters)
}

// Decodes and instance of CompositeFilterSet using the supplied decoder
func (*CompositeFilterSet) Decode(decoder Decoder) (Element, error) {
	return DecodeCompositeFilterSet(decoder)
}

func DecodeCompositeFilterSet(decoder Decoder) (*CompositeFilterSet, error) {
	// Decode CompositeFilterList
	element, err := decoder.DecodeNullableElement(NullCompositeFilterList)
	if err != nil {
		return nil, err
	}
	filters := element.(*CompositeFilterList)

	// Create CompositeFilterList
	compositeFilterSet := &CompositeFilterSet{
		filters,
	}

	return compositeFilterSet, nil
}

// The methods allows the creation of an element in a generic way, i.e., using     the MAL Element polymorphism
func (*CompositeFilterSet) CreateElement() Element {
	return new(CompositeFilterSet)
}

func (c *CompositeFilterSet) IsNull() bool {
	return c == nil
}

func (*CompositeFilterSet) Null() Element {
	return NullCompositeFilterSet
}
