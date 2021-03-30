package actions

import (
	"strconv"
	"testing"
)

// @PASS
func TestUpdateSystemConfigurations(t *testing.T) {
	config := &Configuration{}
	config.AuthMode = "db_auth"
	config.SelfRegistration = false
	config.ReadOnly = false
	config.ProjectCreationRestriction = "adminonly"
	config.StoragePerProject = strconv.Itoa(2 * GB)
	config.CountPerProject = "-1"
	config.TokenExpiration = 60
	if err := UpdateSystemConfigurations(config); err != nil {
		t.Errorf("update configurations: %s", err.Error())
		return
	}
	t.Logf("update configurations successfully.")
}
