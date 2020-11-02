/**
 * MIT License
 *
 * Copyright (c) 2018-2020 CNES
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

import (
	"github.com/CNES/ccsdsmo-malgo/mal"
)

// Constants for all the errors
const (
	ARCHIVE_SERVICE_STORE_LIST_SIZE_ERROR                       mal.String = "ArchiveDetailsList and ElementList must have the same size"
	ARCHIVE_SERVICE_OBJECTTYPE_VALUES_ERROR                     mal.String = "ObjectType's attributes must not be equal to 'O'"
	ARCHIVE_SERVICE_IDENTIFIERLIST_VALUES_ERROR                 mal.String = "Domain's elements must not be equal to '*'"
	ARCHIVE_SERVICE_STORE_ARCHIVEDETAILSLIST_VALUES_ERROR       mal.String = "ArchiveDetailsList elements must not be equal to '0', '*' or NULL"
	ARCHIVE_SERVICE_AREA_OBJECT_INSTANCE_IDENTIFIER_VALUE_ERROR mal.String = "Object instance identifier must not be equal to '0'"
	ARCHIVE_SERVICE_QUERY_LISTS_SIZE_ERROR                      mal.String = "The size of the two lists must be the same"
	ARCHIVE_SERVICE_QUERY_SORT_FIELD_NAME_INVALID_ERROR         mal.String = "SortFieldName parameter doesn't reference a defined field"
	ARCHIVE_SERVICE_QUERY_QUERY_FILTER_ERROR                    mal.String = "QueryFilter contains an error"
	ARCHIVE_SERVICE_UNKNOWN_ELEMENT                             mal.String = "Unknown element, cannot find it in the archive"
)

const (
	LENGTH = 16394
)
