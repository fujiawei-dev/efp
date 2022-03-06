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

func RequestWithProxy(urlString string) ([]byte, error) {
	proxyUrl, err := url.Parse(urlString)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return body, nil
}

func TestSOCKS5Proxy(t *testing.T) {
	body, err := RequestWithProxy("socks5://localhost:1080")
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%s", body)
	}
}
