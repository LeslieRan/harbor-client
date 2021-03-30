package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var config_path = "/configurations"

type Parameter struct {
	DailyTime int `json:"daily_time"`
}

type ScanAllPolicy struct {
	// The type of scan all policy, currently the valid values are "none" and
	// "daily"
	Type string `json:"type"`
	// The parameters of the policy, the values are dependant on the type of the
	// policy.
	Parameter *Parameter `json:"parameter"`
}

type Configuration struct {
	// The auth mode of current system, such as "db_auth", "ldap_auth"
	//
	// Please chose "db_auth".
	AuthMode string `json:"auth_mode"`
	// Whether the Harbor instance supports self-registration. If it's set to false,
	// admin need to add user to the instance.   //
	// Please set it to `false`.
	SelfRegistration bool `json:"self_registration"`
	// The default storage quota for the new created projects.
	StoragePerProject string `json:"storage_per_project"`
	//ScanAllPolicy     *ScanAllPolicy `json:"scan_all_policy"`
	// 'docker push' is prohibited by Harbor if you set it to true.
	//
	// Please it to `false`.
	ReadOnly bool `json:"read_only"`
	// This attribute indicates whether quota per project enabled in harbor
	QuotaPerProjectEnable bool `json:"quota_per_project_enable"`
	// This attribute restricts what users have the permission to create project. It
	// can be "everyone" or "adminonly".
	//
	// Please set to `adminonly`.
	ProjectCreationRestriction string `json:"project_creation_restriction"`
	// Whether or not the certificate will be verified when Harbor tries to access the
	// email server.
	EmailInsecure bool `json:"email_insecure"`
	// The username for authenticate against SMTP server.
	EmailUsername string `json:"email_username"`
	// The default count quota for the new created projects.
	CountPerProject string `json:"count_per_project"`
	// The expiration time of the token for internal Registry, in minutes.
	TokenExpiration int `json:"token_expiration"`
	// When it's set to true the system will access Email server via TLS by default.
	// If it's set to false, it still will handle "STARTTLS" from server side.
	EmailSSL bool `json:"email_ssl"`
	// The port of SMTP server.
	EmailPort int `json:"email_port"`
	// The hostname of SMTP server that sends Email notification.
	EmaiHost string `json:"email_host"`
	// The sender name for Email notification.
	EmailFrom string `json:"email_from"`
}

func UpdateSystemConfigurations(config *Configuration) error {
	r, err := json.Marshal(*config)
	if err != nil {
		return fmt.Errorf("json marshal user: %s", err.Error())
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, config_path), bytes.NewBuffer(r))
	if err != nil {
		return fmt.Errorf("new request: %s", err.Error())
	}
	// set headers
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(basic_auth.username, basic_auth.password) // "Authorization: Basic base64([username]:[password])"

	// send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %s", err.Error())
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
		return fmt.Errorf("update system configurations failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}
