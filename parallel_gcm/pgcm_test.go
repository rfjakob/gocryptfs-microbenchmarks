package main

import (
	"crypto/aes"
	"crypto/cipher"
	"sync"
	"testing"

	"github.com/rfjakob/gocryptfs-microbenchmarks/parallel_gcm/stupidgcm"
)

var gogcm cipher.AEAD
var opensslgcm cipher.AEAD

var nonce = make([]byte, 16)
var ad = make([]byte, 24)
var plaintext = make([]byte, 4096)

const nblocks = 32

func init() {
	// AES-256 with a 32-byte key
	k := make([]byte, 32)
	a, err := aes.NewCipher(k)
	if err != nil {
		panic(err)
	}
	gogcm, err = cipher.NewGCMWithNonceSize(a, 16)
	if err != nil {
		panic(err)
	}
	opensslgcm = stupidgcm.New(k, false)
}

// encrypt "n" 4kiB blocks
func encryptBlocks(gcmimpl cipher.AEAD, n int) {
	for i := 0; i < n; i++ {
		gcmimpl.Seal(nonce, nonce, plaintext, ad)
	}
}

// encrypt 132kiB of 4kiB blocks with p-way parallelism
func encryptBlocksPar(gcmimpl cipher.AEAD, p int) {
	var wg sync.WaitGroup
	for g := 0; g < p; g++ {
		wg.Add(1)
		go func() {
			encryptBlocks(gcmimpl, nblocks/p)
			wg.Done()
		}()
	}
	wg.Wait()
}

// 1 thread
func Benchmark1_gogcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		encryptBlocks(gogcm, nblocks)
	}
}

// 2 threads (main goroutine + new one)
func Benchmark2coop_gogcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			encryptBlocks(gogcm, nblocks/2)
			wg.Done()
		}()
		encryptBlocks(gogcm, nblocks/2)
		wg.Wait()
	}
}

// 2 threads (two new goroutines)
func Benchmark2_gogcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		encryptBlocksPar(gogcm, 2)
	}
}

// 4 threads (new goroutines)
func Benchmark4_gogcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		encryptBlocksPar(gogcm, 4)
	}
}

// 8 threads (new goroutines)
func Benchmark8_gogcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		encryptBlocksPar(gogcm, 8)
	}
}

func Benchmark1_opensslgcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		encryptBlocks(opensslgcm, nblocks)
	}
}
func Benchmark2_opensslgcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		encryptBlocksPar(opensslgcm, 2)
	}
}
func Benchmark2coop_opensslgcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			encryptBlocks(opensslgcm, nblocks/2)
			wg.Done()
		}()
		encryptBlocks(opensslgcm, nblocks/2)
		wg.Wait()
	}
}
func Benchmark4_opensslgcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		encryptBlocksPar(opensslgcm, 4)
	}
}
func Benchmark8_opensslgcm(b *testing.B) {
	b.SetBytes(int64(len(plaintext) * nblocks))
	for i := 0; i < b.N; i++ {
		encryptBlocksPar(opensslgcm, 8)
	}
}
