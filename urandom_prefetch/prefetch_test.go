/*
Test the throughput of reading from /dev/urandom (or, equivalently, the getentropy
syscall) with different block sizes.
This benchmark helped decide the "prefetchN" constant in gocryptfs.

Run: go test -bench .
*/

package main

import (
	"bytes"
	"crypto/rand"
	"sync"
	"testing"
)

// RandBytes gets "n" random bytes from /dev/urandom or panics
func RandBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic("Failed to read random bytes: " + err.Error())
	}
	return b
}

var b bytes.Buffer
var l sync.Mutex

func prefetch(n int, pre int) []byte {
	b2 := make([]byte, n)
	l.Lock()
	have, err := b.Read(b2)
	if have == n && err == nil {
		l.Unlock()
		return b2
	}
	b.Reset()
	b.Write(RandBytes(pre))
	_, err = b.Read(b2)
	if err != nil {
		panic("!!!")
	}
	l.Unlock()
	return b2
}

func Benchmark16(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		RandBytes(16)
	}
}

func Benchmark64(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		prefetch(16, 64)
	}
}

func Benchmark128(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		prefetch(16, 128)
	}
}
func Benchmark256(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		prefetch(16, 256)
	}
}
func Benchmark512(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		prefetch(16, 512)
	}
}
func Benchmark1024(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		prefetch(16, 1024)
	}
}
func Benchmark2048(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		prefetch(16, 2048)
	}
}
func Benchmark4096(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		prefetch(16, 4096)
	}
}

func Benchmark40960(b *testing.B) {
	b.SetBytes(16)
	for i := 0; i < b.N; i++ {
		prefetch(16, 40960)
	}
}
