package model

// CallTree represents the entire call tree structure
type CTree struct {
	SourceFile string     `yaml:"source_file"`
	Language   string     `yaml:"language"`
	Functions  []Function `yaml:"functions"`
	CallGraph  []CallEdge `yaml:"call_graph"`
}

// Function represents a function or method in the source code
type Function struct {
	Name      string   `yaml:"name"`
	File      string   `yaml:"file"`
	Line      int      `yaml:"line"`
	Kind      string   `yaml:"kind"`            // function, method, class, etc.
	Signature string   `yaml:"signature"`       // function signature
	Class     string   `yaml:"class,omitempty"` // class name if it's a method
	Namespace string   `yaml:"namespace,omitempty"`
	Access    string   `yaml:"access,omitempty"` // public, private, protected
	CallsTo   []string `yaml:"calls_to,omitempty"`
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
