/*
 * @Date: 2022.03.02 11:28
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 11:28
 */

package cipher

import (
	"strings"

	netp "efp/net"
)

func PickCipher(cipherType string, key []byte) (netp.ConnCipher, error) {
	cipherType = strings.ToLower(cipherType)

	if choice, ok := aeadList[cipherType]; ok {
		if len(key) != choice.KeySize {
			return nil, ErrKeySize
		}
		aead, err := choice.New(key)
		return aeadConnCipher(aead), err
	}

	if choice, ok := streamList[cipherType]; ok {
		if len(key) != choice.KeySize {
			return nil, ErrKeySize
		}
		cipher, err := choice.New(key)
		return streamConnCipher(cipher), err
	}

	if cipherType == "dummy" {
		return dummyConnCipher(), nil
	}

	return nil, ErrCipherNotSupported
}
