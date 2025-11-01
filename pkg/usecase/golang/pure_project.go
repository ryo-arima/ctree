package golang

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"

	"github.com/ryo-arima/ctree/pkg/config"
	"github.com/ryo-arima/ctree/pkg/entity/model"
	"github.com/ryo-arima/ctree/pkg/entity/request"
	"github.com/ryo-arima/ctree/pkg/repository/golang"
	"gopkg.in/yaml.v3"
)

// GoPureProjectGenerateUsecase handles pure Go project specific generation
type GoPureProjectGenerateUsecase interface {
	Generate(req request.GenerateRequest, format string) (string, error)
}

type goPureProjectGenerateUsecase struct {
	config *config.Config
	repo   golang.GoPureProjectRepository
}

// NewGoPureProjectGenerateUsecase creates new Go pure project analyze usecase
func NewGoPureProjectGenerateUsecase(conf *config.Config) GoPureProjectGenerateUsecase {
	return &goPureProjectGenerateUsecase{
		config: conf,
		repo:   golang.NewGoPureProjectRepository(),
	}
}

// Generate performs Go pure project specific source code generation
func (u *goPureProjectGenerateUsecase) Generate(req request.GenerateRequest, format string) (string, error) {
	// Find all Go files
	goFiles, err := u.repo.FindGoFiles(req.SourcePath, req.Recursive, req.MaxDepth)
	if err != nil {
		return "", fmt.Errorf("failed to find Go files: %w", err)
	}

	if len(goFiles) == 0 {
		return "", fmt.Errorf("no Go files found in %s", req.SourcePath)
	}

	// Parse all files and extract functions
	var allFunctions []model.Function
	var entryPoints []model.Function
	functionCalls := make(map[string][]string) // function name -> called functions
	importMap := make(map[string]string)       // package name -> full import path

	for _, filePath := range goFiles {
		file, fset, err := u.repo.ParseGoFile(filePath)
		if err != nil {
			// Log error but continue with other files
			fmt.Printf("Warning: failed to parse %s: %v\n", filePath, err)
			continue
		}

		// Extract import information
		fileImports := u.repo.ExtractImports(file)
		for pkg, path := range fileImports {
			importMap[pkg] = path
		}

		// Convert to relative path before extracting functions
		relPath := u.getRelativePath(filePath)
		functions, err := u.repo.ExtractFunctions(file, fset, relPath)
		if err != nil {
			fmt.Printf("Warning: failed to extract functions from %s: %v\n", relPath, err)
			continue
		}

		// Find entry points (main functions and init functions)
		for _, fn := range functions {
			if fn.Name == "main" && fn.Package == "main" {
				fn.Kind = "entrypoint"
				entryPoints = append(entryPoints, fn)
			} else if fn.Name == "init" {
				fn.Kind = "initializer"
				entryPoints = append(entryPoints, fn)
			}

			// Extract function calls
			calls := u.extractFunctionCalls(file, fn.Name)
			functionKey := u.getFunctionKey(fn)
			functionCalls[functionKey] = calls
		}

		allFunctions = append(allFunctions, functions...)
	}

	// Log entry points found
	if len(entryPoints) > 0 {
		fmt.Printf("Found %d entry point(s):\n", len(entryPoints))
		for _, ep := range entryPoints {
			fmt.Printf("  - %s in %s:%d\n", ep.Name, ep.File, ep.Line)
		}
	} else {
		fmt.Printf("Warning: No entry points (main or init functions) found\n")
	}

	// Build call graph
	var callGraph []model.CallEdge
	for funcKey, calls := range functionCalls {
		for _, calledFunc := range calls {
			// Find the function details
			for _, fn := range allFunctions {
				if fn.Name == calledFunc || u.getFunctionKey(fn) == calledFunc {
					callGraph = append(callGraph, model.CallEdge{
						From: funcKey,
						To:   u.getFunctionKey(fn),
						File: fn.File,
						Line: fn.Line,
					})
					break
				}
			}
		}
	}

	// Update functions with CallsTo information
	for i := range allFunctions {
		funcKey := u.getFunctionKey(allFunctions[i])
		if calls, ok := functionCalls[funcKey]; ok {
			allFunctions[i].CallsTo = calls
		}
	}

	// Build hierarchical call tree from entry points
	callTreeNodes := u.buildHierarchicalCallTree(entryPoints, allFunctions, functionCalls, importMap)

	// Build call tree visualization text
	callTreeData := u.buildCallTreeVisualization(callTreeNodes)

	// Create call tree
	ctree := model.CTree{
		SourceFile:            u.getRelativePath(req.SourcePath),
		Language:              "go",
		Functions:             allFunctions,
		CallGraph:             callGraph,
		EntryPoints:           entryPoints,
		CallTree:              callTreeNodes,
		CallTreeVisualization: callTreeData,
		Metadata: map[string]interface{}{
			"total_functions": len(allFunctions),
			"entry_points":    len(entryPoints),
			"call_edges":      len(callGraph),
		},
	}

	// Format output
	switch strings.ToLower(format) {
	case "yaml", "yml", "":
		data, err := yaml.Marshal(ctree)
		if err != nil {
			return "", fmt.Errorf("failed to marshal to YAML: %w", err)
		}
		return string(data), nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// extractFunctionCalls extracts function calls from a function body
func (u *goPureProjectGenerateUsecase) extractFunctionCalls(file *ast.File, funcName string) []string {
	var calls []string
	callMap := make(map[string]bool)

	// Find the function
	ast.Inspect(file, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok && fn.Name.Name == funcName {
			// Inspect function body
			ast.Inspect(fn.Body, func(node ast.Node) bool {
				if callExpr, ok := node.(*ast.CallExpr); ok {
					callName := u.getCallName(callExpr.Fun)
					if callName != "" && !callMap[callName] {
						callMap[callName] = true
						calls = append(calls, callName)
					}
				}
				return true
			})
			return false
		}
		return true
	})

	return calls
}

// getCallName extracts the function name from a call expression
func (u *goPureProjectGenerateUsecase) getCallName(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		// For method calls like obj.Method()
		if ident, ok := e.X.(*ast.Ident); ok {
			return ident.Name + "." + e.Sel.Name
		}
		return e.Sel.Name
	}
	return ""
}

// getFunctionKey generates a unique key for a function
func (u *goPureProjectGenerateUsecase) getFunctionKey(fn model.Function) string {
	if fn.Receiver != "" {
		return fmt.Sprintf("%s.%s.%s", fn.Package, fn.Receiver, fn.Name)
	}
	return fmt.Sprintf("%s.%s", fn.Package, fn.Name)
}

// buildFunctionSignature builds a full function signature like "func name(args) returnTypes"
func (u *goPureProjectGenerateUsecase) buildFunctionSignature(fn model.Function) string {
	var sig strings.Builder

	sig.WriteString("func ")

	// Add receiver if it's a method
	if fn.Receiver != "" {
		sig.WriteString("(")
		sig.WriteString(fn.Receiver)
		sig.WriteString(") ")
	}

	sig.WriteString(fn.Name)
	sig.WriteString("(")

	// Add parameters
	if len(fn.Parameters) > 0 {
		var params []string
		for _, p := range fn.Parameters {
			if p.Name != "" {
				params = append(params, fmt.Sprintf("%s %s", p.Name, p.Type))
			} else {
				params = append(params, p.Type)
			}
		}
		sig.WriteString(strings.Join(params, ", "))
	}

	sig.WriteString(")")

	// Add return types
	if len(fn.ReturnTypes) > 0 {
		if len(fn.ReturnTypes) == 1 {
			sig.WriteString(" ")
			sig.WriteString(fn.ReturnTypes[0])
		} else {
			sig.WriteString(" (")
			sig.WriteString(strings.Join(fn.ReturnTypes, ", "))
			sig.WriteString(")")
		}
	}

	return sig.String()
}

// getRelativePath converts absolute path to relative path from current directory
func (u *goPureProjectGenerateUsecase) getRelativePath(absPath string) string {
	// Get current working directory
	cwd, err := filepath.Abs(".")
	if err != nil {
		fmt.Printf("Warning: failed to get current directory: %v\n", err)
		return absPath
	}

	// Try to get relative path from current directory
	if relPath, err := filepath.Rel(cwd, absPath); err == nil {
		return relPath
	}
	// If failed, return absolute path
	return absPath
}

// buildHierarchicalCallTree builds a hierarchical call tree structure from entry points
func (u *goPureProjectGenerateUsecase) buildHierarchicalCallTree(entryPoints []model.Function, allFunctions []model.Function, functionCalls map[string][]string, importMap map[string]string) []model.CallTreeNode {
	// Create function map for quick lookup
	funcMap := make(map[string]model.Function)
	for _, fn := range allFunctions {
		funcMap[u.getFunctionKey(fn)] = fn
	}

	var callTreeNodes []model.CallTreeNode

	// Build tree for each entry point
	for _, ep := range entryPoints {
		visited := make(map[string]bool)
		node := u.buildTreeNodeRecursive(ep, funcMap, functionCalls, importMap, visited, 0, 10) // max depth 10
		callTreeNodes = append(callTreeNodes, node)
	}

	return callTreeNodes
}

// buildTreeNodeRecursive recursively builds a call tree node
func (u *goPureProjectGenerateUsecase) buildTreeNodeRecursive(fn model.Function, funcMap map[string]model.Function, functionCalls map[string][]string, importMap map[string]string, visited map[string]bool, depth int, maxDepth int) model.CallTreeNode {
	funcKey := u.getFunctionKey(fn)

	// Build full function signature for title
	fullSignature := u.buildFunctionSignature(fn)

	// Get relative path
	relativePath := u.getRelativePath(fn.File)

	node := model.CallTreeNode{
		Title:       fullSignature,
		Name:        fn.Name,
		Package:     fn.Package,
		File:        relativePath,
		Line:        fn.Line,
		Kind:        fn.Kind,
		Receiver:    fn.Receiver,
		Parameters:  fn.Parameters,
		ReturnTypes: fn.ReturnTypes,
	}

	// Check for circular reference
	if visited[funcKey] {
		node.IsRecursive = true
		return node
	}

	// Stop if max depth reached
	if depth >= maxDepth {
		return node
	}

	// Mark as visited
	visited[funcKey] = true
	defer func() { visited[funcKey] = false }()

	// Get called functions
	calls, ok := functionCalls[funcKey]
	if !ok || len(calls) == 0 {
		return node
	}

	// Build child nodes
	for _, calledFuncName := range calls {
		// Try to find the function in funcMap
		var childFn model.Function
		found := false

		// Try exact match first
		if fn, exists := funcMap[calledFuncName]; exists {
			childFn = fn
			found = true
		} else {
			// Try partial match (simple function name)
			for key, fn := range funcMap {
				if strings.HasSuffix(key, "."+calledFuncName) || fn.Name == calledFuncName {
					childFn = fn
					found = true
					break
				}
			}
		}

		if found {
			childNode := u.buildTreeNodeRecursive(childFn, funcMap, functionCalls, importMap, visited, depth+1, maxDepth)
			node.Children = append(node.Children, childNode)
		} else {
			// Create a placeholder node for external or unresolved functions
			// Extract package name and lookup in importMap
			packageName := u.extractPackageFromFunctionName(calledFuncName)
			packagePath := ""
			if packageName != "" {
				if path, ok := importMap[packageName]; ok {
					packagePath = path
				}
			}

			node.Children = append(node.Children, model.CallTreeNode{
				Title:       calledFuncName + "()",
				Name:        calledFuncName,
				Package:     packageName,
				PackagePath: packagePath,
				Kind:        "external",
				File:        "",
			})
		}
	}

	return node
}

// buildCallTreeVisualization creates a text visualization of the call tree
func (u *goPureProjectGenerateUsecase) buildCallTreeVisualization(nodes []model.CallTreeNode) string {
	var result strings.Builder
	result.WriteString("Call Tree:\n")
	result.WriteString("==========\n\n")

	for i, node := range nodes {
		if i > 0 {
			result.WriteString("\n")
		}
		result.WriteString(fmt.Sprintf("Entry Point %d: %s\n", i+1, node.Title))
		u.visualizeNodeRecursive(&result, node, "", true)
	}

	return result.String()
}

// visualizeNodeRecursive recursively visualizes a call tree node
func (u *goPureProjectGenerateUsecase) visualizeNodeRecursive(result *strings.Builder, node model.CallTreeNode, prefix string, isLast bool) {
	// Print current node
	if prefix != "" {
		connector := "├─"
		if isLast {
			connector = "└─"
		}
		result.WriteString(prefix)
		result.WriteString(connector)
		result.WriteString(" ")
	}

	// Node info - use Title for full function signature
	result.WriteString(node.Title)
	if node.File != "" {
		result.WriteString(fmt.Sprintf(" (%s:%d)", node.File, node.Line))
	}
	if node.IsRecursive {
		result.WriteString(" [recursive]")
	}
	if node.Kind == "external" {
		result.WriteString(" [external]")
	}
	result.WriteString("\n")

	// Print children
	for i, child := range node.Children {
		childIsLast := i == len(node.Children)-1
		childPrefix := prefix
		if prefix != "" {
			if isLast {
				childPrefix += "   "
			} else {
				childPrefix += "│  "
			}
		}
		u.visualizeNodeRecursive(result, child, childPrefix, childIsLast)
	}
}

// extractPackageFromFunctionName extracts the package name from a qualified function name
// Example: "controlplaneapiserver.BuildGenericConfig" -> "controlplaneapiserver"
func (u *goPureProjectGenerateUsecase) extractPackageFromFunctionName(functionName string) string {
	lastDot := strings.LastIndex(functionName, ".")
	if lastDot > 0 {
		return functionName[:lastDot]
	}
	return ""
}
