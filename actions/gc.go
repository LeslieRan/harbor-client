package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	gc_schedule_path = "/system/gc/schedule"
)

type ScheduleType string

const (
	Hourly ScheduleType = "Hourly"
	Daily  ScheduleType = "Daily"
	Weekly ScheduleType = "Weekly"
	Custom ScheduleType = "Custom"
	Manual ScheduleType = "Manual"
)

type Schedule struct {
	Cron string       `json:"cron"`
	Type ScheduleType `json:"type"`
}

type GCScheduleRequestBody struct {
	Parameters map[string]interface{} `json:"parameters"`
	Schedule   *Schedule              `json:"schedule"`
}

func CreateGCSchedule(schedule *GCScheduleRequestBody) error {
	r, err := json.Marshal(*schedule)
	if err != nil {
		return fmt.Errorf("json marshal user: %s", err.Error())
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, gc_schedule_path), bytes.NewBuffer(r))
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
		return fmt.Errorf("create schedule failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}

func UpdateGCSchedule(schedule *GCScheduleRequestBody) error {
	r, err := json.Marshal(*schedule)
	if err != nil {
		return fmt.Errorf("json marshal user: %s", err.Error())
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, gc_schedule_path), bytes.NewBuffer(r))
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
		return fmt.Errorf("update schedule failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}
