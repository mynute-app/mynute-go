package main

import (
	"fmt"
	"io"
	"mynute-go/core/src/config/db/model"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		// Public schema models (GeneralModels)
		&model.Sector{},
		&model.Company{},
		&model.Holiday{},
		&model.Client{},
		&model.EndPoint{},
		&model.Role{},
		&model.PolicyRule{},
		&model.Resource{},
		&model.Property{},
		&model.Subdomain{},
		&model.ClientAppointment{},

		// Tenant schema models (TenantModels)
		&model.Appointment{},
		&model.AppointmentArchive{},
		&model.BranchServiceDensity{},
		&model.BranchWorkRange{},
		&model.Branch{},
		&model.EmployeeServiceDensity{},
		&model.EmployeeWorkRange{},
		&model.Employee{},
		&model.Service{},
		&model.Payment{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
