package actions

import (
	"testing"
	//"ldblock.com/harbor-client/actions"
)

var (
	admin_username = "admin"
	admin_pwd      = "Harbor12345"
)

// @PASS
func TestLogin(t *testing.T) {
	if err := Login(admin_username, admin_pwd); err != nil {
		t.Error(err.Error())
	}
}
