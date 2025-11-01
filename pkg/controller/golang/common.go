package golang

import (
	"fmt"
	"os"

	"github.com/ryo-arima/ctree/pkg/config"
	"github.com/ryo-arima/ctree/pkg/entity/request"
	"github.com/spf13/cobra"
)

// InitGenerateGolangCmd creates a generate command for Golang projects
func InitGenerateGolangCmd(conf *config.Config) *cobra.Command {
	generateCmd := &cobra.Command{
		Use:   "golang",
		Short: "Generate call tree for Golang project source code",
		Long: `Generate call tree for Golang project source code using ctags and output in YAML format.
Supports various Go frameworks and pure projects.`,
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
				Language:   "golang",
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
			case "gin":
				// result, err = AnalyzeGinProject(conf, req, "yaml")
				err = fmt.Errorf("gin framework support not implemented yet")
			case "echo":
				// result, err = AnalyzeEchoProject(conf, req, "yaml")
				err = fmt.Errorf("echo framework support not implemented yet")
			case "pure", "":
				result, err = GeneratePureProject(conf, req, "yaml")
			default:
				err = fmt.Errorf("unsupported framework: %s", framework)
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Output handling
			if outputPath != "" {
				err = os.WriteFile(outputPath, []byte(result), 0644)
				if err != nil {
					fmt.Printf("Error writing to file %s: %v\n", outputPath, err)
					return
				}
				fmt.Printf("Output written to %s\n", outputPath)
			} else {
				fmt.Print(result)
			}
		},
	}

	// フレームワーク選択オプションを追加
	generateCmd.Flags().String("framework", "pure", "Framework to generate (pure, gin, echo)")
	generateCmd.Flags().StringP("source", "s", ".", "Source directory or file to generate")
	generateCmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	generateCmd.Flags().BoolP("recursive", "r", true, "Recursively analyze subdirectories")
	generateCmd.Flags().IntP("max-depth", "d", 10, "Maximum depth for recursive generation")

	return generateCmd
}

// InitGetGolangCmd creates a get command for Golang projects
func InitGetGolangCmd(conf *config.Config) *cobra.Command {
	getCmd := &cobra.Command{
		Use:   "golang",
		Short: "Get specific information from Golang project",
		Long: `Get specific information from Golang project source code.
Available subcommands: call-tree, functions, classes, variables, imports`,
	}

	// サブコマンドを追加
	getCmd.AddCommand(initGetCallTreeCmd(conf))
	getCmd.AddCommand(initGetFunctionsCmd(conf))
	getCmd.AddCommand(initGetClassesCmd(conf))
	getCmd.AddCommand(initGetVariablesCmd(conf))
	getCmd.AddCommand(initGetImportsCmd(conf))

	return getCmd
}

// InitListGolangCmd creates a list command for Golang projects
func InitListGolangCmd(conf *config.Config) *cobra.Command {
	listCmd := &cobra.Command{
		Use:   "golang",
		Short: "List information from Golang project",
		Long: `List information from Golang project source code using ctags.
Available options: --type (functions, classes, variables, imports)`,
		Run: func(cmd *cobra.Command, args []string) {
			sourcePath, _ := cmd.Flags().GetString("source")
			itemType, _ := cmd.Flags().GetString("type")
			recursive, _ := cmd.Flags().GetBool("recursive")
			format, _ := cmd.Flags().GetString("format")

			if sourcePath == "" && len(args) > 0 {
				sourcePath = args[0]
			}
			if sourcePath == "" {
				sourcePath = "."
			}

			req := request.GenerateRequest{
				SourcePath: sourcePath,
				Recursive:  recursive,
				MaxDepth:   10,
			}

			var result string
			var err error

			switch itemType {
			case "functions", "func":
				result, err = ListFunctions(conf, req, format)
			case "classes", "class", "types":
				result, err = ListClasses(conf, req, format)
			case "variables", "var":
				result, err = ListVariables(conf, req, format)
			case "imports", "import":
				result, err = ListImports(conf, req, format)
			default:
				err = fmt.Errorf("unsupported type: %s (available: functions, classes, variables, imports)", itemType)
			}

			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			fmt.Print(result)
		},
	}

	listCmd.Flags().StringP("source", "s", ".", "Source directory or file to generate")
	listCmd.Flags().StringP("type", "t", "functions", "Type of items to list (functions, classes, variables, imports)")
	listCmd.Flags().StringP("format", "f", "table", "Output format (table, json, yaml)")
	listCmd.Flags().BoolP("recursive", "r", true, "Recursively analyze subdirectories")

	return listCmd
}

// サブコマンド実装
func initGetFunctionsCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "functions [function_name]",
		Short: "Get specific function information",
		Run: func(cmd *cobra.Command, args []string) {
			sourcePath, _ := cmd.Flags().GetString("source")
			if sourcePath == "" {
				sourcePath = "."
			}

			var functionName string
			if len(args) > 0 {
				functionName = args[0]
			}

			req := request.GenerateRequest{
				SourcePath: sourcePath,
				Recursive:  true,
			}

			result, err := GetFunction(conf, req, functionName, "yaml")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Print(result)
		},
	}
}

func initGetClassesCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "classes [class_name]",
		Short: "Get specific class/type information",
		Run: func(cmd *cobra.Command, args []string) {
			sourcePath, _ := cmd.Flags().GetString("source")
			if sourcePath == "" {
				sourcePath = "."
			}

			var className string
			if len(args) > 0 {
				className = args[0]
			}

			req := request.GenerateRequest{
				SourcePath: sourcePath,
				Recursive:  true,
			}

			result, err := GetClass(conf, req, className, "yaml")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Print(result)
		},
	}
}

func initGetVariablesCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "variables [variable_name]",
		Short: "Get specific variable information",
		Run: func(cmd *cobra.Command, args []string) {
			sourcePath, _ := cmd.Flags().GetString("source")
			if sourcePath == "" {
				sourcePath = "."
			}

			var variableName string
			if len(args) > 0 {
				variableName = args[0]
			}

			req := request.GenerateRequest{
				SourcePath: sourcePath,
				Recursive:  true,
			}

			result, err := GetVariable(conf, req, variableName, "yaml")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Print(result)
		},
	}
}

func initGetImportsCmd(conf *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "imports",
		Short: "Get import information",
		Run: func(cmd *cobra.Command, args []string) {
			sourcePath, _ := cmd.Flags().GetString("source")
			if sourcePath == "" {
				sourcePath = "."
			}

			req := request.GenerateRequest{
				SourcePath: sourcePath,
				Recursive:  true,
			}

			result, err := GetImports(conf, req, "yaml")
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}
			fmt.Print(result)
		},
	}
}

// initGetCallTreeCmd creates a get call-tree command
func initGetCallTreeCmd(conf *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "call-tree",
		Short: "Get call tree from generated ctree file",
		Long:  `Extract and display call tree information from a previously generated ctree YAML file`,
		Run: func(cmd *cobra.Command, args []string) {
			ctreePath, _ := cmd.Flags().GetString("ctree")
			framework, _ := cmd.Flags().GetString("framework")
			outputPath, _ := cmd.Flags().GetString("output")
			format, _ := cmd.Flags().GetString("format")
			expandSignature, _ := cmd.Flags().GetBool("expand-signature")

			if ctreePath == "" {
				fmt.Println("Error: --ctree flag is required")
				cmd.Usage()
				return
			}

			req := request.GenerateRequest{
				SourcePath: ctreePath,
				Framework:  framework,
			}

			result, err := GetCallTree(conf, req, format, expandSignature)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				return
			}

			// Output handling
			if outputPath != "" {
				err = os.WriteFile(outputPath, []byte(result), 0644)
				if err != nil {
					fmt.Printf("Error writing to file %s: %v\n", outputPath, err)
					return
				}
				fmt.Printf("Output written to %s\n", outputPath)
			} else {
				fmt.Print(result)
			}
		},
	}

	cmd.Flags().StringP("ctree", "c", "", "Path to ctree YAML file (required)")
	cmd.Flags().String("framework", "pure", "Framework type (pure, gin, echo)")
	cmd.Flags().StringP("output", "o", "", "Output file path (default: stdout)")
	cmd.Flags().String("format", "yaml", "Output format (yaml, json, text)")
	cmd.Flags().Bool("expand-signature", false, "Show function parameters and return values on separate lines")
	cmd.MarkFlagRequired("ctree")

	return cmd
}
