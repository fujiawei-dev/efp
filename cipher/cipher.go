/*
 * @Date: 2022.03.06 10:34
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.06 10:34
 */

package cipher

import (
	"crypto/cipher"
	"efp/net/crypto/stream"
	"errors"
	"fmt"
	"golang.org/x/crypto/chacha20poly1305"
	"io"
	"sort"
)

// ErrKeySize means the key size does not meet the requirement of cipher.
var ErrKeySize = errors.New("key size error")

// ErrCipherNotSupported means the cipher has not been implemented.
var ErrCipherNotSupported = errors.New("cipher not supported")

// List of AEAD ciphers: key size in bytes and constructor
var aeadList = map[string]struct {
	KeySize int
	New     func(key []byte) (cipher.AEAD, error)
}{
	"aes-128-gcm":            {16, aesGCM},
	"aes-192-gcm":            {24, aesGCM},
	"aes-256-gcm":            {32, aesGCM},
	"aes-128-gcm-16":         {16, aesGCM16},
	"aes-192-gcm-16":         {24, aesGCM16},
	"aes-256-gcm-16":         {32, aesGCM16},
	"chacha20-ietf-poly1305": {32, chacha20poly1305.New},
}

// List of stream ciphers: key size in bytes and constructor
var streamList = map[string]struct {
	KeySize int
	New     func(key []byte) (stream.Cipher, error)
}{
	"aes-128-ctr":   {16, aesCTR},
	"aes-192-ctr":   {24, aesCTR},
	"aes-256-ctr":   {32, aesCTR},
	"aes-128-cfb":   {16, aesCFB},
	"aes-192-cfb":   {24, aesCFB},
	"aes-256-cfb":   {32, aesCFB},
	"chacha20-ietf": {32, newChacha20ietf},
}

// Return two lists of sorted cipher names.
func availableCiphers() (aead []string, stream []string) {
	for k := range aeadList {
		aead = append(aead, k)
	}

	for k := range streamList {
		stream = append(stream, k)
	}

	sort.Strings(aead)
	sort.Strings(stream)
	return
}

// PrintCiphers Print available ciphers to w
func PrintCiphers(w io.Writer) {
	fmt.Fprintf(w, "## Available AEAD ciphers (recommended)\n\n")

	aead, stream := availableCiphers()
	for _, name := range aead {
		fmt.Fprintf(w, "%s\n", name)
	}

	fmt.Fprintf(w, "\n## Available stream ciphers\n\n")
	for _, name := range stream {
		fmt.Fprintf(w, "%s\n", name)
	}
}
