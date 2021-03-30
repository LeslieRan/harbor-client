package actions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"
)

const login_path = "/c/login"

func Login(username, pwd string) error {
	formContent := map[string][]string{
		"principal": []string{username},
		"password":  []string{pwd},
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", harbor_server_domain, login_path), strings.NewReader(encode(formContent)))
	if err != nil {
		return fmt.Errorf("login: %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// login
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("login: %s", err.Error())
	}
	if !(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated) {
		var errString string
		if resp.Body != nil {
			var errors *ErrorBody
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("[WARNING]read error message from body: %s\n", err.Error())
			} else {
				err = json.Unmarshal(body, &errors)
				if err != nil {
					fmt.Printf("[WARNING]unmarshal error body: %s\n", err.Error())
					errString = string(body)
				} else {
					errString = fmt.Sprintf(" and error contents are: %v", *errors)
				}

			}
		}
		return fmt.Errorf("login failed. The response code is %v %s", resp.StatusCode, errString)
	}
	// store cookies.
	// url, err := url.Parse(harbor_server_domain)
	// if err != nil {
	// 	return fmt.Errorf("parse url: %s", err.Error())
	// }
	// client.Jar.SetCookies(url, resp.Cookies())
	fmt.Printf("login successfully.")

	return nil

}

func encode(v map[string][]string) string {
	if v == nil {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(v))
	for k := range v {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		vs := v[k]
		keyEscaped := url.QueryEscape(k)
		for _, v := range vs {
			if buf.Len() > 0 {
				buf.WriteByte('&')
			}
			buf.WriteString(keyEscaped)
			buf.WriteByte('=')
			buf.WriteString(url.QueryEscape(v))
		}
	}
	return buf.String()
}
