package golang

import (
	"github.com/ryo-arima/ctree/pkg/entity/model"
)

// GoPureProjectCtagsRepository handles Go pure project specific ctags operations
type GoPureProjectCtagsRepository interface {
	GenerateTags(sourcePath string, recursive bool) ([]model.Tag, error)
	CheckCtagsInstalled() error
}

type goPureProjectCtagsRepository struct {
}

// NewGoPureProjectCtagsRepository creates a new Go pure project ctags repository
func NewGoPureProjectCtagsRepository() GoPureProjectCtagsRepository {
	return &goPureProjectCtagsRepository{}
}

// CheckCtagsInstalled checks if ctags supports Go
func (r *goPureProjectCtagsRepository) CheckCtagsInstalled() error {
	// TODO: Implement actual ctags installation check
	return nil
}

// GenerateTags generates tags for Go pure project files
func (r *goPureProjectCtagsRepository) GenerateTags(sourcePath string, recursive bool) ([]model.Tag, error) {
	// TODO: Implement actual ctags execution for Go
	return []model.Tag{}, nil
}
