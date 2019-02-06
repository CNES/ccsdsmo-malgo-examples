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
package errors

import (
	. "github.com/CNES/ccsdsmo-malgo/mal"
)

// ServiceError : TODO
type ServiceError struct {
	ErrorNumber  *UInteger
	ErrorComment *String
	ErrorExtra   Element
}

func EncodeError(encoder Encoder, errorNumber UInteger, errorComment String, errorExtra Element) (Encoder, error) {
	// Encode UInteger
	err := errorNumber.Encode(encoder)
	if err != nil {
		return nil, err
	}

	// SL This information must not be encoded in the message
	/*
	// Encode String
	err = errorComment.Encode(encoder)
	if err != nil {
		return nil, err
	}
	*/

	// Encode Element
	err = encoder.EncodeAbstractElement(errorExtra)
	if err != nil {
		return nil, err
	}

	return encoder, nil
}

func DecodeError(decoder Decoder) (*ServiceError, error) {
	// Decode UInteger
	errorNumber, err := decoder.DecodeElement(NullUInteger)
	if err != nil {
		return nil, err
	}

	// SL This information must not be encoded in the message
	/*
	// Decode String
	errorComment, err := decoder.DecodeElement(NullString)
	if err != nil {
		return nil, err
	}
	*/
	var errorComment Element = NewString("dummy")

	// Decode Element
	errorExtra, err := decoder.DecodeAbstractElement()
	if err != nil {
		return nil, err
	}

	// Create ServiceError
	serviceError := &ServiceError{
		errorNumber.(*UInteger),
		errorComment.(*String),
		errorExtra,
	}

	return serviceError, nil
}
