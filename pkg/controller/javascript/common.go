package javascript

import (
	"fmt"

	"github.com/ryo-arima/ctree/pkg/config"
	"github.com/ryo-arima/ctree/pkg/entity/request"
	"github.com/spf13/cobra"
)

// InitGenerateJavaScriptCmd creates a generate command for JavaScript projects
func InitGenerateJavaScriptCmd(conf *config.Config) *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "javascript",
		Short: "Generate call tree for JavaScript/TypeScript project source code",
		Long: `Generate call tree for JavaScript/TypeScript project source code using ctags and output in YAML format.
Supports various JavaScript/TypeScript frameworks and pure projects.`,
		Run: func(cmd *cobra.Command, args []string) {
			sourcePath, _ := cmd.Flags().GetString("source")
			outputPath, _ := cmd.Flags().GetString("output")
			recursive, _ := cmd.Flags().GetBool("recursive")
			maxDepth, _ := cmd.Flags().GetInt("max-depth")
			framework, _ := cmd.Flags().GetString("framework")

			if sourcePath == "" && len(args) > 0 {
				sourcePath = args[0]
			}
			if sourcePath == "" {
				sourcePath = "."
			}

			req := request.GenerateRequest{
				Language:   "javascript",
				Framework:  framework,
				SourcePath: sourcePath,
				OutputPath: outputPath,
				Recursive:  recursive,
				MaxDepth:   maxDepth,
			}

			var result string
			var err error

			// フレームワークに応じて適切な関数を呼び出し
			switch framework {
			case "react":
				result, err = GenerateReactProject(conf, req, "yaml")
			case "pure", "":
				result, err = GeneratePureProject(conf, req, "yaml")
			default:
				err = fmt.Errorf("unsupported framework: %s", framework)
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Print(result)
		},
	}

	// フレームワーク選択オプションを追加
	generateCmd.Flags().StringP("framework", "f", "pure", "Framework to generate (pure, react)")
	generateCmd.Flags().StringP("source", "s", ".", "Source directory or file to generate")
	generateCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	generateCmd.Flags().BoolP("recursive", "r", true, "Recursively analyze subdirectories")
	generateCmd.Flags().IntP("max-depth", "d", 10, "Maximum depth for recursive generation")

	return generateCmd
}

// InitGetJavaScriptCmd creates a get command for JavaScript projects
func InitGetJavaScriptCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "javascript",
		Short: "Get information from JavaScript project",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Get command for JavaScript - Not yet implemented")
		},
	}
}

// InitListJavaScriptCmd creates a list command for JavaScript projects
func InitListJavaScriptCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "javascript",
		Short: "List information from JavaScript project",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("List command for JavaScript - Not yet implemented")
		},
	}
}
