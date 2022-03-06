/*
 * @Date: 2022.03.06 10:34
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.06 10:34
 */

package cipher

import (
	"crypto/cipher"
	"crypto/md5"
	"fmt"
	"io"
	"sort"
	"strings"

	"golang.org/x/crypto/chacha20poly1305"

	netp "github.com/fujiawei-dev/efp/net"
	"github.com/fujiawei-dev/efp/net/crypto/stream"
)

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

func New(cipherType string, key []byte, password string) (netp.ConnCipher, error) {
	cipherType = strings.ToLower(cipherType)

	if cipherType == "dummy" || cipherType == "" || (len(key) == 0 && password == "") {
		return dummyConnCipher(), nil
	}

	if choice, ok := aeadList[cipherType]; ok {
		if len(key) == 0 {
			key = kdf(password, choice.KeySize)
		} else if len(key) != choice.KeySize {
			return nil, fmt.Errorf("key size error: need %d-byte key", choice.KeySize)
		}
		c, err := choice.New(key)
		return aeadConnCipher(c), err
	}

	if choice, ok := streamList[cipherType]; ok {
		if len(key) == 0 {
			key = kdf(password, choice.KeySize)
		} else if len(key) != choice.KeySize {
			return nil, fmt.Errorf("key size error: need %d-byte key", choice.KeySize)
		}
		c, err := choice.New(key)
		return streamConnCipher(c), err
	}

	return nil, fmt.Errorf("cipher %q not supported", cipherType)
}

// key-derivation function from original Shadowsocks
func kdf(password string, keyLen int) []byte {
	var b, prev []byte
	h := md5.New()
	for len(b) < keyLen {
		h.Write(prev)
		h.Write([]byte(password))
		b = h.Sum(b)
		prev = b[len(b)-h.Size():]
		h.Reset()
	}
	return b[:keyLen]
}
