package model

type Sector struct {
	BaseModel
	Name        string `gorm:"not null;unique" json:"name"`
	Description string `json:"description"`
}

const SectorTableName = "public.sectors"

func (Sector) TableName() string { return SectorTableName }
func (Sector) SchemaType() string { return "public" }
