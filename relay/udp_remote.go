/*
 * @Date: 2022.03.02 11:12
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 11:12
 */

package relay

import (
	"net"

	netp "efp/net"
	"efp/net/socks"
)

// NewUDPRemoteProxyServer Listen on addr for encrypted packets and basically do UDP NAT.
func NewUDPRemoteProxyServer(addr string, cipher netp.PacketConnCipher) {
	c, err := netp.ListenPacket("udp", addr, cipher)
	if err != nil {
		logf("UDP remote listen error, %v", err)
		return
	}
	defer c.Close()

	nm := NewNATMap(UDPTimeout)
	buf := make([]byte, UDPBufferSize)

	logf("listening UDP on %s", addr)

	for {
		n, raddr, err := c.ReadFrom(buf)
		if err != nil {
			logf("UDP remote read error, %v", err)
			continue
		}

		tgtAddr := socks.SplitAddr(buf[:n])
		if tgtAddr == nil {
			logf("failed to split target address from packet: %q", buf[:n])
			continue
		}

		tgtUDPAddr, err := net.ResolveUDPAddr("udp", tgtAddr.String())
		if err != nil {
			logf("failed to resolve target UDP address, %v", err)
			continue
		}

		payload := buf[len(tgtAddr):n]

		pc := nm.Get(raddr.String())
		if pc == nil {
			pc, err = net.ListenPacket("udp", "")
			if err != nil {
				logf("UDP remote listen error, %v", err)
				continue
			}

			nm.Add(raddr, c, pc)
		}

		_, err = pc.WriteTo(payload, tgtUDPAddr) // accept only UDPAddr despite the signature
		if err != nil {
			logf("UDP remote write error, %v", err)
			continue
		}
	}
}
