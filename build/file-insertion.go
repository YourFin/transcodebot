// Copyright © 2018 Patrick Nuckolls <nuckollsp at gmail>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package build

import (
	"os"
	"io"
	"compress/gzip"
	"sync/atomic"
	"sync"

	"github.com/yourfin/transcodebot/common"
)

type appendedFile struct {
	//TODO: Copy to temp file before opening a reader
	//TODO: CopyToTmp bool
	Compressed bool
	StartFilePtr int64
	ZippedSize int64
}

const METADATA_VERSION string = "0.1"
type appendedMetadata struct {
	Version string
	Files map[string]appendedFile
}

type BinAppender struct {
	fileHandle *os.File
	writer io.Writer
	metadata appendedMetadata
	writeLock *sync.Mutex
}

func MakeCustomAppender(filename string, writeWrapper func(io.Writer) io.Writer) (*BinAppender, error){
	var err error
	output := BinAppender{}
	output.fileHandle, err = os.OpenFile(filename, os.O_RDWR, 0755)
	if err != nil {
		return nil, err
	}
	output.writer = writeWrapper(output.fileHandle)
	output.writeLock = &sync.Mutex{}
	output.metadata = appendedMetadata{}
	output.metadata.Files = make(map[string]appendedFile)
	output.metadata.Version = METADATA_VERSION
	return &output, nil
}

func MakeGzAppender(filename string) (*BinAppender, error) {
	return MakeCustomAppender(
		filename,
		func(input io.Writer) io.Writer {
			return input
		},
	)
}

//func (bb BinAppender) ()

// Procedure:
//  appendFile
// Purpose:
//  To gzip and pack a file onto the end of another file and return
//  the append files byte location
// Parameters:
//  The file to append: toAppend string
//  The file to be appended to: appendee string
// Produces:
//  The "index" of the first byte of the appended file: startByte int
//  The number of bytes in the appended file: size int
//  Any errors that occur: err error
//  Side effects:
//    The file at $appendee will have the binary data of $toAppend at the end of it
// Preconditions:
//  toAppend and apendee exist and are readable in the file system
//  apendee is writeable on the filesystem
// Postconditions:
//  bash equivalent is executed:
//    cat $toAppend | gzip >> $apendee
//
//  $apendee.toByteArray()[$startByte:$size].gunzip() == $toAppend.toByteArray()
//
func AppendFile(toAppend, appendee string) (startByte, size int64, err error) {
	appendeeFile, err := os.OpenFile(appendee, os.O_RDWR, 0755)
	if err != nil {
		return 0, 0, err
	}
	toAppendFile, err := os.Open(toAppend)
	if err != nil {
		return 0, 0, err
	}
	defer toAppendFile.Close()

	startByte, err = appendeeFile.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, 0, err
	}

	gzWriter := gzip.NewWriter(appendeeFile)

	num, err := io.Copy(gzWriter, toAppendFile)
	if err != nil {
		return 0, 0, err
	}
	gzWriter.Close()
	common.Println(num)

	end, err := appendeeFile.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, 0, err
	}
	common.Println("end: ", end)

	err = appendeeFile.Close()
	if err != nil {
		return 0, 0, err
	}

	size = end - startByte

	return startByte, size, nil
}

func ReadStuff(start, size int64, path string) error {
	readFile, err := os.Open(path)
	if err != nil {
		return err
	}
	_, err = readFile.Seek(start, io.SeekStart)
	if err != nil {
		return err
	}

	gzReader, err := gzip.NewReader(io.LimitReader(readFile, size))

	if err != nil {
		return err
	}

	num, err := io.Copy(os.Stdout, gzReader)
	common.Println()
	common.Println(num)
	return err
}

// Type:
//  writeCounter
// Purpose:
//  To count the number of bytes written to an io.writer
//
//  As it turns out, gzip.Writer().Write() returns the
//  number of bytes in, not out. This shim goes between
//  gzip and the filesystem so we can figure out how many
//  bytes were actually written out to pasture
type writeCounter struct {
	Counter int64
	child io.Writer
}
func (w *writeCounter) Write(p []byte) (n int, err error) {
	n, err = w.child.Write(p)
	atomic.AddInt64(&w.Counter, int64(n))
	return
}
func newWriteCounter(child io.Writer) (*writeCounter) {
	return &writeCounter{0, child}
}
