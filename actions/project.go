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

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
)

var (
	placeholder_project = "{project_name_or_id}"
	placeholder_member  = "{mid}"
	projectPaths        = struct {
		create        string
		search        string
		del           string
		create_member string
		search_member string
		del_member    string
	}{
		create:        "/projects",
		search:        "/projects/" + placeholder_project,
		del:           "/projects/" + placeholder_project,                                    // only project name
		create_member: "/projects/" + placeholder_project + "/members",                       // only project id
		search_member: "/projects/" + placeholder_project + "/members",                       // only project id
		del_member:    "/projects/" + placeholder_project + "/members/" + placeholder_member, // only id for project
	}
)

type CreateProjectRequestBody struct {
	Name         string    `json:"project_name"`
	StorageLimit int64     `json:"storage_limit"`
	MetaData     *MetaData `json:"metadata"`

	// other information you care about...
}

type MetaData struct {
	AutoScan   string `json:"auto_scan"`
	Public     string `json:"public"`
	PreventVul string `json:"prevent_vul"`
}

type SearchProjectResponseBody struct {
	Name      string    `json:"name"`
	ProjectID int       `json:"project_id"`
	OwnerName string    `json:"owner_name"`
	OwnerID   int       `json:"owner_id"`
	Deleted   bool      `json:"deleted"`
	RepoCount int       `json:"repo_count"`
	MetaData  *MetaData `json:"metadata"`
}

func CreateProject(project *CreateProjectRequestBody) error {
	r, err := json.Marshal(project)
	if err != nil {
		return fmt.Errorf("json marshal project: %s", err.Error())
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, projectPaths.create), bytes.NewBuffer(r))
	if err != nil {
		return fmt.Errorf("new request: %s", err.Error())
	}
	// set headers
	req.Header.Set("X-Resource-Name-In-Location", "false")
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
		return fmt.Errorf("create project failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}

func SeachProject(projectName string) (*SearchProjectResponseBody, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, strings.Replace(projectPaths.search, placeholder_project, projectName, -1)), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %s", err.Error())
	}
	// set headers
	req.Header.Set("X-Is-Resource-Name", "true")               // NOTE
	req.SetBasicAuth(basic_auth.username, basic_auth.password) // "Authorization: Basic base64([username]:[password])"

	// send request
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
		return nil, fmt.Errorf("search project failed. The response code is %v %s", resp.StatusCode, errString)
	}

	var project *SearchProjectResponseBody
	// read response body
	if resp.Body != nil {
		cb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read data from response body: %s", err.Error())
		}
		if resp.Header.Get("Content-Type") != "application/json" {
			return nil, fmt.Errorf("data in body is not in json type")
		}
		err = json.Unmarshal(cb, &project)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data in body: %s", err.Error())
		}
	}

	return project, nil
}

func DeleteProject(project_name string) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, strings.Replace(projectPaths.del, placeholder_project, project_name, -1)), nil)
	if err != nil {
		return fmt.Errorf("new request: %s", err.Error())
	}
	// set headers
	req.Header.Set("X-Is-Resource-Name", "true")               // NOTE
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
		return fmt.Errorf("delete project failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}

type MemberRoleType int

const (
	ProjectAdmin MemberRoleType = 1
	Developer    MemberRoleType = 2
	Guest        MemberRoleType = 3
	Maintainer   MemberRoleType = 4
	LimitedGues  MemberRoleType = 5
)

type CreateMemberRequestBody struct {
	RoleID     int         `json:"role_id"`
	MemberUser *MemberUser `json:"member_user"`

	// other information you care about...
}

type MemberUser struct {
	UserName string `json:"username"`
	UserID   int    `json:"user_id"`
}

func CreateMember(projectName, username string, role MemberRoleType) error {
	// get project id by project name.
	project, err := SeachProject(projectName)
	if err != nil {
		return fmt.Errorf("search for project %s: %s", projectName, err.Error())
	}

	member := &CreateMemberRequestBody{
		RoleID:     int(role),
		MemberUser: new(MemberUser),
	}
	member.MemberUser.UserName = username

	r, err := json.Marshal(member)
	if err != nil {
		return fmt.Errorf("json marshal project: %s", err.Error())
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, strings.Replace(projectPaths.create_member, placeholder_project, strconv.Itoa(project.ProjectID), -1)), bytes.NewBuffer(r))
	if err != nil {
		return fmt.Errorf("new request: %s", err.Error())
	}
	// set headers.
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
		return fmt.Errorf("create project member failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}

type Member struct {
	ID         int    `json:"id"`
	ProjectID  int    `json:"project_id"`
	EntityID   int    `json:"entity_id"`
	EntityName string `json:"entity_name"`
	EntityType string `json:"entity_type"`
	RoleID     int    `json:"role_id"`
	RoleName   string `json:"role_name"`
}

func SearchMember(projectName, username string) ([]*Member, error) {
	// get project id first
	project, err := SeachProject(projectName)
	if err != nil {
		if err != nil {
			return nil, fmt.Errorf("search members %s: %s", projectName, err.Error())
		}
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, strings.Replace(projectPaths.search_member, placeholder_project, strconv.Itoa(project.ProjectID), -1)), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %s", err.Error())
	}
	// build parameters
	if len(username) > 0 {
		q := url.Values(make(map[string][]string, 1))
		q.Add("entity", username)
		req.URL.RawPath = q.Encode()
	}
	fmt.Printf("[DEBUG]url of searching members is: %s\n", req.URL.String())

	// set headers.
	req.SetBasicAuth(basic_auth.username, basic_auth.password) // "Authorization: Basic base64([username]:[password])"

	// send request
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
		return nil, fmt.Errorf("search members failed. The response code is %v %s", resp.StatusCode, errString)
	}

	var members []*Member
	// read response body
	if resp.Body != nil {
		cb, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read data from response body: %s", err.Error())
		}
		// if resp.Header.Get("Content-Type") != "application/json" {
		// 	return nil, fmt.Errorf("data in body is not in json type")
		// }
		err = json.Unmarshal(cb, &members)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data in body: %s", err.Error())
		}
	}

	return members, nil
}

func DeleteMember(projectName, username string) error {
	// get user id by user name.
	members, err := SearchMember(projectName, username)
	if err != nil {
		return fmt.Errorf("search the member %s: %s", username, err.Error())
	}
	if len(members) == 0 {
		return fmt.Errorf("user %s is not existing.", username)
	} else if len(members) > 1 { // hardly
		return fmt.Errorf("the number of user %s is more than one.", username)
	}
	// get project id by project name.
	project, err := SeachProject(projectName)
	if err != nil {
		return fmt.Errorf("search for project %s: %s", projectName, err.Error())
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s%s", harbor_server_domain, v2api_path, strings.Replace(strings.Replace(projectPaths.del_member, placeholder_project, strconv.Itoa(project.ProjectID), -1), placeholder_member, strconv.Itoa(members[0].ID), -1)), nil)
	if err != nil {
		return fmt.Errorf("new request: %s", err.Error())
	}

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
		return fmt.Errorf("create project member failed. The response code is %v %s", resp.StatusCode, errString)
	}

	return nil
}
