/*
 * @Date: 2022.03.02 11:27
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 11:27
 */

package cipher

import (
	"net"

	netp "efp/net"
)

// Dummy ciphers (no encryption)

func dummyConnCipher() netp.ConnCipher {
	return func(c net.Conn) net.Conn { return c }
}
