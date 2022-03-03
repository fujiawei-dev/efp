/*
 * @Date: 2022.03.02 14:36
 * @Description: Stream Cipher.
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 14:36
 */

package stream

import "crypto/cipher"

// Cipher generates a pair of stream ciphers for encryption and decryption.
type Cipher interface {
	IVSize() int
	Encrypter(iv []byte) cipher.Stream
	Decrypter(iv []byte) cipher.Stream
}
