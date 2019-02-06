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
package utils

import (
	"bytes"
	"strings"

	. "github.com/CNES/ccsdsmo-malgo/com"
	. "github.com/CNES/ccsdsmo-malgo/mal"
	. "github.com/CNES/ccsdsmo-malgo/mal/encoding/binary"
)

// AdaptDomainToString transforms a list of Identifiers to a domain of this
// type: first.second.third.[...]
func AdaptDomainToString(identifierList IdentifierList) String {
	var domain String
	for i := 0; i < identifierList.Size(); i++ {
		domain += String(*identifierList.GetElementAt(i).(*Identifier))
		if i+1 < identifierList.Size() {
			domain += "."
		}
	}
	return domain
}

// AdaptDomainToIdentifierList transforms a domain of this
// type: first.second.third.[...] to a list of Identifiers
func AdaptDomainToIdentifierList(domain string) IdentifierList {
	var identifierList = NewIdentifierList(0)
	var domains = strings.Split(domain, ".")
	for i := 0; i < len(domains); i++ {
		identifierList.AppendElement(NewIdentifier(domains[i]))
	}
	return *identifierList
}

func DecodeObjectID(encodedObjectId []byte) (*ObjectId, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the decoder
	decoder := factory.NewDecoder(encodedObjectId)

	// Decode the ObjectId
	elem, err := decoder.DecodeNullableElement(NullObjectId)
	if err != nil {
		return nil, err
	}
	objectId := elem.(*ObjectId)

	return objectId, nil
}

func DecodeElement(encodedObjectElement []byte) (Element, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the decoder
	decoder := factory.NewDecoder(encodedObjectElement)

	// Decode the Element
	element, err := decoder.DecodeNullableAbstractElement()
	if err != nil {
		return nil, err
	}

	return element, nil
}

func DecodeElements(_objectId []byte, _element []byte) (*ObjectId, Element, error) {
	// Decode the ObjectId
	objectId, err := DecodeObjectID(_objectId)
	if err != nil {
		return nil, nil, err
	}

	// Decode the Element
	element, err := DecodeElement(_element)
	if err != nil {
		return nil, nil, err
	}

	return objectId, element, nil
}

func EncodeElements(_element Element, _objectId *ObjectId) ([]byte, []byte, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the encoder
	encoder := factory.NewEncoder(make([]byte, 0, 8192))

	// Encode Element
	err := encoder.EncodeNullableAbstractElement(_element)
	if err != nil {
		return nil, nil, err
	}
	element := encoder.Body()

	// Reallocate the encoder
	encoder = factory.NewEncoder(make([]byte, 0, 8192))

	// Encode ObjectId
	err = encoder.EncodeNullableElement(_objectId)
	if err != nil {
		return nil, nil, err
	}
	objectId := encoder.Body()

	return element, objectId, nil
}

// This part is useful for type short form conversion (from typeShortForm to listShortForm)
func TypeShortFormToShortForm(objectType ObjectType) Long {
	return objectType.GetMALBodyType()
}

// ConvertToListShortForm converts an ObjectType to a Long (which
// will be used for a List Short Form)
func ConvertToListShortForm(objectType ObjectType) Long {
	return objectType.GetMALBodyListType()
}

func CheckCondition(cond *bool, buffer *bytes.Buffer) {
	if *cond {
		buffer.WriteString(" AND")
	} else {
		*cond = true
	}
}
