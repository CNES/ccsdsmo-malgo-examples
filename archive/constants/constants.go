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
package archive

// Constants for the Archive Service
const (
	ARCHIVE_SERVICE_SERVICE_IDENTIFIER = "Archive"
	ARCHIVE_SERVICE_SERVICE_NUMBER     = 2
)

const (
	ARCHIVE_SERVICE_AREA_IDENTIFIER = "COM"
	DEFAULT_SERVICE_NUMBER          = 0
)

// Constants for the operations
const (
	OPERATION_IDENTIFIER_RETRIEVE = iota + 1
	OPERATION_IDENTIFIER_QUERY
	OPERATION_IDENTIFIER_COUNT
	OPERATION_IDENTIFIER_STORE
	OPERATION_IDENTIFIER_UPDATE
	OPERATION_IDENTIFIER_DELETE
)

const (
	SLEEP_TIME = 250
)
