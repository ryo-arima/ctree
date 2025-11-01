package javascript

import (
	"github.com/ryo-arima/ctree/pkg/config"
	"github.com/ryo-arima/ctree/pkg/entity/request"
	javascript_usecase "github.com/ryo-arima/ctree/pkg/usecase/javascript"
)

// GeneratePureProject generates pure JavaScript/TypeScript projects
func GeneratePureProject(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := javascript_usecase.NewJavaScriptPureProjectGenerateUsecase(conf)
	return uc.Generate(req, format), nil
}

// GetOutputFormat returns current output format
func GetOutputFormat() string {
	return "yaml" // TODO: Import from common controller
}
