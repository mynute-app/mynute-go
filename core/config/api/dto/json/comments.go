package dJSON

import "github.com/google/uuid"

type Comments []Comment

type Comment struct {
	CreatedAt     string          `json:"created_at" example:"2021-01-01T09:00:00Z"`
	UpdatedAt     string          `json:"updated_at" example:"2021-01-01T09:00:00Z"`
	DeletedAt     string          `json:"deleted_at" example:"2021-01-01T09:00:00Z"`
	Comment       string          `json:"comment" example:"Some comment example text"`
	OldVersions   CommentVersions `json:"old_versions"`
	CreatedBy     uuid.UUID       `json:"created_by" example:"00000000-0000-0000-0000-000000000000"`
	LastUpdatedBy uuid.UUID          `json:"last_updated_by" example:"00000000-0000-0000-0000-000000000000"`
	FromClient    bool            `json:"from_client" example:"true"`
	FromEmployee  bool            `json:"from_employee" example:"false"`
	Type          string          `json:"type" example:"internal"` // "internal" or "external"
}

type CommentVersions []CommentVersion

type CommentVersion struct {
	CreatedAt string `json:"created_at" example:"2021-01-01T09:00:00Z"`
	Comment   string `json:"comment" example:"Some different version comment text"`
	CreatedBy uuid.UUID `json:"created_by" example:"00000000-0000-0000-0000-000000000000"`
}
