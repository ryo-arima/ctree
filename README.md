# ctree

> **⚠️ ALPHA VERSION**: This project is under active development. Go call tree generation is functional, but other languages are still in progress.

A command-line tool for generating and visualizing call trees from source code. `ctree` analyzes your codebase using AST parsing and outputs structured call trees in YAML format with beautiful tree visualization, making it easier to understand code dependencies and execution flows.

**Current Status**: ✅ Go fully supported | 🚧 C, C++, Rust, Python in progress

**Supported Languages**: Go • C • C++ • Rust • Python

## Features

- **Multi-language Support**: Generate call trees for Go, C, C++, Rust, and Python
- **Advanced Go AST Parsing**: Full function signature extraction with parameters and return types
- **Import Path Resolution**: Shows complete import paths for external packages (e.g., `k8s.io/kubernetes/pkg/controlplane/apiserver`)
- **Call Tree Visualization**: 
  - Hierarchical tree structure with indentation
  - Color-coded output (internal/external functions)
  - `[internal]` and `[external]` tags with file paths/package names
  - Entry point detection (main, init functions)
- **Flexible Display Options**:
  - `--expand-signature`: Show function parameters and return values on separate lines
  - Text/YAML output formats
- **Language-Specific Features**: 
  - Go: Full AST parsing with import resolution
  - C/C++: Header file tracking
  - Rust: Trait and macro analysis
  - Python: Module and decorator support
- **Recursive Analysis**: Automatically traverse directory structures
- **Configurable Depth**: Control analysis depth to focus on relevant code

## Installation

```bash
# Clone the repository
git clone https://github.com/ryo-arima/ctree.git
cd ctree

# Build the binary
go build -o ctree ./cmd/main.go

# Install to your PATH (optional)
sudo mv ctree /usr/local/bin/
```

## Usage

### Generate Call Tree

Generate a call tree for your source code:

```bash
# Generate for Go project
ctree generate golang --source ./myproject --output call-tree.yaml

# Generate for Kubernetes apiserver
ctree generate golang --source ./kubernetes/cmd/kube-apiserver --output apiserver-tree.yaml

# Generate for C project
ctree generate c --source ./myapp --output c-tree.yaml

# Generate for C++ project
ctree generate cpp --source ./myapp --output cpp-tree.yaml

# Generate for Rust project
ctree generate rust --source ./myapp --output rust-tree.yaml

# Generate for Python project
ctree generate python --source ./myapp --output python-tree.yaml
```

### View Call Tree

Extract and visualize call tree from generated YAML:

```bash
# View as text with tree structure
ctree get golang call-tree --ctree call-tree.yaml --format text

# View with expanded function signatures
ctree get golang call-tree --ctree call-tree.yaml --format text --expand-signature

# View as YAML
ctree get golang call-tree --ctree call-tree.yaml --format yaml
```

### Command Options

#### Generate Command
- `--source, -s`: Source directory or file to analyze (default: current directory)
- `--output, -o`: Output file path (default: stdout)
- `--framework`: Framework to use (pure, react, django, flask, etc.)
- `--recursive, -r`: Recursively analyze subdirectories (default: true)
- `--max-depth, -d`: Maximum depth for recursive analysis (default: 10)

#### Get Call-Tree Command
- `--ctree, -c`: Path to ctree YAML file (required)
- `--format`: Output format (yaml, text) (default: yaml)
- `--expand-signature`: Show function parameters and return values on separate lines
- `--output, -o`: Output file path (default: stdout)

### Examples

```bash
# Analyze current directory
ctree generate golang

# Analyze specific directory and save to file
ctree generate c --source ./src --output call-tree.yaml

# Analyze C++ project
ctree generate cpp --source ./app --output cpp-tree.yaml

# Analyze Rust project
ctree generate rust --source ./app --output rust-tree.yaml

# Control recursion depth
ctree generate golang --source ./pkg --max-depth 5

# View call tree with color-coded output
ctree get golang call-tree --ctree apiserver-tree.yaml --format text
```

### Output Example

Text format with tree structure:

```
Call Tree:
==========

Entry Point 1: func main() [internal] (cmd/apiserver.go:32)
  ├─ func NewAPIServerCommand() *cobra.Command [internal] (app/server.go:70)
  │  ├─ func NewServerRunOptions() *ServerRunOptions [internal] (app/options/options.go:66)
  │  │  ├─ controlplaneapiserver.NewOptions() [external] (k8s.io/kubernetes/pkg/controlplane/apiserver)
  │  │  ├─ time.Duration() [external] (time)
  │  │  └─ append() [external]
  │  ├─ genericapiserver.SetupSignalContext() [external] (k8s.io/apiserver/pkg/server)
  │  └─ func Run(ctx context.Context, opts options.CompletedOptions) error [internal] (app/server.go:148)
  ├─ cli.Run() [external] (k8s.io/component-base/cli)
  └─ os.Exit() [external] (os)
```

With `--expand-signature` flag:

```
Entry Point 1: func main [internal] (cmd/apiserver.go:32)
  Parameters:
    (none)
  Returns:
    (none)
  ├─ func NewAPIServerCommand [internal] (app/server.go:70)
  │  Returns:
  │    - *cobra.Command
```

## Supported Languages

### Go ✅
- Pure Go projects
- Full AST parsing with go/parser
- Function signature extraction (parameters, return types)
- Import path resolution for external packages
- Entry point detection (main, init)
- Call tree construction with parent-child relationships

### C 🚧
- In progress
- Planned features: Function call analysis, header file resolution

### C++ 🚧
- In progress
- Planned features: Class hierarchy analysis, template support

### Rust 🚧
- In progress
- Planned features: Trait resolution, macro expansion

### Python 🚧
- In progress
- Planned features: Import resolution, decorator support

## Configuration

You can configure `ctree` using a configuration file located at:
- `~/.config/ctree/app.toml` (Linux/macOS)
- `%APPDATA%\ctree\app.toml` (Windows)

Example configuration:

```toml
[app]
name = "ctree"
version = "0.1.0"
```

## Project Structure

```
ctree/
├── cmd/
│   └── main.go           # Entry point
├── pkg/
│   ├── base.go           # CLI base commands
│   ├── config/           # Configuration management
│   ├── controller/       # Command handlers
│   │   ├── golang/
│   │   ├── javascript/
│   │   └── python/
│   ├── usecase/          # Business logic
│   │   ├── golang/
│   │   ├── javascript/
│   │   └── python/
│   ├── repository/       # Data access (ctags integration)
│   └── entity/           # Data models
│       ├── model/
│       ├── request/
│       └── response/
└── etc/
    └── app.toml          # Default configuration
```

## Architecture

`ctree` follows Clean Architecture principles:

- **Controller Layer**: Handles CLI commands and user input
- **Usecase Layer**: Contains business logic for code analysis
- **Repository Layer**: Interfaces with code parsing infrastructure
- **Entity Layer**: Defines data structures and models

## Development Status

### Completed ✅
- Project structure and architecture
- CLI framework with Cobra
- Multi-language command structure (generate/get/list)
- Configuration management
- **Go call tree generation**:
  - Full AST parsing with function signatures
  - Entry point detection (main, init)
  - Call graph construction
  - Import path resolution
- **Output formats**:
  - YAML with hierarchical structure
  - Text with tree visualization
  - Color-coded terminal output
- **Display features**:
  - [internal]/[external] function tags
  - File paths and line numbers
  - Full package import paths
  - Expandable function signatures

### In Progress 🚧
- C call tree generation with header file resolution
- C++ call tree generation with class hierarchy
- Rust call tree generation with trait resolution
- Python call tree generation with import analysis
- Additional get subcommands (functions, classes, variables)

### Planned 📋
- Language-specific optimizations
- List command implementations
- Advanced filtering and query capabilities
- Graph visualization output (DOT, Mermaid)
- IDE integration (VS Code extension)

## Requirements

- Go 1.25+ for building

## Contributing

Contributions are welcome! This project is in early development, so there are many opportunities to contribute.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Roadmap

### Phase 1: Core Functionality (Current)
- [x] Go call tree generation with AST parsing
- [x] Import path resolution
- [x] Tree visualization with colors
- [x] YAML and text output formats
- [ ] C call tree generation
- [ ] C++ call tree generation
- [ ] Rust call tree generation
- [ ] Python call tree generation

### Phase 2: Language-Specific Features
- [ ] C: Header file resolution, function pointer tracking
- [ ] C++: Class hierarchy, template instantiation, namespace resolution
- [ ] Rust: Trait resolution, macro expansion, lifetime analysis
- [ ] Python: Import resolution, decorator support, type hints
- [ ] Advanced filtering options (by package, depth, pattern)
- [ ] Query capabilities (find function, trace call path)

### Phase 3: Enhanced Commands
- [ ] Additional get commands (functions, classes, variables, imports)
- [ ] List commands for overview and statistics
- [ ] Diff command to compare call trees
- [ ] Search command with pattern matching

### Phase 4: Visualization & Integration
- [ ] Graph visualization output (DOT, SVG, Mermaid)
- [ ] Interactive HTML visualization
- [ ] VS Code extension
- [ ] Language Server Protocol support
- [ ] Performance optimization for large codebases
- [ ] Caching mechanism for incremental analysis

## Author

[@ryo-arima](https://github.com/ryo-arima)

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
