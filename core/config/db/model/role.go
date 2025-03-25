package model

type Role struct {
	ID        uint   `gorm:"primaryKey"`
	TenantID  uint   `gorm:"index"`
	Name      string
	IsDefault bool
}

type UserRole struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint
	RoleID   uint
	TenantID uint
}

type RolePermission struct {
	ID     uint   `gorm:"primaryKey"`
	RoleID uint
	Method string
	Path   string
}