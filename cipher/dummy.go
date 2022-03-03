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

func dummyConnCipher() netp.ConnCipher {
	return func(c net.Conn) net.Conn { return c }
}

func dummyPacketConnCipher() netp.PacketConnCipher {
	return func(c net.PacketConn) net.PacketConn { return c }
}
