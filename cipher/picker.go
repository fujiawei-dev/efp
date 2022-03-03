/*
 * @Date: 2022.03.02 11:28
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 11:28
 */

package cipher

import (
	"fmt"
	"strings"

	"gitlab.com/yawning/chacha20.git"
	"golang.org/x/crypto/chacha20poly1305"

	netp "efp/net"
	"efp/net/crypto/stream"
)

func PickCipher(cipherType string, key []byte) (netp.ConnCipher, netp.PacketConnCipher, error) {
	switch strings.ToLower(cipherType) {
	case "aes-128-gcm", "aes-192-gcm", "aes-256-gcm":
		aead, err := aesGCM(key, 0) // 0 for standard 12-byte nonce
		return aeadConn(aead), aeadPacketConn(aead), err

	case "aes-128-gcm-16", "aes-192-gcm-16", "aes-256-gcm-16":
		aead, err := aesGCM(key, 16) // 16-byte nonce for better collision avoidance
		return aeadConn(aead), aeadPacketConn(aead), err

	case "chacha20-ietf-poly1305":
		aead, err := chacha20poly1305.New(key)
		return aeadConn(aead), aeadPacketConn(aead), err

	case "aes-128-ctr", "aes-192-ctr", "aes-256-ctr":
		cipher, err := aesCTR(key)
		return streamConn(cipher), streamPacketConn(cipher), err

	case "aes-128-cfb", "aes-192-cfb", "aes-256-cfb":
		cipher, err := aesCFB(key)
		return streamConn(cipher), streamPacketConn(cipher), err

	case "chacha20-ietf":
		if len(key) != chacha20.KeySize {
			return nil, nil, stream.ErrShortPacket
		}
		cipher := chacha20ietfkey(key)
		return streamConn(cipher), streamPacketConn(cipher), nil

	case "dummy": // only for benchmarking and debugging
		return dummyConnCipher(), dummyPacketConnCipher(), nil

	default:
		err := fmt.Errorf("cipher not supported: %s", cipherType)
		return nil, nil, err
	}
}
