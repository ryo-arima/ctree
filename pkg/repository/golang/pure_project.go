package golang

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/ryo-arima/ctree/pkg/entity/model"
)

// GoPureProjectRepository handles Go pure project file operations
type GoPureProjectRepository interface {
	FindGoFiles(sourcePath string, recursive bool, maxDepth int) ([]string, error)
	ParseGoFile(filePath string) (*ast.File, *token.FileSet, error)
	ExtractFunctions(file *ast.File, fset *token.FileSet, filePath string) ([]model.Function, error)
	ExtractImports(file *ast.File) map[string]string // alias/name -> full import path
}

type goPureProjectRepository struct {
}

// NewGoPureProjectRepository creates a new Go pure project repository
func NewGoPureProjectRepository() GoPureProjectRepository {
	return &goPureProjectRepository{}
}

// FindGoFiles finds all Go files in the specified path
func (r *goPureProjectRepository) FindGoFiles(sourcePath string, recursive bool, maxDepth int) ([]string, error) {
	var goFiles []string

	// Get absolute path
	absPath, err := filepath.Abs(sourcePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Check if path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access path: %w", err)
	}

	// If it's a single file
	if !info.IsDir() {
		if strings.HasSuffix(absPath, ".go") && !strings.HasSuffix(absPath, "_test.go") {
			return []string{absPath}, nil
		}
		return []string{}, nil
	}

	// Walk directory
	err = r.walkDir(absPath, absPath, 0, maxDepth, recursive, &goFiles)
	if err != nil {
		return nil, err
	}

	return goFiles, nil
}

// walkDir recursively walks through directories
func (r *goPureProjectRepository) walkDir(basePath, currentPath string, currentDepth, maxDepth int, recursive bool, goFiles *[]string) error {
	if currentDepth > maxDepth {
		return nil
	}

	entries, err := os.ReadDir(currentPath)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", currentPath, err)
	}

	for _, entry := range entries {
		fullPath := filepath.Join(currentPath, entry.Name())

		// Skip hidden files and directories
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		// Skip vendor and test directories
		if entry.IsDir() && (entry.Name() == "vendor" || entry.Name() == "testdata") {
			continue
		}

		if entry.IsDir() {
			if recursive {
				if err := r.walkDir(basePath, fullPath, currentDepth+1, maxDepth, recursive, goFiles); err != nil {
					return err
				}
			}
		} else if strings.HasSuffix(entry.Name(), ".go") && !strings.HasSuffix(entry.Name(), "_test.go") {
			*goFiles = append(*goFiles, fullPath)
		}
	}

	return nil
}

// ParseGoFile parses a Go source file
func (r *goPureProjectRepository) ParseGoFile(filePath string) (*ast.File, *token.FileSet, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}
	return file, fset, nil
}

// ExtractFunctions extracts function information from AST
func (r *goPureProjectRepository) ExtractFunctions(file *ast.File, fset *token.FileSet, filePath string) ([]model.Function, error) {
	var functions []model.Function

	ast.Inspect(file, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.FuncDecl:
			fn := model.Function{
				Name:    x.Name.Name,
				File:    filePath,
				Line:    fset.Position(x.Pos()).Line,
				Package: file.Name.Name,
				Kind:    "function",
			}

			// Extract receiver type for methods
			if x.Recv != nil && len(x.Recv.List) > 0 {
				if starExpr, ok := x.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						fn.Receiver = ident.Name
					}
				} else if ident, ok := x.Recv.List[0].Type.(*ast.Ident); ok {
					fn.Receiver = ident.Name
				}
			}

			// Extract parameters
			if x.Type.Params != nil {
				for _, param := range x.Type.Params.List {
					paramType := formatType(param.Type)
					if len(param.Names) > 0 {
						for _, name := range param.Names {
							fn.Parameters = append(fn.Parameters, model.Parameter{
								Name: name.Name,
								Type: paramType,
							})
						}
					} else {
						fn.Parameters = append(fn.Parameters, model.Parameter{
							Type: paramType,
						})
					}
				}
			}

			// Extract return types
			if x.Type.Results != nil {
				for _, result := range x.Type.Results.List {
					returnType := formatType(result.Type)
					fn.ReturnTypes = append(fn.ReturnTypes, returnType)
				}
			}

			functions = append(functions, fn)
		}
		return true
	})

	return functions, nil
}

// formatType formats an AST type expression to string
func formatType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + formatType(t.X)
	case *ast.ArrayType:
		return "[]" + formatType(t.Elt)
	case *ast.MapType:
		return "map[" + formatType(t.Key) + "]" + formatType(t.Value)
	case *ast.SelectorExpr:
		return formatType(t.X) + "." + t.Sel.Name
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.ChanType:
		return "chan " + formatType(t.Value)
	case *ast.FuncType:
		return "func"
	case *ast.Ellipsis:
		return "..." + formatType(t.Elt)
	default:
		return "unknown"
	}
}

// ExtractImports extracts import information from AST
// Returns a map of package alias/name -> full import path
func (r *goPureProjectRepository) ExtractImports(file *ast.File) map[string]string {
	imports := make(map[string]string)

	for _, imp := range file.Imports {
		if imp.Path == nil {
			continue
		}

		// Remove quotes from import path
		importPath := strings.Trim(imp.Path.Value, "\"")

		// Get the package name/alias
		var packageName string
		if imp.Name != nil {
			// Named import: alias "path"
			packageName = imp.Name.Name
		} else {
			// Default import: use the last part of the path
			parts := strings.Split(importPath, "/")
			packageName = parts[len(parts)-1]
		}

		imports[packageName] = importPath
	}

	return imports
}
