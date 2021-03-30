package actions

import (
	"testing"
)

// @PASS
func TestCreateGCSchedule(t *testing.T) {
	schedule := &GCScheduleRequestBody{
		Parameters: map[string]interface{}{
			"delete_untagged": "true",
		},
		Schedule: &Schedule{
			Type: "Hourly",
			Cron: "0 0 * * * *",
		},
	}
	if err := CreateGCSchedule(schedule); err != nil {
		t.Errorf("create schedule: %s", err.Error())
		return
	}
	t.Logf("create schedule successfully.")
}

// @PASS
func TestUpdateGCSchedule(t *testing.T) {
	schedule := &GCScheduleRequestBody{
		Parameters: map[string]interface{}{
			"delete_untagged": false,
		},
		Schedule: &Schedule{
			Type: "Weekly",
			Cron: "0 0 0 * * *",
		},
	}
	if err := UpdateGCSchedule(schedule); err != nil {
		t.Errorf("update schedule: %s", err.Error())
		return
	}
	t.Logf("update schedule successfully.")
}
