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

	netp "efp/net"
	"efp/net/crypto/aead"
)

func aeadConn(cipher cipher.AEAD) netp.ConnCipher {
	return func(c net.Conn) net.Conn { return aead.NewConn(c, cipher) }
}

func aeadPacketConn(cipher cipher.AEAD) netp.PacketConnCipher {
	return func(c net.PacketConn) net.PacketConn { return aead.NewPacketConn(c, cipher) }
}

func aesGCM(key []byte, nonceSize int) (cipher.AEAD, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if nonceSize > 0 {
		return cipher.NewGCMWithNonceSize(blk, nonceSize)
	}
	return cipher.NewGCM(blk) // standard 12-byte nonce
}
