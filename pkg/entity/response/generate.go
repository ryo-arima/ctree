package response

import "github.com/ryo-arima/ctree/pkg/entity/model"

// GenerateResponse represents the response from generating source code
type GenerateResponse struct {
	Success  bool           `json:"success" yaml:"success"`
	Message  string         `json:"message,omitempty" yaml:"message,omitempty"`
	CallTree *model.CTree   `json:"call_tree,omitempty" yaml:"call_tree,omitempty"`
	Stats    *AnalysisStats `json:"stats,omitempty" yaml:"stats,omitempty"`
}

// AnalysisStats represents statistics about the generation
type AnalysisStats struct {
	TotalFiles     int            `json:"total_files" yaml:"total_files"`
	TotalFunctions int            `json:"total_functions" yaml:"total_functions"`
	TotalCallEdges int            `json:"total_call_edges" yaml:"total_call_edges"`
	LanguagesFound map[string]int `json:"languages_found" yaml:"languages_found"`
	ProcessingTime string         `json:"processing_time" yaml:"processing_time"`
}
