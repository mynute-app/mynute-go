package model

type Property struct {
	BaseModel
	Name         string   `json:"name" gorm:"unique;not null"`
	Description  string   `json:"description"`
	ResourceName string   `json:"resource_name"`
	Resource     Resource `gorm:"foreignKey:ResourceName;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"resource"`
}

func (Property) TableName() string { return "public.properties" }
func (Property) SchemaType() string { return "public" }

