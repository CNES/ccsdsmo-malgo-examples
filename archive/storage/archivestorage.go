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
package storage

import "os"

const (
	FILE_TO_READ_AND_WRITE string = "./archive/storage/archive.txt"
)

// WriteInArchive :
func WriteInArchive(data []byte) error {
	var file *os.File
	var err error
	// Open the file in Append mode
	file, err = os.OpenFile(FILE_TO_READ_AND_WRITE, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		// If we have an error it may mean that the file doesn't exist
		file, err = os.Create(FILE_TO_READ_AND_WRITE)
		if err != nil {
			return err
		}
	}

	_, err = file.WriteString(string(data) + "\n")
	if err != nil {
		return err
	}

	file.Sync()

	return nil
}

// StoreInArchive : store objects in the Archive
func StoreInArchive() error {
	return nil
}
