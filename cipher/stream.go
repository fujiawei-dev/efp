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

	"gitlab.com/yawning/chacha20.git"

	netp "efp/net"
	"efp/net/crypto/stream"
)

func streamConn(cipher stream.Cipher) netp.ConnCipher {
	return func(c net.Conn) net.Conn { return stream.NewConn(c, cipher) }
}

type ctrStream struct{ cipher.Block }

func (b *ctrStream) IVSize() int                       { return b.BlockSize() }
func (b *ctrStream) Encrypter(iv []byte) cipher.Stream { return cipher.NewCTR(b, iv) }
func (b *ctrStream) Decrypter(iv []byte) cipher.Stream { return b.Encrypter(iv) }

func aesCTR(key []byte) (stream.Cipher, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &ctrStream{blk}, nil
}

type cfbStream struct{ cipher.Block }

func (b *cfbStream) IVSize() int                       { return b.BlockSize() }
func (b *cfbStream) Encrypter(iv []byte) cipher.Stream { return cipher.NewCFBEncrypter(b, iv) }
func (b *cfbStream) Decrypter(iv []byte) cipher.Stream { return cipher.NewCFBDecrypter(b, iv) }

func aesCFB(key []byte) (stream.Cipher, error) {
	blk, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return &cfbStream{blk}, nil
}

type chacha20ietfkey []byte

func (k chacha20ietfkey) IVSize() int { return chacha20.INonceSize }
func (k chacha20ietfkey) Encrypter(iv []byte) cipher.Stream {
	cs, err := chacha20.New(k, iv)
	if err != nil {
		panic(err)
	}
	return cs
}
func (k chacha20ietfkey) Decrypter(iv []byte) cipher.Stream { return k.Encrypter(iv) }
