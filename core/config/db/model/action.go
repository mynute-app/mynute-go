package model

import "gorm.io/gorm"

var AllowActionCreation = false

type Action struct {
	gorm.Model
	Name        string `json:"name" gorm:"unique;not null"`
	Description string `json:"description"`
}

var ViewAction = &Action{
	Name:        "view",
	Description: "Can view a resource",
}

var CreateAction = &Action{
	Name:        "create",
	Description: "Can create a resource",
}

var UpdateAction = &Action{
	Name:        "update",
	Description: "Can update a resource",
}

var DeleteAction = &Action{
	Name:        "delete",
	Description: "Can delete a resource",
}

var ActionList = []*Action{
	ViewAction,
	CreateAction,
	UpdateAction,
	DeleteAction,
}

func SeedActions(db *gorm.DB) {
	AllowActionCreation = true
	defer func() { AllowActionCreation = false }()
	for _, action := range ActionList {
		if err := db.Where(Action{Name: action.Name}).FirstOrCreate(action).Error; err != nil {
			panic(err)
		}
	}
}
