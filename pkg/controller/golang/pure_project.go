package golang

import (
	"fmt"
	"os"
	"strings"

	"github.com/ryo-arima/ctree/pkg/config"
	"github.com/ryo-arima/ctree/pkg/entity/model"
	"github.com/ryo-arima/ctree/pkg/entity/request"
	golang_usecase "github.com/ryo-arima/ctree/pkg/usecase/golang"
	"gopkg.in/yaml.v3"
)

// ANSI color codes
const (
	colorReset      = "\033[0m"
	colorRed        = "\033[31m"
	colorGreen      = "\033[32m"
	colorYellow     = "\033[33m"
	colorBlue       = "\033[34m"
	colorMagenta    = "\033[35m"
	colorCyan       = "\033[36m"
	colorWhite      = "\033[37m"
	colorBrightBlue = "\033[94m"
	colorBrightCyan = "\033[96m"
	colorGray       = "\033[90m"
	colorBold       = "\033[1m"
)

// GeneratePureProject generates pure Go projects
func GeneratePureProject(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	return uc.Generate(req, format)
}

// ListFunctions lists all functions in the project
func ListFunctions(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	// TODO: Implement function listing logic using ctags
	_ = uc // suppress unused variable warning for now
	return fmt.Sprintf("Functions list for %s (format: %s)", req.SourcePath, format), nil
}

// ListClasses lists all classes/types in the project
func ListClasses(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	// TODO: Implement class/type listing logic using ctags
	_ = uc // suppress unused variable warning for now
	return fmt.Sprintf("Classes/Types list for %s (format: %s)", req.SourcePath, format), nil
}

// ListVariables lists all variables in the project
func ListVariables(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	// TODO: Implement variable listing logic using ctags
	_ = uc // suppress unused variable warning for now
	return fmt.Sprintf("Variables list for %s (format: %s)", req.SourcePath, format), nil
}

// ListImports lists all imports in the project
func ListImports(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	// TODO: Implement import listing logic using ctags
	_ = uc // suppress unused variable warning for now
	return fmt.Sprintf("Imports list for %s (format: %s)", req.SourcePath, format), nil
}

// GetFunction gets specific function information
func GetFunction(conf *config.Config, req request.GenerateRequest, functionName string, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	// TODO: Implement specific function retrieval logic using ctags
	_ = uc // suppress unused variable warning for now
	if functionName == "" {
		return fmt.Sprintf("All functions in %s (format: %s)", req.SourcePath, format), nil
	}
	return fmt.Sprintf("Function '%s' details from %s (format: %s)", functionName, req.SourcePath, format), nil
}

// GetClass gets specific class/type information
func GetClass(conf *config.Config, req request.GenerateRequest, className string, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	// TODO: Implement specific class/type retrieval logic using ctags
	_ = uc // suppress unused variable warning for now
	if className == "" {
		return fmt.Sprintf("All classes/types in %s (format: %s)", req.SourcePath, format), nil
	}
	return fmt.Sprintf("Class/Type '%s' details from %s (format: %s)", className, req.SourcePath, format), nil
}

// GetVariable gets specific variable information
func GetVariable(conf *config.Config, req request.GenerateRequest, variableName string, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	// TODO: Implement specific variable retrieval logic using ctags
	_ = uc // suppress unused variable warning for now
	if variableName == "" {
		return fmt.Sprintf("All variables in %s (format: %s)", req.SourcePath, format), nil
	}
	return fmt.Sprintf("Variable '%s' details from %s (format: %s)", variableName, req.SourcePath, format), nil
}

// GetImports gets import information
func GetImports(conf *config.Config, req request.GenerateRequest, format string) (string, error) {
	uc := golang_usecase.NewGoPureProjectGenerateUsecase(conf)
	// TODO: Implement import retrieval logic using ctags
	_ = uc // suppress unused variable warning for now
	return fmt.Sprintf("Imports from %s (format: %s)", req.SourcePath, format), nil
}

// GetCallTree extracts call tree from a previously generated ctree YAML file
func GetCallTree(conf *config.Config, req request.GenerateRequest, format string, expandSignature bool) (string, error) {
	// Read the ctree YAML file
	data, err := os.ReadFile(req.SourcePath)
	if err != nil {
		return "", fmt.Errorf("failed to read ctree file: %w", err)
	}

	// Parse YAML
	var ctree model.CTree
	if err := yaml.Unmarshal(data, &ctree); err != nil {
		return "", fmt.Errorf("failed to parse ctree YAML: %w", err)
	}

	// Extract call tree based on format
	switch format {
	case "text", "tree":
		// Return indented tree visualization
		return formatCallTreeAsText(ctree.CallTree, expandSignature), nil
	case "yaml", "":
		// Return call tree nodes as YAML
		if len(ctree.CallTree) == 0 {
			return "call_tree: []\n", nil
		}
		output, err := yaml.Marshal(map[string]interface{}{
			"call_tree": ctree.CallTree,
		})
		if err != nil {
			return "", fmt.Errorf("failed to marshal call tree: %w", err)
		}
		return string(output), nil
	default:
		return "", fmt.Errorf("unsupported format: %s (supported: text, tree, yaml)", format)
	}
}

// formatCallTreeAsText formats call tree nodes as indented text with colors
func formatCallTreeAsText(nodes []model.CallTreeNode, expandSignature bool) string {
	if len(nodes) == 0 {
		return "No call tree available\n"
	}

	var result strings.Builder
	result.WriteString(colorBold + colorCyan + "Call Tree:\n" + colorReset)
	result.WriteString(colorCyan + "==========" + colorReset + "\n\n")

	for i, node := range nodes {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(colorBold + colorYellow + fmt.Sprintf("Entry Point %d: ", i+1) + colorReset)
		if expandSignature {
			// For entry points with expand signature, show function name only
			funcName := "func "
			if node.Receiver != "" {
				funcName += fmt.Sprintf("(%s) ", node.Receiver)
			}
			funcName += node.Name
			result.WriteString(colorGreen + funcName + colorReset)
			result.WriteString(" " + colorGreen + "[internal]" + colorReset)
			if node.File != "" {
				result.WriteString(" " + colorGray + fmt.Sprintf("(%s:%d)", node.File, node.Line) + colorReset)
			}
			// Show parameters and returns on separate lines
			if len(node.Parameters) > 0 {
				result.WriteString("\n")
				result.WriteString(colorGray + "  Parameters:" + colorReset)
				for _, param := range node.Parameters {
					result.WriteString("\n")
					result.WriteString(colorGray + "    - " + colorReset)
					result.WriteString(colorWhite + param.Name + ": " + colorMagenta + param.Type + colorReset)
				}
			}
			if len(node.ReturnTypes) > 0 {
				result.WriteString("\n")
				result.WriteString(colorGray + "  Returns:" + colorReset)
				for _, returnType := range node.ReturnTypes {
					result.WriteString("\n")
					result.WriteString(colorGray + "    - " + colorMagenta + returnType + colorReset)
				}
			}
		} else {
			result.WriteString(colorGreen + node.Title + colorReset)
			result.WriteString(" " + colorGreen + "[internal]" + colorReset)
			if node.File != "" {
				result.WriteString(" " + colorGray + fmt.Sprintf("(%s:%d)", node.File, node.Line) + colorReset)
			}
		}
		result.WriteString("\n")
		formatNodeRecursive(&result, node, "", true, expandSignature)
	}

	return result.String()
}

// formatNodeRecursive recursively formats a call tree node with indentation
func formatNodeRecursive(result *strings.Builder, node model.CallTreeNode, prefix string, isLast bool, expandSignature bool) {
	// Skip the root node itself as we already printed it
	if prefix == "" {
		// Print children directly with initial indentation
		for i, child := range node.Children {
			childIsLast := i == len(node.Children)-1
			formatNodeRecursive(result, child, "  ", childIsLast, expandSignature)
		}
		return
	}

	// Print current node
	connector := "├─"
	if isLast {
		connector = "└─"
	}
	result.WriteString(colorGray + prefix + colorReset)
	result.WriteString(colorBrightBlue + connector + colorReset)
	result.WriteString(" ")

	// Node title with color based on kind
	titleColor := colorWhite
	if node.Kind == "external" {
		titleColor = colorGray
	} else if node.Kind == "function" || node.Kind == "method" {
		titleColor = colorBrightCyan
	}

	if expandSignature && (node.Kind == "function" || node.Kind == "method") {
		// Show expanded signature with parameters and return values
		// formatExpandedTitle handles tag and location display
		formatExpandedTitle(result, node, prefix+getChildPrefix(isLast))
	} else {
		// Show function name
		result.WriteString(titleColor + node.Title + colorReset)

		// Show [internal] or [external] tag after function name
		if node.Kind == "external" {
			result.WriteString(" " + colorGray + "[external]" + colorReset)
		} else if node.File != "" {
			result.WriteString(" " + colorGreen + "[internal]" + colorReset)
		}

		// Location info or package info after tag
		if node.Kind == "external" {
			// For external functions, show full package path if available
			if node.PackagePath != "" {
				result.WriteString(" " + colorGray + fmt.Sprintf("(%s)", node.PackagePath) + colorReset)
			} else if node.Package != "" {
				result.WriteString(" " + colorGray + fmt.Sprintf("(%s)", node.Package) + colorReset)
			} else {
				// Fallback: extract package name from function name
				packageName := extractPackageName(node.Name)
				if packageName != "" {
					result.WriteString(" " + colorGray + fmt.Sprintf("(%s)", packageName) + colorReset)
				}
			}
		} else if node.File != "" {
			// For internal functions, show file path
			result.WriteString(" " + colorGray + fmt.Sprintf("(%s:%d)", node.File, node.Line) + colorReset)
		}
	}

	// Special markers
	if node.IsRecursive {
		result.WriteString(colorYellow + " [recursive]" + colorReset)
	}
	result.WriteString("\n")

	// Print children
	for i, child := range node.Children {
		childIsLast := i == len(node.Children)-1
		childPrefix := prefix
		if isLast {
			childPrefix += "   "
		} else {
			childPrefix += "│  "
		}
		formatNodeRecursive(result, child, childPrefix, childIsLast, expandSignature)
	}
}

// getChildPrefix returns the prefix for child nodes
func getChildPrefix(isLast bool) string {
	if isLast {
		return "   "
	}
	return "│  "
}

// formatExpandedTitle formats a function signature with parameters and return values on separate lines
func formatExpandedTitle(result *strings.Builder, node model.CallTreeNode, childPrefix string) {
	titleColor := colorBrightCyan

	// Function name with receiver if present
	funcName := "func "
	if node.Receiver != "" {
		funcName += fmt.Sprintf("(%s) ", node.Receiver)
	}
	funcName += node.Name
	result.WriteString(titleColor + funcName + colorReset)

	// Show [internal] or [external] tag after function name
	if node.Kind == "external" {
		result.WriteString(" " + colorGray + "[external]" + colorReset)
		// For external functions, show full package path if available
		if node.PackagePath != "" {
			result.WriteString(" " + colorGray + fmt.Sprintf("(%s)", node.PackagePath) + colorReset)
		} else if node.Package != "" {
			result.WriteString(" " + colorGray + fmt.Sprintf("(%s)", node.Package) + colorReset)
		} else {
			// Fallback: extract package name from function name
			packageName := extractPackageName(node.Name)
			if packageName != "" {
				result.WriteString(" " + colorGray + fmt.Sprintf("(%s)", packageName) + colorReset)
			}
		}
	} else if node.File != "" {
		result.WriteString(" " + colorGreen + "[internal]" + colorReset)
		// For internal functions, show file path
		result.WriteString(" " + colorGray + fmt.Sprintf("(%s:%d)", node.File, node.Line) + colorReset)
	}

	// Parameters on separate lines
	if len(node.Parameters) > 0 {
		result.WriteString("\n")
		result.WriteString(colorGray + childPrefix + "│  Parameters:" + colorReset)
		for _, param := range node.Parameters {
			result.WriteString("\n")
			result.WriteString(colorGray + childPrefix + "│    - " + colorReset)
			result.WriteString(colorWhite + param.Name + ": " + colorMagenta + param.Type + colorReset)
		}
	}

	// Return types on separate lines
	if len(node.ReturnTypes) > 0 {
		result.WriteString("\n")
		result.WriteString(colorGray + childPrefix + "│  Returns:" + colorReset)
		for _, returnType := range node.ReturnTypes {
			result.WriteString("\n")
			result.WriteString(colorGray + childPrefix + "│    - " + colorMagenta + returnType + colorReset)
		}
	}
}

// extractPackageName extracts full package path from a function name like "package.subpackage.FunctionName"
func extractPackageName(functionName string) string {
	// Remove () if present
	name := strings.TrimSuffix(functionName, "()")

	// Handle method calls like "receiver.Method" or "package.subpackage.FunctionName"
	// Extract the part before the last dot
	parts := strings.Split(name, ".")
	if len(parts) > 1 {
		// Return all parts except the last one (the function name)
		// This gives us the full package path
		return strings.Join(parts[:len(parts)-1], ".")
	}
	return ""
}

// GetOutputFormat returns current output format - will be imported from common
func GetOutputFormat() string {
	return "yaml" // TODO: Import from common controller
}
