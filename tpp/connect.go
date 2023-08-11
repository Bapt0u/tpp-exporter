package tpp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type oauthReturn struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type VenafiTpp struct {
	AccessToken  string
	Url          string
	ExpireSoon   int
	Certificates []Certificate
}

func (v *VenafiTpp) Connect(conf *Config) {
	var oauthResponse oauthReturn

	v.Url = conf.VenafiTpp.Url

	// Format json payload
	body := fmt.Sprintf(`
		{"username": "%s", "password": "%s", "client_id": "%s", "scope": "%s"}`,
		conf.VenafiTpp.Username, conf.VenafiTpp.Password,
		conf.VenafiTpp.ClientId, conf.VenafiTpp.Scope,
	)

	// Create a Post request to get Access_token
	req, err := http.NewRequest(
		"POST",
		string(conf.VenafiTpp.Url+"/vedauth/authorize/oauth"),
		bytes.NewBuffer([]byte(body)),
	)
	CheckErr(err)
	req.Header.Add("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		CheckErr(err)
	}

	// Read response from API
	response, err := io.ReadAll(resp.Body)
	CheckErr(err)
	err = json.Unmarshal(response, &oauthResponse)
	CheckErr(err)

	resp.Body.Close()
	v.AccessToken = oauthResponse.AccessToken
	v.ExpireSoon = conf.VenafiTpp.ExpireSoon
}

func (v *VenafiTpp) GetStatus() {
	resp, err := http.Get(string(v.Url + "/vedsdk/ServerStatus/Status"))
	CheckErr(err)

	// Read response from API
	response, err := io.ReadAll(resp.Body)
	CheckErr(err)
	log.Print(resp.StatusCode)
	log.Printf("%s ; %s (Status code %d)", v.Url, string(response), resp.StatusCode)

	resp.Body.Close()
}
