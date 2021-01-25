package mcmodel

import "time"

type GlobusRequestFile struct {
	ID              int            `json:"id"`
	UUID            string         `json:"string"`
	ProjectID       int            `json:"project_id"`
	OwnerID         int            `json:"owner_id"`
	GlobusRequestID int            `json:"globus_request_id"`
	GlobusRequest   *GlobusRequest `gorm:"foreignKey:GlobusRequestID;references:ID"`
	Path            string         `json:"path"`
	FileID          int            `json:"file_id"`
	File            *File          `gorm:"foreignKey:FileID;references:ID"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (GlobusRequestFile) TableName() string {
	return "globus_request_files"
}