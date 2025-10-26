package model

type Sector struct {
	BaseModel
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}

func (Sector) TableName() string { return "public.sectors" }
func (Sector) SchemaType() string { return "public" }
