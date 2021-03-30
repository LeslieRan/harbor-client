package actions

import (
	"fmt"
	"testing"
)

var (
	projectName = "project-leslieran-test"
	memberName  = "leslie-test"
)

// @PASS
func TestCreateProject(t *testing.T) {
	project := &CreateProjectRequestBody{
		Name:         projectName,
		StorageLimit: 2 * GB,
		MetaData:     new(MetaData),
	}
	project.MetaData.Public = "false"
	project.MetaData.AutoScan = "true"
	project.MetaData.PreventVul = "true"
	if err := CreateProject(project); err != nil {
		t.Errorf("create project %s: %s", projectName, err.Error())
		return
	}
	fmt.Printf("create project %s successfully.", projectName)
}

// @PASS
func TestSearchProject(t *testing.T) {
	project, err := SeachProject(projectName)
	if err != nil {
		t.Errorf("search project %s: %s", projectName, err.Error())
		return
	}
	fmt.Println("project info: ", *project)
}

// @PASS
func TestDeleteProject(t *testing.T) {
	if err := DeleteProject(projectName); err != nil {
		t.Errorf("delete project %s: %s", projectName, err.Error())
		return
	}
	fmt.Printf("delete project %s successfully.", projectName)
}

// @PASS
func TestCreateMember(t *testing.T) {
	err := CreateMember(projectName, memberName, Developer)
	if err != nil {
		t.Errorf("create member: %s", err.Error())
		return
	}
	fmt.Println("create member successfully.")
}

// @PASS
func TestSearchMember(t *testing.T) {
	members, err := SearchMember(projectName, memberName)
	if err != nil {
		t.Errorf("search member %s of project %s: %s\n", memberName, projectName, err.Error())
		return
	}
	if len(members) > 0 {
		fmt.Println("member info: ", *members[0])
	}
}

// @PASS
func TestDeleteMember(t *testing.T) {
	if err := DeleteMember(projectName, memberName); err != nil {
		t.Errorf("delete member: %s", err.Error())
		return
	}
	fmt.Println("delete member successfully.")
}
