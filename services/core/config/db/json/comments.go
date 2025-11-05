package mJSON

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Comments []Comment

type Comment struct {
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"deleted_at"`
	Comment       string          `json:"comment"`
	OldVersions   CommentVersions `json:"old_versions"` // List of old versions of the comment
	CreatedBy     uuid.UUID       `gorm:"type:uuid;not null;index" json:"created_by"`
	LastUpdatedBy *uuid.UUID      `gorm:"type:uuid;index" json:"last_updated_by"`
	FromClient    bool            `json:"from_client"`   // true if the comment is from the client
	FromEmployee  bool            `json:"from_employee"` // true if the comment is from the employee
	Type          string          `json:"type"`          // "internal" or "external"
}

type CommentVersions []CommentVersion

type CommentVersion struct {
	CreatedAt time.Time `json:"created_at"`
	Comment   string    `json:"comment"`
	CreatedBy uuid.UUID `gorm:"type:uuid;not null;index" json:"created_by"`
}

func (cv *CommentVersions) IsEmpty() bool {
	if cv == nil {
		return true
	}
	return len(*cv) == 0
}

func (c *Comment) IsEmpty() bool {
	return len(c.Comment) == 0
}

func (c *Comment) Edit(newCommentStr string, editor uuid.UUID) error {
	if c == nil {
		return errors.New("comment is nil")
	}
	old_version := CommentVersion{
		Comment: c.Comment,
	}

	if c.UpdatedAt.IsZero() {
		old_version.CreatedAt = c.CreatedAt
	} else {
		old_version.CreatedAt = c.UpdatedAt
	}

	if c.OldVersions.IsEmpty() {
		c.OldVersions = make(CommentVersions, 0)
		old_version.CreatedBy = c.CreatedBy
	}

	c.OldVersions = append(c.OldVersions, old_version)
	c.Comment = newCommentStr
	c.UpdatedAt = time.Now()
	c.LastUpdatedBy = &editor

	return nil
}

// --- Implement Scanner/Valuer for Comments ---
func (ac *Comments) Value() (driver.Value, error) {
	if ac == nil || len(*ac) == 0 {
		// Return empty JSON array `[]` which is valid JSON
		return json.Marshal([]Comment{})
	}
	return json.Marshal(ac)
}

func (ac *Comments) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		// Handle nil from DB
		if value == nil {
			*ac = []Comment{} // Initialize to empty slice
			return nil
		}
		return errors.New("failed to scan Comments: expected []byte")
	}
	// Handle empty JSON array or null from DB
	if len(bytes) == 0 || string(bytes) == "null" {
		*ac = []Comment{} // Initialize to empty slice
		return nil
	}
	// Important: Unmarshal into the pointer *ac
	return json.Unmarshal(bytes, ac)
}

// Optional: Add helper methods directly to the type
func (ac *Comments) Add(c Comment) {
	if ac == nil {
		ac = &Comments{}
	}
	*ac = append(*ac, c)
}

