package e2e_test

import (
	"agenda-kaki-go/tests/handlers"
	"agenda-kaki-go/tests/lib"
	"fmt"
	"testing"
)

var branch = handlers.Tester{
	Entity:  "branch",
	BaseURL: baseURL,
	PostBody: map[string]interface{}{
		"name": lib.GenerateRandomName("Branch"),
	},
	PatchBody: map[string]interface{}{"name": lib.GenerateRandomName("Branch")},
}

func TestBranchFlow(t *testing.T) {
	t.Run("CreateCompanyType", companyType.ExpectStatus(201).POST)
	company.PostBody["company_types"] = []map[string]interface{}{
		{"id": companyType.EntityID, "name": companyType.PostBody["name"]},
	}
	t.Run("CreateCompany", company.ExpectStatus(201).POST)
	branch.RelatedPath = fmt.Sprintf("company/%d/branch", company.EntityID)
	t.Run("CreateBranch", branch.ExpectStatus(201).POST)
	t.Run("UpdateBranch", branch.ExpectStatus(200).PATCH)
	t.Run("GetBranch", branch.ExpectStatus(200).GET)
	t.Run("DeleteBranch", branch.ExpectStatus(204).DELETE)
	t.Run("GetBranchWith404", branch.ExpectStatus(404).GET)
	t.Run("ForceGetBranch", branch.ExpectStatus(200).ForceGET)
	t.Run("ForceDeleteBranch", branch.ExpectStatus(204).ForceDELETE)
	t.Run("ForceDeleteCompany", company.ExpectStatus(204).ForceDELETE)
	t.Run("ForceDeleteCompanyType", companyType.ExpectStatus(204).ForceDELETE)
}
