package pkg

import (
	"github.com/ryo-arima/ctree/pkg/config"
	golang_controller "github.com/ryo-arima/ctree/pkg/controller/golang"
	javascript_controller "github.com/ryo-arima/ctree/pkg/controller/javascript"
	python_controller "github.com/ryo-arima/ctree/pkg/controller/python"
	"github.com/spf13/cobra"
)

// BaseCmdForCtree represents the base commands for ctree
type BaseCmdForCtree struct {
	Generate *cobra.Command
	Get      *cobra.Command
	List     *cobra.Command
	Version  *cobra.Command
}

// InitRootCmdForCtree creates the root command for ctree
func InitRootCmdForCtree() *cobra.Command {
	var output string
	var rootCmd = &cobra.Command{
		Use:   "ctree",
		Short: "ctree is a CLI tool to generate call trees from source code",
		Long: `ctree generates source code using ctags and generates call trees in YAML format.
It supports all programming languages that ctags supports.

Examples:
  ctree generate golang --framework pure    # generate call tree for Go project
  ctree generate javascript --framework react  # generate for React project
  ctree get golang functions                # get function information
  ctree list golang --type functions        # list all functions`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// TODO: Set output format if needed
		},
	}
	rootCmd.PersistentFlags().StringVarP(&output, "output-format", "f", "yaml", "Output format: yaml|json|table")
	return rootCmd
}

// InitBaseCmdForCtree creates the base commands for ctree
func InitBaseCmdForCtree(conf *config.Config) BaseCmdForCtree {
	// Create generate command
	generateCmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate call tree from source code",
		Long:  `Generate call tree from source code using ctags`,
	}
	generateCmd.AddCommand(golang_controller.InitGenerateGolangCmd(conf))
	generateCmd.AddCommand(javascript_controller.InitGenerateJavaScriptCmd(conf))
	generateCmd.AddCommand(python_controller.InitGeneratePythonCmd(conf))

	// Create get command
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get specific information from source code",
		Long:  `Get specific information like functions, classes, variables from source code`,
	}
	getCmd.AddCommand(golang_controller.InitGetGolangCmd(conf))
	getCmd.AddCommand(javascript_controller.InitGetJavaScriptCmd(conf))
	getCmd.AddCommand(python_controller.InitGetPythonCmd(conf))

	// Create list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List information from source code",
		Long:  `List information like functions, classes, variables from source code`,
	}
	listCmd.AddCommand(golang_controller.InitListGolangCmd(conf))
	listCmd.AddCommand(javascript_controller.InitListJavaScriptCmd(conf))
	listCmd.AddCommand(python_controller.InitListPythonCmd(conf))

	// Create version command
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("ctree version 0.1.0")
		},
	}

	return BaseCmdForCtree{
		Generate: generateCmd,
		Get:      getCmd,
		List:     listCmd,
		Version:  versionCmd,
	}
}

// ClientForCtree initializes and executes the ctree CLI
func ClientForCtree(conf *config.Config) {
	rootCmd := InitRootCmdForCtree()
	rootCmd.CompletionOptions.HiddenDefaultCmd = true

	baseCmd := InitBaseCmdForCtree(conf)

	// Add commands to root
	rootCmd.AddCommand(baseCmd.Generate)
	rootCmd.AddCommand(baseCmd.Get)
	rootCmd.AddCommand(baseCmd.List)
	rootCmd.AddCommand(baseCmd.Version)

	// Execute the root command
	rootCmd.Execute()
}
