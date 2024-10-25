package e2e_test

import (
	"agenda-kaki-go/tests/e2e"
	"agenda-kaki-go/tests/lib"
	"testing"
)

func TestBranchFlow(t *testing.T) {
	branch := &e2e.BaseE2EActions{T: t}
	postBody := map[string]interface{}{"name": lib.GenerateRandomName("Branch")}
	patchBody := map[string]interface{}{"name": lib.GenerateRandomName("Branch")}
	branch.GenerateTesters(5, "branch", "branch", postBody, patchBody).RunAll()
}

// func (b *Branch) Init() *branchActios {
// 	return &branchActios{Branch: b}
// }

// func (b *branchActios) LoadTester(n int) *branchActios {
// 	for i := 0; i < n; i++ {
// 		branch := handlers.Tester{
// 			Entity:  "branch",
// 			BaseURL: namespace.GeneralKey.BaseURL,
// 			PostBody: map[string]interface{}{
// 				"name": lib.GenerateRandomName("Branch"),
// 			},
// 			PatchBody: map[string]interface{}{"name": lib.GenerateRandomName("Branch")},
// 		}
// 		b.Testers = append(b.Testers, branch)
// 	}
// 	return b
// }

// func TestBranchFlow(t *testing.T) {
// 	t.Run("CreateCompanyType", companyType.ExpectStatus(201).POST)
// 	company.PostBody["company_types"] = []map[string]interface{}{
// 		{"id": companyType.EntityID, "name": companyType.PostBody["name"]},
// 	}
// 	t.Run("CreateCompany", company.ExpectStatus(201).POST)
// 	branch.RelatedPath = fmt.Sprintf("company/%d/branch", company.EntityID)
// 	t.Run("GetCompany", company.ExpectStatus(200).GET)
// 	t.Run("CreateBranch", branch.ExpectStatus(201).POST)
// 	t.Run("UpdateBranch", branch.ExpectStatus(200).PATCH)
// 	t.Run("GetBranch", branch.ExpectStatus(200).GET)
// 	t.Run("DeleteCompany", company.ExpectStatus(204).DELETE)
// 	t.Run("GetBranch", branch.ExpectStatus(400).GET)
// 	t.Run("DeleteBranch", branch.ExpectStatus(400).DELETE)
// 	t.Run("ForceGetBranch", branch.ExpectStatus(400).ForceGET)
// 	t.Run("ForceDeleteCompany", company.ExpectStatus(204).ForceDELETE)
// 	t.Run("ForceDeleteBranch", branch.ExpectStatus(400).ForceDELETE)
// 	t.Run("CreateCompany", company.ExpectStatus(201).POST)
// 	t.Run("GetCompany", company.ExpectStatus(200).GET)
// 	branch.RelatedPath = fmt.Sprintf("company/%d/branch", company.EntityID)
// 	t.Run("CreateBranch", branch.ExpectStatus(201).POST)
// 	t.Run("GetBranch", branch.ExpectStatus(200).GET)
// 	t.Run("UpdateBranch", branch.ExpectStatus(200).PATCH)
// 	t.Run("DeleteBranch", branch.ExpectStatus(204).DELETE)
// 	t.Run("GetBranch", branch.ExpectStatus(404).GET)
// 	t.Run("ForceGetBranch", branch.ExpectStatus(200).ForceGET)
// 	t.Run("ForceDeleteBranch", branch.ExpectStatus(204).ForceDELETE)
// 	t.Run("ForceDeleteCompany", company.ExpectStatus(204).ForceDELETE)
// 	t.Run("ForceDeleteCompanyType", companyType.ExpectStatus(204).ForceDELETE)
// }
