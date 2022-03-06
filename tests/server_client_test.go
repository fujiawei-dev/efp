/*
 * @Date: 2022.03.06 11:53
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.06 11:53
 */

package tests

import (
	"efp/cipher"
	"efp/relay"
	"encoding/hex"
	"testing"
)

func TestServerAndClient(t *testing.T) {
	cipherType := "aes-128-gcm"
	keyString := "1234567890abcdef1234567890abcdef"
	clientAddr := ":1080"
	serverAddr := ":8488"

	key, err := hex.DecodeString(keyString)
	if err != nil {
		t.Error(err)
		return
	}

	connCipher, err := cipher.New(cipherType, key)
	if err != nil {
		t.Error(err)
		return
	}

	go relay.NewSOCKS5ProxyClient(clientAddr, serverAddr, connCipher)
	go relay.NewTCPRemoteProxyServer(serverAddr, connCipher)

	body, err := RequestWithProxy("socks5://" + clientAddr)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%s", body)
	}
}
