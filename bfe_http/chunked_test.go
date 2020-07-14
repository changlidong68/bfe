// Copyright (c) 2019 Baidu, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bfe_http

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"runtime"
	"testing"
)

func TestChunk(t *testing.T) {
	var b bytes.Buffer

	w := newChunkedWriter(&b)
	const chunk1 = "hello, "
	const chunk2 = "world! 0123456789abcdef"
	w.Write([]byte(chunk1))
	w.Write([]byte(chunk2))
	w.Close()

	if g, e := b.String(), "7\r\nhello, \r\n17\r\nworld! 0123456789abcdef\r\n0\r\n"; g != e {
		t.Fatalf("chunk writer wrote %q; want %q", g, e)
	}

	r := newChunkedReader(&b)
	data, err := ioutil.ReadAll(r)
	if err != nil {
		t.Logf(`data: "%s"`, data)
		t.Fatalf("ReadAll from reader: %v", err)
	}
	if g, e := string(data), chunk1+chunk2; g != e {
		t.Errorf("chunk reader read %q; want %q", g, e)
	}
}

func TestChunkReaderAllocs(t *testing.T) {
	// temporarily set GOMAXPROCS to 1 as we are testing memory allocations
	defer runtime.GOMAXPROCS(runtime.GOMAXPROCS(1))
	var buf bytes.Buffer
	w := newChunkedWriter(&buf)
	a, b, c := []byte("aaaaaa"), []byte("bbbbbbbbbbbb"), []byte("cccccccccccccccccccccccc")
	w.Write(a)
	w.Write(b)
	w.Write(c)
	w.Close()

	r := newChunkedReader(&buf)
	readBuf := make([]byte, len(a)+len(b)+len(c)+1)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m0 := ms.Mallocs

	n, err := io.ReadFull(r, readBuf)

	runtime.ReadMemStats(&ms)
	mallocs := ms.Mallocs - m0
	if mallocs > 1 {
		t.Errorf("%d mallocs; want <= 1", mallocs)
	}

	if n != len(readBuf)-1 {
		t.Errorf("read %d bytes; want %d", n, len(readBuf)-1)
	}
	if err != io.ErrUnexpectedEOF {
		t.Errorf("read error = %v; want ErrUnexpectedEOF", err)
	}
}

func TestParseHexUint(t *testing.T) {
	for i := uint64(0); i <= 1234; i++ {
		line := []byte(fmt.Sprintf("%x", i))
		got, err := parseHexUint(line)
		if err != nil {
			t.Fatalf("on %d: %v", i, err)
		}
		if got != i {
			t.Errorf("for input %q = %d; want %d", line, got, i)
		}
	}
	_, err := parseHexUint([]byte("bogus"))
	if err == nil {
		t.Error("expected error on bogus input")
	}
}
