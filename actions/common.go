package actions

import (
	"crypto/tls"
	"net/http"
	//"net/http/cookiejar"
)

const harbor_server_domain = "https://core.harbor.domain"

const v2api_path = "/api/v2.0"

var (
	// basic auth for server
	basic_auth = struct {
		username string
		password string
	}{
		username: "admin",
		password: "Harbor12345",
	}

	client *http.Client
)

// Response body returned from harbor sever on error occasions.
type ErrorBody struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func init() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client = &http.Client{
		Transport: tr,
	}

}
