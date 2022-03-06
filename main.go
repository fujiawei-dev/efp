/*
 * @Date: 2022.03.02 11:35
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 11:35
 */

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"efp/cipher"
	"efp/relay"
)

var flags struct {
	Cipher    string
	Client    string
	Key       string
	Key64     string
	Keygen    int
	Password  string
	Server    string
	Socks     string
	TCPTunnel string
	Verbose   bool
}

func main() {
	flag.StringVar(&flags.Cipher, "cipher", "", "cipher to encrypt/decrypt")
	flag.StringVar(&flags.Client, "c", "", "client connect address")
	flag.StringVar(&flags.Key, "key", "", "secret key in hexadecimal (derive from key64 or password if empty)")
	flag.StringVar(&flags.Key64, "key64", "", "base64url-encoded key")
	flag.IntVar(&flags.Keygen, "keygen", 0, "generate a base64url-encoded random key of given length in byte")
	flag.StringVar(&flags.Password, "password", "", "password")
	flag.StringVar(&flags.Server, "s", "", "server listen address")
	flag.StringVar(&flags.Socks, "socks", ":1080", "(client-only) SOCKS listen address")
	flag.StringVar(&flags.TCPTunnel, "tcptunnel", "", "(client-only) TCP tunnel (laddr1=raddr1,laddr2=raddr2,...)")
	flag.BoolVar(&flags.Verbose, "v", false, "verbose mode")

	flag.Parse()

	if flags.Keygen > 0 {
		key := make([]byte, flags.Keygen)
		if _, err := rand.Read(key); err != nil {
			log.Fatalf("failed to generate key, %v", err)
		}
		fmt.Println(base64.URLEncoding.EncodeToString(key))
		return
	}

	if flags.Cipher == "" {
		cipher.PrintCiphers(os.Stderr)
		return
	}

	if flags.Client == "" && flags.Server == "" {
		flag.Usage()
		return
	}

	var (
		key []byte
		err error
	)

	if flags.Key != "" {
		key, err = hex.DecodeString(flags.Key)
		if err != nil {
			log.Fatalf("failed to parse key, %v", err)
		}
	} else if flags.Key64 != "" {
		key, err = base64.URLEncoding.DecodeString(flags.Key64)
		if err != nil {
			log.Fatalf("failed to parse key, %v", err)
		}
	}

	connCipher, err := cipher.New(flags.Cipher, key, flags.Password)
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
