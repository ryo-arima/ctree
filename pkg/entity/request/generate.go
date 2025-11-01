package request

import "fmt"

// GenerateRequest represents the request to generate source code
type GenerateRequest struct {
	Language     string   `json:"language,omitempty" yaml:"language,omitempty"`
	Framework    string   `json:"framework,omitempty" yaml:"framework,omitempty"`
	SourcePath   string   `json:"source_path" yaml:"source_path"`
	OutputPath   string   `json:"output_path" yaml:"output_path"`
	Languages    []string `json:"languages,omitempty" yaml:"languages,omitempty"`
	Recursive    bool     `json:"recursive,omitempty" yaml:"recursive,omitempty"`
	ExcludeFiles []string `json:"exclude_files,omitempty" yaml:"exclude_files,omitempty"`
	IncludeFiles []string `json:"include_files,omitempty" yaml:"include_files,omitempty"`
	MaxDepth     int      `json:"max_depth,omitempty" yaml:"max_depth,omitempty"`
}

// Validate validates the generate request
func (r *GenerateRequest) Validate() error {
	if r.SourcePath == "" {
		return fmt.Errorf("source_path is required")
	}
	if r.OutputPath == "" {
		return fmt.Errorf("output_path is required")
	}
	return nil
}
