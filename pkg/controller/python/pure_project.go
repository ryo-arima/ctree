package python

import (
	"github.com/ryo-arima/ctree/pkg/config"
	"github.com/ryo-arima/ctree/pkg/entity/request"
	python_usecase "github.com/ryo-arima/ctree/pkg/usecase/python"
)

// GeneratePureProject generates pure Python projects
func GeneratePureProject(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := python_usecase.NewPythonPureProjectGenerateUsecase(conf)
	return uc.Generate(req, format), nil
}

// GetOutputFormat returns current output format
func GetOutputFormat() string {
	return "yaml" // TODO: Import from common controller
}
