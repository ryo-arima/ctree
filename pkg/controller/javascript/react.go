package javascript

import (
	"github.com/ryo-arima/ctree/pkg/config"
	"github.com/ryo-arima/ctree/pkg/entity/request"
	javascript_usecase "github.com/ryo-arima/ctree/pkg/usecase/javascript"
)

// GenerateReactProject generates React projects
func GenerateReactProject(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := javascript_usecase.NewJavaScriptReactGenerateUsecase(conf)
	return uc.Generate(req, format), nil
}
