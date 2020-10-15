package net

import (
	"io/ioutil"
	"net/http"

	"github.com/mapkn3/gompose/util"
)

// DoRequestWithBasicAuth does request and returns response like byte array
func DoRequestWithBasicAuth(url string, username string, password string) []byte {
	req, err := http.NewRequest("GET", url, nil)
	util.Check(err)
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(req)
	util.Check(err)
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	util.Check(err)
	return body
}
