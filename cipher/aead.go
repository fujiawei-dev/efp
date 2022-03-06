/*
 * @Date: 2022.03.02 15:28
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 15:28
 */

package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"net"

	netp "github.com/fujiawei-dev/efp/net"
	"github.com/fujiawei-dev/efp/net/crypto/aead"
)

// AEAD ciphers

func aeadConnCipher(cipher cipher.AEAD) netp.ConnCipher {
	return func(c net.Conn) net.Conn { return aead.NewConn(c, cipher) }
}

// AES-GCM with standard 12-byte nonce
func aesGCM(key []byte) (cipher.AEAD, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(blk)
}

// AES-GCM with 16-byte nonce for better collision avoidance
func aesGCM16(key []byte) (cipher.AEAD, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCMWithNonceSize(blk, 16)
}
