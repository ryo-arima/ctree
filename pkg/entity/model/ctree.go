package model

// CallTree represents the entire call tree structure
type CTree struct {
	SourceFile            string                 `yaml:"source_file"`
	Language              string                 `yaml:"language"`
	Functions             []Function             `yaml:"functions,omitempty"`
	CallGraph             []CallEdge             `yaml:"call_graph,omitempty"`
	EntryPoints           []Function             `yaml:"entry_points,omitempty"`
	CallTree              []CallTreeNode         `yaml:"call_tree,omitempty"`
	CallTreeVisualization string                 `yaml:"call_tree_visualization,omitempty"`
	ImportMap             map[string]string      `yaml:"import_map,omitempty"` // package name -> full import path
	Metadata              map[string]interface{} `yaml:"metadata,omitempty"`
}

// CallTreeNode represents a node in the hierarchical call tree
type CallTreeNode struct {
	Title       string         `yaml:"title"`
	Name        string         `yaml:"name,omitempty"`
	Package     string         `yaml:"package,omitempty"`
	PackagePath string         `yaml:"package_path,omitempty"` // Full import path for external packages
	File        string         `yaml:"file"`
	Line        int            `yaml:"line"`
	Kind        string         `yaml:"kind,omitempty"`
	Receiver    string         `yaml:"receiver,omitempty"`
	Signature   string         `yaml:"signature,omitempty"`
	Parameters  []Parameter    `yaml:"parameters,omitempty"`
	ReturnTypes []string       `yaml:"return_types,omitempty"`
	Children    []CallTreeNode `yaml:"children,omitempty"`
	IsRecursive bool           `yaml:"is_recursive,omitempty"`
}

// Function represents a function or method in the source code
type Function struct {
	Name        string      `yaml:"name"`
	File        string      `yaml:"file"`
	Line        int         `yaml:"line"`
	Kind        string      `yaml:"kind"`            // function, method, class, etc.
	Signature   string      `yaml:"signature"`       // function signature
	Class       string      `yaml:"class,omitempty"` // class name if it's a method
	Namespace   string      `yaml:"namespace,omitempty"`
	Access      string      `yaml:"access,omitempty"` // public, private, protected
	CallsTo     []string    `yaml:"calls_to,omitempty"`
	Package     string      `yaml:"package,omitempty"`      // Go package name
	Receiver    string      `yaml:"receiver,omitempty"`     // Go method receiver
	Parameters  []Parameter `yaml:"parameters,omitempty"`   // Function parameters
	ReturnTypes []string    `yaml:"return_types,omitempty"` // Return types
}

// Parameter represents a function parameter
type Parameter struct {
	Name string `yaml:"name,omitempty"`
	Type string `yaml:"type"`
}

// CallEdge represents a call relationship between functions
type CallEdge struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
	File string `yaml:"file"`
	Line int    `yaml:"line"`
}

// Tag represents a ctags tag entry
type Tag struct {
	Name      string
	File      string
	Address   string
	Kind      string
	Language  string
	Line      int
	Signature string
	Class     string
	Namespace string
	Access    string
}
