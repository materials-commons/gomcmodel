package mcmodel

import (
	"encoding/json"
	"fmt"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
	"path/filepath"
	"time"
)

type Dataset struct {
	ID            int
	Name          string
	UUID          string
	ProjectID     int
	License       string
	LicenseLink   string
	Description   string
	Summary       string
	DOI           string
	Authors       string
	Files         []File `gorm:"many2many:dataset2file"`
	FileSelection string
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PublishedAt   time.Time
}

func (d Dataset) ZipfileDir(mcfsDir string) string {
	return filepath.Join(mcfsDir, "zipfiles", d.UUID)
}

func (d Dataset) ZipfilePath(mcfsDir string) string {
	return filepath.Join(d.ZipfileDir(mcfsDir), fmt.Sprintf("%s.zip", slug.Make(d.Name)))
}

func (d Dataset) GetFileSelection() (*FileSelection, error) {
	var fs FileSelection
	err := json.Unmarshal([]byte(d.FileSelection), &fs)
	return &fs, err
}

type FileSelection struct {
	IncludeFiles []string `json:"include_files"`
	ExcludeFiles []string `json:"exlude_files"`
	IncludeDirs  []string `json:"include_dirs"`
	ExcludeDirs  []string `json:"exclude_dirs"`
}

func (d Dataset) GetFiles(db *gorm.DB) *gorm.DB {
	return db.Preload("Directory").Where("id in (?)",
		db.Table("dataset2file").
			Select("file_id").
			Where("dataset_id = ?", d.ID))
}

func (d Dataset) GetEntitiesFromTemplate(db *gorm.DB) ([]Entity, error) {
	experimentIdsSubSubquery := db.Table("item2entity_selection").
		Select("experiment_id").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset")

	entityIdsFromExperimentSubquery := db.Table("experiment2entity").
		Select("entity_id").
		Where("experiment_id in (?)", experimentIdsSubSubquery)

	entityNamesFromExperimentSubquery := db.Table("item2entity_selection").
		Select("entity_name").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset").
		Where("experiment_id in (?)", experimentIdsSubSubquery)

	entityIdSubquery := db.Table("item2entity_selection").
		Select("entity_id").
		Where("item_id = ?", d.ID).
		Where("item_type = ?", "App\\Models\\Dataset")

	var entities []Entity
	result := db.Preload("Files.Directory").
		Where("id in (?)", entityIdsFromExperimentSubquery).
		Where("name in (?)", entityNamesFromExperimentSubquery).
		Or("id in (?)", entityIdSubquery).
		Find(&entities)

	return entities, result.Error
}
