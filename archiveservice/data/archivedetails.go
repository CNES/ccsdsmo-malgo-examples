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
	. "github.com/CNES/ccsdsmo-malgo/com"
	. "github.com/CNES/ccsdsmo-malgo/mal"
)

// ArchiveDetails structure is used to hold information about a single entry in an Archive
type ArchiveDetails struct {
	InstId    Long
	Details   ObjectDetails
	Network   *Identifier
	Timestamp *FineTime
	Provider  *URI
}

var (
	NullArchiveDetails *ArchiveDetails = nil
)

const (
	COM_ARCHIVE_DETAILS_TYPE_SHORT_FORM Integer = 0x01
	COM_ARCHIVE_DETAILS_SHORT_FORM      Long    = 0x2000201000001
)

func NewArchiveDetails(instId Long, details ObjectDetails, network *Identifier, timestamp *FineTime, provider *URI) *ArchiveDetails {
	archiveDetails := &ArchiveDetails{
		instId,
		details,
		network,
		timestamp,
		provider,
	}
	return archiveDetails
}

// ----- Defines COM ArchiveDetails as a MAL Composite -----
func (a *ArchiveDetails) Composite() Composite {
	return a
}

// ================================================================================
// Defines COM ArchiveDetails type as a MAL Element
// ================================================================================
// Registers COM ArchiveDetails type for polymorpsism handling
func init() {
	RegisterMALElement(COM_ARCHIVE_DETAILS_SHORT_FORM, NullArchiveDetails)
}

// ----- Defines COM ArchiveDetails as a MAL Element -----
// Returns the absolute short form of the element type
func (*ArchiveDetails) GetShortForm() Long {
	return COM_ARCHIVE_DETAILS_SHORT_FORM
}

// Returns the number of the area this element belongs to
func (a *ArchiveDetails) GetAreaNumber() UShort {
	return a.Details.GetAreaNumber()
}

// Returns the version of the area this element belongs to
func (a *ArchiveDetails) GetAreaVersion() UOctet {
	return a.Details.GetAreaVersion()
}

func (a *ArchiveDetails) GetServiceNumber() UShort {
	return a.Details.GetServiceNumber()
}

// Returns the relative short form of the element type
func (*ArchiveDetails) GetTypeShortForm() Integer {
	return COM_ARCHIVE_DETAILS_TYPE_SHORT_FORM
}

// ----- Encoding and Decoding -----
// Encodes this element using the supplied encoder
func (a *ArchiveDetails) Encode(encoder Encoder) error {
	specific := encoder.LookupSpecific(COM_ARCHIVE_DETAILS_SHORT_FORM)
	if specific != nil {
		return specific(a, encoder)
	}
	
	// Encode instId (Long)
	err := encoder.EncodeElement(&a.InstId)
	if err != nil {
		return err
	}

	// Encode details (ObjectDetails)
	err = encoder.EncodeElement(&a.Details)
	if err != nil {
		return err
	}

	// Encode network (NullableIdentifier)
	err = encoder.EncodeNullableElement(a.Network)
	if err != nil {
		return err
	}

	// Encode timestamp (NullableFineTime)
	err = encoder.EncodeNullableElement(a.Timestamp)
	if err != nil {
		return err
	}

	// Encode provider (NullableURI)
	return encoder.EncodeNullableElement(a.Provider)
}

// Decodes and instance of ArchiveDetails using the supplied decoder
func (*ArchiveDetails) Decode(decoder Decoder) (Element, error) {
	specific := decoder.LookupSpecific(COM_ARCHIVE_DETAILS_SHORT_FORM)
	if specific != nil {
		return specific(decoder)
	}
	return DecodeArchiveDetails(decoder)
}

func DecodeArchiveDetails(decoder Decoder) (*ArchiveDetails, error) {
	// Decode instId (Long)
	element, err := decoder.DecodeElement(NullLong)
	if err != nil {
		return nil, err
	}
	instId := element.(*Long)

	// Decode details (ObjectDetails)
	element, err = decoder.DecodeElement(NullObjectDetails)
	if err != nil {
		return nil, err
	}
	details := element.(*ObjectDetails)

	// Decode network (NullableIdentifier)
	element, err = decoder.DecodeNullableElement(NullIdentifier)
	if err != nil {
		return nil, err
	}
	network := element.(*Identifier)

	// Decode timestamp (NullableFineTime)
	element, err = decoder.DecodeNullableElement(NullFineTime)
	if err != nil {
		return nil, err
	}
	timestamp := element.(*FineTime)

	// Decode provider (NullableURI)
	element, err = decoder.DecodeNullableElement(NullURI)
	if err != nil {
		return nil, err
	}
	provider := element.(*URI)

	archiveDetails := &ArchiveDetails{
		*instId,
		*details,
		network,
		timestamp,
		provider,
	}

	return archiveDetails, nil
}

// The methods allows the creation of an element in a generic way, i.e., using     the MAL Element polymorphism
func (*ArchiveDetails) CreateElement() Element {
	return new(ArchiveDetails)
}

func (a *ArchiveDetails) IsNull() bool {
	return a == nil
}

func (*ArchiveDetails) Null() Element {
	return NullArchiveDetails
}
