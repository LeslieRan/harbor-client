package actions

import (
	"fmt"
	"testing"
)

var (
	testUser = &User{
		Username:     "leslie-test",
		Password:     "Rf123456",
		Email:        "1234@google.com",
		Realname:     "LeslieRan",
		Comment:      "Aloha, harbor",
		SysAdminFlag: false,
	}
)

// @PASS
func TestUserCreate(t *testing.T) {
	if err := CreateUser(testUser); err != nil {
		t.Error(err.Error())
		return
	}
	fmt.Printf("create user %s successfully.", testUser.Username)
}

// @PASS
func TestUserSearch(t *testing.T) {
	users, err := SearchUsers(testUser.Username, 1, 2)
	if err != nil {
		t.Errorf("search user %s: %s", testUser.Username, err.Error())
		return
	}
	if len(users) > 0 {
		fmt.Printf("users: %v", users[0])
	}
}

// @PASS
func TestUserList(t *testing.T) {
	users, err := ListUsers(testUser.Username, "", 1, 1)
	if err != nil {
		t.Errorf("list user %s: %s", testUser.Username, err.Error())
		return
	}
	if len(users) > 0 {
		fmt.Printf("users: %v", users[0])
	}
}

// @PASS
func TestChangePassword(t *testing.T) {
	if err := ChangePassword(testUser.Username, testUser.Password, "Rf452342"); err != nil {
		t.Errorf("%s", err.Error())
		return
	}
	fmt.Printf("change password of user %s successfully.", testUser.Username)
}

// @PASS
func TestUserDelete(t *testing.T) {
	if err := DeleteUser(testUser.Username); err != nil {
		t.Errorf("delete user %s: %s", testUser.Username, err.Error())
		return
	}
	fmt.Printf("delete user %s successfully.", testUser.Username)
}
