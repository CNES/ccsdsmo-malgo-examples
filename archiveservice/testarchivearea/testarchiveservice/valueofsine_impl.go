/**
 * MIT License
 *
 * Copyright (c) 2018-20 CNES
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
package testarchiveservice

import (
	"github.com/CNES/ccsdsmo-malgo-examples/archiveservice/testarchivearea"
	"github.com/CNES/ccsdsmo-malgo/com"
	"github.com/CNES/ccsdsmo-malgo/mal"
	"github.com/CNES/ccsdsmo-malgo/mal/debug"
)

var (
	logger debug.Logger = debug.GetLogger("archive.data")
)

func init() {
	// We need to register the ObjectType mapping of the COM object
	// this code cannot currently be generated

	// In the tests the short form is also used as the COM type number
	// the fields values are those used in the test
	comObjType := com.ObjectType{
		Area:    testarchivearea.AREA_NUMBER,
		Service: SERVICE_NUMBER,
		Version: testarchivearea.AREA_VERSION,
		Number:  mal.UShort(VALUEOFSINE_TYPE_SHORT_FORM),
	}
	err := comObjType.RegisterMALBodyType(VALUEOFSINE_SHORT_FORM)
	if err != nil {
		logger.Errorf("ValueOfSine.init, cannot register COM object: %s", err.Error())
	}
}
