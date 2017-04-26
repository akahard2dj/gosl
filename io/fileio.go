// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package io

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cpmech/gosl/chk"
)

// functions to handle filenames //////////////////////////////////////////////////////////////////

// FnKey returns the file name key (without path and extension, if any)
func FnKey(fn string) string {
	base := filepath.Base(fn)
	return base[:len(base)-len(filepath.Ext(base))]
}

// FnExt returns the extension of a file name.
// The extension is the suffix beginning at the final dot in the final element of path; it is empty if there is no dot.
func FnExt(fn string) string {
	return filepath.Ext(fn)
}

// PathKey returs the full path except the extension
func PathKey(fn string) string {
	return fn[:len(fn)-len(filepath.Ext(fn))]
}

// RemoveAll deletes all files matching filename specified by key (be careful)
func RemoveAll(key string) {
	fns, _ := filepath.Glob(os.ExpandEnv(key))
	for _, fn := range fns {
		os.RemoveAll(fn)
	}
}

// functions to write files ///////////////////////////////////////////////////////////////////////

// AppendToFile appends data to an existent (or new) file
func AppendToFile(fn string, buffer ...*bytes.Buffer) {
	fil, err := os.OpenFile(os.ExpandEnv(fn), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		chk.Panic("cannot create file <%s>", fn)
	}
	defer fil.Close()
	for k, _ := range buffer {
		if buffer[k] != nil {
			fil.Write(buffer[k].Bytes())
		}
	}
}

// WriteFile writes data to a new file with given bytes.Buffer(s)
func WriteFile(fn string, buffer ...*bytes.Buffer) {
	fil, err := os.Create(os.ExpandEnv(fn))
	if err != nil {
		chk.Panic("cannot create file <%s>", fn)
	}
	defer fil.Close()
	for k, _ := range buffer {
		if buffer[k] != nil {
			fil.Write(buffer[k].Bytes())
		}
	}
}

// WriteFileD writes data to a new file after creating a directory
func WriteFileD(dirout, fn string, buffer ...*bytes.Buffer) {
	os.MkdirAll(dirout, 0777)
	WriteFile(filepath.Join(dirout, fn), buffer...)
}

// WriteFileV writes data to a new file, and shows message
func WriteFileV(fn string, buffer ...*bytes.Buffer) {
	WriteFile(fn, buffer...)
	Pf("file <%s> written\n", fn)
}

// WriteFileVD writes data to a new file, and shows message after creating a directory
func WriteFileVD(dirout, fn string, buffer ...*bytes.Buffer) {
	WriteFileD(dirout, fn, buffer...)
	Pf("file <%s> written\n", filepath.Join(dirout, fn))
}

// WriteStringToFile writes string to a new file
func WriteStringToFile(fn, data string) {
	fil, err := os.Create(os.ExpandEnv(fn))
	if err != nil {
		chk.Panic("cannot create file <%s>", fn)
	}
	defer fil.Close()
	fil.WriteString(data)
}

// WriteStringToFileD writes string to a new file after creating a directory
func WriteStringToFileD(dirout, fn, data string) {
	os.MkdirAll(dirout, 0777)
	WriteStringToFile(filepath.Join(dirout, fn), data)
}

// WriteBytesToFile writes slice of bytes to a new file
func WriteBytesToFile(fn string, b []byte) {
	fil, err := os.Create(os.ExpandEnv(fn))
	if err != nil {
		chk.Panic("cannot create file <%s>", fn)
	}
	defer fil.Close()
	if _, err = fil.Write(b); err != nil {
		chk.Panic("%v", err)
	}
}

// WriteBytesToFileD writes slice of bytes to a new file after creating a directory
func WriteBytesToFileD(dirout, fn string, b []byte) {
	os.MkdirAll(dirout, 0777)
	WriteBytesToFile(filepath.Join(dirout, fn), b)
}

// WriteBytesToFileVD writes slice of bytes to a new file, and print message, after creating a directory
func WriteBytesToFileVD(dirout, fn string, b []byte) {
	os.MkdirAll(dirout, 0777)
	WriteBytesToFile(filepath.Join(dirout, fn), b)
	Pf("file <%s> written\n", filepath.Join(dirout, fn))
}

// functions to read files ////////////////////////////////////////////////////////////////////////

// OpenFileR opens a file for reading data
func OpenFileR(fn string) (fil *os.File, err error) {
	fil, err = os.Open(os.ExpandEnv(fn))
	return
}

// ReadFile reads bytes from a file
func ReadFile(fn string) (b []byte, err error) {
	return ioutil.ReadFile(os.ExpandEnv(fn))
}

// ReadLinesCallback is used in ReadLines to process line by line during reading of a file
type ReadLinesCallback func(idx int, line string) (stop bool)

// ReadLines reads lines from a file and calls ReadLinesCallback to process each line being read
func ReadLines(fn string, cb ReadLinesCallback) {
	fil, err := os.Open(os.ExpandEnv(fn))
	if err != nil {
		chk.Panic("could not open file <%s>", fn)
	}
	defer fil.Close()
	r := bufio.NewReader(fil)
	idx := 0
	for {
		lin, prefix, errl := r.ReadLine()
		if prefix {
			chk.Panic("cannot read long line. file = <%s>", fn)
		}
		if errl == io.EOF {
			break
		}
		if errl != nil {
			chk.Panic("cannot read line. file = <%s>", fn)
		}
		stop := cb(idx, string(lin))
		if stop {
			break
		}
		idx++
	}
}

// ReadLinesFile reads lines from a file and calls ReadLinesCallback to process each line being read
func ReadLinesFile(fil *os.File, cb ReadLinesCallback) (oserr error) {
	r := bufio.NewReader(fil)
	idx := 0
	for {
		lin, prefix, errl := r.ReadLine()
		if prefix {
			return chk.Err("cannot read long line. file = <%s>", fil.Name())
		}
		if errl == io.EOF {
			break
		}
		if errl != nil {
			return chk.Err("cannot read line. file = <%s>", fil.Name())
		}
		stop := cb(idx, string(lin))
		if stop {
			break
		}
		idx++
	}
	return
}

//  ReadTableOrPanic reads text file as ReadTable; but panic on errors
func ReadTableOrPanic(fn string) (keys []string, T map[string][]float64) {
	var err error
	keys, T, err = ReadTable(fn)
	if err != nil {
		chk.Panic("cannot read table:\n%v", err)
	}
	return
}

// ReadTable reads a text file in which the first line contains the headers and the next lines the float64
// type of numeric values. The number of columns must be equal, including for the headers
func ReadTable(fn string) (keys []string, T map[string][]float64, err error) {
	f, err := OpenFileR(fn)
	if err != nil {
		return
	}
	header := true
	ReadLinesFile(f, func(idx int, line string) (stop bool) {
		r := strings.Fields(line)
		if len(r) == 0 { // skip empty lines
			return
		}
		if r[0] == "#" { // skip comments
			return
		}
		ncol := len(r)
		if ncol < 1 {
			err = chk.Err("ReadTable: number of columns must be at least 1")
			return true
		}
		if header {
			T = make(map[string][]float64)
			keys = make([]string, ncol)
			for i := 0; i < ncol; i++ {
				keys[i] = r[i]
				T[r[i]] = make([]float64, 0)
			}
			header = false
		} else {
			for i := 0; i < ncol; i++ {
				T[keys[i]] = append(T[keys[i]], Atof(r[i]))
			}
		}
		return
	})
	return
}

// ReadMatrix reads a text file in which the float64 type of numeric values represent
// a matrix of data. The number of columns must be equal, including for the headers
func ReadMatrix(fn string) (M [][]float64, err error) {
	f, err := OpenFileR(fn)
	if err != nil {
		return
	}
	ncol_fix := 0
	ReadLinesFile(f, func(idx int, line string) (stop bool) {
		r := strings.Fields(line)
		if len(r) == 0 { // skip empty lines
			return
		}
		if r[0] == "#" { // skip comments
			return
		}
		ncol := len(r)
		if ncol < 1 {
			err = chk.Err("ReadMatrix: number of columns must be at least 1")
			return true
		}
		if M == nil {
			M = make([][]float64, 0)
			ncol_fix = ncol
		}
		if ncol != ncol_fix {
			err = chk.Err("ReadMatrix: nubmer of columns must be equal for all lines (%d != %d)", ncol, ncol_fix)
			return
		}
		vals := make([]float64, ncol)
		for i := 0; i < ncol; i++ {
			vals[i] = Atof(r[i])
		}
		M = append(M, vals)
		return
	})
	return
}
