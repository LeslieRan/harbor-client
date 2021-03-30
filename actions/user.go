package actions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	placeholder_user = "{user_id}"

	userPaths = struct {
		create     string
		search     string
		list       string
		del        string
		change_pwd string
	}{
		create:     "/users",
		search:     "/users/search",
		list:       "/users",
		del:        "/users/" + placeholder_user,
		change_pwd: "/users/" + placeholder_user + "/password",
	}
)

// User holds the details of a user.
type User struct {
	UserID       int    `json:"user_id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	Realname     string `json:"realname"`
	Comment      string `json:"comment"`
	SysAdminFlag bool   `json:"sysadmin_flag"`

	// other information you care about...
}

func CreateUser(user *User) error {
	r, err := json.Marshal(*user)
	if err != nil {
		return fmt.Errorf("json marshal user: %s", err.Error())
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, userPaths.create), bytes.NewBuffer(r))
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
		return fmt.Errorf("create user failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}

// NOTE page and page_size are both > 0. If you're searching for single user,
// and then you should set page to 1.

func SearchUsers(username string, page, page_size int32) ([]*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, userPaths.search), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %s", err.Error())
	}
	// build parameters in url
	q := url.Values(make(map[string][]string, 3))
	q.Add("username", username)
	q.Add("page", strconv.Itoa(int(page)))
	q.Add("page_size", strconv.Itoa(int(page_size)))
	req.URL.RawQuery = q.Encode()
	fmt.Printf("[DEBUG]search url: %s\n", req.URL.String())
	// set headers
	req.SetBasicAuth(basic_auth.username, basic_auth.password) // "Authorization: Basic base64([username]:[password])"
	// do search
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %s", err.Error())
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
		return nil, fmt.Errorf("search users failed. The response code is %v %s", resp.StatusCode, errString)
	}

	var users []*User
	// read response body
	if resp.Body != nil {
		cb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read data from response body: %s", err.Error())
		}
		// if resp.Header.Get("Content-Type") != "application/json" {
		// 	return nil, fmt.Errorf("data in body is not in json type")
		// }
		usersN, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))
		users = make([]*User, usersN, usersN) // FIXME
		err = json.Unmarshal(cb, &users)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data in body: %s", err.Error())
		}
	}

	return users, nil
}

// NOTE page and page_size are both > 0. If you're searching for single user,
// and then you should set page to 1.

func ListUsers(username, email string, page, page_size int32) ([]*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, userPaths.list), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %s", err.Error())
	}
	// build parameters in url
	q := url.Values(make(map[string][]string, 3))
	q.Add("username", username)
	q.Add("email", email)
	q.Add("page", strconv.Itoa(int(page)))
	q.Add("page_size", strconv.Itoa(int(page_size)))
	req.URL.RawQuery = q.Encode()
	fmt.Printf("[DEBUG]list url: %s\n", req.URL.String())
	// set headers
	req.SetBasicAuth(basic_auth.username, basic_auth.password) // "Authorization: Basic base64([username]:[password])"
	// do search
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %s", err.Error())
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
		return nil, fmt.Errorf("list users failed. The response code is %v %s", resp.StatusCode, errString)
	}

	var users []*User
	// read response body
	if resp.Body != nil {
		cb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read data from response body: %s", err.Error())
		}
		// if resp.Header.Get("Content-Type") != "application/json" {
		// 	return nil, fmt.Errorf("data in body is not in json type")
		// }
		usersN, _ := strconv.Atoi(resp.Header.Get("X-Total-Count"))
		users = make([]*User, usersN, usersN) // FIXME
		err = json.Unmarshal(cb, &users)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data in body: %s", err.Error())
		}
	}

	return users, nil
}

func DeleteUser(username string) error {
	// search the user first.
	users, err := SearchUsers(username, 1, 1)
	if err != nil {
		return fmt.Errorf("search the user %s: %s", username, err.Error())
	}
	if len(users) == 0 {
		return fmt.Errorf("user %s is not existing.", username)
	} else if len(users) > 1 { // hardly
		return fmt.Errorf("the number of user %s is more than one.", username)
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, strings.Replace(userPaths.del, placeholder_user, strconv.Itoa(users[0].UserID), -1)), nil)
	if err != nil {
		return fmt.Errorf("new request: %s", err.Error())
	}
	fmt.Println("delete url: ", req.URL.String())
	// set headers
	req.SetBasicAuth(basic_auth.username, basic_auth.password) // "Authorization: Basic base64([username]:[password])"
	// do search
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
		return fmt.Errorf("delete user failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}

func ChangePassword(username string, old_pwd, new_pwd string) error {
	// get user id by searching the user info.
	users, err := SearchUsers(username, 1, 1)
	if err != nil {
		return fmt.Errorf("search the user %s: %s", username, err.Error())
	}
	if len(users) == 0 {
		return fmt.Errorf("user %s is not existing.", username)
	} else if len(users) > 1 { // hardly
		return fmt.Errorf("the number of user %s is more than one.", username)
	}

	req_body := &struct {
		Old_pwd string `json:"old_password"`
		New_pwd string `json:"new_password"`
	}{
		Old_pwd: old_pwd,
		New_pwd: new_pwd,
	}
	r, err := json.Marshal(*req_body)
	if err != nil {
		return fmt.Errorf("json marshal request body: %s", err.Error())
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, strings.Replace(userPaths.change_pwd, placeholder_user, strconv.Itoa(users[0].UserID), -1)), bytes.NewBuffer(r))
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
		return fmt.Errorf("change the password of user %s failed. The response code is %v %s", username, resp.StatusCode, errString)
	}

	return nil
}
