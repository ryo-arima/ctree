# ctree

> **⚠️ WARNING: This project is currently under active development and is not yet ready for production use.**

A command-line tool for generating call trees from source code. `ctree` analyzes your codebase and outputs structured call trees in YAML format, making it easier to understand code dependencies and execution flows.

## Features

- **Multi-language Support**: Generate call trees for Go, JavaScript, Python, and more
- **Framework-Specific Analysis**: Special handling for popular frameworks:
  - JavaScript: React
  - Python: Django, Flask
- **Flexible Output**: YAML format for easy parsing and integration
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
ctree generate golang --source ./myproject --framework pure

# Generate for JavaScript/React project
ctree generate javascript --source ./myapp --framework react

# Generate for Python project
ctree generate python --source ./myapp --framework pure
```

### Command Options

- `--source, -s`: Source directory or file to analyze (default: current directory)
- `--output, -o`: Output file path (default: stdout)
- `--framework, -f`: Framework to use (pure, react, django, flask, etc.)
- `--recursive, -r`: Recursively analyze subdirectories (default: true)
- `--max-depth, -d`: Maximum depth for recursive analysis (default: 10)

### Examples

```bash
# Analyze current directory
ctree generate golang

# Analyze specific directory and save to file
ctree generate javascript --source ./src --output call-tree.yaml

# Analyze with specific framework
ctree generate python --source ./app --framework django

# Control recursion depth
ctree generate golang --source ./pkg --max-depth 5
```

## Supported Languages and Frameworks

### Go
- Pure Go projects

### JavaScript
- Pure JavaScript
- React

### Python
- Pure Python
- Django (coming soon)
- Flask (coming soon)

### Coming Soon
- Java
- C/C++
- Ruby
- Rust
- PHP

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

- ✅ Project structure and architecture
- ✅ CLI framework with Cobra
- ✅ Multi-language command structure (generate/get/list)
- ✅ Configuration management
- 🚧 Call tree generation logic (in progress)
- 📋 Framework-specific analyzers (planned)
- 📋 Output formatting options (planned)
- 📋 Filtering and query capabilities (planned)

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

- [ ] Implement call tree generation logic
- [ ] Add support for more languages (Java, C/C++, Ruby, Rust, PHP)
- [ ] Framework-specific analyzers (Django, Flask, Spring, Rails)
- [ ] Advanced filtering and query options
- [ ] Graph visualization output
- [ ] Integration with IDEs
- [ ] Performance optimization for large codebases

## Author

[@ryo-arima](https://github.com/ryo-arima)

## Acknowledgments

- [Cobra](https://github.com/spf13/cobra) - CLI framework
