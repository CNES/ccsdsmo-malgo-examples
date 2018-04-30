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

	. "github.com/ccsdsmo/malgo/com"
	. "github.com/ccsdsmo/malgo/mal"
	. "github.com/ccsdsmo/malgo/mal/encoding/binary"

	. "github.com/etiennelndr/archiveservice/data"
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
	elem, err := decoder.DecodeElement(NullObjectId)
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
	element, err := decoder.DecodeAbstractElement()
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

func EncodeElements(_element Element, _objectId ObjectId) ([]byte, []byte, error) {
	// Create the factory
	factory := new(FixedBinaryEncoding)

	// Create the encoder
	encoder := factory.NewEncoder(make([]byte, 0, 8192))

	// Encode Element
	err := encoder.EncodeAbstractElement(_element)
	if err != nil {
		return nil, nil, err
	}
	element := encoder.Body()

	// Reallocate the encoder
	encoder = factory.NewEncoder(make([]byte, 0, 8192))

	// Encode ObjectId
	err = _objectId.Encode(encoder)
	if err != nil {
		return nil, nil, err
	}
	objectId := encoder.Body()

	return element, objectId, nil
}

// This part is useful for type short form conversion (from typeShortForm to listShortForm)
func TypeShortFormToShortForm(objectType ObjectType) Long {
	var typeShortForm = Long(objectType.Number) | 0xFFFFFFF000000
	var areaVersion = (Long(objectType.Version) << 24) | 0xFFFFF00FFFFFF
	var serviceNumber = (Long(objectType.Service) << 32) | 0xF0000FFFFFFFF
	var areaNumber = (Long(objectType.Area) << 48) | 0x0FFFFFFFFFFFF

	return areaNumber & serviceNumber & areaVersion & typeShortForm
}

// ConvertToListShortForm converts an ObjectType to a Long (which
// will be used for a List Short Form)
func ConvertToListShortForm(objectType ObjectType) Long {
	var listByte []byte
	listByte = append(listByte, byte(objectType.Area), byte(objectType.Service>>8), byte(objectType.Service), byte(objectType.Version))
	typeShort := TypeShortFormToShortForm(objectType)
	quatuor5 := (typeShort & 0x0000F0) >> 4
	if quatuor5 == 0x0 {
		var b byte
		for i := 2; i >= 0; i-- {
			b = byte(typeShort>>uint(i*8)) ^ 255
			if i == 0 {
				b++
			}
			listByte = append(listByte, b)
		}

		var byte0 = Long(listByte[6]) | 0xFFFFFFFFFFF00
		var byte1 = (Long(listByte[5]) << 8) | 0xFFFFFFFFF00FF
		var byte2 = (Long(listByte[4]) << 16) | 0xFFFFFFF00FFFF
		var byte3 = (Long(listByte[3]) << 24) | 0xFFFFF00FFFFFF
		var byte4 = (Long(listByte[2]) << 32) | 0xFFF00FFFFFFFF
		var byte5 = (Long(listByte[1]) << 40) | 0xF00FFFFFFFFFF
		var byte6 = (Long(listByte[0]) << 48) | 0x0FFFFFFFFFFFF

		return byte6 & byte5 & byte4 & byte3 & byte2 & byte1 & byte0
	}

	// Force bytes 2, 3 to 1
	return typeShort | 0x0000000FFFF00
}

func CheckCondition(cond *bool, buffer *bytes.Buffer) {
	if *cond {
		buffer.WriteString(" AND")
	} else {
		*cond = true
	}
}

// WhichExpressionOperatorIsIt transforms an ExpressionOperator to a string
func WhichExpressionOperatorIsIt(expressionOperator ExpressionOperator) string {
	switch expressionOperator {
	case COM_EXPRESSIONOPERATOR_EQUAL:
		return "="
	case COM_EXPRESSIONOPERATOR_DIFFER:
		return "!="
	case COM_EXPRESSIONOPERATOR_GREATER:
		return ">"
	case COM_EXPRESSIONOPERATOR_GREATER_OR_EQUAL:
		return ">="
	case COM_EXPRESSIONOPERATOR_LESS:
		return "<"
	case COM_EXPRESSIONOPERATOR_LESS_OR_EQUAL:
		return "<="
	case COM_EXPRESSIONOPERATOR_CONTAINS:
		return "LIKE '%"
	case COM_EXPRESSIONOPERATOR_ICONTAINS:
		return "NOT LIKE '%"
	default:
		return ""
	}
}
