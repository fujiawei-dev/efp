/*
 * @Date: 2022.03.03 14:13
 * @Description: Omit
 * @LastEditors: Rustle Karl
 * @LastEditTime: 2022.03.03 14:13
 */

package tests

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestTCPTunnel(t *testing.T) {
	client := &http.Client{}
	defer client.CloseIdleConnections()

	request, err := http.NewRequest(
		"GET",
		"http://localhost:80/s?wd=OK",
		nil,
	)
	// 否则有可能对方服务器拒绝访问
	request.Host = "www.baidu.com"

	response, err := client.Do(request)
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
