/*
 * @Date: 2022.03.02 14:01
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.02 14:01
 */

package tests

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestSOCK5Proxy(t *testing.T) {
	proxyUrl, err := url.Parse("socks5://localhost:1080")
	if err != nil {
		t.Error(err)
		return
	}

	proxyClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	defer proxyClient.CloseIdleConnections()
	// In principle, it should be actively closed, otherwise the server will throw an error:
	// an existing connection was forcibly closed by the remote host.

	response, err := proxyClient.Get("https://www.baidu.com")
	if err != nil {
		t.Error(err)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Error(err)
		return
	}

	defer response.Body.Close()

	t.Logf("%s", body)
}
