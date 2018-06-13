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
)

// ArchiveQuery structure is used to specify filters on the
// common parts of an object in an archive
type ArchiveQuery struct {
	Domain        *IdentifierList
	Network       *Identifier
	Provider      *URI
	Related       Long
	Source        *ObjectId
	StartTime     *FineTime
	EndTime       *FineTime
	SortOrder     *Boolean
	SortFieldName *String
}

var (
	NullArchiveQuery *ArchiveQuery = nil
)

const (
	COM_ARCHIVE_QUERY_TYPE_SHORT_FORM Integer = 0x02
	COM_ARCHIVE_QUERY_SHORT_FORM      Long    = 0x2000201000002
)

// NewArchiveQuery TODO:
func NewArchiveQuery(domain *IdentifierList,
	network *Identifier,
	provider *URI,
	related Long,
	source *ObjectId,
	startTime *FineTime,
	endTime *FineTime,
	sortOrder *Boolean,
	sortFieldName *String) *ArchiveQuery {
	archiveQuery := &ArchiveQuery{
		domain,
		network,
		provider,
		related,
		source,
		startTime,
		endTime,
		sortOrder,
		sortFieldName,
	}
	return archiveQuery
}

// ----- Defines COM ArchiveQuery as a MAL Composite -----
func (a *ArchiveQuery) Composite() Composite {
	return a
}

// ================================================================================
// Defines COM ArchiveQuery type as a MAL Element
// ================================================================================
// Registers COM ArchiveQuery type for polymorpsism handling
func init() {
	RegisterMALElement(COM_ARCHIVE_QUERY_SHORT_FORM, NullArchiveQuery)
}

// ----- Defines COM ArchiveQuery as a MAL Element -----
// Returns the absolute short form of the element type
func (*ArchiveQuery) GetShortForm() Long {
	return COM_ARCHIVE_QUERY_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (a *ArchiveQuery) GetAreaNumber() UShort {
	return a.Source.GetAreaNumber()
}

// Returns the version of the area this element belongs to
func (a *ArchiveQuery) GetAreaVersion() UOctet {
	return a.Source.GetAreaVersion()
}

// Returns the number of the service this element belongs to
func (a *ArchiveQuery) GetServiceNumber() UShort {
	return a.Source.GetServiceNumber()
}

// Returns the relative short form of the element type
func (*ArchiveQuery) GetTypeShortForm() Integer {
	return COM_ARCHIVE_QUERY_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (a *ArchiveQuery) Encode(encoder Encoder) error {
	// Encode domain (NullableIdentifierList)
	err := encoder.EncodeNullableElement(a.Domain)
	if err != nil {
		return err
	}

	// Encode network (NullableIdentifier)
	err = encoder.EncodeNullableIdentifier(a.Network)
	if err != nil {
		return err
	}

	// Encode provider (NullableURI)
	err = encoder.EncodeNullableURI(a.Provider)
	if err != nil {
		return err
	}

	// Encode related (Long)
	err = encoder.EncodeLong(&a.Related)
	if err != nil {
		return err
	}

	// Encode source (NullableObjectId)
	err = encoder.EncodeNullableElement(a.Source)
	if err != nil {
		return err
	}

	// Encode startTime (NullableFineTime)
	err = encoder.EncodeNullableFineTime(a.StartTime)
	if err != nil {
		return err
	}

	// Encode endTime (NullableFineTime)
	err = encoder.EncodeNullableFineTime(a.EndTime)
	if err != nil {
		return err
	}

	// Encode sortOrder (NullableBoolean)
	err = encoder.EncodeNullableBoolean(a.SortOrder)
	if err != nil {
		return err
	}

	// Encode sortFieldName (NullableString)
	return encoder.EncodeNullableString(a.SortFieldName)
}

// Decodes an instance of ObjectDetails using the supplied decoder
func (*ArchiveQuery) Decode(decoder Decoder) (Element, error) {
	return DecodeArchiveQuery(decoder)
}

func DecodeArchiveQuery(decoder Decoder) (*ArchiveQuery, error) {
	// Encode domain (NullableIdentifierList)
	element, err := decoder.DecodeNullableElement(NullIdentifierList)
	if err != nil {
		return nil, err
	}
	domain := element.(*IdentifierList)

	// Encode network (NullableIdentifier)
	element, err = decoder.DecodeNullableElement(NullIdentifier)
	if err != nil {
		return nil, err
	}
	network := element.(*Identifier)

	// Encode provider (NullableURI)
	element, err = decoder.DecodeNullableElement(NullURI)
	if err != nil {
		return nil, err
	}
	provider := element.(*URI)

	// Encode related (Long)
	element, err = decoder.DecodeElement(NullLong)
	if err != nil {
		return nil, err
	}
	related := element.(*Long)

	// Encode source (NullableObjectId)
	element, err = decoder.DecodeNullableElement(NullObjectId)
	if err != nil {
		return nil, err
	}
	source := element.(*ObjectId)

	// Encode startTime (NullableFineTime)
	element, err = decoder.DecodeNullableElement(NullFineTime)
	if err != nil {
		return nil, err
	}
	startTime := element.(*FineTime)

	// Encode endTime (NullableFineTime)
	element, err = decoder.DecodeNullableElement(NullFineTime)
	if err != nil {
		return nil, err
	}
	endTime := element.(*FineTime)

	// Encode sortOrder (NullableBoolean)
	element, err = decoder.DecodeNullableElement(NullBoolean)
	if err != nil {
		return nil, err
	}
	sortOrder := element.(*Boolean)

	// Encode sortFieldName (NullableString)
	element, err = decoder.DecodeNullableElement(NullString)
	if err != nil {
		return nil, err
	}
	sortFieldName := element.(*String)

	archiveQuery := &ArchiveQuery{
		domain,
		network,
		provider,
		*related,
		source,
		startTime,
		endTime,
		sortOrder,
		sortFieldName,
	}

	return archiveQuery, nil
}

// The methods allows the creation of an element in a generic way, i.e., using the MAL Element polymorphism
func (*ArchiveQuery) CreateElement() Element {
	return new(ArchiveQuery)
}

func (a *ArchiveQuery) IsNull() bool {
	return a == nil
}

func (*ArchiveQuery) Null() Element {
	return NullArchiveQuery
}
