package lib

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var token string
var client = http.DefaultClient

// Login attempts to login and get API token from worker
func Login(workerURL, username, password string) error {
	body := &bytes.Buffer{}
	buf, _ := json.Marshal(struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{Username: username, Password: password})
	body.Write(buf)
	res, err := client.Post(workerURL, "application/json", body)
	if err != nil {
		return err
	}
	tokenRes := struct {
		Authorization string
	}{}
	respBody, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(respBody, &tokenRes)
	token = tokenRes.Authorization
	return nil
}
