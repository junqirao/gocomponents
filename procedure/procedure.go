package procedure

import (
	"context"
	"text/template"
)

// NodeLifeCycle defines node life cycle
type NodeLifeCycle interface {
	// BeforeExecute defines aspect before execute
	// input is copy of Execute input and filtered.
	BeforeExecute(ctx context.Context, node *Node, input map[string]any) (err error)
	// Execute defines how to execute a node.
	// input is copy of Execute input and
	// filtered. after execute, output will be
	// saved to Execute output in key node.Name
	Execute(ctx context.Context, node *Node, input map[string]any) (output any, err error)
	// AfterExecute defines aspect after execute.
	AfterExecute(ctx context.Context, node *Node, output any, err error)
	// BeforeScript defines parameters to build the
	// template.Template, returns the template.FuncMap
	// and the template execute input.
	// this method called after AfterExecute only if
	// node.Script is not empty.
	BeforeScript(ctx context.Context, node *Node, input map[string]any, output any) (fm template.FuncMap)
}
