package pkg

import (
	"github.com/ryo-arima/ctree/pkg/config"
	// c_controller "github.com/ryo-arima/ctree/pkg/controller/c"
	// cpp_controller "github.com/ryo-arima/ctree/pkg/controller/cpp"
	golang_controller "github.com/ryo-arima/ctree/pkg/controller/golang"
	python_controller "github.com/ryo-arima/ctree/pkg/controller/python"

	// rust_controller "github.com/ryo-arima/ctree/pkg/controller/rust"
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
  ctree generate golang --source ./myproject    # generate call tree for Go project
  ctree generate c --source ./myapp             # generate for C project
  ctree generate cpp --source ./myapp           # generate for C++ project
  ctree generate rust --source ./myapp          # generate for Rust project
  ctree get golang functions                    # get function information
  ctree list golang --type functions            # list all functions`,
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
	// generateCmd.AddCommand(c_controller.InitGenerateCCmd(conf))          // TODO: Implement C support
	// generateCmd.AddCommand(cpp_controller.InitGenerateCppCmd(conf))      // TODO: Implement C++ support
	generateCmd.AddCommand(golang_controller.InitGenerateGolangCmd(conf))
	// generateCmd.AddCommand(rust_controller.InitGenerateRustCmd(conf))    // TODO: Implement Rust support
	generateCmd.AddCommand(python_controller.InitGeneratePythonCmd(conf))

	// Create get command
	getCmd := &cobra.Command{
		Use:   "get",
		Short: "Get specific information from source code",
		Long:  `Get specific information like functions, classes, variables from source code`,
	}
	// getCmd.AddCommand(c_controller.InitGetCCmd(conf))          // TODO: Implement C support
	// getCmd.AddCommand(cpp_controller.InitGetCppCmd(conf))      // TODO: Implement C++ support
	getCmd.AddCommand(golang_controller.InitGetGolangCmd(conf))
	// getCmd.AddCommand(rust_controller.InitGetRustCmd(conf))    // TODO: Implement Rust support
	getCmd.AddCommand(python_controller.InitGetPythonCmd(conf))

	// Create list command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List information from source code",
		Long:  `List information like functions, classes, variables from source code`,
	}
	// listCmd.AddCommand(c_controller.InitListCCmd(conf))          // TODO: Implement C support
	// listCmd.AddCommand(cpp_controller.InitListCppCmd(conf))      // TODO: Implement C++ support
	listCmd.AddCommand(golang_controller.InitListGolangCmd(conf))
	// listCmd.AddCommand(rust_controller.InitListRustCmd(conf))    // TODO: Implement Rust support
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
