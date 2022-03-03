/*
 * @Date: 2022.03.02 11:35
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 11:35
 */

package main

import (
	"encoding/hex"
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"efp/cipher"
	"efp/relay"
)

var flags struct {
	Client    string
	Server    string
	Cipher    string
	Key       string
	Socks     string
	TCPTunnel string
	Verbose   bool
}

func main() {
	flag.StringVar(&flags.Client, "c", "", "client connect address")
	flag.StringVar(&flags.Server, "s", "", "server listen address")
	flag.StringVar(&flags.Cipher, "cipher", "dummy", "cipher to encrypt/decrypt")
	flag.StringVar(&flags.Key, "key", "", "secret key in hexadecimal")
	flag.StringVar(&flags.Socks, "socks", ":1080", "(client-only) SOCKS listen address")
	flag.StringVar(&flags.TCPTunnel, "tcptunnel", "", "(client-only) TCP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	flag.BoolVar(&flags.Verbose, "v", false, "verbose mode")

	flag.Parse()

	if flags.Client == "" && flags.Server == "" {
		flag.Usage()
		return
	}

	key, err := hex.DecodeString(flags.Key)
	if err != nil {
		log.Fatalf("failed to parse key, %v", err)
	}

	connCipher, err := cipher.PickCipher(flags.Cipher, key)
	if err != nil {
		log.Fatalf("failed to create cipher %s, %v", flags.Cipher, err)
	}

	relay.SetVerboseMode(flags.Verbose)

	// Proxy Client
	if flags.Client != "" {
		if flags.TCPTunnel != "" {
			for _, tun := range strings.Split(flags.TCPTunnel, ",") {
				p := strings.Split(tun, "=")
				go relay.NewTCPTunnel(p[0], flags.Client, p[1], connCipher)
			}
		}

		if flags.Socks != "" {
			go relay.NewSOCKS5ProxyClient(flags.Socks, flags.Client, connCipher)
		}
	}

	// Proxy Server
	if flags.Server != "" {
		go relay.NewTCPRemoteProxyServer(flags.Server, connCipher)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
