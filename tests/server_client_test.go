/*
 * @Date: 2022.03.06 11:53
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.06 11:53
 */

package tests

import (
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/fujiawei-dev/efp/cipher"
	"github.com/fujiawei-dev/efp/relay"
)

func StartServerAndClient(clientAddr, serverAddr, cipherType string, key []byte, password string) error {
	connCipher, err := cipher.New(cipherType, key, password)
	if err != nil {
		return err
	}

	go relay.NewSOCKS5ProxyClient(clientAddr, serverAddr, connCipher)
	go relay.NewTCPRemoteProxyServer(serverAddr, connCipher)

	return nil
}

func TestServerAndClientWithKey(t *testing.T) {
	var (
		cipherType = "aes-128-gcm"
		keyString  = "1234567890abcdef1234567890abcdef"
		clientAddr = ":1080"
		serverAddr = ":8488"

		key  []byte
		err  error
		body []byte
	)

	if key, err = hex.DecodeString(keyString); err != nil {
		t.Error(err)
		return
	}

	if err = StartServerAndClient(clientAddr, serverAddr, cipherType, key, ""); err != nil {
		t.Error(err)
		return
	}

	if body, err = RequestWithProxy("socks5://" + clientAddr); err != nil {
		t.Error(err)
	} else {
		t.Logf("%s", body)
	}
}

func TestServerAndClientWithKey64(t *testing.T) {
	var (
		cipherType  = "aes-128-gcm"
		key64String = "mUOlCtqTTE9r4qkNNXkpoA=="
		clientAddr  = ":1080"
		serverAddr  = ":8488"

		key  []byte
		err  error
		body []byte
	)

	// use key64
	{
		if key, err = base64.URLEncoding.DecodeString(key64String); err != nil {
			t.Error(err)
			return
		}

		if err = StartServerAndClient(clientAddr, serverAddr, cipherType, key, ""); err != nil {
			t.Error(err)
			return
		}

		if body, err = RequestWithProxy("socks5://" + clientAddr); err != nil {
			t.Error(err)
		} else {
			t.Logf("%s", body)
		}
	}
}

func TestServerAndClientWithPassword(t *testing.T) {
	var (
		cipherType = "aes-128-gcm"
		password   = "password"
		clientAddr = ":1080"
		serverAddr = ":8488"

		key  []byte
		err  error
		body []byte
	)

	// use key64
	{
		if err = StartServerAndClient(clientAddr, serverAddr, cipherType, key, password); err != nil {
			t.Error(err)
			return
		}

		if body, err = RequestWithProxy("socks5://" + clientAddr); err != nil {
			t.Error(err)
		} else {
			t.Logf("%s", body)
		}
	}
}
